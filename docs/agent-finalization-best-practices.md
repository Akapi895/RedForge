# Agent Final-Response Governance Best Practices

[Back to English documentation](README.md)

Research date: 2026-07-28

This document focuses on a specific problem: an agent emits natural language that looks like a conclusion while tool calls, reasoning, planning, or sub-agent collaboration are not actually complete, and the frontend or orchestration layer displays it as the final response. The conclusion is straightforward: mature agent systems do not determine completion from the most recent assistant text. They use runtime state, tool state, verification results, and explicit terminal events together.

## 1. Core Conclusions

1. **A final response is a runtime event, not natural-language content.** Text such as “obtained,” “next step,” or “Huge breakthrough” can only be treated as a candidate observation or progress update, never as a completion signal.
2. **The process surface and delivery surface must be separated.** `thinking`, `reasoning_chain`, `planning`, `response_delta`, sub-agent replies, and tool output belong to the process surface. Only a `response` / `final` event that passes the final gate may be written to the main message bubble and `messages.content`.
3. **Complex tasks need a verifier, not a longer prompt.** A prompt can remind the model to be careful, but completion must be decided in code: whether tools, background executions, required plan steps, evidence, recorded facts/vulnerabilities, or cleanup actions remain pending or unaccounted for.
4. **Agent modes differ, but final-governance principles are consistent.** Single-agent, Deep, Plan-Execute, and Supervisor modes all need a final gate. Only the evidence source differs: tool traces for single-agent, sub-agent results for Deep, Replanner termination for Plan-Execute, and `exit` plus Supervisor aggregation for Supervisor mode.

## 2. Public Practices in Mature Agent Systems

| System | Public practice | Final-governance implication |
|---|---|---|
| Codex | OpenAI's Codex prompting guide recommends not forcing an upfront plan, preamble, or status updates in the prompt because this can cause a rollout to stop before completion. | Do not treat model-generated plans or status text as completion; the agent harness should own the execution loop and closeout. |
| Claude Code | Claude Code provides `PreToolUse`, `PostToolUse`, and `Stop` hooks; `PostToolUse` explicitly occurs after a tool succeeds. | Lifecycle events are more reliable than natural language. Attach verification, audit, and blocking to definite phase boundaries. |
| Claude Code Subagents | Sub-agents have separate contexts, custom system prompts, specific tool permissions, and independent permissions. They are useful for isolating large retrieval, log-reading, and file-reading tasks. | Sub-agent output is evidence material, not the main task's final conclusion. The main agent must aggregate, accept, and then finalize it. |
| Claude Code Plan Mode | Plan mode reads files and produces a plan; it does not edit before approval. | Planning and execution are different states; completing a plan does not complete the task. |
| Cursor Plan Mode | Cursor Plan Mode researches the codebase, asks clarifying questions, produces a reviewable plan, and waits for user confirmation before building. | Separating plan, review, and build in the UI prevents a plan from being mistaken for delivery. |
| OpenCode | OpenCode separates Build, Plan, Review, Debug, Docs, and other agents by tool permissions and purpose; the Plan agent analyzes and plans without modifying files. | Capability boundaries reduce accidental actions: an agent that can plan is not necessarily an agent that can execute to completion. |
| Eino ADK | Eino ADK provides event-driven output, Runner callbacks, interruption, checkpoints, and Supervisor and Plan-Execute collaboration primitives. Plan-Execute coordinates Planner, Executor, and Replanner. | The current project's direction is sound; event-driven capabilities should be further formalized as a finalization contract. |

Main references:

- OpenAI Codex Prompting Guide: https://developers.openai.com/cookbook/examples/gpt-5/codex_prompting_guide
- Claude Code Hooks: https://docs.anthropic.com/en/docs/claude-code/hooks
- Claude Code Subagents: https://docs.anthropic.com/en/docs/claude-code/sub-agents
- Claude Code Common Workflows: https://docs.anthropic.com/en/docs/claude-code/common-workflows
- Cursor Agent Best Practices: https://cursor.com/blog/agent-best-practices
- OpenCode Agents: https://opencode.ai/docs/agents/
- CloudWeGo Eino ADK: https://www.cloudwego.io/docs/eino/core_modules/eino_adk/
- CloudWeGo Eino ADK Patterns: https://www.cloudwego.io/docs/eino/overview/eino_adk0_1/

