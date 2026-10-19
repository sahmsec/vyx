# Component decisions and feasibility evidence

## 2026-10-19: reject the first WASM wrapper for this prototype

Candidate: `github.com/malivvan/tree-sitter` v0.0.1, upstream master inspected during Phase 1.

Its bundled `_languages` list in const.go includes JavaScript but not TypeScript or TSX. Its Tree API also lacks an exported tree/module close operation and ParseString defers an allocation cleanup before checking its error. These are implementation inspection findings, not benchmark results. Do not adopt this wrapper as the planned JS/TS parser without a separately reviewed repair and grammar build pipeline.

Source: https://github.com/malivvan/tree-sitter/blob/5aa45b3d087ac4f9014faa31e2ea7a5eeabe834b/const.go

## 2026-10-19: evaluate a pure-Go Tree-sitter implementation

Candidate pinned: `github.com/odvcencio/gotreesitter` v0.56.0 (MIT).

- Uses embedded Tree-sitter grammar tables with a pure-Go runtime; no C compiler or Node runtime required.
- Import only the JavaScript, TypeScript, and TSX blob accessors plus runtime support. Avoid the all-language registry.
- Load our embedded bytes explicitly through runtime.LoadLanguage; do not consult external grammar overrides.
- Use strict parsing, bounded per-file time/memory, syntax error detection, tree release, and sanitized failure reporting.
- This is a different candidate from upstream's C runtime. Passing local fixtures does not establish full upstream compatibility or production suitability. Broader parity and adversarial tests remain a gate.

Source: https://github.com/odvcencio/gotreesitter/tree/v0.56.0

## Secret detector

Pinned `github.com/zricethezav/gitleaks/v8` v8.30.1 (repository is gitleaks/gitleaks). Use the detection package and embedded rules, not its CLI, Git/file discovery, executable configuration, or remote extension loading. Construct a private Viper instance and silence upstream logs before use. Never serialize raw upstream findings: adjacent secrets could survive single-match redaction. Our model retains location/rule identity and a fixed redaction marker only.

## Dependency scanner

OSV-Scanner v2.7.0 remains a feasibility candidate, not an implemented dependency. Its Go package imports the broad Scalibr pipeline and container helpers; its dependency footprint must be measured. Its ScannerActions exposes local DB, plugin network disabling, call-analysis states, and transitive-scanning controls. Source comments warn that not all CLI offline transformations occur inside DoScan. Do not assume calling the library with CompareOffline alone reproduces CLI offline behavior.

Source: https://github.com/google/osv-scanner/blob/v2.7.0/pkg/osvscanner/osvscanner.go

The next prototype must prove supported artifact matching using a local advisory snapshot with egress blocked and no project execution before adding this package to the main executable.
