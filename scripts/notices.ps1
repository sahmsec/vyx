# Generate a development dependency inventory and retain upstream root notices.
# This is evidence for review, not a complete redistribution/legal clearance.
$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path -Parent $PSScriptRoot
$goPath = Join-Path $projectRoot '.tools\go\bin\go.exe'
if (-not (Test-Path -LiteralPath $goPath)) { $goPath = (Get-Command go -ErrorAction Stop).Source }
$saved = @{}
foreach ($name in @('CGO_ENABLED','GOTOOLCHAIN','GOTELEMETRY','GOPATH','GOCACHE')) {
    $saved[$name] = [Environment]::GetEnvironmentVariable($name, 'Process')
}
Push-Location $projectRoot
try {
    $env:CGO_ENABLED = '0'
    $env:GOTOOLCHAIN = 'local'
    $env:GOTELEMETRY = 'off'
    $env:GOPATH = Join-Path $projectRoot '.tools\gopath'
    $env:GOCACHE = Join-Path $projectRoot '.tools\gocache'
    $mainModule = & $goPath list -m -f '{{.Path}}'
    if ($LASTEXITCODE -ne 0) { throw 'Cannot determine the main Go module.' }
    $rows = & $goPath list -deps -f '{{if .Module}}{{.Module.Path}}|{{.Module.Version}}|{{.Module.Dir}}{{end}}' ./cmd/vyx
    if ($LASTEXITCODE -ne 0) { throw 'Cannot enumerate linked-package modules.' }
    $inventory = @()
    $notice = [System.Text.StringBuilder]::new()
    [void]$notice.AppendLine('Vyx development third-party notices')
    [void]$notice.AppendLine('Generated from modules providing imported production packages. Not a complete release audit.')
    [void]$notice.AppendLine('Vyx original code is MIT-licensed; see LICENSE. Third-party licenses remain separate.')
    foreach ($row in ($rows | Where-Object { $_ -and $_ -notlike ($mainModule + '|*') } | Sort-Object -Unique)) {
        $parts = $row -split '\|', 3
        $modulePath, $version, $directory = $parts
        $files = @(Get-ChildItem -LiteralPath $directory -File | Where-Object { $_.Name -match '^(LICENSE|LICENCE|COPYING|NOTICE)(\..*)?$' })
        $retained = @()
        foreach ($file in $files) {
            [void]$notice.AppendLine("`n===== $modulePath $version / $($file.Name) =====")
            [void]$notice.AppendLine([IO.File]::ReadAllText($file.FullName))
            $retained += [PSCustomObject]@{ file=$file.Name; sha256=(Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant() }
        }
        $inventory += [PSCustomObject]@{ module=$modulePath; version=$version; root_notices=$retained; review= $(if ($files.Count) { 'root notices retained; artifact audit pending' } else { 'missing root notice; investigate' }) }
    }
    $grammarRoot = Join-Path $env:GOPATH 'pkg\mod\github.com\odvcencio\gotreesitter@v0.56.0'
    $audit = Get-Content -LiteralPath (Join-Path $grammarRoot 'licenses\grammars.json') -Raw | ConvertFrom-Json
    $grammars = @($audit.entries | Where-Object name -In @('javascript','typescript','tsx'))
    foreach ($grammar in $grammars) {
        [void]$notice.AppendLine("`n===== Embedded grammar: $($grammar.name) =====")
        [void]$notice.AppendLine("$($grammar.repo) @ $($grammar.ref)")
        foreach ($holder in $grammar.copyright_holders) { [void]$notice.AppendLine($holder) }
        [void]$notice.AppendLine([IO.File]::ReadAllText((Join-Path $grammarRoot "licenses\texts\$($grammar.spdx).txt")))
    }
    [IO.File]::WriteAllText((Join-Path $projectRoot 'THIRD_PARTY_NOTICES.txt'), $notice.ToString(), [Text.UTF8Encoding]::new($false))
    $snapshot = [PSCustomObject]@{ schema_version=1; scope='production import graph; compiled-artifact audit pending'; modules=$inventory; embedded_grammars=$grammars }
    $snapshot | ConvertTo-Json -Depth 8 | Set-Content -LiteralPath 'docs\DEPENDENCY_INVENTORY.json' -Encoding utf8NoBOM
    Write-Output "Retained root notices for $($inventory.Count) imported modules and three embedded grammars. Release audit remains open."
} finally {
    Pop-Location
    foreach ($name in $saved.Keys) { [Environment]::SetEnvironmentVariable($name, $saved[$name], 'Process') }
}
