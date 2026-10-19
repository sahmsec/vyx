# Vyx: accepted plan

Decision date: 2026-10-19. The user authorized starting Phase 1 and requested persistent plans plus AGENTS.md and CLAUDE.md. The repository was empty and was not a Git repository at the initial audit.

The user selected **Vyx** (pronounced "vicks") after reviewing naming options. It supersedes the original working name, Vibehound. The command is `vyx`, the canonical Go module is `github.com/sahmsec/vyx`, and first-party rule IDs use `VYX-`. The user supplied https://github.com/sahmsec/vyx as the repository. Broader name/domain availability has not been established.

## Product

`vyx scan .`: a local, token-free security scanner with actionable evidence. First framework: Next.js JS/TS. AI-assisted developers are the initial audience; code provenance does not change vulnerability semantics. Accuracy and repeated-scan speed matter more than rule count.

## Distribution and dependencies

- One standalone Go executable, prebuilt for Windows, Linux, and macOS.
- Go alone for source installation; `CGO_ENABLED=0` is a release gate.
- No external Node/Python/Docker/scanner runtime. Embedded WASM is allowed; it is an implementation detail.
- Local advisory data is separate from the executable and updated explicitly. Missing data must not become a clean dependency result.
- Git is optional for future history scans, not required for directory scans.
- Canonical repository/module: `github.com/sahmsec/vyx`. Origin: `https://github.com/sahmsec/vyx.git`. The user authorized initial source publication and explicitly chose MIT. This does not authorize tags/releases, live-system testing, or license changes to third-party code.

## Architecture

Go CLI -> deterministic file inventory -> embedded analysis components -> normalized findings/coverage -> policy and deduplication -> terminal/JSON/SARIF/compact agent output.

Preferred feasibility candidates, not unconditional production selections:

1. Original preference: Tree-sitter WASM through wazero. Phase 1 inspection rejected the first wrapper because its bundle lacks TypeScript/TSX. Current prototype candidate: gotreesitter v0.56.0, a pure-Go runtime with embedded JavaScript/TypeScript/TSX grammar tables. Validate correctness, memory, startup, cancellation, and maintenance before production adoption. A syntax tree does not provide TypeScript's type checker. See DECISIONS.md for evidence.
2. Gitleaks detection as a Go library, behind a small adapter. Pin its API/rules, redact all results, disable validation/network behavior. Upstream is security-patches-only; keep replaceable.
3. OSV-Scanner's Go package for supported lockfile/artifact scanning using local advisories. Explicitly disable online resolution, call analysis, build execution, and telemetry. Audit transitive dependencies and binary footprint.
4. Original, bounded Go framework rules. No new universal engine or handwritten general parser/version matcher.

Do not silently replace a failed feasibility candidate with an external runtime. Record the evidence and choose a compatible alternative or bring the architecture tradeoff to the user.

## MVP: six rule families

| ID | Check | Evidence/limits |
| --- | --- | --- |
| VYX-SEC-001 | Credential/private-key candidates | Strong patterns/context; redact; no claim credentials are active |
| VYX-SEC-002 | Secrets in public environment references | Secret classification plus a supported public exposure path; names alone insufficient |
| VYX-SEC-003 | Secrets exported through Next.js env config | Parse supported static config without executing it |
| VYX-CFG-002 | Sensitive files under public/ | Public asset path plus sensitive contents; examples need handling |
| VYX-SB-001 | Privileged Supabase key in client code | Privileged key and supported client boundary; publishable/anon keys are public by design |
| VYX-SCA-001 | Known dependency vulnerabilities | Exact resolved version, advisory identity/freshness; version match is not reachability |

Express, Go, and other frameworks initially receive only accurately described generic secret/dependency coverage. Do not claim deployed Supabase/RLS coverage from repository files.

## Findings and CLI policy

Preserve rule/engine versions, severity and confidence with their basis, relative source ranges, redacted evidence, related locations, remediation, references, assumptions, limitations, fingerprint, and suppression state. Include scan completeness, exclusions, unsupported syntax, timeouts, and advisory metadata.

Exit 0: required checks complete, policy not triggered. Exit 1: required checks complete, policy triggered. Exit 2: invalid invocation or incomplete required analysis. Failure takes precedence over findings. Default eventual failure policy: unsuppressed high/critical high-confidence findings. Review items are separate.

Planned commands: scan, rules list, db status/update, version; finding show and local MCP mode come with agent integration. Only implement commands backed by actual behavior. Machine output goes to stdout and diagnostics to stderr.

Configuration: small versioned TOML. Suppressions require rule-specific scope and a reason; preserve visibility, support expiry, and allow trusted CI policy. No executable configuration.

## Privacy and threat model

No scan-time network, uploads, telemetry, credential validation, project builds/scripts/config execution, or package-manager hooks. Contain paths and avoid symlinks/junctions outside the root. Apply resource limits and sanitize outputs. Keep caches user-local and minimal. Test blocked-egress behavior, not just flags. Updates require explicit commands. Dynamic or remote/AI features need separate opt-in and are deferred.

## Performance and validation

Inventory once; reuse parsed content; bound concurrency. File-based content-addressed caching is enough initially. Key caches by content, rules, config, parser, and advisory snapshot; invalidate affected dependents. No SQLite/daemon/cloud/plugin system initially.

Measure cold/warm/incremental scans, shared-helper/config/lockfile changes, peak full-process memory, CPU time, binary size, and advisory footprint. Compare against identical upstream components and fixtures. Set latency budgets after realistic baselines.

Use vulnerable/secure pairs for every enabled rule, framework fixtures, malformed/unsupported syntax, alias/guard cases, redaction, traversal/symlink cases, and incremental/full equivalence. Report per-rule precision/recall, false-positive opportunities, coverage, and abstentions. Proposed MVP actionable precision gate: 95% on an independently reviewed sufficiently sized pilot, with sample size and uncertainty; not a measured guarantee.

## Agent integration (accepted addition)

Phase 2: compact agent output plus a short skill teaching scan -> inspect selected evidence -> fix -> rescan. Keep full results local and retrieve evidence on demand; never truncate silently. Phase 3: local stdio MCP mode in the same executable with scan/list_findings/get_finding, plus sound incremental workflows. No extra runtime or persistent daemon required.

Deterministic scans use no model tokens; agent interpretation/fixes still do. Measure token savings against a comparable agent-only task, preserving detection quality. A cloud-backed agent can transmit returned findings to its provider: redact and disclose that boundary.

## Primary references

- https://github.com/malivvan/tree-sitter (prerelease cgo-free wrapper)
- https://github.com/wazero/wazero (embedded WASM runtime)
- https://github.com/gitleaks/gitleaks (MIT, maintenance status and library code)
- https://pkg.go.dev/github.com/google/osv-scanner/v2/pkg/osvscanner
- https://github.com/google/osv-scanner/blob/main/docs/offline-mode.md
- https://github.com/google/osv-scanner/blob/main/docs/scan-source.md (call-analysis execution risks)
- https://nextjs.org/docs/app/api-reference/config/next-config-js/env
- https://supabase.com/docs/guides/getting-started/api-keys
- https://semgrep.dev/legal/rules-license/ (do not redistribute restricted rules)
- https://learn.chatgpt.com/docs/build-skills
- https://learn.chatgpt.com/docs/extend/mcp?surface=cli

Check pinned artifact licenses, grammars, advisories, transitive code, notices, and redistribution independently. The project's original code is MIT-licensed; dependencies retain their own terms. See FUTURE_PLAN.md for the ordered handoff queue.
