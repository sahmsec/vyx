package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/sahmsec/vyx/internal/model"
	"github.com/sahmsec/vyx/internal/scan"
)

const Version = "0.0.1-dev"

func Run(ctx context.Context, args []string, out, diagnostics io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(out, "Vyx — Phase 1 feasibility prototype\n\nUsage:\n  vyx scan [directory] [--format text|json] [--max-file-bytes N] [--max-files N]\n  vyx version\n\nImplemented: embedded secret detection and JS/TS syntax validation only.\nDependency and framework security checks are not implemented. Scans make no network requests.")
		return 0
	}
	if args[0] == "version" {
		fmt.Fprintln(out, "vyx "+Version)
		return 0
	}
	if args[0] != "scan" {
		fmt.Fprintln(diagnostics, "Unknown command; run vyx --help.")
		return 2
	}
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	format := flags.String("format", "text", "text or json")
	options := scan.DefaultOptions()
	flags.Int64Var(&options.MaxFileBytes, "max-file-bytes", options.MaxFileBytes, "maximum bytes per file")
	flags.IntVar(&options.MaxFiles, "max-files", options.MaxFiles, "maximum files inspected")
	// Accept the advertised scan . --format json order as well as flags before
	// the path. Consume flag values even when they begin with a dash.
	var positionals, flagArgs []string
	for i := 1; i < len(args); i++ {
		a := args[i]
		if a == "--format" || a == "-format" || a == "--max-file-bytes" || a == "-max-file-bytes" || a == "--max-files" || a == "-max-files" {
			flagArgs = append(flagArgs, a)
			if i+1 < len(args) {
				i++
				flagArgs = append(flagArgs, args[i])
			}
		} else if len(a) > 0 && a[0] == '-' {
			flagArgs = append(flagArgs, a)
		} else {
			positionals = append(positionals, a)
		}
	}
	if err := flags.Parse(flagArgs); err != nil || len(positionals) > 1 || options.MaxFileBytes <= 0 || options.MaxFileBytes > 64<<20 || options.MaxFiles <= 0 || options.MaxFiles > 1_000_000 {
		fmt.Fprintln(diagnostics, "Invalid scan arguments; run vyx --help.")
		return 2
	}
	if *format != "text" && *format != "json" {
		fmt.Fprintln(diagnostics, "Unsupported format; use text or json. SARIF and agent output are planned.")
		return 2
	}
	path := "."
	if len(positionals) == 1 {
		path = positionals[0]
	}
	report, err := scan.Scan(ctx, path, options)
	if err != nil {
		fmt.Fprintln(diagnostics, err.Error())
		return 2
	}
	if *format == "json" {
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(report); err != nil {
			fmt.Fprintln(diagnostics, "Cannot write report.")
			return 2
		}
	} else if err := writeText(out, report); err != nil {
		fmt.Fprintln(diagnostics, "Cannot write report.")
		return 2
	}
	return report.ExitCode()
}

func writeText(out io.Writer, r model.Report) error {
	status := "completed"
	if !r.Complete {
		status = "incomplete"
	}
	if _, err := fmt.Fprintf(out, "Vyx feasibility scan: %s\n%d files scanned; %d JS/TS files parsed; %d credential candidates.\n", status, r.Coverage.FilesScanned, r.Coverage.SourcesParsed, len(r.Findings)); err != nil {
		return err
	}
	for _, f := range r.Findings {
		if _, err := fmt.Fprintf(out, "\n%s [%s confidence] %s (%s)\n  %s:%d:%d\n  %s\n  %s\n  Fix: %s\n", f.Severity, f.Confidence, f.RuleID, f.UpstreamRuleID, strconv.QuoteToASCII(f.Path), f.Line, f.Column, f.Explanation, f.Evidence, f.Remediation); err != nil {
			return err
		}
	}
	for _, d := range r.Diagnostics {
		if _, err := fmt.Fprintf(out, "\n%s %s: %s\n", d.Code, strconv.QuoteToASCII(d.Path), d.Message); err != nil {
			return err
		}
	}
	if len(r.Exclusions) > 0 {
		if _, err := fmt.Fprintf(out, "\n%d paths excluded (use JSON for the complete list).\n", len(r.Exclusions)); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(out, "\nScope: secrets and parser validation only. Dependency/framework security checks are not implemented.\nNo conclusion about overall application security.")
	return err
}
