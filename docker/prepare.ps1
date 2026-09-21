# ============================================================================
# CyberStrikeAI - Windows environment preparation helper
#
# Windows-native equivalent of ./docker/prepare.sh
#
# Some antivirus products (e.g. Windows Defender) flag legitimate
# offensive-security source files in this repo as malware and delete them
# from the working tree (false positives, e.g. internal/c2/payload_oneliner.go
# and knowledge_base/SQL Injection/*.md). This restores them from git so the
# Docker build works on Windows (no WSL required).
#
# Usage (from the repo root), in PowerShell:
#   .\docker\prepare.ps1
# or just double-click / run in a terminal. You may need a global execution
# policy override once:  Set-ExecutionPolicy -Scope Process Bypass
#
# Requires: a git clone of CyberStrikeAI with a clean history containing
# these files.
# ============================================================================

$ErrorActionPreference = 'Stop'

# Locate the repo root (parent of the docker\ directory this script lives in).
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir   = Split-Path -Parent $scriptDir
Set-Location $rootDir

Write-Host "[prepare] working in: $rootDir"

# Avoid git "dubious ownership" when the repo is shared/mounted (e.g. accessed
# over a UNC/WSL path or a different user). Failure is fine if it already exists.
& git config --global --add safe.directory $rootDir 2>$null | Out-Null

# Files known to be flagged by antivirus false positives.
$flaggedFiles = @(
    'internal\c2\payload_oneliner.go',
    'knowledge_base\SQL Injection\MySQL Injection.md',
    'knowledge_base\SQL Injection\SQLite Injection.md'
)

$restored    = 0
$stillMissing = 0

foreach ($f in $flaggedFiles) {
    if (Test-Path -LiteralPath $f -PathType Leaf) {
        Write-Host "[prepare] OK    $f"
        continue
    }

    # Restore from git (use full path so git resolves spaces correctly).
    & git checkout -- $f
    if ($LASTEXITCODE -eq 0 -and (Test-Path -LiteralPath $f -PathType Leaf)) {
        Write-Host "[prepare] RESTORED  $f"
        $restored++
    } else {
        Write-Host "[prepare] FAILED to restore  $f  (not tracked by git?)"
        $stillMissing++
    }
}

Write-Host ""
Write-Host "[prepare] summary: restored=$restored still_missing=$stillMissing"

if ($stillMissing -gt 0) {
    Write-Host ""
    Write-Host "[prepare] Some files could not be kept. Add this repo path to your"
    Write-Host "          antivirus exclusion list, then re-run: .\docker\prepare.ps1"
    exit 1
}

Write-Host "[prepare] done. You can now build: docker compose up -d --build"
exit 0
