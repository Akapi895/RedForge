# Deployment Guide


CyberStrikeAI can run as a local testing tool, an internal team service, or a production red-team platform. Treat it as a high-privilege security system: it can execute commands, call MCP tools, manage WebShell connections, and optionally run C2 listeners.

## Prerequisites

- Go for source runs and binary builds.
- Python for some MCP servers and tool scripts.
- SQLite files under `data/`; no external DB is required by default.
- Actual security tools installed in PATH. YAML files under `tools/` only describe commands.
- At least one `ai.channels` entry. Use `provider: openai_compatible` for OpenAI-compatible endpoints, or `provider: claude` for Eino's native Claude component.

Important persistent paths:

```text
config.yaml
data/
tools/
roles/
skills/
agents/
knowledge_base/
chat_uploads/
```

Back these up before upgrades.

## Quick Start

Local quick start:

```bash
chmod +x run.sh && ./run.sh
```

`run.sh` is the most common startup path for local use, development, small temporary internal deployments, and quick post-upgrade verification.

The default configuration normally enables TLS with a self-signed certificate. Open `https://127.0.0.1:8080/`; a browser security warning is expected for local self-signed testing. Use a real certificate in production.

For long-running service, boot-time startup, managed logs, and crash recovery, prefer a binary managed by systemd.

## Source Run

```bash
go run ./cmd/server --config config.yaml
```

If dependency downloads are slow, configure a Go proxy first:

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

## Binary Build

```bash
go build -o cyberstrike-ai ./cmd/server
./cyberstrike-ai --config config.yaml
```

The binary still needs `web/templates`, `web/static`, and the runtime resource directories.

When distributing a binary, also include `tools/`, `roles/`, `skills/`, `agents/`, and `config.yaml`.

## HTTPS

For local testing, self-signed HTTPS is acceptable:

```yaml
server:
  tls_enabled: true
  tls_auto_self_sign: true
```

For production, use certificate files:

```yaml
server:
  host: 0.0.0.0
  port: 8080
  tls_enabled: true
  tls_cert_path: /etc/letsencrypt/live/example.com/fullchain.pem
  tls_key_path: /etc/letsencrypt/live/example.com/privkey.pem
```

When TLS is enabled, HTTP requests on the same port redirect to HTTPS with status 308. If a reverse proxy terminates TLS, disable application TLS and handle HTTPS at the proxy.

## Reverse Proxy

For production, use real certificates or terminate TLS at a reverse proxy. If the proxy terminates TLS and forwards HTTP to the app, avoid enabling app-side TLS on the same upstream unless `proxy_pass` uses HTTPS.

Nginx must not buffer SSE:

```nginx
proxy_buffering off;
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";
```

Example Nginx configuration:

```nginx
server {
    listen 443 ssl http2;
    server_name cyberstrike.example.com;

    ssl_certificate /etc/letsencrypt/live/cyberstrike.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/cyberstrike.example.com/privkey.pem;
    client_max_body_size 200m;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_buffering off;
    }
}
```

`proxy_buffering off` is important for SSE output and WebSocket terminals.

## systemd

```ini
[Unit]
Description=CyberStrikeAI
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/CyberStrikeAI
ExecStart=/opt/CyberStrikeAI/cyberstrike-ai --config /opt/CyberStrikeAI/config.yaml
Restart=on-failure
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

Enable it with:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now cyberstrikeai
sudo journalctl -u cyberstrikeai -f
```

## Data and Backups

Back up `config.yaml`, both database files, `data/eino-checkpoints/`, custom resource directories, and `chat_uploads/`. For SQLite hot backups, stop the service first when possible, or copy the matching `*.db`, `*.db-wal`, and `*.db-shm` files together.

## Upgrades

1. Stop the service.
2. Back up `config.yaml`, `data/`, and custom directories.
3. Pull or replace the new code or binary.
4. Preserve the existing configuration and add new fields from the new sample.
5. Start the service and check login, model test, tool list, and knowledge-base status.

The repository provides `upgrade.sh` for quick upgrades without known compatibility issues; production deployments should still back up first.

## Rollback

Restore the previous binary or code, pre-upgrade `config.yaml`, and pre-upgrade `data/` together. If the new version changed the database schema, restore the database backup; replacing only the binary may not be enough.

## Deployment Decision Table

| Scenario | Recommended setup | Key settings | Avoid |
| --- | --- | --- | --- |
| Personal testing | `./run.sh` + self-signed HTTPS | `tls_auto_self_sign: true` | Public exposure |
| Internal team | Binary + systemd + internal HTTPS | strong password, audit, backup, IP restrictions | Shared weak password |
| Production red-team platform | Reverse proxy + dedicated OS user + log collection | real certs, proxy auth, C2 only when needed | Direct public admin UI |
| Chat/KB only | Disable C2 and unnecessary MCP | `c2.enabled: false` | All tools enabled by default |
| Tool automation | Isolated workspace + HITL | `workspace_root_dir`, `hitl`, `monitor` | Shell tools globally allowlisted |

## Runtime File Layers

- Replaceable: binary, `web/`, default docs/resources.
- Preserve: `config.yaml`, `data/`, custom tools/roles/skills/agents, `knowledge_base`, uploads.
- Cleanup candidates: checkpoints, temporary workspaces, stale payloads, old tool execution records.

When replacing the whole deployment directory, move custom resources and `data/` out first. Many “configuration lost after upgrade” incidents are caused by overwriting runtime files as if they were part of the release package.

## Acceptance Checklist

After startup:

1. Open `/` and verify no HTTP/HTTPS redirect loop.
2. Login and validate `/api/auth/validate`.
3. Run model test in settings.
4. Check tool list and schemas.
5. If KB is enabled, check index status.
6. If external MCP is enabled, verify connection and tool visibility.
7. If C2 is enabled, start and stop a test listener only in an authorized network.
8. Check audit logs for login and config activity.

Startup success does not prove that the deployment is usable. Also verify that the `tools/` schemas load, the knowledge-base `index-status` is healthy, external MCP tools are visible, and an authorized C2 test listener can be stopped and removed.

## Common Reverse-Proxy Pitfalls

- Buffered SSE makes the Agent appear silent until the request ends; disable `proxy_buffering`.
- WebSocket failures usually indicate missing `Upgrade` or `Connection` headers.
- When app-side and Nginx TLS are both enabled, the `proxy_pass` protocol must match the upstream.
- Upload failures may require a larger `client_max_body_size` and application-side upload limit.
- A 308 loop occurs when app-side same-port HTTPS redirect is combined with HTTP proxying; disable app TLS or proxy to HTTPS.

## Source Anchors

- App wiring and routes: `internal/app/app.go`
- TLS bootstrap: `internal/app/main_server_tls.go`
- HTTP to HTTPS redirect: `internal/app/main_server_http_redirect.go`
- Config structs: `internal/config/config.go`
- Config apply: `internal/handler/config.go`
