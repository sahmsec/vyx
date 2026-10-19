// Package secrets adapts Gitleaks without invoking its CLI, source discovery,
// executable configuration, remote extensions, or Git commands.
package secrets

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	"github.com/zricethezav/gitleaks/v8/logging"

	"github.com/sahmsec/vyx/internal/model"
)

const Version = "v8.30.1"

var quietOnce sync.Once

type Detector struct{ detector *detect.Detector }

func New() (*Detector, error) {
	// Upstream logging is process-global. Disable it once before constructing
	// detectors: only our normalized, redacted model may leave this adapter.
	quietOnce.Do(func() { logging.Logger = zerolog.Nop() })
	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(strings.NewReader(config.DefaultConfig)); err != nil {
		return nil, fmt.Errorf("load embedded secret rules")
	}
	var raw config.ViperConfig
	if err := v.Unmarshal(&raw); err != nil {
		return nil, fmt.Errorf("decode embedded secret rules")
	}
	cfg, err := raw.Translate()
	if err != nil {
		return nil, fmt.Errorf("compile embedded secret rules")
	}
	d := detect.NewDetector(cfg)
	d.Redact = 100
	d.IgnoreGitleaksAllow = true // suppression support will be explicit and audited
	d.Verbose = false
	return &Detector{detector: d}, nil
}

func (d *Detector) Scan(ctx context.Context, path string, source []byte) []model.Finding {
	upstream := d.detector.DetectContext(ctx, detect.Fragment{Raw: string(source), FilePath: path, StartLine: 1})
	findings := make([]model.Finding, 0, len(upstream))
	for _, f := range upstream {
		confidence := "high"
		if f.RuleID == "generic-api-key" {
			confidence = "medium"
		}
		// Deliberately omit all source snippets, Match, Secret, and Fragment.
		// Redacting one match is insufficient when a line contains other secrets.
		identity := fmt.Sprintf("%s:%s:%d:%d", f.RuleID, path, f.StartLine, f.StartColumn)
		hash := sha256.Sum256([]byte(identity))
		findings = append(findings, model.Finding{
			ID: fmt.Sprintf("%x", hash[:12]), RuleID: "VYX-SEC-001", UpstreamRuleID: f.RuleID,
			Engine: "gitleaks", EngineVersion: Version, Severity: "high", Confidence: confidence,
			Path: path, Line: f.StartLine, Column: f.StartColumn,
			Explanation: "Credential-like content matches an embedded secret rule; validity and exposure are unverified.",
			Evidence:    "[REDACTED] (" + f.RuleID + ")",
			Remediation: "Determine whether this is a real credential; remove it from source and rotate it if exposed. Use server-side secret storage.",
		})
	}
	return findings
}
