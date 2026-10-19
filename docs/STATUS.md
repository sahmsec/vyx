# Current status

Updated: 2026-10-19 (Asia/Dhaka).

## Authorization and scope

User approved starting implementation and requested AGENTS.md, CLAUDE.md, and persistent future plans. Current detection scope: Phase 1 feasibility. The user subsequently explicitly authorized adding a license and instructions and committing/pushing to origin/main, and selected MIT. No tags/releases or real-system testing were requested.

User selected **Vyx** as the project name. Updated the CLI, command directory (`cmd/vyx`), module/imports, first-party rule prefix (`VYX-`), development scripts, and saved instructions/plans. Broader name/domain availability remains unverified.

The user supplied `https://github.com/sahmsec/vyx`, initially verified empty. Canonical Go module/imports: `github.com/sahmsec/vyx`; branch: main; origin: `https://github.com/sahmsec/vyx.git`. Initial source publication is authorized. Git history and remote refs record the actual publication commit; no release/tag is part of this task.

The user requested normal command usage after installation. Added and ran a Windows per-user installer at `%LOCALAPPDATA%\Vyx\bin\vyx.exe`, with one entry added to user PATH. Machine PATH was not changed. Existing terminals/apps may need restarting to inherit the user PATH update. Reinstall via `scripts/install.ps1`; uninstall via `scripts/install.ps1 -Uninstall`.

## Completed

- Confirmed the workspace was empty and had no Git repository or existing project instructions.
- Saved accepted architecture in PLAN.md and fixed phases/deferred additions in ROADMAP.md.
- Added shared AGENTS.md and Claude entry instructions.
- Added MIT LICENSE (owner selected), expanded installation/usage/development instructions, CONTRIBUTING.md, and a dedicated FUTURE_PLAN.md linked from both AGENTS.md and CLAUDE.md. It preserves the ordered next step, all deferred phases, and agent skill/MCP plans. Installer retains the project license separately from third-party notices.
- Installed Go 1.27.2 windows/amd64 under `.tools/go`, without changing system PATH. Verified archive SHA-256: `1314008898bd40df77af4b014f777f08873dbdfbcd3d92308728ee03304fe04f` against the official download manifest.
- Added Go module, pinned Gitleaks v8.30.1 and gotreesitter v0.56.0, CLI entry point, deterministic inventory, normalized findings/completeness, text/JSON output, and development helper scripts.
- Embedded secret detector returns fixed redaction markers, upstream rule identity, and accurate one-based locations; generic matches retain medium confidence.
- Parser candidate handles initial JS/JSX, TS, and TSX fixtures and extracts direct client/server directives/import counts internally. This is not framework vulnerability detection.
- Discovery includes environment files; explicitly excludes generated/dependency/tooling directories, binary/non-UTF-8 data, and symlinks/nonregular files. Scan-root reads use Go's os.Root. Limits/parse failures/cancellation mark scans incomplete.
- Saved parser rejection/alternative evidence in DECISIONS.md.
- Generated DEPENDENCY_INVENTORY.json and THIRD_PARTY_NOTICES.txt for 63 modules supplying imported production packages and three embedded grammars. No missing root notices in that snapshot; full artifact/license review remains open.
- Added `scripts/install.ps1`: installs the existing build and notices per user (builds only if missing), records file hashes and PATH ownership, supports isolated `-InstallDirectory ... -NoPath` use, and removes only managed unmodified files on uninstall. It refuses existing unmanaged targets and linked installation paths.

## In progress

- Phase 1 is started, not complete. Current increment is a working secret-scanning and parser-feasibility CLI.
- Linux amd64 and macOS arm64 cgo-disabled cross-compilation passed; runtime validation on those operating systems is still pending.

## Not implemented

Dependency advisory scanning/updates, the five other MVP rule families, framework security models, SARIF, configuration/suppressions, caches/incremental scans, compact agent mode, agent skill, and MCP remain unimplemented. No live-system testing or release/tag creation is part of this task. The initial source snapshot is being published under the user's explicit instruction; see Git history for its commit.

## Verification

