# Configuration Reference


The main configuration file is `config.yaml`. Many fields are editable through the Web settings page, but not every field has the same hot-apply behavior.

## Core Configuration

```yaml
server:
  host: 0.0.0.0
  port: 8080
  tls_enabled: true
  # Optional: other trusted Web integrations; Chromium extensions need no entry.
  # cors_allowed_origins:
  #   - https://trusted-integration.example
auth:
  session_duration_hours: 12
ai:
  default_channel: openai-main
  channels:
    openai-main:
      name: OpenAI Main
      provider: openai_compatible
      base_url: https://api.openai.com/v1
      api_key: sk-...
      model: gpt-4.1
agent:
  max_iterations: 12000
  tool_timeout_minutes: 60
```

- `version`: version displayed by the frontend. Use the current release value from `config.example.yaml`.
- `server.host/port`: address and port listened to by the Web service.
- `server.tls_*`: HTTPS settings. In production, use `tls_cert_path` and `tls_key_path`.
- `auth.session_duration_hours`: login-session lifetime in hours. Login passwords are managed by RBAC users; the initial `admin` password is printed to the console on first startup.
- `log.output`: `stdout`, `stderr`, or a file path.

Change the initial `admin` password from the Web UI after first login. Use HTTPS or a trusted reverse proxy in any shared environment.

Valid Chromium `chrome-extension://<32-character-extension-id>` origins are recognized automatically. The extension must still obtain host permission and authenticate with a password and Bearer token. `server.cors_allowed_origins` remains available as an exact allowlist for other trusted Web integrations; wildcards are not accepted, and changing it requires a restart.

## AI Channels

`ai` is the recommended model configuration entry. In the Web UI, use **System Settings → Basic Settings → AI Channel Configuration**. Saving that form writes `ai.default_channel` and `ai.channels`. The legacy `openai` field remains as a backward-compatible runtime field; on load, CyberStrikeAI ensures a default channel exists and synchronizes the resolved `ai.default_channel` into runtime `openai`.

```yaml
ai:
  default_channel: openai-main
  channels:
    openai-main:
      name: OpenAI Main
      provider: openai_compatible
      base_url: https://api.openai.com/v1
      api_key: sk-...
      model: gpt-4.1
      max_total_tokens: 120000
      max_completion_tokens: 16384
      reasoning:
        mode: on
        effort: high
        allow_client_reasoning: true
        profile: openai_compat
    claude-main:
      name: Claude Main
      provider: claude
      base_url: https://api.anthropic.com/v1
      api_key: sk-ant-...
      model: claude-sonnet-4-5
```

| Field | Meaning |
| --- | --- |
| `ai.default_channel` | Default channel ID for new conversations and requests without an explicit channel. |
| `ai.channels.<id>` | Channel config. IDs are normalized to lowercase letters, digits, and hyphens. |
| `name` | Display name in the Web UI; falls back to the ID. |
| `provider` | `openai_compatible` or `claude`. OpenAI-compatible channels map to runtime `openai`; Claude channels use Eino's native Anthropic Messages API component. |
| `base_url/api_key/model` | Required. Base URL usually includes a version path such as `/v1`. |
| `max_total_tokens` | Shared context budget for compression, attack-chain generation, multi-agent summaries, and similar paths. |
| `max_completion_tokens` | Per-response output cap; default is used when empty. |
| `reasoning` | Default reasoning fields for the channel. Gateway support varies; try `mode: off` first when a provider rejects requests. |

The chat page reads saved channels into the “AI Channel” selector. A non-empty request `aiChannelId` selects a channel for that run/session without sending API credentials through the prompt path. Empty `aiChannelId` follows `ai.default_channel`.

Common Web UI operations:

- Add: click `+`, fill required fields, then save.
- Set default: select a channel, click **Set as default**, then save/apply.
- Copy: duplicate the current form, useful for the same provider with a different model.
- Delete: keep at least one channel; the default channel is protected from bulk delete.
- Probe: use **Test connection** or **Bulk probe** to validate API key, Base URL, and model.

