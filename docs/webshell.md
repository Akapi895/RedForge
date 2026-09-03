# WebShell Management


WebShell management stores authorized WebShell connections and allows command/file operations through the UI and Agent tools.

## Workflow

1. Add a connection.
2. Fill URL, parameter/password, and metadata.
3. Test connectivity.
4. Run read-only identification commands first.
5. Let AI assist only after selecting the correct connection.

Connections are stored in SQLite.

## API

- `GET /api/webshell/connections`
- `POST /api/webshell/connections`
- `PUT /api/webshell/connections/:id`
- `DELETE /api/webshell/connections/:id`
- `GET /api/webshell/connections/:id/state`
- `PUT /api/webshell/connections/:id/state`
- `POST /api/webshell/exec`
- `POST /api/webshell/file`
- `GET /api/webshell/connections/:id/ai-history`
- `GET /api/webshell/connections/:id/ai-conversations`

## MCP Tools

Typical tools include:

- `webshell_exec`: execute a command on a connection.
- `webshell_file_list`: list a directory.
- `webshell_file_read`: read a file.
- `webshell_file_write`: write a file.
- WebShell connection-management tools.

These tools require `connection_id`; the frontend normally injects the selected connection into context.

## Command Execution

Before running a command, confirm that the connection belongs to an authorized target, that the command will not disrupt the business, that output may contain sensitive information, and that long or interactive commands are unsuitable for the WebShell channel.

Start with read-only environment checks:

```bash
whoami
pwd
uname -a
id
```

On Windows targets:

```cmd
whoami
cd
ver
ipconfig
```

## File Operations

File operations include directory listing, reading, and writing:

- Back up a file before writing it.
- Do not write unconfirmed scripts or binaries to production targets.
- Use dedicated download/upload channels for large files.
- Account for target encoding and line endings.

## AI Assistance

AI can identify the operating system and current privilege, plan read-only enumeration, analyze command output, and summarize risks and remediation.

Do not let AI automatically delete files, modify business configuration, establish persistence, collect credentials, or scan a broad internal network. Require human confirmation and HITL for these actions.

## Security Recommendations

- Store connections only for authorized targets.
- Include project, environment, and target in connection names.
- Delete connections after an exercise.
- Keep WebShell write tools out of the global approval-free allowlist.
- Record important output as project facts or reports, then clean up sensitive raw data.

## Troubleshooting

Connection failures may be caused by an unreachable URL, incorrect parameter/password, target WAF blocking, or proxy/TLS configuration. For garbled command output, check target encoding and try a different output encoding or base64 wrapping. If the Agent cannot find a connection, confirm that the frontend selected it, that it still exists, and that `connection_id` is correct.

## Operation Tiers

| Tier | Operation | Risk | Guidance |
| --- | --- | --- | --- |
| Identify | `whoami`, `pwd`, OS version | low | may automate |
| Enumerate | dirs, processes, env vars | medium | constrain path/command |
| Read | config, logs, source | medium-high | human confirms sensitivity |
| Write/execute | write, run script, delete | high | human approval and rollback |

Having a WebShell does not make follow-up operations low risk.

## Naming

Use:

```text
<project>-<environment>-<target>-<privilege>-<date>
```

Example:

```text
acme-staging-web01-www-20260707
```

Avoid vague names like `test`, `shell1`, or `customer machine`.

```text
test
shell1
customer machine
```

## AI Guardrail Prompt

```text
Before using WebShell, confirm connection_id, target name, current directory, and privilege. Default to read-only commands. Any write, delete, upload, permission change, persistence, credential access, or internal probing requires purpose, impact, rollback plan, and approval.
```

## Source Anchors

- Handler: `internal/handler/webshell.go`
- Context: `internal/handler/webshell_context.go`
- Probe: `internal/handler/webshell_probe.go`
- Encoding/OS tests: `internal/handler/webshell_encoding_test.go`, `internal/handler/webshell_os_test.go`
- Tool registration: `internal/app/app.go`
