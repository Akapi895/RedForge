# Eino Multi-Agent Notes


CyberStrikeAI uses CloudWeGo Eino ADK for the current single-agent and multi-agent execution paths. The native legacy ReAct path has been removed.

## Overall Conclusion

The refactor is ready for production trials: streaming chat, MCP tool bridging, configuration switches, and frontend mode selection are implemented. Single-agent requests use `/api/eino-agent/stream`; multi-agent requests use `/api/multi-agent/stream` with `orchestration` set to `deep`, `plan_execute`, or `supervisor`.

- **Deep** is for complex security testing and task-sub-agent collaboration.
- **Plan-Execute** is for explicit goals requiring a plan → execute → replan loop.
- **Supervisor** is for dynamically routing work among specialist sub-agents.
- Robots default to `robot_default_agent_mode: eino_single`; batch queues default to `eino_single` and require `multi_agent.enabled` for multi-agent modes.

**Entrypoints**

- Single-agent: `/api/eino-agent` and `/api/eino-agent/stream`
- Multi-agent: `/api/multi-agent` and `/api/multi-agent/stream`

Multi-agent orchestration is selected by request body:

- `deep`
- `plan_execute`
- `supervisor`

Robots default to `robot_default_agent_mode`, and batch tasks can opt into multi-agent through config.

**Agent Definitions**

Markdown agents live under `agents/`.

Typical files include `agents/orchestrator.md`, `agents/orchestrator-plan-execute.md`, `agents/orchestrator-supervisor.md`, and other `agents/*.md` files.

Front matter controls name, id, description, tools, bound role, max iterations, and optional orchestrator kind.

**Middleware**

Important Eino middleware:

- tool search: exposes a small visible tool set and unlocks others on demand;
- patch tool calls: repairs interrupted histories;
- plan task: structured task board;
- reduction: truncates or persists large tool outputs;
- summarization: compresses long contexts;
- checkpoint: resume after crash/OOM.

These settings live under `multi_agent.eino_middleware`.

**Skills**

Eino Skills support progressive disclosure. The Agent initially sees names and descriptions; details are loaded only when needed through the configured skill tool.

**Operational Notes**

- Tool visibility is not the same as tool availability in the UI.
- Running streams keep their startup context even if config changes mid-run.
- Summarization can write transcripts under `data/conversation_artifacts/...`.
- High-risk tools should still be constrained by roles and HITL.

## Completed Items

| Area | Completed implementation |
| --- | --- |
| Dependencies and agents | `go.mod` directly depends on Eino and its OpenAI extensions; Go proxy and bootstrap guidance are documented. |
| Configuration | `agent.max_iterations` is the shared ReAct limit; `multi_agent` controls enablement, robot/batch defaults, sub-agents, Skills, and middleware. |
| Markdown agents | `agents_dir` contains `*.md`; mode-specific orchestrators are selected by fixed filenames or `kind: orchestrator`, with body instructions taking precedence over YAML. |
| MCP bridge | `internal/einomcp` bridges MCP tool definitions and conversation execution IDs. |
| Orchestration | Single-agent, Deep, Supervisor, and Plan-Execute paths use typed Eino Agentic messages through adapters preserving Runner, TurnLoop, SSE, and MCP boundaries. |
| HTTP | `/api/multi-agent` and `/api/multi-agent/stream` are always registered; `multi_agent.enabled` controls runtime availability and emits SSE `error` plus `done` when disabled. |
| Session preparation | `prepareMultiAgentSession` preserves WebShell context and tool allowlists consistently with single-agent behavior. |
| Frontend and automation | Frontend, WebShell, robots, and batch queues consistently select Eino single-agent or the requested orchestration mode. |
| Streaming | The shared boundary supports conversation, progress, response start/delta, thinking, tool events, response, and done events. |
| Configuration API | APIs expose middleware, user-input budgets, model retry/failover, and persistent tool allowlists without overwriting sub-agent definitions. |
| Middleware | Typed patch-tool-calls, tool search, plan task, reduction, checkpoint, retry/failover, filesystem, Skill, and summarization are available. Plan-Execute keeps the official session contract while Planner/Replanner remain classic until typed constructors exist. |
| Agentic adapters | Message and event adapters preserve text, reasoning, tool calls, tool results, checkpoint resume, stream rendering, persistence, and MCP display updates. |
| TurnLoop and Runner | Bridges provide run IDs, interrupt/continue, cancellation, checkpoint resume, recovery, transient retry, context-overflow compression, pending-tool flushing, partial/final results, and usage summaries. |

The implementation details behind these items are intentionally recorded here because they define compatibility boundaries. `internal/multiagent/eino_agentic_message.go` maps text, reasoning, function calls, and function results between `schema.Message` and `schema.AgenticMessage`. `eino_agentic_event_adapter.go` maps typed assistant, tool-result, stream, and error events back to the existing SSE/MCP drain, preserving tool names and call IDs. `eino_agentic_agent_adapter.go` wraps typed agents as `adk.Agent` and `adk.ResumableAgent`, so Runner, TurnLoop, and checkpoint resume do not need a new outer contract. `eino_agentic_chat_model_agent.go` assembles the typed chat-model Agent.

