# Vyx: future work and agent handoff

Read this file with [AGENTS.md](AGENTS.md) or [CLAUDE.md](CLAUDE.md) before continuing development. Keep completed/verified work in [docs/STATUS.md](docs/STATUS.md), architecture in [docs/PLAN.md](docs/PLAN.md), and milestone gates in [docs/ROADMAP.md](docs/ROADMAP.md). This file is the ordered future-work queue, not implemented coverage.

## Current position

Phase 1 (feasibility) is in progress. Implemented: standalone Go CLI, embedded credential detection, JS/JSX/TS/TSX parsing, text/JSON reporting, bounded discovery, explicit incomplete results, and a Windows per-user installer. Syntax parsing is not framework security analysis.

Canonical repository/module: `github.com/sahmsec/vyx`. Project license: MIT, selected by the owner. Source publication was explicitly authorized; tags/releases and real-system testing were not requested.

## Next: finish Phase 1

1. **Offline dependencies:** prototype the OSV-Scanner v2.7.0 Go API in isolation. Use a supported lockfile and local advisory fixture. Disable downloads, network-backed plugins, transitive resolution, call analysis, project configuration execution, and builds. Inspect which CLI offline transformations the library does not apply. Verify vulnerable and secure versions, missing-database errors, no egress, and no project execution. Measure its footprint before linking it into the main CLI. Reuse narrower upstream packages if the high-level package is unsuitable; do not write a new advisory/version engine.
2. **Parser qualification:** expand secure/vulnerable syntax fixtures and compare JS/TS/TSX behavior against upstream parser expectations. Exercise malformed input, resource budgets, cancellation, and complete/partial-tree handling. Keep gotreesitter replaceable until validated.
3. **Privacy qualification:** run the full executable under blocked egress. Current empty-PATH tests demonstrate standalone execution, not network isolation. Avoid running application code or real endpoints.
4. **Real performance baselines:** pin representative small apps and larger repositories; measure cold/repeated/incremental scenarios, CPU, process memory, binary size, and database storage. Repair Windows peak-memory accounting. Existing synthetic smoke timings are not product benchmarks.
5. **Portability and redistribution:** native Linux/macOS runtime checks, not just cross-compilation; audit linked code, grammars, advisories, transitive licenses, and required notices.

Record results and unresolved gates before expanding scope.

## Then: Phase 2 MVP

Finish the six rule families in PLAN.md: credential candidates, public environment secret exposure, Next.js env config exposure, sensitive public files, privileged Supabase keys reaching client code, and known dependency advisories.

Add terminal/JSON/SARIF parity, a documented support/version matrix, configuration, reasoned suppressions with expiry, deduplication, and a reviewed accuracy corpus. Target 95% precision for the default actionable set as a proposed gate, with sample counts and uncertainty disclosed; do not claim it is already achieved.

**Agent integration starts here:** compact agent output and a short skill for scan -> inspect selected evidence -> fix -> rescan. Keep complete reports local and retrieve details on demand. Always expose incomplete coverage/truncation. Scanning itself consumes no model tokens; agent interpretation/fixes still do. Measure savings rather than promising a percentage.

## Then: Phase 3 performance and integration

Content-hash file caching; dependency-aware invalidation; incremental/full result equivalence; changed-code/baseline workflows; Git hooks and CI examples.

Add local stdio MCP mode inside the same executable with a small scan/list_findings/get_finding interface. No separate Node/Python runtime or background daemon. Cloud-backed agents can send returned evidence to their provider; redact and disclose this boundary.

## Later: Phase 4 expansion

- Validate injection checks individually: SQL/NoSQL/command injection, XSS, path traversal, SSRF, uploads, mass assignment, and unsafe evaluation/deserialization.
- Add Express routing/middleware models and later Go analysis by reusing established components.
- Analyze supported Supabase migrations: RLS, grants, exposed schemas, combined policies, storage, and RPC/SECURITY DEFINER. Repository files never prove deployed state.
- Research bounded cross-file authentication/authorization/ownership paths and role/recovery workflows. Unknown guards mean unresolved analysis, not proof of vulnerability.
- Add explicit local Git-history scope and reliable package deprecation/EOL metadata if justified.

Each capability needs secure/vulnerable fixtures, evidence quality, false-positive review, performance results, and a separate support contract.

## Preserve these constraints

One Go-installable standalone executable; no cgo or external developer runtime; no scan-time uploads/telemetry; explicit advisory updates; no untrusted application execution; evidence and confidence separate from severity. No GUI/cloud service, automatic fixes, universal engine, arbitrary plugins, or live/dynamic probing unless separately requested.

Use fixed phases 0–4. Do not reintroduce 1A/1B or silently mark deferred features complete. Update this queue and STATUS when work changes.
