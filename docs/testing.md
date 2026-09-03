# Testing Guide


Testing CyberStrikeAI means more than running Go tests. Agent, MCP, HITL, C2, WebShell, and frontend streaming all have different failure modes.

## Go Unit Tests

```bash
go test ./internal/...
```

Run focused package tests while developing and the full internal suite before release:

```bash
go test ./internal/workflow
go test ./internal/multiagent
go test ./internal/handler
```

Important packages include `internal/security`, `internal/mcp`, `internal/multiagent`, `internal/workflow`, `internal/knowledge`, `internal/project`, `internal/handler`, and `internal/c2`. Add tests for new behavior, migration compatibility, error paths, and security boundaries rather than only happy paths.

## Build Tests

When entrypoints, dependencies, or build assets change, run:

```bash
go test ./cmd/...
go build -o cyberstrike-ai ./cmd/server
```

## Configuration Validation

Before startup, verify YAML indentation, model settings, a writable database path, existing `tools_dir`, `roles_dir`, `skills_dir`, and `agents_dir`, and valid HTTPS certificate paths. After startup, use the Web Settings page to test an OpenAI-compatible model, the vision model, the tool list, and external MCP status. Start with a temporary configuration, load it, apply it through the Web API, and confirm that old configurations still start.

## Manual API Tests

Use `/api-docs` to obtain the current request and response schemas. Test login, token validation, representative read/write endpoints, expected 401/403/404/409 responses, and SSE `error`/`done` behavior.

At minimum, exercise `/api/eino-agent/stream`, `/api/config`, `/api/config/tools`, `/api/knowledge/search`, and `/api/monitor`. When a streaming endpoint is behind a reverse proxy, verify that output is delivered in real time.

```text
login -> streaming chat -> settings apply -> tool list -> HITL -> knowledge base -> external MCP
```

## Tool Tests

After adding or changing `tools/*.yaml`, verify the schema in the tool list, execute it with harmless parameters, confirm readable errors and timeout behavior, and confirm that HITL blocks the expected operations. Do not test new tools against production targets. Also verify input schema validation, no-output timeout, audit behavior, output persistence, cancellation, and duplicate-event handling.

## MCP Tests

For external MCP, run stdio commands independently in a terminal, use curl for HTTP/SSE connectivity, verify `/api/external-mcp/stats` after starting the provider from the Web page, and confirm that conversation-time `tool_search` can find the tools. Test built-in MCP and both external transports for tool-list retrieval, invocation, process/network failure, reconnection, stop, and removal.

## Knowledge-Base Tests

Use a small Markdown corpus to test scan, index rebuild, keyword and synonym search, chunking, embedding failures, reranking, empty results, thresholds, retrieval logs, and index rebuilds. When using a real embedding API, account for quota and rate limits.

## Frontend Smoke Tests

Manually verify login/logout, the conversation sidebar, new conversations and streaming replies, settings save/apply, tool listing, related CRUD pages, language switching, upload/file management, and browser-console errors after frontend changes.

## High-Risk Module Tests

Test C2, WebShell, Terminal, external MCP write tools, and batch operations only in an authorized lab. Confirm that the target is a local machine, training target, or exercise environment, commands are non-destructive, HITL matches policy, and sessions, payloads, uploaded files, and task results are cleaned up afterward. Also confirm audit, cancellation, disabled-mode responses, and network isolation.

## Test Pyramid

| Layer | Goal | Example |
| --- | --- | --- |
| Unit | pure logic | expressions, chunking, sanitization |
| Handler | HTTP behavior | validation, auth, status codes |
| Integration | module cooperation | external MCP, KB indexing, HITL |
| Smoke | user path | login, chat, tools, settings |
| Authorized lab | high-risk features | C2, WebShell, terminal |

Do not use end-to-end manual testing as a substitute for unit tests, or unit tests as a substitute for high-risk lab validation.

## Regression Focus

Expand testing when changing:

- `internal/handler/config.go`: model, KB, MCP, C2, robot apply paths;
- `internal/multiagent/`: streaming, tool calls, summarization, retry, HITL;
- `internal/security/`: auth, shell, timeout, no-output;
- `internal/database/`: old data compatibility;
- `web/static/js/chat.js`: chat, process details, attack chain, groups.

Also regression-test finalization gates, model retry/failover, tool visibility, role boundaries, workflow resume, and frontend handling of non-final `response` events.

## Test Data

Avoid real customer data. Prepare:

- small Markdown KB sample;
- fake local MCP server;
- controlled local HTTP target;
- harmless WebShell simulator;
- temporary SQLite DB.

Delete temporary databases and uploaded files after testing so the development environment is not polluted.

## Failure Cases

Cover:

- model API 401/429/500;
- MCP startup failure;
- tool timeout;
- HITL rejection;
- interrupted KB indexing;
- unwritable database;
- WebShell non-200 response;
- C2 disabled endpoint access.

Failure cases are more valuable than success-only cases: a system that fails safely is easier to operate and review than one that is merely fast on the happy path.

## Source Anchors

Existing tests live across:

- `internal/handler/*_test.go`
- `internal/multiagent/*_test.go`
- `internal/workflow/*_test.go`
- `internal/knowledge/*_test.go`
- `internal/security/*_test.go`
- `internal/mcp/*_test.go`
- `internal/c2/*_test.go`