## Agent

```yaml
agent:
  max_iterations: 12000
  tool_timeout_minutes: 60
  shell_no_output_timeout_seconds: 1200
  workspace_root_dir: ""
  system_prompt_path: ""
```

- `max_iterations`: default iteration limit for the single-agent executor, multi-agent main executor, and sub-agents.
- `tool_timeout_minutes`: maximum runtime for one tool call.
- `shell_no_output_timeout_seconds`: terminates a Shell process that produces no output for too long.
- `workspace_root_dir`: session workspace root; avoid setting it to the system `/tmp`.
- `system_prompt_path`: override file for the single-agent system prompt.

## HITL

```yaml
hitl:
  default_reviewer: audit_agent
  retention_days: 90
  tool_whitelist: [read_file, list_dir, glob, grep, tool_search]
  audit_model:
    provider: ""
    base_url: ""
    api_key: ""
    model: ""
```

- `default_reviewer`: `human` or `audit_agent`.
- `tool_whitelist`: global approval-free tool list; it is merged with the session allowlist.
- `audit_model`: separate model for the audit agent; empty values reuse the main model.
- `audit_agent_prompt` / `audit_agent_prompt_review_edit`: override the default approval policy.

See [HITL Best Practices](hitl-best-practices.md) for more policy details.

## Multi-Agent

```yaml
multi_agent:
  enabled: true
  robot_default_agent_mode: eino_single
  batch_use_multi_agent: false
  eino_skills:
    disable: false
    filesystem_tools: true
    skill_tool_name: skill
```

Supported modes:

- `eino_single`: Eino single-agent mode.
- `deep`: DeepAgent-style multi-agent mode.
- `plan_execute`: planning, execution, and replanning.
- `supervisor`: a supervisor agent delegates to sub-agents.

`agents_dir` points to the Markdown sub-agent directory. An individual agent can set `tools`, `bind_role`, and `max_iterations` in front matter.

## Tools and MCP

```yaml
security:
  tools_dir: tools
  tool_description_mode: full
mcp:
  enabled: false
  host: 0.0.0.0
  port: 8081
  auth_header: "X-MCP-Token"
  auth_header_value: ""
external_mcp:
  servers: {}
```

- `security.tools_dir`: directory containing built-in tool YAML files.
- `tool_description_mode`: `short` saves tokens; `full` provides more detail.
- `mcp.enabled`: whether to start the standalone HTTP MCP service.
- `mcp.auth_header_value`: shared secret for external MCP callers; it must be set in production.
- `external_mcp.servers`: external MCP federation configuration.

See `tools/README.md` for tool YAML rules.

## Knowledge Base

```yaml
knowledge:
  enabled: false
  base_path: knowledge_base
  embedding:
    provider: openai
    model: text-embedding-v4
    base_url: ""
    api_key: ""
  retrieval:
    top_k: 5
    similarity_threshold: 0.4
  indexing:
    chunk_size: 512
    chunk_overlap: 50
    batch_size: 10
```

When enabled, the knowledge-base retrieval tools and management endpoints are registered. See [Knowledge Base](knowledge-base.md) for details.

## Database

```yaml
database:
  path: data/conversations.db
  knowledge_db_path: data/knowledge.db
```

SQLite is used by default. When `knowledge_db_path` is empty, the knowledge base may reuse the conversation database; a separate file makes knowledge-base migration easier.

## Audit and Monitoring

```yaml
audit:
  enabled: true
  retention_days: 15
  max_detail_bytes: 8192
monitor:
  retention_days: 90
```

- `audit` records platform operations, but not conversation bodies or the body of every tool call.
- `monitor` controls retention for tool-execution records.

## C2, WebShell, and Projects

```yaml
c2:
  enabled: true
project:
  enabled: true
  fact_index_max_runes: 65000
```

