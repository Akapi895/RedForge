# Audit and Monitoring


CyberStrikeAI has separate observability streams:

- Audit: who performed platform management actions.
- Monitor: how tool executions ran.
- HITL logs: why a tool call was approved, edited, or rejected.
- Process details: how an Agent chained reasoning, tools, and outputs.

Use them together during review.

## Audit

Config:

```yaml
audit:
  enabled: true
  retention_days: 15
  max_detail_bytes: 8192
  auth_failure_cooldown_seconds: 60
```

Endpoints:

- `GET /api/audit/meta`
- `GET /api/audit/summary`
- `GET /api/audit/logs`
- `GET /api/audit/logs/:id`
- `GET /api/audit/logs/export`

Watch for login failures, password changes, config updates, external MCP changes, WebShell/C2 actions, and HITL rejections.

Audit covers login, configuration, and resource-management operations. It does not store complete conversation bodies or the body of every tool call.

## Tool Monitoring

Config:

```yaml
monitor:
  retention_days: 90
```

Endpoints:

- `GET /api/monitor`
- `GET /api/monitor/execution/:id`
- `POST /api/monitor/execution/:id/cancel`
- `DELETE /api/monitor/execution/:id`
- `DELETE /api/monitor/executions`
- `GET /api/monitor/stats`
- `GET /api/monitor/calls-timeline`
- `POST /api/monitor/executions/names`

Monitoring is for execution state, duration, cancellation, and result review. It is not a substitute for platform audit.

## Notification Summaries

Endpoints:

- `GET /api/notifications/summary`
- `POST /api/notifications/read`

Notifications highlight pending actions, unread state, and a summary of running tasks. The exact presentation depends on the frontend page.

## HITL Logs

HITL decision logs are managed separately:

- `GET /api/hitl/pending`
- `GET /api/hitl/logs`
- `GET /api/hitl/logs/:id`
- `DELETE /api/hitl/logs`
- `POST /api/hitl/decision`
- `POST /api/hitl/dismiss`

Review HITL logs together with platform audit to understand why an Agent did or did not execute a tool.

## Retention Guidance

Security-tool logs can include targets, paths, commands, and sensitive outputs. Longer retention is not always safer.

- Short engagements: 15-30 days.
- Continuous red-team platform: 90-180 days.
- Compliance archive: export and encrypt.

Recommended retention practices:

- Keep audit logs for 15 to 90 days, adjusted to organizational requirements.
- Keep tool-execution records for 30 to 180 days.
- Clean C2, WebShell, uploaded attachments, and task results separately according to the project lifecycle.
- Redact sensitive data and restrict access when exporting logs.

## Review Checklist

Weekly:

- failed logins and unusual IPs;
- config changes;
- long-running or frequently failing tools;
- external MCP state;
- DB size and disk.

After engagement:

- export required evidence;
- delete stale WebShell/C2 resources;
- clean uploads and temporary workspaces;
- archive reports, vulnerabilities, and project facts.

## Audit and Monitoring Boundaries

These streams are often confused, but they answer different questions:

- Audit answers: “Who performed which management action on the platform?”
- Monitoring answers: “How did a tool call run?”
- HITL logs answer: “Why was a tool call approved, edited, or rejected?”
- Conversation process details answer: “How did the Agent reason and chain steps?”

A security review usually combines all four. Audit alone omits concrete tool output; monitoring alone omits who changed configuration.

## Key Event Interpretation

Pay particular attention to:

| Event | Why it matters |
| --- | --- |
| Login failure | Brute-force attempts, leaked passwords, or misconfiguration |
| Password change | All old sessions are revoked and active users may be affected |
| Configuration update | May change models, tools, C2, knowledge-base, or audit policy |
| External MCP change | A new tool may gain local or remote execution capability |
| C2 listener/task | Directly affects authorized targets and network exposure |
| WebShell connection change | May introduce an execution channel to a real business system |
| HITL rejection | Shows that an Agent or user request reached a risk boundary |

## Log Retention Is Not Always Better

Security-tool logs often contain targets, vulnerabilities, paths, command output, and internal organizational information. Retention should balance review value against disclosure risk:

- Short exercises: 15–30 days.
- Continuous red-team platform: 90–180 days.
- Compliance requirements: archive according to organizational policy and encrypt exports.

Without a dedicated logging platform, do not retain every SQLite detail indefinitely.

## Source Anchors

- Audit service: `internal/audit/service.go`
- Sanitization: `internal/audit/sanitize.go`
- Retention: `internal/audit/retention.go`
- Audit handler: `internal/handler/audit.go`
- Monitor: `internal/monitor/reconcile.go`
- Monitor handler: `internal/handler/monitor.go`
- HITL logs: `internal/handler/hitl_logs.go`
