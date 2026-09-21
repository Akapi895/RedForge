@echo off
rem ============================================================================
rem CyberStrikeAI - Windows preparation launcher
rem
rem Runs docker\prepare.ps1 with execution-policy bypass so that regular
rem Windows users (no WSL required) can restore antivirus-flagged files before
rem building with Docker. Double-click this file, or run it from a terminal:
rem     docker\prepare.bat
rem ============================================================================
setlocal

rem Locate the repo root (directory above this script's docker\ folder).
set "SCRIPT_DIR=%~dp0"
for %%I in ("%SCRIPT_DIR%..") do set "ROOT_DIR=%%~fI"

echo [prepare.bat] working in: %ROOT_DIR%
powershell -NoProfile -ExecutionPolicy Bypass -File "%SCRIPT_DIR%prepare.ps1"
set "EXITCODE=%ERRORLEVEL%"

if not "%EXITCODE%"=="0" (
    echo.
    echo [prepare.bat] Preparation did not complete successfully.
    echo              Add this repo path to your antivirus exclusion list, then re-run.
) else (
    echo [prepare.bat] Done. You can now build:  docker compose up -d --build
)

exit /b %EXITCODE%
