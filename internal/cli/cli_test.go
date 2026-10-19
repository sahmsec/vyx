package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sahmsec/vyx/internal/model"
)

func TestJSONAfterDirectoryAndIncompleteExit(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "broken.ts"), []byte("export const = {"), 0600); err != nil {
		t.Fatal(err)
	}
	var out, diagnostics bytes.Buffer
	code := Run(context.Background(), []string{"scan", root, "--format", "json"}, &out, &diagnostics)
	if code != 2 {
		t.Fatalf("got exit %d, want 2", code)
	}
	var report model.Report
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal("stdout is not a JSON report")
	}
	if report.Complete || diagnostics.Len() != 0 || len(report.Limitations) == 0 {
		t.Fatal("incomplete scope missing or machine output contaminated")
	}
}

func TestInvalidArgumentsDoNotEchoPotentialSecrets(t *testing.T) {
	for _, args := range [][]string{{"scan", "--format", "private-value"}, {"scan", "--unknown=private-value"}, {"scan", "--max-files", "-1"}, {"scan", "one", "two"}} {
		var out, diagnostics bytes.Buffer
		if code := Run(context.Background(), args, &out, &diagnostics); code != 2 {
			t.Fatalf("invalid arguments accepted: %d", code)
		}
		if strings.Contains(diagnostics.String(), "private-value") {
			t.Fatal("argument value echoed into diagnostics")
		}
	}
}

func TestEmptyDirectoryReportsLimitedScope(t *testing.T) {
	var out, diagnostics bytes.Buffer
	if code := Run(context.Background(), []string{"scan", t.TempDir()}, &out, &diagnostics); code != 0 {
		t.Fatalf("empty scan failed: %d", code)
	}
	if !strings.Contains(out.String(), "Dependency/framework security checks are not implemented") {
		t.Fatal("prototype presented as complete scanner")
	}
}
