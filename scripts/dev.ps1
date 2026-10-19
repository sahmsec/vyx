param([ValidateSet('check','build','test','vet')][string]$Task = 'check')

$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path -Parent $PSScriptRoot
$goPath = Join-Path $projectRoot '.tools\go\bin\go.exe'
if (-not (Test-Path -LiteralPath $goPath)) {
    $goCommand = Get-Command go -ErrorAction Stop
    $goPath = $goCommand.Source
}
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
    if ($Task -in @('check','test')) {
        & $goPath test ./...
        if ($LASTEXITCODE -ne 0) { throw 'Go tests failed.' }
    }
    if ($Task -in @('check','vet')) {
        & $goPath vet ./...
        if ($LASTEXITCODE -ne 0) { throw 'Go vet failed.' }
    }
    if ($Task -in @('check','build')) {
        $output = if ($IsWindows) { 'dist/vyx.exe' } else { 'dist/vyx' }
        & $goPath build -trimpath -o $output ./cmd/vyx
        if ($LASTEXITCODE -ne 0) { throw 'Go build failed.' }
    }
} finally {
    Pop-Location
    foreach ($name in $saved.Keys) {
        [Environment]::SetEnvironmentVariable($name, $saved[$name], 'Process')
    }
}