- `c2.enabled`: when disabled, C2 listeners do not start and C2 MCP tools are not registered.
- WebShell connection settings are stored in SQLite and have no separate main-config switch.
- `project` controls the injection budget for the cross-conversation fact blackboard.

## Robots

`robots` supports personal WeChat iLink, WeCom, DingTalk, Lark, Telegram, Slack, Discord, and QQ. See [Robot Guide](robot.md) for detailed configuration steps.

## Configuration Change Recommendations

- Validate models, MCP, the knowledge base, and high-risk tools in a test environment first.
- After changing `tools_dir`, `roles_dir`, `skills_dir`, or `agents_dir`, check that the corresponding resources are listed in the Web UI.
- In production, avoid enabling unneeded C2, WebShell, terminal, and external MCP capabilities.
- After changing sensitive configuration, check the audit page for unexpected logins or configuration-change records.

## Hot-Apply Boundaries

`POST /api/config/apply` coordinates model config, tool description mode, MCP tool registration, knowledge components, robot restarts, and C2 runtime reconciliation. It does not make every field instantly effective.

| Section | Usually hot-applies | Extra action |
| --- | --- | --- |
| `ai.default_channel` / `ai.channels` | new requests use the resolved default or selected channel | running streams keep their current state; reload config for the frontend channel list |
| `openai` | compatibility field, usually synchronized from the default AI channel | prefer maintaining new config in `ai.channels` |
| `agent.max_iterations` | new tasks | existing tasks continue |
| `security.tool_description_mode` | when tools are re-exposed | existing model context is not rolled back |
| `hitl.tool_whitelist` | new approval checks | pending approvals are not re-decided |
| `knowledge.enabled` | initializes/updates components | scan and index are still required |
| `knowledge.embedding` | updates retriever/indexer config | rebuild index for existing vectors |
| `robots` | restarts long-lived connections | platform callback settings must still match |
| `c2.enabled` | reconciles C2 runtime | verify existing listeners/sessions manually |
| `server.port/tls` | usually needs process restart | listener settings are not ordinary hot state |

## Fallback Relationships

- `vision.api_key/base_url/provider` can inherit from the resolved default AI channel.
- `hitl.audit_model` can inherit from the resolved default AI channel.
- `knowledge.embedding.base_url/api_key` can inherit from model settings.
- `knowledge.retrieval.rerank.base_url/api_key` can inherit from embedding/OpenAI settings.
- `database.knowledge_db_path` can reuse the main conversation database when empty, but a separate file is easier to back up.

When debugging, inspect both the child config and the fallback parent.

## Recommended Values

| Field | Conservative | Aggressive | Decide by |
| --- | --- | --- | --- |
| `agent.tool_timeout_minutes` | 10-30 | 60+ | long scanners |
| `shell_no_output_timeout_seconds` | 300-600 | 1200+ | quiet tools |
| `knowledge.indexing.batch_size` | 5-10 | 20+ | embedding API limits |
| `knowledge.indexing.rate_limit_delay_ms` | 300-800 | 0-100 | 429 frequency |
| `retrieval.top_k` | 3-5 | 8-12 | context budget |
| `similarity_threshold` | 0.35-0.45 | 0.5+ | recall vs precision |
| `audit.retention_days` | 15-30 | 90+ | compliance and disk |
| `monitor.retention_days` | 30-90 | 180+ | need for long-term review |

## Change Template

Before changing config, write down:

```text
Purpose:
Sections:
Expected impact:
Rollback:
Validation endpoints:
```

After changing, validate the specific subsystem rather than trusting the save message.

```bash
curl -k https://127.0.0.1:8080/api/auth/validate \
  -H "Authorization: Bearer <token>"
```

Then validate the model, tools, knowledge base, C2, or robot according to the changed configuration type.

## Source Anchors

- Config structs: `internal/config/config.go`
- Env expansion: `internal/config/envexpand.go`
- Config API and apply: `internal/handler/config.go`
- Route registration: `internal/app/app.go`
- C2 reconciliation: `internal/app/c2_lifecycle.go`
