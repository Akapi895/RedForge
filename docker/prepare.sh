#!/bin/sh
# ============================================================================
# CyberStrikeAI — environment preparation helper (portable)
#
# Some antivirus / EDR products (e.g. Windows Defender) flag legitimate
# offensive-security source files in this repo as malware and delete them
# from the working tree (false positives). This script restores them from git
# so the Docker build (or ./run.sh) works on any machine.
#
# Usage (from the repo root):
#   ./docker/prepare.sh
#
# Requires: a git clone of CyberStrikeAI with a clean history containing
# these files. If a file is still detected after this runs, add the repo path
# to your antivirus exclusion list, then re-run this script.
# ============================================================================

set -u

# Resolve repo root (parent of the docker/ directory this script lives in).
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR" || exit 1

# Files that are known to be flagged by antivirus false positives.
restored=0
missing_still=0

echo "[prepare] working in: $ROOT_DIR"

# Avoid "dubious ownership" when the repo is shared/mounted (e.g. WSL -> Windows
# drives, or vice versa). Ignore failures: the config entry may already exist.
git config --global --add safe.directory "$ROOT_DIR" 2>/dev/null || true

# Iterate line-by-line (not word-split) so paths with spaces are handled.
# Feed from a here-doc (not a pipe) so the counter variables above propagate.
while IFS= read -r f; do
  [ -n "$f" ] || continue
  if [ -f "$f" ]; then
    echo "[prepare] OK    $f"
    continue
  fi
  if git checkout -- "$f" 2>/dev/null; then
    if [ -f "$f" ]; then
      echo "[prepare] RESTORED  $f"
      restored=$((restored + 1))
    else
      echo "[prepare] RESTORED but still missing (antivirus re-deleted?)  $f"
      missing_still=$((missing_still + 1))
    fi
  else
    echo "[prepare] FAILED to restore  $f  (not tracked by git?)"
    missing_still=$((missing_still + 1))
  fi
done <<'EOF'
internal/c2/payload_oneliner.go
knowledge_base/SQL Injection/MySQL Injection.md
knowledge_base/SQL Injection/SQLite Injection.md
EOF

echo ""
echo "[prepare] summary: restored=$restored still_missing=$missing_still"

if [ "$missing_still" -gt 0 ]; then
  echo ""
  echo "[prepare] Some files could not be kept. Add this repo path to your"
  echo "          antivirus exclusion list, then re-run: ./docker/prepare.sh"
  exit 1
fi

echo "[prepare] done. You can now build: docker compose up -d --build"
exit 0
