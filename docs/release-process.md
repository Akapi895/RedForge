# Release Process

This guide is for CyberStrikeAI maintainers and operators who release, upgrade, or roll back the system.

## Pre-Release Checklist

Before release, check:

- `README.md` and the intended localized README describe current features.
- `docs/` documents new capabilities.
- The sample `config.yaml` includes new fields.
- OpenAPI includes new endpoints.
- English and frontend i18n resources are synchronized.
- High-risk features have security documentation.

## Testing

At minimum, run:

```bash
go test ./internal/...
```

If entrypoints, builds, or commands changed, also run:

```bash
go test ./cmd/...
go build -o cyberstrike-ai ./cmd/server
```

If frontend code changed, manually verify login, streaming conversation output, settings save/apply, the tool list, and the absence of browser-console errors on relevant pages.

## Build

Build the release binary with:

```bash
go build -o cyberstrike-ai ./cmd/server
```

A release package should contain:

- `cyberstrike-ai`;
- `web/templates/`;
- `web/static/`;
- a sample `config.yaml`;
- `tools/`, `roles/`, `skills/`, `agents/`, and `docs/`;
- `README.md` and any intended localized README files;
- `LICENSE`.

Do not include local `data/`, real configuration secrets, uploaded attachments, or logs in a public release package.

## Upgrade Checklist

Before upgrading:

- stop the service;
- back up `config.yaml`;
- back up `data/`;
- back up custom `tools/roles/skills/agents/knowledge_base`;
- record the current version and startup method.

After upgrading:

- start the service;
- log in;
- test the model;
- check the tool list;
- check knowledge-base status;
- check external MCP;
- confirm C2/WebShell are enabled or disabled as expected;
- review logs and audit records.

## Rollback

Common rollback triggers are failure to start, failed database migration, unavailable core conversation functionality, or abnormal high-risk behavior.

Rollback steps:

1. Stop the new version.
2. Restore the previous binary or code.
3. Restore the pre-upgrade `config.yaml`.
4. Restore the pre-upgrade `data/`.
5. Start the old version and validate it.

If the new version changed the database schema, restore the database backup; replacing only the binary is insufficient.

## Changelog Recommendations

Record the following for every release:

- new features;
- behavior changes;
- configuration changes;
- database changes;
- security fixes;
- compatibility notes;
- upgrade notes.

Mark changes to high-risk modules separately, including C2, WebShell, Terminal, external MCP, and HITL.

## Release Risk Tiers

| Change | Risk | Must test |
| --- | --- | --- |
| Documentation and images | low | links and rendering |
| Frontend pages | medium | login, page states, API errors |
| Handler/API | medium | OpenAPI, permissions, error codes |
| Configuration structure | high | old-config compatibility, ApplyConfig |
| Database structure | high | old-database migration, rollback |
| Agent/MCP/HITL | high | tool calls, approvals, streaming interruption |
| C2/WebShell/Terminal | critical | authorized environment, auditing, disable switch |

Release notes should call out risk level rather than listing only feature names.

## Configuration Compatibility

New configuration fields should:

- have safe defaults when omitted;
- allow old configurations to start;
- be documented in the sample `config.yaml`;
- not be accidentally deleted by the Web Settings page;
- be tested through both hot-apply and restart paths.

If a new field enables a high-risk feature by default, reconsider the default.

## Database Changes

SQLite migrations must account for:

- users upgrading directly from a very old version;
- idempotence after an interrupted migration;
- whether new fields allow null values;
- indexes locking tables for too long;
- whether data backfill is required.

Release notes must clearly state that `data/` should be backed up before upgrading.

## Release Acceptance Script Outline

Minimum automation:

```bash
go test ./internal/...
go test ./cmd/...
go build -o cyberstrike-ai ./cmd/server
```

Manual smoke test:

```text
login -> model test -> new conversation -> tool list -> HITL -> knowledge base -> external MCP -> C2 off/on
```

For high-risk modules, perform an additional authorized-lab test instead of relying on unit tests alone.