## 3. General Best Practices

### 1. Establish a Finalization Contract

Every execution entry point should produce one structured closeout object, and only that object may trigger a final response.

```go
type FinalizationDecision struct {
    Status             string   // in_progress | completed | blocked | failed | cancelled
    Finalizable        bool
    CompletionReason   string   // verified | user_cancelled | timeout | blocked | failed
    FinalText          string
    EvidenceVerified   bool
    EvidenceRefs       []string
    PendingToolRuns    []string
    PendingPlanSteps   []string
    PendingApprovals   []string
    MissingChecks      []string
}
```

Hard rules:

- When `Finalizable=false`, do not send a terminal `response` event.
- When `Status=in_progress`, only process events such as `progress`, `planning`, `tool_*`, and `reasoning_chain` may be sent.
- `FinalText` must not be empty, but non-empty text does not mean that finalization is allowed.
- If any of `PendingToolRuns`, `PendingPlanSteps`, or `PendingApprovals` is non-empty, the status cannot be `completed`.
- When `EvidenceVerified=false`, do not write a candidate output as a verified conclusion.

### 2. Fix SSE Event Semantics

Recommended event layers:

| Event | Display location | May write to `messages.content` | Description |
|---|---|---:|---|
| `progress` | Task status/timeline | No | Short progress update |
| `planning` | Execution details | No | Main-agent plan and interim judgment |
| `reasoning_chain` / `thinking` | Execution details | No | Reasoning/thinking summary |
| `tool_call` / `tool_result` | Execution details | No | Tool events |
| `eino_agent_reply` | Execution details | No | Sub-agent return material |
| `finalization_check` | Execution details | No | Verifier result |
| `finalization_auto_continue` | Execution details | No | Engineering continuation triggered by the verifier; `contextInjection=false` |
| `response` | Main message bubble | Yes | Use only when `data.finalized=true` |
| `done` | Stream close | No | Only means the stream ended, not that the task succeeded |
| `error` / `cancelled` | Main bubble or system notice | Yes, as a terminal failure state | Must include a reason |

### 3. Separate the Final Candidate from the Final Response

The model may produce a candidate conclusion, but it must first enter `final_candidate` or `planning`; the verifier then decides whether to promote it:

```text
assistant text
  -> candidate
  -> finalization gate
  -> response(finalized=true)
```

Do not do this:

```text
assistant text
  -> response
```

### 4. Stop-Time Verification

Following the Claude Code hook model, perform a deterministic check when an agent run stops:

- Every tool call has a corresponding tool result.
- Every background execution is terminal, or is explicitly recorded as still running with task status `in_progress` / `blocked`.
- Plan-Execute has no unexecuted required step.
- Supervisor has no unaggregated sub-agent result.
- Under an evidence-required policy, at least one queryable completed tool execution provides evidence.

### 5. Treat Sub-Agent Output as Evidence Only

Sub-agent returns cannot directly become the user's final response. The main agent must complete:

- Deduplication and conflict merging.
- Evidence-strength ranking.
- Uncertainty labeling.
- Scope-boundary confirmation.
- User-readable delivery.

### 6. Prompts Provide Soft Constraints; Code Provides Hard Constraints

The prompt may say:

```text
Interim observations must be marked as progress, not final.
Do not produce a final answer until verification is complete.
```

However, backend fields and the state machine must make the actual final decision. Otherwise, a natural-language conclusion that looks final can still be misclassified by the UI.

## 4. Current CyberStrikeAI Implementation Status

The project currently has an explicit final gate:

