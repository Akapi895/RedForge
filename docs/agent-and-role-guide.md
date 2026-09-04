# Agent and Role Guide

CyberStrikeAI Agent behavior is jointly determined by roles, sub-agents, and Skills. Roles define the task identity and available tools; sub-agents define multi-agent specialization; Skills provide on-demand domain knowledge and procedures.

## Roles

Role files are YAML files under `roles/` and normally contain a name, description, system prompt, and available tool list. Bind only necessary tools, state authorization boundaries clearly, require approval before high-risk operations, and keep output formats stable for reporting and review.

Design roles for reconnaissance, Web application scanning, API security testing, cloud security auditing, digital forensics, binary analysis, or CTF work. A specialized role should bind only the tools it needs, and its prompt should state what it may and may not do.

## Single-Agent

The single-Agent endpoints are `POST /api/eino-agent` and `POST /api/eino-agent/stream`. Use them for quick questions, single-target testing, short tool chains, and interactive analysis that needs stable context. The Agent uses the selected role, visible tools, Skills, project facts, and conversation history; final responses still pass the finalization gate.

## Multi-Agent Modes

The multi-agent endpoints are `POST /api/multi-agent` and `POST /api/multi-agent/stream`.

- `deep` decomposes work dynamically and delegates specialized tasks.
- `plan_execute` follows planning, execution, and replanning stages.
- `supervisor` routes work among specialist sub-agents and aggregates their results.

Use multi-agent orchestration for multi-stage penetration tests, broad reconnaissance, parallel role specialization, long-running work, and batch tasks. Choose the mode by task shape, risk, required evidence, and interruption needs.

## Markdown Sub-Agent

Sub-agents live in `agents/*.md`. Example front matter:

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

The body is the system prompt. It should define responsibility boundaries, expected inputs, tool order, output format, and prohibited actions. A sub-agent should return evidence and a confidence level so the main Agent can continue orchestration.

## Main Agent

Available main-agent files include `agents/orchestrator.md`, `agents/orchestrator-plan-execute.md`, and `agents/orchestrator-supervisor.md`. A Markdown file can also be marked with `kind: orchestrator` in front matter. Define only one main-agent implementation for each orchestration mode. Keep main-agent instructions separate from specialist sub-agent instructions and define handoff and finalization behavior.

## Tool Selection

The recommended selection order is the role's minimal tool set, task-specific tools supplied by the sub-agent, dynamic discovery through `tool_search`, and HITL approval for high-risk tools. Binding every tool to every role increases context cost and misuse risk.

With `tool_search`, the model initially sees only a subset of tools:

- visible in the UI does not mean visible in the current model context;
- `tool_search_always_visible_tools` are easier to call;
- clear tool descriptions improve search hits;
- sub-agent tool constraints still matter.

## Prompt Recommendations

Prompts should state that the Agent acts only within the authorized scope, confirm the target and constraints first, request approval for writes, deletion, brute force, persistence, C2, and WebShell operations, provide reviewable evidence, and label assumptions instead of inventing results. Output formats should be stable and suitable for reporting and review.

## Debugging

When behavior is unexpected, inspect the selected role, mode, Markdown front matter, available and visible tools, Skill loading, HITL state, pending executions, and finalization decision. A tool visible in the UI may still be absent from the model context.

If the Agent selects the wrong tool, narrow the role tool list, improve `short_description`, or tune `tool_search_always_visible_tools` and the prompt's tool order. If a multi-agent run drifts, inspect overly broad sub-agent descriptions, `sub_agent_user_context_max_runes`, the orchestrator prompt, process details, and tool-execution monitoring.

## Role, Sub-Agent, and Skill Boundaries

| Resource | Purpose | Not for |
| --- | --- | --- |
| Role | identity, tone, tool boundary, authorization rules | large reference material |
| Agent Markdown | multi-agent specialization, handoff format, local strategy | one-off facts |
| Skill | reusable procedures, checklists, templates, references | permission control |

Roles own identity, authorization, and tool boundaries. Sub-agents own specialized delegation and handoff behavior. Skills own reusable procedures and references. Do not use a Skill to grant permissions or use a role as a replacement for a knowledge base; authorization boundaries belong in roles and HITL first.

## Orchestration Modes

| Mode | Good for | Poor fit |
| --- | --- | --- |
| `eino_single` | short tasks, interactive analysis | large multi-stage work |
| `deep` | dynamic task decomposition | strict sequential workflows |
| `plan_execute` | plan, execute, replan loops | frequent user interruption |
| `supervisor` | expert routing | vague or too many sub-agents |

Start with `eino_single`; use `plan_execute` for structured projects; use `deep` or `supervisor` when specialist Agents matter.

## Tool Visibility and Behavior

When a tool is not used, check role tools, sub-agent tools, `tool_search` configuration, and the tool description. The tool list shown in the UI is not proof that the model currently has the tool in context. Check role tools, sub-agent tools, visibility settings, and descriptions together.

## Sub-Agent Output Format

Sub-agents should return structured results:

```markdown
## Conclusion
## Evidence
- Tool:
- Key output:
- Confidence:
## Risks
## Suggested next step
```

This helps the orchestrator continue and supports reporting, attack-chain construction, and project-fact persistence.

## Source Anchors

- Markdown Agent parser: `internal/agents/markdown.go`
- Multi-agent preparation: `internal/handler/multi_agent_prepare.go`
- Orchestration: `internal/multiagent/eino_orchestration.go`
- Tool search middleware: `internal/multiagent/eino_middleware.go`
- Sub-agent context: `internal/multiagent/sub_agent_context_test.go`
