[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$Target,

    [ValidateSet("saas-web", "fullstack-saas", "expo-mobile")]
    [string]$Profile = "fullstack-saas",

    [Parameter(Mandatory = $true)]
    [string]$Name,

    [Parameter(Mandatory = $true)]
    [string]$Idea,

    [ValidateSet("claude", "codex", "both")]
    [string]$Agent = "both",

    [switch]$Force
)

$ErrorActionPreference = "Stop"
$RepoRoot = Split-Path -Parent $PSScriptRoot
$Args = @(
    "run", "./cmd/autoco", "init",
    "--target", $Target,
    "--profile", $Profile,
    "--name", $Name,
    "--idea", $Idea,
    "--agent", $Agent
)
if ($Force) {
    $Args += "--force"
}

Push-Location $RepoRoot
try {
    & go @Args
    if ($LASTEXITCODE -ne 0) {
        throw "Auto Company initialization failed with exit code $LASTEXITCODE"
    }
}
finally {
    Pop-Location
}
