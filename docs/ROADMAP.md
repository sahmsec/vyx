# Roadmap and deferred work

Keep this fixed numbering. Earlier discussion used 1A/1B; those labels are superseded.

The ordered next-work queue is [FUTURE_PLAN.md](../FUTURE_PLAN.md), linked from both agent instruction files. The owner chose MIT and authorized initial source publication; that publishing work is separate from unfinished Phase 1 detection gates.

| Phase | Status | Deliverables | Acceptance gate |
| --- | --- | --- | --- |
| 0 — Planning | Completed | Audit, research, accepted architecture and constraints | Saved plan, explicit unknowns |
| 1 — Feasibility | In progress | Minimal CLI; embedded JS/TS parser, secret detector, offline dependency prototype; packaging/license inventory and baselines | Go-only/cgo-disabled build; standalone clean-machine execution; paired fixtures; blocked-network checks; measured footprint; report failures honestly |
| 2 — MVP | Planned | Six PLAN rule families; terminal/JSON/SARIF; policy/config/suppressions; compact agent output and agent skill | Evidence and redaction tests; per-rule supported-version fixtures; visible incomplete analysis; reviewed pilot accuracy |
| 3 — Performance and integration | Planned | Persistent cache, dependency-aware incremental scans, baseline/diff workflows, Git hooks/CI docs, local stdio MCP | Incremental/full equivalence; measured improvement; on-demand bounded agent evidence; token-savings evaluation |
| 4 — Expansion | Deferred | Selected injection checks, Express framework model, Supabase migration analysis, Go analysis, bounded cross-file authorization | Each addition has a support contract, secure/vulnerable fixtures, licensing review, and performance/accuracy evidence |

## Phase 1 checklist

- [x] Workspace-local Go toolchain and verified checksums.
- [x] Minimal CLI and normalized completeness/finding model.
- [x] Deterministic source inventory with explicit exclusions/limits.
- [x] Embedded JS/TS/TSX parser passes initial syntax fixtures with no cgo (broader parity evaluation pending).
- [x] Embedded secret detector validates redacted positive/negative examples.
- [ ] Offline dependency prototype identifies a known advisory from a local snapshot.
- [ ] Audit transitive dependencies and redistribute-required notices.
- [x] Windows runtime checks and Linux amd64/macOS arm64 compilation (native Linux/macOS runtime validation pending).
- [ ] Cold/warm timing, peak memory, CPU, and binary size baseline.
- [ ] Blocked-egress and no-project-execution checks.
- [ ] Record component decisions and remaining gates before expanding.

## Deferred catalog: preserve for later

- SQL/NoSQL/OS command injection; XSS; path traversal; SSRF; unsafe upload/processing; mass assignment; unsafe evaluation/deserialization.
- JWT misuse, session cookies, CSRF, CORS, redirects, debug endpoints/configuration. Missing headers/rate limiting remain review signals without effective deployment context.
- Express routing/middleware semantics; Go gosec/govulncheck/go/analysis reuse.
- Supabase migration/schema modeling: RLS, grants, exposed schemas, combined policies, public storage, RPC/SECURITY DEFINER. Local migrations never prove deployed state.
- Authentication/authorization/ownership flow, role escalation, related-route inconsistencies, and recovery/session workflows. Absence of a recognized guard is never proof.
- Optional explicitly scoped Git history scanning; supported dependency deprecation/EOL metadata without treating age as a vulnerability.
- Benchmark alternative secret/SAST engines only if they preserve distribution and license requirements.
- Optional authorized sandbox dynamic verification, only after a separate threat model and user opt-in.

## Not part of current work

GUI/web dashboard; cloud service; automatic code changes; default AI/remote scanning; arbitrary executable plugins; SQLite; background daemon; universal engine; production endpoint/database probing; tags/releases. Initial source commit/push was explicitly authorized separately.

## Decisions still open

Production parser binding; supported Next.js versions; advisory freshness limits; performance budgets; actual pilot corpus; final library versions and maintenance tradeoffs. Do not turn these into undocumented assumptions. Repository/module is `github.com/sahmsec/vyx`; project license is MIT, both selected by the user.
