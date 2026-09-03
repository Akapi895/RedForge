// Package projectprompt contains shared system-prompt fragments for project
// facts, evidence, and vulnerability recording.
package projectprompt

import (
	"strings"

	"cyberstrike-ai/internal/mcp/builtin"
)

const (
	factRhythmCore              = "Do not wait until the end of the engagement to persist findings. After confirming each new observation (open port/service version, entry path, authentication state or credential characteristic, exploitable condition, or attack-surface change), immediately call `upsert_project_fact` and update the same fact_key. After verifying a reproducible vulnerability with proof and impact, immediately call `record_vulnerability`; the same discovery may be recorded in both places. Persist before continuing so details survive context compaction. If no project is bound, state that the blackboard cannot be updated and retain an evidence summary in the current response."
	factRhythmCoordinatorSuffix = " When a delegated task returns a new observation or vulnerability, the coordinator must persist it promptly; never assume that a sub-agent already recorded it."
	factRhythmSubAgentSuffix    = " If the tool set does not include the persistence tools, finish with a structured pending-record item containing a suggested fact_key, summary, and body/PoC details for the coordinator to store immediately."
)

// FactRecordingIncrementalRhythmMarkdown returns the incremental recording
// guidance used by agent files and documentation.
func FactRecordingIncrementalRhythmMarkdown(coordinator, subAgent bool) string {
	var b strings.Builder
	b.WriteString("- **Record continuously (mandatory cadence):** ")
	b.WriteString(factRhythmCore)
	if coordinator {
		b.WriteString(factRhythmCoordinatorSuffix)
	}
	if subAgent {
		b.WriteString(factRhythmSubAgentSuffix)
	}
	return b.String()
}

func factRecordingIncrementalRhythmBuiltin(coordinator, subAgent bool) string {
	var b strings.Builder
	b.WriteString("- **Record continuously (mandatory cadence):** Do not wait until the end of the engagement to persist findings. After confirming each new observation (open port/service version, entry path, authentication state or credential characteristic, exploitable condition, or attack-surface change), immediately call ")
	b.WriteString(builtin.ToolUpsertProjectFact)
	b.WriteString(" and update the same fact_key. After verifying a reproducible vulnerability with proof and impact, immediately call ")
	b.WriteString(builtin.ToolRecordVulnerability)
	b.WriteString("; the same discovery may be recorded in both places. Persist before continuing so details survive context compaction. If no project is bound, state that the blackboard cannot be updated and retain an evidence summary in the current response.")
	if coordinator {
		b.WriteString(factRhythmCoordinatorSuffix)
	}
	if subAgent {
		b.WriteString(factRhythmSubAgentSuffix)
	}
	return b.String()
}

func factEdgeRecordingGuidance() string {
	return "" + `### Fact relationship edges (links)

- When writing a **finding / chain / exploit / poc**, provide ` + "`links`" + ` in ` + "`upsert_project_fact`" + `. Prefer ` + "`from`" + ` to point from the source fact to the current fact: ` + "`from`" + ` → current ` + "`fact_key`" + `.
- **Minimum requirement:** every finding must have at least one ` + "`from=target/*`" + ` with ` + "`type=discovered_on`" + `; an exploit recorded on a finding should use ` + "`from=exploit/*`" + ` with ` + "`type=exploits`" + `.
- **Common types:** ` + "`discovered_on`" + `, ` + "`depends_on`" + `, ` + "`leads_to`" + `, ` + "`enables`" + `, ` + "`exploits`" + `, ` + "`contains`" + `, ` + "`part_of`" + `, and ` + "`supports`" + `.
- Omitting ` + "`links`" + ` during an update preserves existing edges; providing it replaces all edges for the current fact.
- A human-readable dependency section may coexist with ` + "`links`" + `; structured relationships are defined by ` + "`links`" + `.
`
}

func factRecordingGuidanceBlock() string {
	return "" + `### Fact recording standard (reproducibility and knowledge retention)

- **summary:** one searchable line containing what was observed, where it was observed, and how it was triggered or verified. Do not write a conclusion only, such as “SQL injection exists”.
- **body:** the complete reproducible context stored in ` + "`upsert_project_fact`" + `. The index contains only the summary; later sessions must call ` + "`get_project_fact`" + ` to retrieve the body.
- **category / fact_key suggestions:**
  - Environment observations: ` + "`target/`" + `, ` + "`auth/`" + `, ` + "`infra/`" + `, ` + "`business/`" + `.
  - Discovery and exploitation: ` + "`finding/`" + `, ` + "`chain/`" + `, ` + "`exploit/`" + `, ` + "`poc/`" + `. The body must include the entry point, attack steps, raw request/response or command, evidence, and related vulnerability ID.
- **Separation from vulnerability records:** ` + "`record_vulnerability`" + ` stores deliverable findings; facts store all context needed to reproduce them, including failed attempts, bypasses, and session dependencies. Both records may be needed.
- Keep the same ` + "`fact_key`" + ` when updating a discovery. Do not scatter one finding across multiple keys.
`
}