- [internal/agentfinalizer/decision.go](../internal/agentfinalizer/decision.go) is the single final-response decision contract.
- [internal/handler/finalization_helpers.go](../internal/handler/finalization_helpers.go) writes decisions to `process_details` and calls `UpdateAssistantMessageFinalize` only when `Finalizable=true`.
- [internal/handler/eino_single_agent.go](../internal/handler/eino_single_agent.go), [internal/handler/multi_agent.go](../internal/handler/multi_agent.go), [internal/handler/workflow_integration.go](../internal/handler/workflow_integration.go), and [internal/handler/batch_queue_executor.go](../internal/handler/batch_queue_executor.go) all connect the finalizer at closeout.
- [web/static/js/monitor.js](../web/static/js/monitor.js) treats only `response` events with `data.finalized === true` as final; non-finalized text is shown as a failed final-response check.
- [web/static/js/webshell.js](../web/static/js/webshell.js) marks streamed body text as candidate output and switches to completed state only for `response(finalized=true)`.
- [internal/agentfinalizer/decision_test.go](../internal/agentfinalizer/decision_test.go) covers pending tools, HITL, empty output, required evidence without execution evidence, failed evidence, and finalization with completed evidence.
- [internal/handler/finalization_auto_continue.go](../internal/handler/finalization_auto_continue.go) automatically continues for at most two segments when completed execution evidence is missing; continuation restores the existing model trace without injecting new user/system text.

Core rules of the current contract:

1. Model language is only a candidate. `RunResult.Response` cannot be promoted directly and must pass `agentfinalizer.Decide`.
2. Every `response` event must include terminal fields such as `finalized`, `finalizable`, `status`, `completionReason`, `evidenceVerified`, `evidenceRefs`, `pendingExecutionIds`, and `missingChecks`.
3. Incomplete tools block finalization; queued or running executions produce `in_progress/pending_tool_executions`.
4. Execution evidence must be declared by a structured policy. Chat requests use `finalization.requireExecutionEvidence`; WebShell, Workflow, batch, and robot entry points pass policy explicitly. A required policy needs at least one queryable completed tool execution; failed or cancelled records are insufficient.
5. Missing execution evidence first triggers engineering continuation and then blocking. Eino single-agent and multi-agent paths continue using the existing trace without additional context; after the limit, they write `blocked`.
6. HITL and empty output never finalize. Workflow waiting for approval, empty assistant text, and empty Eino placeholders produce blocking text rather than a successful summary.

## 5. Recommended Architecture for This Project

```text
Agent / Eino ADK events
  -> event normalizer
  -> process_details
  -> finalization verifier
  -> response(finalized=true)
  -> messages.content
```

### 1. Unified Backend Finalizer

Responsibilities:

- Receive `RunResult` / candidate text, `mcpExecutionIds`, session and assistant-message IDs, HITL state, and orchestration mode.
- Query tool execution state from the database and identify pending, completed, failed, and cancelled evidence.
- Return `FinalizationDecision`.
- Never call high-risk tools; only inspect state and evidence.

### 2. Terminal Fields in `RunResult`

[internal/multiagent/runner.go](../internal/multiagent/runner.go) contains these terminal fields:

```go
type RunResult struct {
    Response             string
    MCPExecutionIDs      []string
    LastAgentTraceInput  string
    LastAgentTraceOutput string

    Finalized           bool
    Status              string
    CompletionReason    string
    EvidenceVerified    bool
    EvidenceRefs        []string
    PendingExecutionIDs []string
    MissingChecks       []string
}
```

### 3. Conditions for Sending `response`

Single-agent, multi-agent, workflow, and batch closeout should use the same gate:

```go
decision := h.finalizeAgentRunForDelivery(...)
if !decision.Finalizable {
    sendEvent("finalization_check", "The task has not met the final-response conditions yet.", decision)
    sendEvent("response", finalizationBlockedMessage(decision), finalizationResponsePayload(decision, extra))
    return
}

sendEvent("response", decision.FinalText, finalizationResponsePayload(decision, extra))
```

### 4. Frontend Trusts Only `finalized=true`

In `case 'response'` in [web/static/js/monitor.js](../web/static/js/monitor.js), apply a hard check:

```js
const responseFinalized = isFinalizedResponseData(responseData);
const bubbleText = responseFinalized
  ? resolvedResponseText
  : (event.message || 'The task has not met the final-response conditions; no successful conclusion is generated yet.');
markAssistantFinalizationState(assistantIdFinal, responseData);
```

The WebShell side follows the same rule: `response_delta` may provide a live preview, but the UI must mark it as execution output; only `response(finalized=true)` is a completed state.

### 5. Final Gates by Mode

