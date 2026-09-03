# MCP Federation


CyberStrikeAI uses MCP as the primary tool protocol. Tools can be built-in, YAML-backed, Skill-local, or provided by external MCP servers.

## Built-In MCP

The internal MCP server registers:

- YAML command tools;
- security execution tools;
- knowledge tools;
- project fact tools;
- C2 tools;
- WebShell tools;
- batch task tools;
- vision analysis.

Agents usually call these internally without extra setup.

## HTTP MCP

```yaml
mcp:
  enabled: true
  host: 0.0.0.0
  port: 8081
  auth_header: "X-MCP-Token"
  auth_header_value: "random-secret"
```

Always set an auth value and restrict network access.

## Web MCP Endpoint

The standalone HTTP MCP endpoint is controlled by `mcp.enabled`, `mcp.host`, `mcp.port`, `auth_header`, and `auth_header_value`. Treat it as an administrative integration surface: use a random secret, HTTPS or a trusted network boundary, and do not expose it publicly without a specific authorization design.

```text
POST /api/mcp
Authorization: X-MCP-Token: random-secret
```

## External MCP Lifecycle

External MCP configuration can be set to:

```yaml
external_mcp:
  servers: {}
```

It can also be created, started, stopped, and deleted from the Web MCP management page. Available endpoints include:

- `GET /api/external-mcp`
- `GET /api/external-mcp/stats`
- `GET /api/external-mcp/:name`
- `PUT /api/external-mcp/:name`
- `POST /api/external-mcp/:name/start`
- `POST /api/external-mcp/:name/stop`
- `DELETE /api/external-mcp/:name`

1. Register config: name, type, command/URL, environment.
2. Start connection: stdio process or HTTP/SSE client.
3. Pull tool list: names, descriptions, schemas.
4. Expose to Agent: affected by role, tool_search, HITL.
5. Execute: validate args, call, monitor.
6. Recover: handle process/network failure.
7. Stop/delete: remove runtime and config.

Debug by locating the failed step.

## stdio

stdio providers run as local processes. Validate the executable path, working directory, environment variables, startup output, process ownership, and shutdown behavior. Keep credentials out of command-line arguments and logs.

## HTTP / SSE

HTTP/SSE providers need a reachable URL, correct authentication, compatible protocol behavior, and network timeout/retry handling. Verify both the SSE stream and message endpoint when the provider exposes them.

For example, an SSE provider can be represented as `external_mcp.servers.trusted-service` with `transport: sse`, a reachable `url`, and an `auth_header`.

## Tool Exposure Policy

External tools are affected by role boundaries, `tool_search`, and HITL. Keep tools that are not needed out of the visible set, and explicitly keep high-risk tools behind approval.

```yaml
multi_agent:
  eino_middleware:
    tool_search_enable: true
    tool_search_min_tools: 20
    tool_search_always_visible: 12
    tool_search_always_visible_tools:
      - read_file
      - glob
      - grep
      - tool_search
```

Keep frequently used tools always visible and dynamically unlock the rest through `tool_search`.

**Tool Naming**

Good names are stable, specific, and action-object oriented, such as `burp_send_to_repeater`, `asset_lookup_domain`, and `cloud_list_public_buckets`. Avoid generic names such as `run`, `execute`, `scan`, and `tool1`.

Avoid generic names:

```text
run
execute
scan
tool1
```

Prefer names such as:

```text
burp_send_to_repeater
asset_lookup_domain
cloud_list_public_buckets
```

Specific names improve tool_search and reduce misuse.

## Security Review

Before connecting an external MCP, ask:

- Can it read/write local files?
- Can it execute commands?
- What network does it access?
- Does it send data to third parties?
- Are tool descriptions trustworthy?
- Can output contain prompt injection?
- Should it run under a separate OS user or container?

## Debugging

Check these in order:

1. Inspect `/api/external-mcp/stats` for provider state.
2. Check service logs.
3. Run the stdio command separately.
4. Use curl to test the HTTP/SSE address.
5. Check whether roles or `tool_search` hide the tool.

Then inspect schema validity, invocation arguments, execution monitoring, and recovery events. Compare the external provider's logs with CyberStrikeAI's MCP and monitor logs.

## MCP Lifecycle

The lifecycle is register configuration, start a stdio or HTTP/SSE connection, retrieve tools, expose approved tools to the Agent, validate and execute calls, recover failures, and stop/delete the provider. Each transition should be observable and auditable.

## Tool Naming Conventions

Use stable, specific, action-object names such as `burp_send_to_repeater`, `asset_lookup_domain`, and `cloud_list_public_buckets`. Avoid generic names such as `run`, `execute`, `scan`, or `tool1`; precise names improve `tool_search` and reduce misuse.

## External MCP Security Review Checklist

Before connecting a provider, record its local file and command access, network destinations, third-party data flow, trustworthiness of descriptions, prompt-injection exposure, operating-system isolation, credential handling, timeout/retry behavior, and removal procedure.

If any answer is unclear, do not keep the provider in the production tool pool.

## Source Anchors

- External manager: `internal/mcp/external_manager.go`
- Recovery: `internal/mcp/connection_recovery.go`
- Tool adapter: `internal/einomcp/mcp_tools.go`
- Handler: `internal/handler/external_mcp.go`
- Invoke notification: `internal/einomcp/tool_invoke_notify.go`
