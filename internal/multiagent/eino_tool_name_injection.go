package multiagent

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/components/tool"
)

// injectToolNamesOnlyInstruction prepends a compact tool-name-only section into
// the system instruction so the model can reference current callable names.
// toolSearchMiddlewareActive must be true when prependEinoMiddlewares mounted toolsearch (dynamic tools); do not infer this
// by scanning tool names — tool_search is injected by middleware and is usually absent from the pre-split tools list.
func injectToolNamesOnlyInstruction(ctx context.Context, instruction string, tools []tool.BaseTool, toolSearchMiddlewareActive bool) string {
	names := collectToolNames(ctx, tools)
	if len(names) == 0 {
		return strings.TrimSpace(instruction)
	}
	hasToolSearch := toolSearchMiddlewareActive
	if !hasToolSearch {
		for _, n := range names {
			if strings.EqualFold(strings.TrimSpace(n), "tool_search") {
				hasToolSearch = true
				break
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("The following is the tool-name index bound to this conversation (names only; no parameter JSON Schema).\n")
	sb.WriteString("If tool_search is enabled, the list may include non-resident tools that are not included in the tool definitions sent in the current turn. Before seeing a tool's complete schema, never infer parameters from its name.\n")
	for _, name := range names {
		sb.WriteString("- ")
		sb.WriteString(name)
		sb.WriteByte('\n')
	}
	sb.WriteString("\nUsage rules:\n")
	sb.WriteString("1) The list above is a name index only; it contains no parameter definitions. Never guess parameter names, types, enum values, or whether a field is required.\n")
	if hasToolSearch {
		sb.WriteString("MANDATORY / HIGHEST PRIORITY: tool_search (the dynamic tool pool) is enabled for this conversation. If a tool appears in the name index but its complete parameter schema is absent from the tools definition attached to the current request, you must call tool_search first. Skipping tool_search to save tokens or time and calling the business tool directly is explicitly prohibited.\n")
		sb.WriteString("2) Default policy: whenever any part of the target tool's parameter definition is uncertain, call tool_search first. One extra tool_search call is preferable to blindly calling a business tool without its schema.\n")
		sb.WriteString("3) Order: call tool_search first (its only required argument is regex_pattern, a regex matching tool names such as nuclei or ^exact_tool_name$); in a later turn confirm that the target tool appears in the tools list and read its schema; then make the real tool call.\n")
		sb.WriteString("4) tool_search returns only the matching tool-name list; the schema is sent in the next turn after the tool is unlocked. Never fabricate JSON arguments before the schema appears.\n")
		sb.WriteString("5) Never invent a tool name that does not exist.\n\n")
	} else {
		sb.WriteString("2) Before calling a specific tool, confirm its parameter requirements from the tool definition in the current request; ask for clarification if uncertain.\n")
		sb.WriteString("3) Never invent a tool name that does not exist.\n\n")
	}
	if s := strings.TrimSpace(injectShellToolGuidance("", names)); s != "" {
		sb.WriteString(s)
		sb.WriteString("\n\n")
	}
	if s := strings.TrimSpace(instruction); s != "" {
		sb.WriteString(s)
	}
	return sb.String()
}

func collectToolNames(ctx context.Context, tools []tool.BaseTool) []string {
	if len(tools) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(tools))
	out := make([]string, 0, len(tools))
	for _, t := range tools {
		if t == nil {
			continue
		}
		info, err := t.Info(ctx)
		if err != nil || info == nil {
			continue
		}
		name := strings.TrimSpace(info.Name)
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, name)
	}
	return out
}