| Mode | Who produces the final candidate | Who decides finalization | Required checks |
|---|---|---|---|
| Eino single-agent | Last single-agent assistant text | Finalizer | No pending tools, complete evidence references, terminal task state |
| Deep | Main-agent aggregate text | Finalizer | Sub-agent results aggregated; sub-agent text cannot finalize directly; tools terminal |
| Plan-Execute | Replanner's final aggregate text | Replanner + Finalizer | Executor step output cannot finalize; plan steps completed or explicitly blocked |
| Supervisor | Supervisor `exit` / aggregate text | Supervisor + Finalizer | Transfer returned; no unprocessed expert result; Supervisor provides the final unified response |

### 6. Evidence Gate for Security Testing

For security testing, WebShell, batch verification, Workflow, and multi-agent execution paths, a final response must satisfy at least:

- A clear target and authorization-scope marker.
- Reviewable evidence references such as tool execution IDs, request/response summaries, screenshot paths, command-output summaries, or fact/vulnerability record IDs.
- Identity or impact verification, rather than marker text alone.
- A project-blackboard or vulnerability-library record, or an explicit explanation that no project binding prevented recording.
- High-risk actions cleaned up, rolled back, or cancelled, or an explicit explanation of why cleanup was not performed.
- Scans, commands, WebShell, and C2 tasks that are still running must not be implicitly treated as complete.

This gate is a governance rule and does not require the final report to expose sensitive exploit details; an evidence summary and internal references are sufficient.

## 6. Implementation Status and Further Enhancements

### P0: Fix False Finalization (Implemented)

1. `FinalizationDecision` has been introduced.
2. Main agent SSE `response` events include `data.finalized`, `finalizable`, `status`, and `completionReason`.
3. `monitor.js` and `webshell.js` distinguish candidate output from final responses using `finalized=true`.
4. `RunResult.Response` remains as a compatibility field, but the finalizer now owns its promotion semantics. It may later be split into `CandidateResponse` and `FinalResponse` to reduce misuse.
5. Plan-Execute, Deep, Supervisor, and Eino Single modes all close through the unified handler gate.

### P1: Complete the Evidence Chain (Partially Implemented)

1. `mcp_execution:<id>` is used as the base evidence reference.
2. The `finalization_check` event displays pending executions and missing checks.
3. Execution entry points use an explicit execution-evidence policy; when required evidence lacks a completed tool record, the Eino main path continues without injection before blocking after the limit.
4. Future work: add granular evidence references for `record_vulnerability`, `upsert_project_fact`, and project-blackboard records.
5. Future work: make final report templates consistently include conclusion, evidence, risks/uncertainty, and next actions.

### P2: Experience and Observability (Future Enhancement)

1. Show `in_progress / verifying / finalizing / completed / blocked` on task cards.
2. Add finalizer logs and metrics: false-block rate, missing-evidence types, and pending-tool count.
3. Support a “Continue verification” button that generates the next input from `FinalizationDecision.MissingChecks`.

## 7. Acceptance-Test Recommendations

At minimum, add these regression tests:

1. **Reasoning text does not finalize:** a seemingly complete candidate appears in `reasoning_chain` without completed tool evidence; the main bubble must not show a successful conclusion.
2. **Interim main-agent output does not finalize:** `response_start/delta` emits “continue verification next”; it belongs only in the planning timeline.
3. **Pending background tools do not finalize:** an execution ID is running; even with a summary, the result remains `in_progress`.
4. **Plan-Execute Executor output does not finalize:** the Executor says “breakthrough succeeded” while the Replanner has not ended; `messages.content` must not finalize.
5. **Supervisor sub-agent output does not finalize:** a sub-agent returns a definitive conclusion while Supervisor has not exited; it belongs only in `eino_agent_reply`.
6. **Final events must include `finalized`:** an old-format `response` without `finalized=true` is treated as candidate/warning content, and the main bubble shows a blocked state.
7. **Failure and cancellation are valid terminal states:** `error` / `cancelled` can update the assistant message, but `completionReason` must be `failed` / `user_cancelled`, never a successful completion.

## 8. Recommended Default Strategy

For CyberStrikeAI:

```text
eino_single: suitable for lightweight tasks, but the final gate must be enabled
deep: recommended by default for complex security testing
plan_execute: recommended for explicit goals requiring strict plan-execute-replan cycles
supervisor: use for multi-expert routing; not the default general-purpose mode
```

In one sentence:

```text
messages.content may come only from FinalizationDecision.FinalText;
process_details may show every process event;
the frontend may treat only response(finalized=true) as the final response.
```
