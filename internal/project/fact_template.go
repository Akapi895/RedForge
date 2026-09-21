package project

import (
	"fmt"
	"strings"

	"cyberstrike-ai/internal/projectprompt"
)

// Fact category constants written to the category field of upsert_project_fact.
const (
	FactCategoryTarget   = "target"
	FactCategoryAuth     = "auth"
	FactCategoryInfra    = "infra"
	FactCategoryBusiness = "business"
	FactCategoryFinding  = "finding"
	FactCategoryChain    = "chain"
	FactCategoryExploit  = "exploit"
	FactCategoryPOC      = "poc"
	FactCategoryNote     = "note"
)

// RequiresAttackChainBody reports whether a fact should include reproducible attack-chain/exploit details in body rather than only summary.
func RequiresAttackChainBody(category, factKey string) bool {
	c := strings.ToLower(strings.TrimSpace(category))
	switch c {
	case FactCategoryFinding, FactCategoryChain, FactCategoryExploit, FactCategoryPOC, "vuln":
		return true
	}
	key := strings.ToLower(strings.TrimSpace(factKey))
	for _, prefix := range []string{"finding/", "chain/", "exploit/", "poc/"} {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// IsSparseFactBody returns true when an attack-chain fact body is too short or lacks key sections (soft validation that does not block writes).
func IsSparseFactBody(category, factKey, body string) bool {
	if !RequiresAttackChainBody(category, factKey) {
		return false
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return true
	}
	lower := strings.ToLower(body)
	// Include at least one reproducibility clue: steps, requests, commands, or code blocks
	hasSteps := strings.Contains(lower, "attack chain") || strings.Contains(lower, "## attack") ||
		strings.Contains(lower, "\u653b\u51fb\u94fe") || strings.Contains(lower, "## \u653b\u51fb") ||
		strings.Contains(lower, "## exploit") || strings.Contains(lower, "## poc")
	hasHTTP := strings.Contains(lower, "```http") || strings.Contains(lower, "```bash") ||
		strings.Contains(lower, "curl ") || strings.Contains(lower, "get ") || strings.Contains(lower, "post ")
	hasReq := strings.Contains(lower, "request") || strings.Contains(lower, "response") || strings.Contains(lower, "payload") ||
		strings.Contains(lower, "\u8bf7\u6c42") || strings.Contains(lower, "\u54cd\u5e94")
	// Without structural clues such as an attack chain, POC, or request, treat the body as conclusion-only regardless of length
	return !(hasSteps || hasHTTP || hasReq)
}

// FactBodyTemplate returns a recommended body Markdown skeleton by category for the Agent to populate with real content.
func FactBodyTemplate(category, factKey string) string {
	if RequiresAttackChainBody(category, factKey) {
		return attackChainFactBodyTemplate
	}
	return envFactBodyTemplate
}

const attackChainFactBodyTemplate = `## Conclusion (verifiable, one sentence)
<Do not merely write "a vulnerability exists"; state the type, location, and trigger conditions>

## Target and Entry Point
- Target: <URL / IP:Port / hostname>
- Entry point: <path / endpoint / parameter>
- Preconditions: <anonymous / role / Cookie / other dependencies>

## Attack Chain (step-by-step reproducible)
1. <Reconnaissance/discovery>
2. <Exploitation/trigger>
3. <Impact proof (file read, RCE output, unauthorized data, etc.)>

## Exploit / POC
### Request
` + "```http\n<METHOD> <path> HTTP/1.1\nHost: ...\n...\n\n<body>\n```" + `

### Response / Observation
<Key response excerpt, status code, and differences>

### Command / Script (if applicable)
` + "```bash\n<command>\n```" + `

## Key Evidence
- <Tool-output summary / screenshot path / conversation or message ID>

## Relationships
- related_vulnerability_id: <optional; corresponding record_vulnerability ID>
- links (upsert argument): [{ "from": "<fact_key>", "type": "discovered_on|..." }] (from → current fact)
- Dependency fact (human-readable mirror in body): <fact_key, such as auth/session_cookie>

## Notes and Uncertainty
<Hypotheses awaiting validation, environmental differences, and bypass-attempt records>`

const envFactBodyTemplate = `## Summary
<Core knowledge represented by this fact>

## Details
<Ports, versions, paths, credential characteristics, business rules, etc.>

## Sources and Evidence
<Command output, response excerpts, and discovery time>

## Relationships
- Related fact_key: <optional>`

// FactRecordingGuidanceBlock provides system guidance requiring facts to preserve attack-chain context rather than conclusions alone.
func FactRecordingGuidanceBlock() string {
	return projectprompt.FactRecordingGuidanceBlock()
}

// SparseBodyWarning returns a tool warning when an attack-chain fact body is insufficient without blocking the save.
func SparseBodyWarning(category, factKey string) string {
	if !IsSparseFactBody(category, factKey, "") {
		return ""
	}
	return fmt.Sprintf(
		"\n\n⚠ Notice: category=%q / fact_key=%q identifies an attack-chain fact, but body is empty or too brief. Add the complete attack chain and POC (using the template) so subsequent audits can reproduce it.\nRecommended body skeleton:\n%s",
		category, factKey, FactBodyTemplate(category, factKey),
	)
}

// SparseBodyWarningIfNeeded determines from the actual body whether to append a warning.
func SparseBodyWarningIfNeeded(category, factKey, body string) string {
	if !IsSparseFactBody(category, factKey, body) {
		return ""
	}
	return SparseBodyWarning(category, factKey)
}
