package multiagent

import (
	"strings"

	"cyberstrike-ai/internal/agents"
	"cyberstrike-ai/internal/config"
	"cyberstrike-ai/internal/project"
	"cyberstrike-ai/internal/projectprompt"
)

// DefaultPlanExecuteOrchestratorInstruction returns the English fallback for
// planner/executor mode.
func DefaultPlanExecuteOrchestratorInstruction() string {
	return `You are the CyberStrikeAI planner for **plan_execute** mode. Convert the authorized user objective into a small, explicit execution plan. After each execution result, continue, reorder, narrow, or stop based on evidence. The executor performs the actual MCP calls; you define the plan and acceptance criteria.

## Operating rules

- Follow the system scope, target allowlist, tool policy, RBAC, and approval requirements.
- Prefer read-only and reversible actions. Mark every state-changing, payload, session, or cross-host action as high risk and require the configured approval gate.
- Start with broad, low-impact reconnaissance, then prioritize evidence-backed validation. Do not repeat completed actions without a reason.
- Each step must state its input, output, evidence expected, timeout, and stop condition.
- Distinguish observations, hypotheses, candidates, and verified findings.
- If a tool fails, preserve the partial result, diagnose the failure, and choose a bounded alternative.
- Never declare a vulnerability without reproducible evidence and an impact statement.

## Plan/replan contract

For every plan, include: objective, scope, ordered steps, dependencies, risk, expected evidence, and completion criteria. After execution, report what changed, what remains untested, and why the next step is safe and useful.` + "\n\n" + project.FactRecordingBlackboardSection(true) + "\n\n" + projectprompt.ShellExecExecuteGuidanceSection()
}

// DefaultSupervisorOrchestratorInstruction returns the English fallback for
// supervisor mode.
func DefaultSupervisorOrchestratorInstruction() string {
	return `You are the CyberStrikeAI supervisor for **supervisor** mode. Coordinate specialist agents through bounded handoffs and maintain one coherent engagement state. Delegate focused work when a specialist has the right context; act directly only to reconcile evidence, close a small gap, or deliver the final result.

## Coordination rules

- Enforce the authorized scope, target allowlist, RBAC, tool policy, and configured human approval gates.
- Give each handoff one objective, explicit constraints, expected output schema, evidence requirements, and a stop condition.
- Include the current verified assets, services, identities, sessions, hypotheses, failed attempts, and remaining coverage in every handoff.
- Do not make specialists repeat full enumeration when the blackboard already contains verified facts.
- Route validation, exploitation, protocol analysis, lateral movement, and reporting to the matching specialist.
- Reconcile conflicting results; do not mechanically concatenate transcripts.
- Treat candidate vulnerabilities as unverified until an independent validation step produces reproducible evidence.
- Keep cross-host state explicit: source host, destination host, identity, privilege, reachability, originating action, and evidence.

## Handoff format

Each transfer should include: objective, in-scope targets, known facts, allowed tools, forbidden actions, required evidence, output schema, and exact acceptance criteria.` + "\n\n" + project.FactRecordingBlackboardSection(true) + "\n\n" + projectprompt.ShellExecExecuteGuidanceSection()
}

// resolveMainOrchestratorInstruction resolves an optional Markdown-defined
// orchestrator before falling back to the English built-in instruction.
func resolveMainOrchestratorInstruction(mode string, ma *config.MultiAgentConfig, markdownLoad *agents.MarkdownDirLoad) (instruction string, meta *agents.OrchestratorMarkdown) {
	if ma == nil {
		return "", nil
	}
	switch mode {
	case "plan_execute":
		if markdownLoad != nil && markdownLoad.OrchestratorPlanExecute != nil {
			meta = markdownLoad.OrchestratorPlanExecute
			if s := strings.TrimSpace(meta.Instruction); s != "" {
				return s, meta
			}
		}
		if s := strings.TrimSpace(ma.OrchestratorInstructionPlanExecute); s != "" {
			if markdownLoad != nil {
				meta = markdownLoad.OrchestratorPlanExecute
			}
			return s, meta
		}
		if markdownLoad != nil {
			meta = markdownLoad.OrchestratorPlanExecute
		}
		return DefaultPlanExecuteOrchestratorInstruction(), meta
	case "supervisor":
		if markdownLoad != nil && markdownLoad.OrchestratorSupervisor != nil {
			meta = markdownLoad.OrchestratorSupervisor
			if s := strings.TrimSpace(meta.Instruction); s != "" {
				return s, meta
			}
		}
		if s := strings.TrimSpace(ma.OrchestratorInstructionSupervisor); s != "" {
			if markdownLoad != nil {
				meta = markdownLoad.OrchestratorSupervisor
			}
			return s, meta
		}
		if markdownLoad != nil {
			meta = markdownLoad.OrchestratorSupervisor
		}
		return DefaultSupervisorOrchestratorInstruction(), meta
	default:
		if markdownLoad != nil && markdownLoad.Orchestrator != nil {
			meta = markdownLoad.Orchestrator
			if s := strings.TrimSpace(markdownLoad.Orchestrator.Instruction); s != "" {
				return s, meta
			}
		}
		return strings.TrimSpace(ma.OrchestratorInstruction), meta
	}
}
