# Vyx agent instructions

Read [FUTURE_PLAN.md](FUTURE_PLAN.md), `docs/PLAN.md`, `docs/ROADMAP.md`, and `docs/STATUS.md` before changing code. They preserve the accepted architecture, deferred work, and current progress across sessions. Update STATUS after meaningful work; keep the ordered future-work queue in FUTURE_PLAN and milestone gates in ROADMAP rather than silently expanding the current phase.

## Product contract

- A local-first security scanner for developers, starting with Next.js JavaScript/TypeScript.
- One standalone executable. No required Node.js, Python, Docker, external scanner executable, or C compiler. Source installation must support `CGO_ENABLED=0` with Go alone. Verify this; do not merely claim it.
- Default scans must never upload source, secrets, findings, dependency inventory, or telemetry. Advisory downloads are an explicit, separate operation.
- Do not execute the scanned application's code, scripts, builds, configuration, or package-manager hooks.
- Reuse established parsers, secret detection, advisory matching, and version semantics. Do not build a universal analysis engine.
- Findings require evidence. An unknown authentication guard is not proof of missing authorization. Unsupported analysis is not a clean result.
- Redact secrets before any output, logging, or persistence. Keep severity separate from confidence.

## Work conventions

- Follow the fixed phase numbering in ROADMAP: 0 planning, 1 feasibility, 2 MVP, 3 performance/integrations, 4 expansion. Do not introduce 1A/1B.
- Start with read-only inspection; preserve unrelated user changes. Do not commit, push, publish, or test real systems without authorization.
- User authorized project setup, local development dependencies, and starting Phase 1 on 2026-10-19. Prefer workspace-local tools; do not alter system PATH or machine-wide installations.
- User requested normal `vyx scan .` use after installation. `scripts/install.ps1` manages an optional per-user Windows installation and user PATH entry; it must never alter machine PATH. Keep its ownership manifest and reversible uninstall behavior intact.
- No sub-agent delegation unless the user explicitly requests it.
- Keep Go code straightforward and interfaces small. Pin direct dependencies and retain upstream licenses/notices.
- Canonical Go module and repository: `github.com/sahmsec/vyx`; local origin is `https://github.com/sahmsec/vyx.git`. The user explicitly authorized the initial commit and push of the prototype, MIT license, instructions, and future-plan handoff. No tags/releases or real-system testing were requested. Follow later explicit user instructions for additional publication.
- Project license: MIT, explicitly selected by the owner. Preserve third-party licenses independently; do not relicense dependencies or restricted rules.
- Go is currently a workspace-local tool under `.tools/go/bin`. Build artifacts and tooling are ignored.
- Use `gofmt`, `go test ./...`, `go vet ./...`, and `CGO_ENABLED=0 go build ./cmd/vyx` as appropriate. Scan changes need vulnerable/secure pairs and evidence/redaction checks.
- Cross-compilation verifies compilation only. Do not claim runtime validation on an OS not actually tested.
- Benchmarks must disclose versions, corpus, exclusions, cold/warm conditions, and correctness. Do not invent performance or token-savings results.
- The agent-facing mode must return compact summaries and on-demand evidence without hiding incomplete coverage. Connecting a cloud agent may transmit returned evidence to that provider.

## Completion

Record changes, commands/tests, measured results, unresolved risks, and the exact next step in STATUS. Report incomplete gates honestly. Do not mark a planned feature implemented based on a stub or a successful build alone.
