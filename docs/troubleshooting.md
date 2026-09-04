# Troubleshooting


Debug by layer. Do not change random config before locating the failing layer.

## Page Inaccessible

Check whether the service is running, whether the port is occupied, whether HTTPS is enabled, and whether the access protocol is correct. The usual default URL is:

```text
https://127.0.0.1:8080/
```

For a self-signed certificate, continue past the browser trust warning manually.

## Login Failure

Check the RBAC password (the initial `admin` password is printed at first startup), whether an old session was invalidated after a password change, browser cookies in a private window, and audit throttling for repeated failures.

### Recover a forgotten `admin` password

If another administrator with `rbac:write` is available, reset the password under **Platform permissions → User management**.

If no administrator session is available, change to the project root and run:

```bash
./run.sh --reset-admin-password
```

Enter and confirm the new password when prompted. The script hides input and stores a bcrypt hash. If the service is running, restart it afterward to invalidate existing login sessions.

If `run.sh` is not available, use `sqlite3` and `htpasswd` to update the built-in account. If `database.path` is not the default, replace `data/conversations.db`. Password input should remain hidden and must not be written to shell history.

```bash
HASH=$(htpasswd -nBC 10 '' | cut -d: -f2 | tr -d '\n') && sqlite3 data/conversations.db "UPDATE rbac_users SET password_hash='$HASH', updated_at=CURRENT_TIMESTAMP WHERE id='admin' AND username='admin' AND is_builtin=1; SELECT changes();"
```

Output `1` means that the row was updated.

## Model No Response

Check the selected channel, `base_url` path, API key, model name, and whether the provider accepts the channel's reasoning fields. Use the model test in System Settings; if the gateway returns 400, try:

```yaml
ai:
  channels:
    your-channel:
      reasoning:
        mode: off
```

## Interrupted Streaming Output

Nginx or another reverse proxy must not buffer SSE. For Nginx, use:

```nginx
proxy_buffering off;
proxy_http_version 1.1;
```

Also check model gateway timeouts, network interruptions, and oversized context.

## Tool Execution Failure

Check that the command is installed in `PATH`, the `tools/*.yaml` schema is valid, HITL did not reject the call, `agent.tool_timeout_minutes` was not exceeded, and `shell_no_output_timeout_seconds` did not terminate a quiet Shell process.

## MCP Connection Failure

For built-in MCP, check `mcp.enabled`, `mcp.port`, `auth_header`, and `auth_header_value`. For external MCP, check stdio command paths, working directories, environment variables, or HTTP/SSE URLs, authentication, and network reachability. Inspect `/api/external-mcp/stats`.

## Knowledge Base Unavailable

Check `knowledge.enabled: true`, embedding settings, whether scanning and index rebuild completed, whether `data/knowledge.db` is writable, and whether the embedding service returns 429 or timeout errors. For many indexing failures, lower `batch_size` and increase `rate_limit_delay_ms`:

```yaml
knowledge:
  indexing:
    batch_size: 5
    rate_limit_delay_ms: 600
```

## Robot Does Not Reply

Check `robots.<platform>.enabled`, the platform callback URL, token/secret/verify-token values, platform reachability to the server, and whether group messages require an @ mention. See [Robot Guide](robot.md).

## C2 Listener Startup Failure

Check `c2.enabled`, port conflicts, firewall/security-group rules, and whether binding a low port requires administrator privileges. When C2 is disabled, `/api/c2/*` returning 503 is expected.

## WebShell Output Encoding

Confirm the target encoding, try a shorter command, use base64-wrapped output where appropriate, and check the code page on Windows targets.

## Database Lock or Write Failure

Check whether `data/` is writable, whether multiple instances share one SQLite file, whether the disk is full, and whether WAL/SHM files were copied inconsistently. Do not let multiple production processes write the same SQLite database.

## Frontend Page Failure

Check browser-console errors, static-resource loading, stale browser cache after frontend changes, and missing i18n keys. When an API fails, compare the request body with `/api-docs`.

## Diagnostic Order

1. Process: is the service alive, any panic?
2. Network: port, HTTPS, reverse proxy, browser console.
3. Auth: does `/api/auth/validate` return 200?
4. Config: can `/api/config` be read and applied?
5. Model: does model test pass?
6. Tools: do tool list and schemas look right?
7. Database: is `data/` writable, any lock?
8. Subsystem: KB, MCP, C2, WebShell minimal action.

## Minimal Commands

```bash
# Process and ports
lsof -i :8080

# Local HTTPS check
curl -k -I https://127.0.0.1:8080/

# Static-resource check
curl -k -I https://127.0.0.1:8080/static/logo.png

# Database files
ls -lh data/
```

If a reverse proxy is involved, test both proxy address and upstream address.

## Common Misdiagnoses

- "Model is broken": HITL is waiting.
- "Tool missing": tool_search hides it from current context.
- "Knowledge base useless": index not rebuilt or risk type too narrow.
- "Config saved but ineffective": listener/TLS changes need restart.
- "Robot silent": platform callback or signature config wrong.

## Issue Template

```text
Version:
Startup method:
Access path:
Relevant config:
Steps:
Expected:
Actual:
Server logs:
Browser console:
API response:
```
