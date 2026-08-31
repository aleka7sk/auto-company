[CmdletBinding()]
param(
    [string]$InstallRoot = "$HOME\.auto-company",
    [switch]$AddToPath,
    [switch]$InstallCodexSkills
)

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent $PSScriptRoot
$BinDir = Join-Path $InstallRoot "bin"
$Binary = Join-Path $BinDir "autoco.exe"

Write-Host "Validating Auto Company..."
Push-Location $RepoRoot
try {
    go test ./...
    node scripts/validate-plugin.mjs
    node scripts/test-guard.mjs

    New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
    go build -trimpath -o $Binary ./cmd/autoco
}
finally {
    Pop-Location
}

Write-Host "Installed CLI: $Binary"

if ($AddToPath) {
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $Entries = @($UserPath -split ";" | Where-Object { $_ })
    if ($Entries -notcontains $BinDir) {
        $NewPath = (($Entries + $BinDir) | Select-Object -Unique) -join ";"
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        Write-Host "Added $BinDir to user PATH. Restart the terminal."
    }
    else {
        Write-Host "$BinDir is already in user PATH."
    }
}
else {
    Write-Host "PATH was not changed. Add this directory manually when ready: $BinDir"
}

if ($InstallCodexSkills) {
    $CodexSkills = Join-Path $HOME ".agents\skills"
    New-Item -ItemType Directory -Force -Path $CodexSkills | Out-Null

    Get-ChildItem (Join-Path $RepoRoot "skills") -Directory | ForEach-Object {
        $Destination = Join-Path $CodexSkills ("auto-company-" + $_.Name)
        if (Test-Path $Destination) {
            Remove-Item -Recurse -Force $Destination
        }
        Copy-Item -Recurse -Force $_.FullName $Destination
    }
    Write-Host "Installed Auto Company skills for Codex under $CodexSkills"
}

Write-Host ""
Write-Host "Claude Code local plugin:"
Write-Host "  claude --plugin-dir `"$RepoRoot`""
Write-Host ""
Write-Host "Persistent Claude marketplace install (run inside Claude Code):"
Write-Host "  /plugin marketplace add aleka7sk/auto-company"
Write-Host "  /plugin install auto-company@auto-company-marketplace"
