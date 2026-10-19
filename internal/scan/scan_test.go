package scan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestIncompleteScanPreservesFindingsAndEnvironmentFiles(t *testing.T) {
	root := t.TempDir()
	hash := sha256.Sum256([]byte("synthetic inventory credential"))
	writeFixture(t, root, ".env.local", "GITHUB_TOKEN=ghp_"+hex.EncodeToString(hash[:])[:36])
	writeFixture(t, root, ".gitignore", ".env.local\n")
	writeFixture(t, root, "broken.ts", "export const = {")
	writeFixture(t, root, "node_modules/ignored.ts", "export const = {")
	r, err := Scan(context.Background(), root, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if r.Complete || r.ExitCode() != 2 {
		t.Fatal("parser failure must take precedence over findings")
	}
	if len(r.Findings) == 0 || r.Findings[0].Path != ".env.local" {
		t.Fatal("ignored environment file credential not detected")
	}
	if r.Coverage.FilesScanned != 3 || len(r.Exclusions) != 1 {
		t.Fatalf("unexpected inventory: %+v", r.Coverage)
	}
}

func TestLimitsAreIncomplete(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "large.txt", "a file larger than the chosen budget")
	options := DefaultOptions()
	options.MaxFileBytes = 4
	r, err := Scan(context.Background(), root, options)
	if err != nil {
		t.Fatal(err)
	}
	if r.Complete || len(r.Diagnostics) != 1 || r.Diagnostics[0].Code != "file_limit" {
		t.Fatal("oversize file silently skipped")
	}
	options = DefaultOptions()
	options.MaxFiles = 1
	writeFixture(t, root, "second.txt", "safe")
	r, err = Scan(context.Background(), root, options)
	if err != nil {
		t.Fatal(err)
	}
	if r.Complete || r.Diagnostics[len(r.Diagnostics)-1].Code != "file_count_limit" {
		t.Fatal("file count limit silently skipped")
	}
}

func TestDoesNotFollowSymlinks(t *testing.T) {
	root, outside := t.TempDir(), t.TempDir()
	writeFixture(t, outside, "outside.ts", "export const = {")
	if err := os.Symlink(filepath.Join(outside, "outside.ts"), filepath.Join(root, "link.ts")); err != nil {
		t.Skip("OS does not permit symlink creation")
	}
	r, err := Scan(context.Background(), root, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if r.Coverage.FilesScanned != 0 || len(r.Exclusions) != 1 {
		t.Fatal("symlink followed")
	}
}

func TestCancellationIsIncomplete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r, err := Scan(ctx, t.TempDir(), DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if r.Complete || r.ExitCode() != 2 {
		t.Fatal("cancelled scan reported completion")
	}
}

func TestScannedConfigurationIsNeverExecuted(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "next.config.js", "require('node:fs').writeFileSync('SHOULD_NOT_EXIST', 'executed'); module.exports = {};")
	r, err := Scan(context.Background(), root, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if !r.Complete || r.Coverage.SourcesParsed != 1 {
		t.Fatal("configuration was not parsed")
	}
	if _, err := os.Stat(filepath.Join(root, "SHOULD_NOT_EXIST")); !os.IsNotExist(err) {
		t.Fatal("unexpected project side effect")
	}
}
