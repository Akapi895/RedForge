# Developer Guide

This guide describes the repository layout, local development, and the complete path for adding a business capability.

## Project Layout

```text
cmd/server/              Web server entrypoint
internal/app/            application assembly, route and MCP-tool registration
internal/handler/        HTTP handlers
internal/database/       SQLite data access
internal/security/       authentication, rate limiting, Shell execution
internal/mcp/             MCP Server and external MCP management
internal/multiagent/     Eino single-agent, multi-agent, and middleware
internal/workflow/       workflow runtime
internal/knowledge/      knowledge-base indexing and retrieval
internal/c2/             built-in C2
internal/project/        project-fact blackboard
web/static/              frontend JavaScript/CSS/resources
web/templates/           HTML templates
tools/                   command-tool YAML
roles/                   role YAML
agents/                  multi-agent Markdown definitions
skills/                  Agent Skills
docs/                    project documentation
```

- `cmd/server/`: server entrypoint.
- `internal/app/`: application assembly and route registration.
- `internal/handler/`: HTTP handlers.
- `internal/agent/`, `internal/multiagent/`, `internal/agents/`: Agent execution and definitions.
- `internal/mcp/`, `internal/einomcp/`: MCP services and adapters.
- `internal/workflow/`: workflow graph and execution.
- `internal/knowledge/`: knowledge-base management and retrieval.
- `internal/database/`: SQLite access and migrations.
- `web/`: templates, static JavaScript/CSS, and i18n.
- `tools/`, `roles/`, `skills/`, `agents/`: runtime resources.
- `config.yaml`: configuration example.

## Development Startup

Use a local configuration and keep runtime data separate from source. Common startup options are:

```bash
go run ./cmd/server --config config.yaml
```

Before testing, confirm model credentials, writable `data/`, resource directories, and any TLS settings. Do not use production targets or secrets for local development.

## Routing

Routes are registered centrally in `internal/app/app.go` and implemented by handlers under `internal/handler/`. Keep authentication, permission checks, request validation, and response errors consistent with existing handlers. Update OpenAPI for every public endpoint.

## Database

SQLite is the default. When adding tables or fields:

- put migration logic in database initialization or the module's migration function;
- preserve compatibility with existing `data/conversations.db` files;
- add tests for migrations and important queries;
- consider locks, nullable fields, defaults, and interrupted startup.

## Adding Tools

Prefer `tools/*.yaml` for command tools. Use Go built-in tools when a tool needs internal state or structured integration. Built-in and YAML tools should define clear input schemas, handle timeouts and errors, and respect HITL for risky actions.

## Adding Roles

Roles are managed through `roles/*.yaml`. Common fields include name, description, system prompt, and tool list. Follow the least-tool principle; do not give every tool to a specialized role by default. Put authorization boundaries in roles and HITL, not only in a Skill.

## Adding Sub-Agents

Put multi-agent sub-agents in `agents/*.md`. Example front matter:

```yaml
---
name: Attack Surface Enumeration
id: attack-surface-enumeration
description: Enumerate exposed attack surface and organize verifiable leads
tools:
  - subfinder
  - nmap
  - http-framework-test
bind_role: Information Gathering
max_iterations: 300
---
```

The body is the system prompt. The main Agent may use a fixed orchestrator filename or `kind: orchestrator`.

## Adding Skills

Put Skills in `skills/<name>/SKILL.md`. They provide specialized capabilities, process guidance, or supporting references. See [Skills Guide](skills-guide.md). Keep secrets and one-off target facts out of Skills.

## Frontend Changes

Use existing helpers such as `apiFetch`, modal utilities, notifications, and i18n. Update the frontend language resources for every new visible string. Avoid putting secrets or provider keys in frontend code. High-risk buttons need confirmation and clear state feedback. See [Frontend i18n](frontend-i18n.md) for conventions.

## OpenAPI

`internal/handler/openapi.go` maintains the built-in OpenAPI output. For every new public endpoint, add its path, method, summary/description, request body, responses, and security requirements so `/api-docs` stays current.

## Development Habits

- preserve existing module boundaries;
- consider timeout and error paths for model, external API, filesystem, and Shell changes;
- connect high-risk capabilities to HITL or at least clear auditing;
- run relevant package tests after code changes;
- update documentation and configuration examples with behavior changes.

## Complete Business-Module Recipe

Do not add only a handler. A complete module usually needs:

1. A data model and SQLite migration.
2. A handler with parameters, errors, pagination, and filtering.
3. Audit records for management actions.
4. Monitor state for long-running execution.
5. An MCP decision: whether Agents should call it.
6. A HITL approval boundary for MCP tools.
7. An OpenAPI update at `/api/openapi/spec`.
8. Frontend i18n, loading, empty, and error states.
9. Database, handler, and edge-case tests.
10. Configuration, usage, troubleshooting, and safety documentation.

Missing one of these often becomes a later usability or safety bug.

## Handler Error Design

Prefer stable JSON responses:

```json
{
  "error": "machine_readable_code",
  "message": "human-readable explanation"
}
```

Frontend code needs stable fields, users need actionable messages, and logs need detailed internal errors without exposing secrets.

## Long-Running Task Design

For scanning, indexing, batch tasks, C2, or external operations, answer:

- Can it be cancelled?
- Can progress be queried?
- Can it be retried?
- Where is the result stored?
- Does state survive page refresh?
- Does it block the HTTP request?

If not, use task tables, event streams, or monitoring. Define cleanup and partial-result behavior as well.

## Test Priority

High-value tests include:

- configuration hot-apply;
- HITL branches;
- Shell timeout and no-output handling;
- external MCP recovery;
- knowledge-base indexing and post-processing;
- WebShell OS and encoding detection;
- SQLite migration compatibility;
- SSE error, done, cancellation, and resume behavior.
