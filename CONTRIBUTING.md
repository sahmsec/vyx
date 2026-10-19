# Contributing to Vyx

Vyx is a Phase 1 feasibility prototype. Start with [AGENTS.md](AGENTS.md), [FUTURE_PLAN.md](FUTURE_PLAN.md), and [docs/STATUS.md](docs/STATUS.md). Discuss a proposed new detector against the current milestone before broadening framework coverage.

## Build and verify

Use Go 1.27 or newer. No Node.js, Python, Docker, or C compiler is required for source builds.

```sh
go mod download
CGO_ENABLED=0 go build -trimpath -o vyx ./cmd/vyx
go test ./...
go vet ./...
```

On Windows PowerShell:

```powershell
.\scripts\dev.ps1 check
.\scripts\install.ps1
vyx version
```

`dev.ps1` prefers `.tools/go/bin/go.exe` when available, then Go on PATH. Dependencies/toolchains may need network access during development; that is separate from scanner runtime privacy.

## Changes and tests

- Format Go changes with gofmt. Use paired vulnerable/secure fixtures for security rules, preserving accurate locations, redaction, confidence, and incomplete-analysis behavior.
- Keep default scans offline and never execute the target application's code, configuration, install hooks, or builds.
- Keep the single-executable/cgo-disabled contract. Validate native runtime behavior separately from cross-compilation.
- Do not check in generated binaries, local tooling, scan reports, credentials, caches, or real private source fixtures.
- Pin dependency changes, regenerate notices with `scripts/notices.ps1`, and review licenses and grammar artifacts. Generated root notices are evidence, not a complete artifact audit.
- Describe the concrete behavior change, validation, coverage limits, and any outstanding gates. Do not claim general authorization detection from a missing guard pattern.
- Update STATUS for verified progress, ROADMAP for milestone gates, and FUTURE_PLAN for deferred work and the next action.

## License

Vyx's original code is MIT-licensed; see [LICENSE](LICENSE). Third-party code and assets retain their own licenses and notices. Do not redistribute rules whose licenses do not permit it.
