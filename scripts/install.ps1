[CmdletBinding()]
param(
    [string]$InstallDirectory = (Join-Path ([Environment]::GetFolderPath('LocalApplicationData')) 'Vyx\bin'),
    [switch]$NoPath,
    [switch]$Uninstall
)

$ErrorActionPreference = 'Stop'
if ([Environment]::OSVersion.Platform -ne [PlatformID]::Win32NT) {
    throw 'This installer supports Windows. See README.md for other platforms.'
}
$projectRoot = Split-Path -Parent $PSScriptRoot
$installRoot = [IO.Path]::GetFullPath($InstallDirectory).TrimEnd('\')
if (-not [IO.Path]::IsPathRooted($InstallDirectory) -or $installRoot -eq [IO.Path]::GetPathRoot($installRoot).TrimEnd('\')) {
    throw 'Choose an absolute installation directory below a drive root.'
}
$executable = Join-Path $installRoot 'vyx.exe'
$manifestPath = Join-Path $installRoot 'vyx-install.json'
$noticesPath = Join-Path $installRoot 'THIRD_PARTY_NOTICES.txt'
$licensePath = Join-Path $installRoot 'LICENSE'

function Get-NormalizedPathEntry([string]$Entry) {
    return [Environment]::ExpandEnvironmentVariables($Entry.Trim().Trim('"')).TrimEnd('\')
}

function Test-PathEntry([string]$PathValue, [string]$Directory) {
    foreach ($entry in ($PathValue -split ';')) {
        if ((Get-NormalizedPathEntry $entry) -ieq $Directory) { return $true }
    }
    return $false
}

# Refuse to operate through an existing link or overwrite an unrelated program.
foreach ($path in @($installRoot, $executable, $manifestPath, $noticesPath, $licensePath)) {
    if (Test-Path -LiteralPath $path) {
        if ((Get-Item -LiteralPath $path -Force).Attributes -band [IO.FileAttributes]::ReparsePoint) {
            throw 'The installation location contains a link; choose a regular directory.'
        }
    }
}
$previous = $null
if (Test-Path -LiteralPath $manifestPath) {
    $previous = Get-Content -LiteralPath $manifestPath -Raw | ConvertFrom-Json
    if ($previous.product -ne 'Vyx' -or $previous.install_directory -ine $installRoot -or $previous.schema_version -ne 1) {
        throw 'The installation manifest does not identify this Vyx installation.'
    }
}
$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')

if ($Uninstall) {
    if (-not $previous) { throw 'No Vyx-managed installation exists at this location.' }
    $artifacts = @(@{path=$executable;hash=$previous.binary_sha256}, @{path=$noticesPath;hash=$previous.notices_sha256})
    if ($previous.license_sha256) { $artifacts += @{path=$licensePath;hash=$previous.license_sha256} }
    foreach ($artifact in $artifacts) {
        if (Test-Path -LiteralPath $artifact.path) {
            if ((Get-FileHash -LiteralPath $artifact.path -Algorithm SHA256).Hash -ine $artifact.hash) {
                throw 'An installed file was changed after installation; refusing to delete it.'
            }
        }
    }
    if ($previous.path_added) {
        $entries = @($userPath -split ';' | Where-Object { (Get-NormalizedPathEntry $_) -ine $installRoot })
        [Environment]::SetEnvironmentVariable('Path', ($entries -join ';'), 'User')
        $env:Path = (@($env:Path -split ';' | Where-Object { (Get-NormalizedPathEntry $_) -ine $installRoot }) -join ';')
    }
    # Only these known files are removed. No recursive directory removal.
    $managedPaths = @($executable, $noticesPath, $manifestPath)
    if ($previous.license_sha256) { $managedPaths += $licensePath }
    foreach ($path in $managedPaths) {
        if (Test-Path -LiteralPath $path) { Remove-Item -LiteralPath $path }
    }
    Write-Output "Uninstalled Vyx from $installRoot. Other files were preserved."
    return
}

if (-not $previous -and ((Test-Path -LiteralPath $executable) -or (Test-Path -LiteralPath $noticesPath) -or (Test-Path -LiteralPath $licensePath))) {
    throw 'The target contains unmanaged files; refusing to overwrite them.'
}
if ($previous -and -not $previous.license_sha256 -and (Test-Path -LiteralPath $licensePath)) {
    throw 'The target contains an unmanaged LICENSE; refusing to overwrite it.'
}
$sourceBinary = Join-Path $projectRoot 'dist\vyx.exe'
if (-not (Test-Path -LiteralPath $sourceBinary)) {
    & (Join-Path $PSScriptRoot 'dev.ps1') build
}
if (-not (Test-Path -LiteralPath $sourceBinary)) { throw 'Build did not produce dist/vyx.exe.' }
$sourceNotices = Join-Path $projectRoot 'THIRD_PARTY_NOTICES.txt'
if (-not (Test-Path -LiteralPath $sourceNotices)) { throw 'Third-party notices are missing.' }
$sourceLicense = Join-Path $projectRoot 'LICENSE'
if (-not (Test-Path -LiteralPath $sourceLicense)) { throw 'Project license is missing.' }
$version = & $sourceBinary version
if ($LASTEXITCODE -ne 0 -or $version -notmatch '^vyx ') { throw 'The build is not a valid Vyx executable.' }

New-Item -ItemType Directory -Path $installRoot -Force | Out-Null
Copy-Item -LiteralPath $sourceBinary -Destination $executable -Force
Copy-Item -LiteralPath $sourceNotices -Destination $noticesPath -Force
Copy-Item -LiteralPath $sourceLicense -Destination $licensePath -Force
$addPath = -not $NoPath -and -not (Test-PathEntry $userPath $installRoot)
$pathOwned = [bool]($addPath -or ($previous -and $previous.path_added))
$manifest = [ordered]@{
    schema_version = 1
    product = 'Vyx'
    version = [string]$version
    install_directory = $installRoot
    binary_sha256 = (Get-FileHash -LiteralPath $executable -Algorithm SHA256).Hash.ToLowerInvariant()
    notices_sha256 = (Get-FileHash -LiteralPath $noticesPath -Algorithm SHA256).Hash.ToLowerInvariant()
    license_sha256 = (Get-FileHash -LiteralPath $licensePath -Algorithm SHA256).Hash.ToLowerInvariant()
    path_added = $pathOwned
}
[IO.File]::WriteAllText($manifestPath, ($manifest | ConvertTo-Json), [Text.UTF8Encoding]::new($false))
if ($addPath) {
    $updatedPath = if ([string]::IsNullOrEmpty($userPath)) { $installRoot } else { $userPath.TrimEnd(';') + ';' + $installRoot }
    [Environment]::SetEnvironmentVariable('Path', $updatedPath, 'User')
}
if (-not $NoPath -and -not (Test-PathEntry $env:Path $installRoot)) {
    $env:Path = $installRoot + ';' + $env:Path
}
Write-Output "Installed $version to $executable."
if ($NoPath) {
    Write-Output 'PATH was not changed (-NoPath).'
} else {
    Write-Output 'Run: vyx scan .'
    Write-Output 'This PowerShell session is ready. For existing terminals or apps, restart them to pick up the updated user PATH.'
}
