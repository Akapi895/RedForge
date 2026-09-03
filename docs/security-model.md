# Security Model


CyberStrikeAI is not a generic chatbot. It is a high-privilege security automation system with command execution, MCP tools, WebShell management, optional C2, batch tasks, and multi-agent orchestration.

## Trust Boundaries

Main actors:

- Web user: can chat, change settings, manage resources, and trigger tools.
- Agent: selects tools based on role, context, and middleware.
- MCP tools: may access files, run commands, call services, or touch targets.
- External MCP: third-party local or remote tool providers.
- Robot callbacks: platform-authenticated message ingress outside Web login.

Anyone who can log into the Web UI should be treated as an operator of the instance.

## Authentication and Sessions

Web credentials are managed by RBAC users, including the built-in `admin` account. Recommended practices:

- Change the initial `admin` password immediately after deployment; it is printed on the first console startup.
- Use a long random password and restrict who may receive it.
- Place the service behind an internal network, VPN, bastion, or reverse-proxy authentication.
- Enable HTTPS in production to avoid transmitting cookies in cleartext.
- Use a reverse proxy to restrict source IPs.

Session lifetime is controlled by `auth.session_duration_hours`.

## Tool-Execution Risk

Tool sources include:

- Built-in security execution tools.
- YAML command tools under `tools/`.
- Eino Skills filesystem tools such as `read_file`, `write_file`, `edit_file`, and `execute`.
- Tools exposed by external MCP providers.
- C2 and WebShell MCP tools.

Recommended controls:

- Enable only the tools needed for the current task.
- Keep high-risk command tools out of the global allowlist.
- Bind the smallest practical tool set to each role.
- Require human approval for destructive, persistent, lateral-movement, and credential operations.
- Prefer trusted local processes for external MCP; remote MCP requires authentication and network isolation.

## HITL

HITL is the approval layer before a tool call. Common modes are:

- `human`: human approval.
- `audit_agent`: automatic approval by the audit Agent.
- `review_edit`: the audit Agent may edit parameters before approval.

Recommended policy:

- Use human approval by default in new environments.
- Allowlist only read-only, low-risk, stable tools.
- Scope scanning tools by target and role prompt.
- Review writes, deletion, payload execution, C2, WebShell, and account changes carefully.

See [HITL Best Practices](hitl-best-practices.md).

## Audit

Platform audit is controlled by `audit`. It records login, configuration, and resource-management operations, but not complete conversation bodies and is not a forensic log.

Tool-execution records are maintained by the monitoring module, with retention controlled by `monitor.retention_days`.

Recommended practices:

- Enable `audit.enabled`.
- Export audit logs regularly.
- Set an appropriate retention period.
- Review failed logins, configuration changes, and C2/WebShell operations carefully.

## C2 Risk

Built-in C2 starts listeners, generates payloads, receives sessions, and executes tasks. Enable it only in explicitly authorized labs, internal exercises, or red-team environments.

- Set `c2.enabled: false` when it is not needed.
- Do not expose C2 listener ports publicly without explicit authorization and isolation.
- Control access to payload files, callback addresses, and task output.
- Prefer HITL for C2 tasks.

## WebShell Risk

WebShell management allows commands and file operations on registered connections:

- Add only authorized targets.
- Give connections clear names, labels, and notes.
- Do not store real production WebShell connections in shared environments.
- Confirm the target and command before the Agent uses WebShell.
- Remove expired or no-longer-authorized connections.

## Data and Privacy

Local data includes:

- Conversations and messages: `data/conversations.db`.
- Knowledge-base indexes: `data/knowledge.db`.
- WebShell, C2, vulnerability, project, and task data.
- Uploaded attachments: `chat_uploads/`.

Recommended practices:

- Restrict file permissions.
- Encrypt backups.
- Do not upload sensitive data without authorization.
- Remove obsolete conversations, attachments, C2 output, and audit logs.

## Threat Model

| Threat | Path | Impact | Controls |
| --- | --- | --- | --- |
| Password leak | login, then use terminal/WebShell/C2 | platform takeover | strong password, HTTPS, internal network, audit |
| Prompt injection | target content instructs Agent to misuse tools | unauthorized actions | role boundaries, HITL, least tools |
| Malicious MCP | external tool lies or has side effects | host/target impact | trusted MCP only, isolation |
| Tool YAML tampering | command template changed | malicious execution | file permissions, review |
| C2 misuse | payload or task against unauthorized target | legal and business risk | disabled by default, approvals |
| WebShell misuse | destructive command on business host | outage/data loss | naming, read-only first, HITL |
| DB leak | copy `data/*.db` or uploads | sensitive target data | permissions, encrypted backups |

## HITL Is Not Magic

HITL sees a tool name, arguments, and context. It does not always see real-world impact. Be conservative when:

- a harmless-looking command wraps `bash -c` or base64;
- the MCP tool description is untrusted;
- WebShell target identity is vague;
- C2 payload delivery happens outside the platform;
- a read-only tool can still create traffic or side effects.

Audit Agent is useful for routine checks, not for replacing humans on destructive operations.

## Data Minimization

Avoid long-term storage of:

- real customer credentials;
- raw production data;
- long-lived cookies;
- unrelated scan output;
- stale WebShell or C2 sessions.

Project closeout should include cleanup of uploads, WebShell connections, C2 payloads, temporary workspaces, and bulky execution logs.

## Production Baseline

- Strong password and HTTPS.
- Internal/VPN/proxy restricted access.
- `audit.enabled: true`.
- Random `mcp.auth_header_value` when HTTP MCP is exposed.
- `c2.enabled: false` unless required.
- Minimal external MCP.
- No high-risk tools in global allowlist.

Back up and update the system regularly.

## Authorization Boundaries

Every high-risk role should state its authorization boundary in its prompt rather than relying only on a user's verbal statement. For example:

```text
Act only within the target scope explicitly provided by the user. Before any write, delete, brute-force, persistence, credential-access, C2, WebShell, or lateral-movement operation, explain the purpose, impact, target, and rollback method, then wait for HITL approval.
```

This does not replace technical controls, but it reduces the chance that an Agent expands its scope during an ambiguous task.

## Source Anchors

- Sessions: `internal/security/auth_manager.go`
- Auth middleware: `internal/security/auth_middleware.go`
- Rate limiting: `internal/security/ratelimit.go`
- Shell execution: `internal/security/executor.go`
- HITL execution: `internal/handler/hitl_execution.go`
- Audit service: `internal/audit/service.go`
