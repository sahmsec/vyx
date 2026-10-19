# Vyx

A local-first security scanner, currently in **Phase 1: feasibility**.

The target is one standalone Go executable with no Node.js, Python, Docker, or external scanner requirement. This repository is an early prototype, not a completed security audit product.

## Project memory

- [Accepted plan](docs/PLAN.md)
- [Future work and agent handoff](FUTURE_PLAN.md)
- [Roadmap and deferred additions](docs/ROADMAP.md)
- [Current progress and next step](docs/STATUS.md)
- [Shared agent instructions](AGENTS.md)
- [Claude instructions](CLAUDE.md)
- [Contributing and development instructions](CONTRIBUTING.md)

Canonical repository: [sahmsec/vyx](https://github.com/sahmsec/vyx). Go module: `github.com/sahmsec/vyx`. Original code is [MIT-licensed](LICENSE).

## Run the current prototype

Install on Windows from this checkout:

```powershell
git clone https://github.com/sahmsec/vyx.git
cd vyx
.\scripts\install.ps1
vyx version
vyx scan .
```

The installer copies the built executable, MIT license, and third-party notices to `%LOCALAPPDATA%\Vyx\bin` and adds that directory to your **user PATH**. Administrator access is not required. It also updates the PowerShell session running the installer; other already-open terminals/apps may need restarting. It uses `dist/vyx.exe` when present and builds from source only if that file is missing. A fresh clone needs Go 1.27 or newer to build; running the installed binary needs no Go or other developer runtime.

Uninstall with `.\scripts\install.ps1 -Uninstall`. It removes only the managed installation files and the PATH entry it added. Modified installed files are preserved by refusing removal. `-InstallDirectory <absolute-path> -NoPath` supports a separate installation without PATH changes.

Alternatively, install from source with Go 1.27 or newer:

```sh
go install github.com/sahmsec/vyx/cmd/vyx@latest
```

Go installs commands to GOBIN (or GOPATH/bin by default), which must be on PATH. Run `go env GOBIN GOPATH` to locate it. Installation downloads Go modules; normal scans do not update or download components. Prebuilt releases are not published yet.

On this workspace, a verified portable Go toolchain is under `.tools/go`; it is ignored and is not a runtime dependency.

```powershell
.\scripts\dev.ps1 check
.\dist\vyx.exe version
.\dist\vyx.exe scan D:\path\to\app
.\dist\vyx.exe scan D:\path\to\app --format json
```

With Go already installed, build on any target host:

```sh
CGO_ENABLED=0 go build -trimpath -o vyx ./cmd/vyx
go test ./...
go vet ./...
```

PowerShell: set `$env:CGO_ENABLED='0'` before invoking `go build`.

Current checks:

- Embedded Gitleaks credential-candidate rules, returning redacted metadata only.
- Embedded JavaScript/JSX, TypeScript, and TSX parsing. Recognizing syntax is **not** a framework security check or TypeScript type check.
- Deterministic discovery, explicit exclusions, bounded file sizes/counts, cancellation between work items, and incomplete-scan diagnostics.

## Usage

```sh
vyx version
vyx --help
vyx scan .
vyx scan ./my-app --format json
vyx scan . --max-file-bytes 2097152 --max-files 20000
```

Only `text` and `json` formats are currently implemented. Text findings show the rule, severity/confidence, path, one-based line/column, redacted evidence, and suggested remediation. JSON includes the complete finding list, exclusions, diagnostics, and coverage limits. Example `VYX-SEC-001` findings are credential candidates, not verified active credentials.

Default limits: 1 MiB per file and 10,000 regular files. Override with `--max-file-bytes` or `--max-files`. Limit breaches produce an incomplete scan. Ignored directories: `.git`, `.tools`, `node_modules`, `.next`, `dist`, `vendor`, `coverage`. Symlinks/nonregular files and binary/non-UTF-8 files are excluded and listed in JSON; source files with unsupported encoding mark the scan incomplete. `.gitignore` is not interpreted yet; environment files are intentionally scanned.

Exit 0 means the requested implemented checks completed without a high-confidence high-severity credential candidate. Exit 1 means such a candidate was found. Exit 2 means invalid input or incomplete analysis; errors take precedence over findings. Neither exit 0 nor parsed syntax establishes application security. Generic credential matches use medium confidence and do not trigger exit 1.

Framework secret-exposure checks, dependency advisories, SARIF, configuration/suppressions, caches, agent output, skills, and MCP are planned, not implemented. Scans do not execute application configuration or install packages. There is no scan-time network operation in Vyx's code; blocked-egress verification of the complete dependency graph remains an open Phase 1 gate.

## Dependency evidence

`scripts/notices.ps1` regenerates the production-import module inventory and retained upstream root notices. See [component decisions](docs/DECISIONS.md), [dependency inventory](docs/DEPENDENCY_INVENTORY.json), and [third-party notices](THIRD_PARTY_NOTICES.txt). This is not a completed release/license audit.

Vyx's original code is licensed under the [MIT License](LICENSE). Third-party components and grammar assets retain their own licenses. See [CONTRIBUTING.md](CONTRIBUTING.md) for change and verification instructions.