The production paths in `eino_single_runner.go`, Deep, and Supervisor use AgenticModel. `eino_model_resilience.go` builds the Agentic OpenAI-compatible model and connects native ModelRetry and ModelFailover. `eino_middleware.go` provides the typed patch-tool-calls, tool-search, plan-task, and reduction middleware; `eino_agentic_summarize.go` connects the domain summary policy to typed summarization; and `eino_agentic_chat_model_tail_middleware.go` provides a protocol-neutral tail for system merging, continuation de-duplication, typed summarization, model-facing traces, and output guards.

The TurnLoop bridge is assembled through `WithAgentTurnLoopInterruptRegistrar`. `eino_run_trace.go` creates a shared `runId` for progress events and observation callbacks. The runtime and starter helpers manage session items, interrupt/continue, preemption timeouts, idle stop, checkpoint parameters, cancellation registrars, fresh-run IDs, and resume options. Event-bridge and runtime-session helpers forward TurnLoop events to the existing SSE/MCP drain, suppress framework-level cancellation caused by preemption, and coordinate startup, resume, restart, recovery, stream errors, completion, cancellation cleanup, and final-result construction.

The run handlers cover non-retryable errors, timeout and iteration-limit progress, native retry, context-overflow restart, transient run retry, assistant-stream receive errors, cancellation, orphan pending-tool flush, and checkpoint cleanup. Shared accumulators and emitters track pending calls, stream IDs, assistant output, run messages, tool results, main and sub-agent replies, reasoning, duplicate tool calls, filesystem monitoring, and partial/final results. The legacy model-output recovery marker remains only as a compatibility guard; the default path does not rewrite model tool calls before execution. Context-overflow and transient-retry restarts are TurnLoop-aware, so abnormal continuation does not fall back to a plain Runner.

## In Progress / Backlog

| Priority | Item | Description |
| --- | --- | --- |
| P1 | Typed Plan-Execute Planner/Replanner | The Executor uses a typed Agentic agent, but official Planner/Replanner constructors still accept classic ChatModel. Migrate when typed constructors or a custom typed root are available. |
| P2 | Observability and cost | Run-level IDs and usage summaries are implemented; add model pricing and cost fields. |
| P3 | Tests | Add integration tests for `internal/multiagent` and `einomcp` using mock models or recorded event playback. |

## Key File Index

- `internal/multiagent/runner.go`: Deep, Plan-Execute, Supervisor assembly and event loop.
- `internal/multiagent/eino_orchestration.go`: Plan-Execute root and Agentic Executor middleware.
- `internal/handler/multi_agent.go`: SSE and synchronous HTTP endpoints.
- `internal/handler/multi_agent_prepare.go`: session preparation, including WebShell.
- `internal/einomcp/`: MCP-to-Eino tool adapter.
- `config.yaml`: `multi_agent` example configuration.
- `web/static/js/chat.js`: mode selection and stream URL.
- `web/static/js/webshell.js`: WebShell stream URL and main-chat mode alignment.
- `web/static/js/settings.js`: multi-agent and Eino retry/failover settings.

## Version History

| Date | Description |
| --- | --- |
| 2026-03-22 | Initial Eino DeepAgent, stream, frontend switch, and GOPROXY guidance. |
| 2026-03-22 | Added progress documentation, session preparation, WebShell alignment, multi-agent HTTP, and OpenAPI entries. |
| 2026-03-22 | Added persistent routes, disabled-stream SSE errors, robot/batch multi-agent execution, `bind_role`, and Skill/tool inheritance. |
| 2026-03-22 | Added `tool_result.toolCallId`, mapped `ReasoningContent` to thinking streams, and connected batch multi-agent settings to Eino execution. |
| 2026-03-22 | Added stable-signature de-duplication for streaming tool events, final-response de-duplication, and `task` display for built-in scheduling. |
| 2026-03-22 | Added `agents/*.md`, `agents_dir`, Deep integration, the frontend Agents menu, and CRUD APIs. |
| 2026-03-22 | Added `orchestrator.md`, `kind: orchestrator`, main/sub-agent markers, and instruction precedence. |
| 2026-04-19 | Added conversation-mode Deep, Plan-Execute, and Supervisor flows with `orchestration`; robot and batch defaults use `deep`. |
| 2026-04-21 | Removed role `skills`; Skills now load on demand through the Eino `skill` tool. |
| 2026-06-02 | Removed native ReAct endpoints and unified Eino ADK single-agent and multi-agent paths. |
| 2026-07-06 | Aligned Deep, Plan-Execute, and Supervisor guidance with their intended scenarios and tightened transfer/exit constraints. |
| 2026-08-14 | Operationalized native model retry/failover, Agentic typed summarization, streaming tool-result replay, context-overflow recovery, and Agentic production-path adapters. |
