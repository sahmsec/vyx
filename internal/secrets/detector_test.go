package secrets

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
)

func testToken(label string) string {
	hash := sha256.Sum256([]byte("vyx synthetic test credential: " + label))
	return "ghp_" + hex.EncodeToString(hash[:])[:36]
}

func TestCredentialEvidenceIsRedacted(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	first, second := testToken("one"), testToken("two")
	source := "const token = '" + first + "'; const other = '" + second + "';"
	findings := d.Scan(context.Background(), "src/server.ts", []byte(source))
	if len(findings) == 0 {
		t.Fatal("synthetic provider tokens were not detected")
	}
	encoded, err := json.Marshal(findings)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), first) || strings.Contains(string(encoded), second) || strings.Contains(string(encoded), source) {
		t.Fatal("raw secret/source leaked into normalized findings")
	}
	for _, f := range findings {
		if f.Line != 1 || f.Column < 1 || f.Path != "src/server.ts" || f.UpstreamRuleID == "" {
			t.Fatalf("missing provenance/location: %+v", f)
		}
	}
}

func TestSafePublicConfiguration(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	source := []byte("NEXT_PUBLIC_SUPABASE_URL=https://example.invalid\nNEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY=sb_publishable_example\nconst password = process.env.PASSWORD;\n")
	if findings := d.Scan(context.Background(), "config.ts", source); len(findings) != 0 {
		t.Fatalf("safe configuration produced %d findings", len(findings))
	}
}

func TestSourceCommentCannotSilentlySuppress(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	findings := d.Scan(context.Background(), "server.ts", []byte("const token = '"+testToken("allow-comment")+"'; // gitleaks:allow"))
	if len(findings) == 0 {
		t.Fatal("upstream suppression comment unexpectedly hid credential")
	}
}

func TestMultilineLocation(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	findings := d.Scan(context.Background(), "server.ts", []byte("// header\n\nconst token = '"+testToken("third-line")+"';\n"))
	if len(findings) == 0 {
		t.Fatal("synthetic credential not found")
	}
	for _, f := range findings {
		if f.Line != 3 {
			t.Fatalf("got line %d, want 3", f.Line)
		}
	}
}