- Initial publication checks: `scripts/dev.ps1 check` passed (tests, vet, cgo-disabled Windows build); dependency notices regenerated; installer syntax and MIT license copy/hash/reinstall/uninstall checks passed in an isolated NoPath directory. The per-user installation was updated to retain LICENSE. A full current repository credential scan completed with zero findings; this is a prepublication secret check, not a claim of complete application security. Source/tooling/build exclusions remain documented. The source snapshot is authorized for origin/main; use Git history/remote refs for the actual publication commit.
- Initial canonical repository setup: remote was empty and local main had no commits at that checkpoint. `scripts/dev.ps1 check` passed after changing the module/imports to `github.com/sahmsec/vyx`. Notices regenerated for 63 dependency modules with the main module correctly excluded. `CGO_ENABLED=0 go install ./cmd/vyx` into an isolated workspace GOBIN passed; binary build metadata identifies `github.com/sahmsec/vyx/cmd/vyx`. Updated the per-user installed binary/notices and verified its hash matches dist/vyx.exe and its version command works. That earlier checkpoint did not test remote installation.
- Installer syntax parsed without errors. An isolated install/reinstall/version/uninstall smoke passed using `-NoPath`, with user/process PATH unchanged. Default installation then passed `vyx version` and `vyx scan . --format json` from the synthetic app directory outside the repository root; the expected credential finding produced exit 1, three sources parsed, and complete requested analysis. Confirmed exactly one user PATH entry and unchanged machine PATH. No scanner code changed, so the earlier Go checks were not repeated for this installer-only change.
- Rename verification: `scripts/dev.ps1 check` passed under the local `vyx` module; notices were regenerated. `dist/vyx.exe version` prints `vyx 0.0.1-dev`; help advertises `vyx scan`. The renamed Windows binary passed the synthetic scan with an empty PATH, three parsed sources, one redacted `VYX-SEC-001` finding, and exit 1. Linux amd64/macOS arm64 cross-builds passed using `./cmd/vyx`. Old generated binary names were removed from dist/. Only the naming history and actual workspace path retain the old name.
- The sizes/timings below describe the earlier pre-rename artifact; they have not been remeasured as performance baselines for Vyx.
- `scripts/dev.ps1 check`: `go test ./...`, `go vet ./...`, and `CGO_ENABLED=0 go build -trimpath` passed on Windows.
- Paired secret tests cover synthetic provider credentials, safe public configuration, adjacent-secret redaction, one-based multiline locations, and preventing unreviewed upstream suppression comments.
- Parser tests cover JS/JSX/TS/TSX, type annotations/generics, directives vs nested/late strings, malformed syntax, and ignoring external grammar-blob overrides.
- Scan/CLI tests cover environment files, exclusions, file/count limits, incomplete-error precedence, cancellation, configuration nonexecution, valid JSON, argument validation, and symlink nonfollowing. The targeted Windows symlink test passed without being skipped.
- Standalone Windows executable launched five times with an empty PATH and a cleared environment except SystemRoot. Four synthetic files: one credential in .env.local and three JS/TS/TSX source files. Every run parsed three sources, returned one redacted candidate, completed requested checks, and exited 1. No Go/Node/Python/scanner executable was available through PATH.
- Unstripped trimpath binary size: 26,250,240 bytes. Initial synthetic launch 2619.07 ms; subsequent launches 133.51, 156.49, 203.98, 155.56 ms. These are small process smoke checks with OS caches uncontrolled and concurrent benchmark work, not representative application benchmarks or incremental-cache results. Peak memory collection was unavailable/null and must be fixed before claiming memory results.
- Parser microbenchmark: 100 iterations, 186436 ns/op, 5784 B/op, 34 allocs/op on Windows amd64, Ryzen 7 4800H, Go 1.27.2. One synthetic TypeScript function, initialized parser. This does not measure full scan speed or broader accuracy.
- `go tool nm` shows embedded blob variables for JavaScript, TypeScript, and TSX only in the inspected Windows artifact; the all-language registry is not imported.
- `CGO_ENABLED=0` cross-builds passed: linux/amd64 (25,866,009 bytes) and darwin/arm64 (24,624,594 bytes). These are compilation checks, not native runtime certification. Artifacts are under ignored dist/.
- `gofmt -l cmd internal` returned no unformatted source.

## Exact next step

Prototype OSV-Scanner v2.7.0's Go API in isolation using a small supported lockfile plus a local advisory fixture. Explicitly disable online plugins, downloads, call analysis, and transitive resolution; confirm which CLI offline transformations must be reproduced. Measure its build/dependency/binary footprint before linking it into the CLI. Prove no egress/project execution and an actual positive/negative advisory match. Do not build a new advisory/version engine if the high-level package is unsuitable; evaluate narrower upstream packages and record the tradeoff.

Other open Phase 1 gates: blocked-egress tests of the full executable, parser parity/resource/cancellation evaluation, real project cold/warm baselines and proper peak-memory accounting, and native Linux/macOS runtime checks. See ROADMAP.md; do not start Phase 2's broad rules or Phase 3's MCP ahead of these gates.
