package scan

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/sahmsec/vyx/internal/model"
	"github.com/sahmsec/vyx/internal/secrets"
	"github.com/sahmsec/vyx/internal/syntax"
)

type Options struct {
	MaxFileBytes int64
	MaxFiles     int
}

func DefaultOptions() Options { return Options{MaxFileBytes: 1 << 20, MaxFiles: 10_000} }

var excludedDirs = map[string]bool{
	".git": true, ".tools": true, "node_modules": true, ".next": true,
	"dist": true, "vendor": true, "coverage": true,
}

var errFileLimit = errors.New("file limit reached")

func Scan(ctx context.Context, path string, options Options) (model.Report, error) {
	report := model.Report{
		SchemaVersion: "0.1", Stage: "phase-1-feasibility", Complete: true,
		Checks:   []string{"secrets", "js-ts-syntax-validation"},
		Findings: []model.Finding{}, Diagnostics: []model.Diagnostic{}, Exclusions: []model.Exclusion{},
		Limitations: []string{
			"Only embedded secret detection and JS/TS syntax validation are implemented.",
			"Framework security rules, dependency advisories, suppressions, history, and incremental scans are not implemented.",
			"Syntax validation is a parser feasibility check, not a security audit or TypeScript type check.",
			"Gitignore files are not interpreted yet; environment files are scanned. See explicit exclusions.",
		},
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		return report, errors.New("cannot open scan directory")
	}
	defer root.Close()
	detector, err := secrets.New()
	if err != nil {
		return report, err
	}
	parser := syntax.New()
	visited := 0
	err = fs.WalkDir(root.FS(), ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if walkErr != nil {
			report.Complete = false
			report.Diagnostics = append(report.Diagnostics, model.Diagnostic{Code: "read_failed", Path: path, Message: "Cannot inspect path."})
			return nil
		}
		if entry.IsDir() {
			if path != "." && excludedDirs[strings.ToLower(entry.Name())] {
				report.Exclusions = append(report.Exclusions, model.Exclusion{Path: path, Reason: "generated, dependency, tooling, or VCS directory"})
				return fs.SkipDir
			}
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 || !entry.Type().IsRegular() {
			report.Exclusions = append(report.Exclusions, model.Exclusion{Path: path, Reason: "symlink or nonregular file; not followed"})
			return nil
		}
		visited++
		if visited > options.MaxFiles {
			return errFileLimit
		}
		file, openErr := root.Open(path)
		if openErr != nil {
			report.Complete = false
			report.Diagnostics = append(report.Diagnostics, model.Diagnostic{Code: "read_failed", Path: path, Message: "Cannot read file within scan root."})
			return nil
		}
		data, readErr := io.ReadAll(io.LimitReader(file, options.MaxFileBytes+1))
		closeErr := file.Close()
		if readErr != nil || closeErr != nil {
			report.Complete = false
			report.Diagnostics = append(report.Diagnostics, model.Diagnostic{Code: "read_failed", Path: path, Message: "Cannot read complete file."})
			return nil
		}
		if int64(len(data)) > options.MaxFileBytes {
			report.Complete = false
			report.Diagnostics = append(report.Diagnostics, model.Diagnostic{Code: "file_limit", Path: path, Message: "File exceeds configured size limit; not scanned."})
			return nil
		}
		if bytes.ContainsRune(data, 0) || !utf8.Valid(data) {
			report.Exclusions = append(report.Exclusions, model.Exclusion{Path: path, Reason: "binary or non-UTF-8 content"})
			if syntax.Language(path) != "" {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, model.Diagnostic{Code: "unsupported_encoding", Path: path, Message: "Source encoding is unsupported."})
			}
			return nil
		}
		report.Coverage.FilesScanned++
		report.Coverage.BytesScanned += int64(len(data))
		report.Findings = append(report.Findings, detector.Scan(ctx, path, data)...)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if syntax.Language(path) != "" {
			_, parseErr := parser.Parse(path, data)
			if parseErr != nil {
				report.Complete = false
				report.Diagnostics = append(report.Diagnostics, model.Diagnostic{Code: "parse_failed", Path: path, Message: parseErr.Error()})
			} else {
				report.Coverage.SourcesParsed++
			}
		}
		return nil
	})
	if err != nil {
		report.Complete = false
		code, message := "scan_failed", "File discovery stopped before completion."
		if errors.Is(err, errFileLimit) {
			code, message = "file_count_limit", "Configured file count limit reached."
		}
		if ctx.Err() != nil {
			code, message = "cancelled", "Scan cancelled before completion."
		}
		report.Diagnostics = append(report.Diagnostics, model.Diagnostic{Code: code, Message: message})
	}
	sort.Slice(report.Findings, func(i, j int) bool {
		a, b := report.Findings[i], report.Findings[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Column != b.Column {
			return a.Column < b.Column
		}
		return a.UpstreamRuleID < b.UpstreamRuleID
	})
	return report, nil
}