// FactRecordingBlackboardSection returns the complete project-blackboard and
// vulnerability-recording prompt fragment for the primary agent.
func FactRecordingBlackboardSection(coordinatorDelegate bool) string {
	var b strings.Builder
	b.WriteString("## Project blackboard (facts) and vulnerability records\n\n")
	b.WriteString("When the conversation is bound to a project, the system injects a project-blackboard index containing fact keys and summaries. **If the summary is insufficient, call ")
	b.WriteString(builtin.ToolGetProjectFact)
	b.WriteString("(fact_key) to retrieve the body; never invent missing details.**\n\n")
	b.WriteString(factRecordingIncrementalRhythmBuiltin(coordinatorDelegate, false))
	b.WriteString("\n\n")
	b.WriteString("- **Environment, target, and authentication observations:** use ")
	b.WriteString(builtin.ToolUpsertProjectFact)
	b.WriteString(" with a `category/slug` fact_key such as `target/primary_domain`; update the same key and include ports, versions, credential characteristics, and evidence sources.\n")
	b.WriteString("- **Discovery and exploitation context:** use `finding/`, `chain/`, `exploit/`, or `poc/` prefixes. The body must contain the complete attack chain: entry point → steps → raw request/response or command → observed result → related_vulnerability_id. Do not record a conclusion only.\n")
	b.WriteString("- **Deliverable vulnerabilities:** use ")
	b.WriteString(builtin.ToolRecordVulnerability)
	b.WriteString(" with title, severity, type, target, proof/PoC, impact, and remediation. Before recording, use ")
	b.WriteString(builtin.ToolListVulnerabilities)
	b.WriteString(" to check for duplicates; use ")
	b.WriteString(builtin.ToolGetVulnerability)
	b.WriteString("(id) for details.\n")
	b.WriteString("- The same discovery may need both a reproducible fact/chain record and a formal finding. Mark false positives with ")
	b.WriteString(builtin.ToolDeprecateProjectFact)
	b.WriteString(" or the vulnerability status `false_positive`.\n")
	b.WriteString("- For many facts, use ")
	b.WriteString(builtin.ToolListProjectFacts)
	b.WriteString(" / ")
	b.WriteString(builtin.ToolSearchProjectFacts)
	b.WriteString(".\n\n")
	b.WriteString(factEdgeRecordingGuidance())
	b.WriteString("\n\n")
	b.WriteString(factRecordingGuidanceBlock())
	b.WriteString("\n\nSeverity values: critical / high / medium / low / info. Proof must contain sufficient evidence such as request/response data, screenshots, or command output.")
	return b.String()
}

// FactRecordingSubAgentSection returns the incremental recording guidance for
// sub-agents.
func FactRecordingSubAgentSection() string {
	return "## Continuous recording\n\n" + factRecordingIncrementalRhythmBuiltin(false, true) + "\n"
}

// FactRecordingBlackboardSectionMarkdown returns the same guidance with
// literal tool names for agent Markdown files.
func FactRecordingBlackboardSectionMarkdown(coordinatorDelegate bool) string {
	var b strings.Builder
	b.WriteString("## Project blackboard (facts) and vulnerability records\n\n")
	b.WriteString("When the conversation is bound to a project, the system injects a project-blackboard index containing fact keys and summaries. **If the summary is insufficient, call `get_project_fact(fact_key)` to retrieve the body; never invent missing details.**\n\n")
	b.WriteString(FactRecordingIncrementalRhythmMarkdown(coordinatorDelegate, false))
	b.WriteString("\n\n")
	b.WriteString("- **Environment, target, and authentication observations:** use **`upsert_project_fact`** with a `category/slug` fact_key such as `target/primary_domain`; update the same key and include ports, versions, credential characteristics, and evidence sources.\n")
	b.WriteString("- **Discovery and exploitation context:** use `finding/`, `chain/`, `exploit/`, or `poc/` prefixes. The body must contain the complete attack chain: entry point → steps → raw request/response or command → observed result → related_vulnerability_id. Do not record a conclusion only.\n")
	b.WriteString("- **Deliverable vulnerabilities:** use **`record_vulnerability`** with title, description, severity, type, target, proof/PoC, impact, and remediation. Severity values are critical / high / medium / low / info.\n")
	b.WriteString("- The same discovery may need both a reproducible chain record and a formal finding. Mark false positives with **`deprecate_project_fact`** or vulnerability status `false_positive`.\n")
	b.WriteString("- For many facts, use **`list_project_facts`** / **`search_project_facts`**.\n\n")
	b.WriteString(factEdgeRecordingGuidance())
	b.WriteString("\n\n")
	b.WriteString(factRecordingGuidanceBlock())
	b.WriteString("\n\nProof must contain sufficient evidence such as request/response data, screenshots, or command output.")
	return b.String()
}

// FactEdgeRecordingGuidance returns the relationship-edge guidance.
func FactEdgeRecordingGuidance() string { return factEdgeRecordingGuidance() }

// FactRecordingGuidanceBlock returns the fact-recording standard.
func FactRecordingGuidanceBlock() string { return factRecordingGuidanceBlock() }
