#!/bin/sh
set -e

# CyberStrikeAI container entrypoint
#
# - Bootstraps /app/config.yaml from config.example.yaml on first run.
# - If CYBERSTRIKE_ADMIN_PASSWORD is set, applies it as the built-in admin
#   password (initial bootstrap on a fresh data dir, reset otherwise).
# - Then starts the server, forwarding the container CMD/args.

cd /app

# 1) Ensure a config file exists.
if [ ! -f /app/config.yaml ]; then
  echo "[entrypoint] config.yaml not found; creating from config.example.yaml"
  cp /app/config.example.yaml /app/config.yaml
fi

# 2) Apply the default admin password when requested.
if [ -n "${CYBERSTRIKE_ADMIN_PASSWORD:-}" ]; then
  echo "[entrypoint] applying admin password from CYBERSTRIKE_ADMIN_PASSWORD"
  /app/set-admin-password -config /app/config.yaml -password "$CYBERSTRIKE_ADMIN_PASSWORD"
else
  echo "[entrypoint] CYBERSTRIKE_ADMIN_PASSWORD not set; leaving admin password as-is"
fi

# 3) Start the server.
echo "[entrypoint] starting CyberStrikeAI: ./cyberstrike-ai -config /app/config.yaml $*"
exec /app/cyberstrike-ai -config /app/config.yaml "$@"
