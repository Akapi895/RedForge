# Plugin Development


Plugins live under `plugins/`. The repo ships two reference implementations: **Burp Suite extension** and **Chromium DevTools extension**. Integrations typically use HTTP APIs, MCP servers, or resource packs (tools, roles, Skills, agents).

## Layout

```text
plugins/
  README.md
  burp-suite/
    cyberstrikeai-burp-extension/
      src/main/java/burp/
      README.md
      README.zh-CN.md
      build.gradle
      pom.xml
  browser-extension/
    cyberstrikeai-browser-extension/
      manifest.json
      devtools.js
      background/service-worker.js
      panel/
      popup/
      lib/
      README.md
      README.zh-CN.md
      package.sh
```

## Plugin Types

Common integration forms are browser/security-tool extensions calling the CyberStrikeAI API, MCP Servers exposing tools, file-based extensions providing `tools`, `roles`, `skills`, or `agents`, and Webhook/robot callbacks.

Do not start with MCP unless the Agent must actively call your capability. If the requirement is only to send an HTTP request from Burp or the browser to AI for analysis, an API plugin is more direct; use MCP when the Agent must actively invoke scanning or query tools.

Other common integrations include Webhook/robot callbacks and file-based extensions that provide `tools`, `roles`, `skills`, or `agents`.

## Burp Suite Extension

Java extension under `plugins/burp-suite/cyberstrikeai-burp-extension/`. Typical flow: read HTTP from Burp → format prompt → call CyberStrikeAI SSE → show Progress/Final in a Burp tab.

The plugin reads HTTP requests and responses from Burp, formats messages, calls the CyberStrikeAI API, and displays AI analysis in a Burp tab. Before building, verify that JDK and Gradle or Maven are available and that the service URL and authentication settings are correct. Build with the repository's Gradle/Maven script and confirm the generated JAR in `dist/`.

## Browser Extension (Chromium DevTools)

MV3 DevTools extension under `plugins/browser-extension/cyberstrikeai-browser-extension/`. Aligned with the Burp plugin: capture Network traffic → HTTP/1.1 prompt → SSE output. Full installation and UI documentation is in `plugins/browser-extension/cyberstrikeai-browser-extension/README.md`.

Load unpacked at `chrome://extensions/`, or `bash package.sh` → `dist/cyberstrikeai-browser-extension.zip`.

The extension can capture XHR/Fetch traffic from the DevTools Network panel, pause capture, keep raw HAR in memory, normalize displayed and prompt content to HTTP/1.1, call login/validate/Agent Stream APIs, and show Progress/Final in the DevTools panel. The popup only displays connection status.

### Auth best practices (browser)

Server `POST /api/auth/login` returns `{ token, expires_at }`. There is **no refresh token** — do not assume silent renewal. Reference: `lib/auth-session.js`, `lib/api.js`, `panel/panel.js`.

| Practice | Description |
| --- | --- |
| Session storage | Store token + `expires_at` in `chrome.storage.session`; never persist password |
| Remaining time | Show `OK · 11h 30m left`; warn when <30min |
| Local check | Re-check `expires_at` + `GET /api/auth/validate` every 30s |
| Server probe | Immediate probe when DevTools panel becomes visible |
| Unreachable | Show warning; keep token during transient outage |
| 401/403 | Clear token (server restart clears in-memory sessions) |
| Before Send | `ensureAuthReady()` before SSE |
| Permissions | `optional_host_permissions` — request origin on Validate |

After extension reload, close DevTools completely and reopen F12 (stale panel context).

### Data and performance (browser)

- Caps: 200 captures/tab, 20 tabs, 512KB progress/run.
- Default XHR/Fetch only; use pause toggle when not capturing.
- Truncate or summarize large bodies before sending to Agent.

## API Integration

- Login: `POST /api/auth/login`, then `GET /api/auth/validate`.
- Persist `expires_at`; re-login when expired (no silent refresh).
- Prefer `/api/eino-agent/stream` or `/api/multi-agent/stream` (SSE).
- Large files: `/api/chat-uploads`, then reference in message.
- Findings can be written through `/api/vulnerabilities`, and project facts through `/api/projects/:id/facts`.
- Full spec: `/api-docs` or `/api/openapi/spec`.

Use `/api-docs` as the authoritative interface reference.

## MCP Plugins

When a plugin must add tools that the Agent actively calls, implement an MCP Server and connect it through external MCP management:

- stdio for a locally started process;
- HTTP/SSE for a long-running service.

MCP tools should use explicit schemas, minimal parameters, stable structured output, readable errors, and separate high-risk actions so HITL can review them.

## File-Based Extensions

A plugin can also ship:

- `tools/*.yaml`;
- `roles/*.yaml`;
- `skills/<name>/SKILL.md`;
- `agents/*.md`.

This is simple and reliable for preserving internal methodology or toolchains.

## Release Checklist

Before releasing a plugin, confirm that it contains no API keys, cookies, or target information; includes installation, configuration, and uninstall instructions; displays clear errors; matches the current CyberStrikeAI API; and clearly documents high-risk capabilities.

## Version Compatibility

Avoid depending on unpublished frontend internals. Prefer `/api/openapi/spec`, stable HTTP APIs, the MCP protocol, and documented directory conventions. If an internal endpoint is unavoidable, document the compatible version in the plugin README.

## Plugin Layers

| Layer | Example | Benefit | Cost |
| --- | --- | --- | --- |
| API plugin | Burp / browser extension calling Agent Stream | simple UI integration | depends on API/auth |
| MCP plugin | exposes tools to Agent | Agent can call it | needs schema and safety design |
| Resource pack | ships tools/roles/skills/agents | simple and versionable | less interactive |

If the requirement is only to send an HTTP request to AI, an API plugin is enough. If the Agent must actively call a scan or query capability, use an MCP plugin.

## API Plugin Payload

Include the source tool and context, target URL and method, key headers, request/response truncation policy, user intent, and authorization boundary. Upload or summarize large responses instead of pasting them whole into the prompt.

## MCP Schema Design

Bad:

```json
{"cmd":{"type":"string"}}
```

Better:

```json
{
  "target_url": {"type":"string","description":"authorized target URL"},
  "scan_profile": {"type":"string","enum":["passive","active-safe"]},
  "max_requests": {"type":"integer","description":"request limit"}
}
```

Specific schemas make HITL and Agent behavior safer.

## Security Boundaries

Plugins should not bypass platform controls:

- no hidden destructive local commands;
- no plaintext long-lived credentials (password only for login; token in session storage);
- no default third-party data exfiltration;
- no dependency on browser state to bypass login;
- on 401/403, clear session and require re-auth — do not silently retry.

## Source Anchors

- Burp plugin: `plugins/burp-suite/cyberstrikeai-burp-extension/src/main/java/burp/`
- Browser extension: `plugins/browser-extension/cyberstrikeai-browser-extension/`
  - Auth: `lib/auth-session.js`, `lib/api.js`, `lib/storage.js`
  - UI: `panel/panel.js`
  - Capture: `devtools.js`, `background/service-worker.js`
- OpenAPI: `internal/handler/openapi.go`
- External MCP: `internal/handler/external_mcp.go`
- Web auth reference: `web/static/js/auth.js`
