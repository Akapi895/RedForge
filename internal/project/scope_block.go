package project

import (
	"encoding/json"
	"fmt"
	"strings"

	"cyberstrike-ai/internal/config"
	"cyberstrike-ai/internal/database"
)

// projectScopePayload parses projects.scope_json (conventional, extensible fields).
type projectScopePayload struct {
	Targets []string `json:"targets"`
	Exclude []string `json:"exclude"`
	Notes   string   `json:"notes"`
}

// BuildScopeBlock formats project scope_json as an authorization-scope block readable by the Agent.
func BuildScopeBlock(proj *database.Project) string {
	if proj == nil {
		return ""
	}
	raw := strings.TrimSpace(proj.ScopeJSON)
	if raw == "" {
		return ""
	}

	var payload projectScopePayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return fmt.Sprintf("## Project Testing Scope (project: %s)\n(scope_json is not valid JSON; review the configuration manually)\n```\n%s\n```\n"+
			"Test only explicitly authorized targets; stop and explain before going out of scope.\n", proj.Name, truncateRunes(raw, 800))
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("## Project Testing Scope (project: %s, id: %s)\n", proj.Name, proj.ID))
	b.WriteString("The following authorization boundaries **must be observed**: test only the listed targets, avoid excluded entries, and do not expand the scope without permission.\n")

	if len(payload.Targets) > 0 {
		b.WriteString("\n**Authorized targets**:\n")
		for _, t := range payload.Targets {
			t = strings.TrimSpace(t)
			if t != "" {
				b.WriteString("- " + t + "\n")
			}
		}
	}
	if len(payload.Exclude) > 0 {
		b.WriteString("\n**Explicit exclusions**:\n")
		for _, t := range payload.Exclude {
			t = strings.TrimSpace(t)
			if t != "" {
				b.WriteString("- " + t + "\n")
			}
		}
	}
	if n := strings.TrimSpace(payload.Notes); n != "" {
		b.WriteString("\n**Notes**:\n" + n + "\n")
	}
	if len(payload.Targets) == 0 && len(payload.Exclude) == 0 && strings.TrimSpace(payload.Notes) == "" {
		b.WriteString("\n(scope_json is configured but contains no recognized targets/exclude/notes fields; the original content is included for reference)\n```json\n")
		b.WriteString(truncateRunes(raw, 1200))
		b.WriteString("\n```\n")
	}
	b.WriteString("\nIf a target is not listed in targets or matches exclude, do not scan or exploit it; continue only after the user explicitly expands authorization.\n")
	return b.String()
}

func truncateRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + "…"
}

// BuildProjectBlackboardBlock combines the testing scope with the fact-blackboard index.
func BuildProjectBlackboardBlock(db *database.DB, projectID string, cfg config.ProjectConfig) (string, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return "", nil
	}
	proj, err := db.GetProject(projectID)
	if err != nil {
		return "", err
	}
	parts := []string{}
	if scope := strings.TrimSpace(BuildScopeBlock(proj)); scope != "" {
		parts = append(parts, scope)
	}
	index, err := BuildFactIndexBlock(db, projectID, cfg)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(index) != "" {
		parts = append(parts, index)
	}
	return strings.Join(parts, "\n\n"), nil
}
