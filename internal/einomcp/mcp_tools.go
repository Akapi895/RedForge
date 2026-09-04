package einomcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cyberstrike-ai/internal/agent"
	"cyberstrike-ai/internal/security"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/eino-contrib/jsonschema"
)

// ExecutionRecorder is optional and is called when an MCP tool succeeds and returns an execution ID (used to aggregate mcpExecutionIds).
// toolCallID comes from Eino compose.GetToolCallID and associates the invocation with the displayed result after reduction.
type ExecutionRecorder func(executionID, toolCallID string)

// ToolErrorPrefix propagates the IsError marker from an internal MCP execution result to the multi-agent layer.
// Eino's tool channel currently returns strings only, so a prefix identifies errors and the multi-agent runner later parses them into success/isError.
const ToolErrorPrefix = "__CYBERSTRIKE_AI_TOOL_ERROR__\n"

// ToolsFromDefinitions converts OpenAI-style tool definitions used by a single agent into Eino InvokableTools that execute through the agent's MCP path.
// invokeNotify is optional and shared with runEinoADKAgentLoop; when InvokableRun returns, it triggers UI and pending-state cleanup (deduplicated against ADK Tool events).
// einoAgentName is the Name of the ChatModelAgent owning this tool set (main-agent or sub-agent ID), used for the einoAgent field in SSE.
func ToolsFromDefinitions(
	ag *agent.Agent,
	holder *ConversationHolder,
	defs []agent.Tool,
	rec ExecutionRecorder,
	toolOutputChunk func(toolName, toolCallID, chunk string),
	invokeNotify *ToolInvokeNotifyHolder,
	einoAgentName string,
) ([]tool.BaseTool, error) {
	out := make([]tool.BaseTool, 0, len(defs))
	for _, d := range defs {
		if d.Type != "function" || d.Function.Name == "" {
			continue
		}
		info, err := toolInfoFromDefinition(d)
		if err != nil {
			return nil, fmt.Errorf("tool %q: %w", d.Function.Name, err)
		}
		out = append(out, &mcpBridgeTool{
			info:          info,
			name:          d.Function.Name,
			agent:         ag,
			holder:        holder,
			record:        rec,
			chunk:         toolOutputChunk,
			invokeNotify:  invokeNotify,
			einoAgentName: strings.TrimSpace(einoAgentName),
		})
	}
	return out, nil
}

func toolInfoFromDefinition(d agent.Tool) (*schema.ToolInfo, error) {
	fn := d.Function
	raw, err := json.Marshal(fn.Parameters)
	if err != nil {
		return nil, err
	}
	var js jsonschema.Schema
	if len(raw) > 0 && string(raw) != "null" && string(raw) != "{}" {
		if err := json.Unmarshal(raw, &js); err != nil {
			return nil, err
		}
	}
	if js.Type == "" {
		js.Type = string(schema.Object)
	}
	if js.Properties == nil && js.Type == string(schema.Object) {
		// Empty argument object
	}
	return &schema.ToolInfo{
		Name:        fn.Name,
		Desc:        fn.Description,
		ParamsOneOf: schema.NewParamsOneOfByJSONSchema(&js),
	}, nil
}

type mcpBridgeTool struct {
	info          *schema.ToolInfo
	name          string
	agent         *agent.Agent
	holder        *ConversationHolder
	record        ExecutionRecorder
	chunk         func(toolName, toolCallID, chunk string)
	invokeNotify  *ToolInvokeNotifyHolder
	einoAgentName string
}

func (m *mcpBridgeTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	_ = ctx
	return m.info, nil
}

func (m *mcpBridgeTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (out string, err error) {
	_ = opts
	toolCallID := compose.GetToolCallID(ctx)
	defer func() {
		if m.invokeNotify == nil {
			return
		}
		tid := strings.TrimSpace(toolCallID)
		if tid == "" {
			return
		}
		success := err == nil && !strings.HasPrefix(out, ToolErrorPrefix)
		body := out
		if err != nil {
			success = false
		} else if strings.HasPrefix(out, ToolErrorPrefix) {
			success = false
			body = strings.TrimPrefix(out, ToolErrorPrefix)
		}
		m.invokeNotify.Fire(tid, m.name, m.einoAgentName, success, body, err)
	}()
	return runMCPToolInvocation(ctx, m.agent, m.holder, m.name, argumentsInJSON, m.record, m.chunk)
}

// runMCPToolInvocation is shared with mcpBridgeTool.InvokableRun.
func runMCPToolInvocation(
	ctx context.Context,
	ag *agent.Agent,
	holder *ConversationHolder,
	toolName string,
	argumentsInJSON string,
	record ExecutionRecorder,
	chunk func(toolName, toolCallID, chunk string),
) (string, error) {
	var args map[string]interface{}
	if argumentsInJSON != "" && argumentsInJSON != "null" {
		if err := json.Unmarshal([]byte(argumentsInJSON), &args); err != nil {
			// Return soft error (nil error) so the eino graph continues and the LLM can self-correct,
			// instead of a hard error that terminates the iteration loop.
			return ToolErrorPrefix + fmt.Sprintf(
				"Invalid tool arguments JSON: %s\n\nPlease ensure the arguments are a valid JSON object "+
					"(double-quoted keys, matched braces, no trailing commas) and retry.\n\n"+
					" (Failed to parse tool-argument JSON: %s. Ensure that arguments is a valid JSON object and try again.)",
				err.Error(), err.Error()), nil
		}
	}
	if args == nil {
		args = map[string]interface{}{}
	}

	if chunk != nil {
		toolCallID := compose.GetToolCallID(ctx)
		if toolCallID != "" {
			if existing, ok := ctx.Value(security.ToolOutputCallbackCtxKey).(security.ToolOutputCallback); ok && existing != nil {
				ctx = context.WithValue(ctx, security.ToolOutputCallbackCtxKey, security.ToolOutputCallback(func(c string) {
					existing(c)
					if strings.TrimSpace(c) == "" {
						return
					}
					chunk(toolName, toolCallID, c)
				}))
			} else {
				ctx = context.WithValue(ctx, security.ToolOutputCallbackCtxKey, security.ToolOutputCallback(func(c string) {
					if strings.TrimSpace(c) == "" {
						return
					}
					chunk(toolName, toolCallID, c)
				}))
			}
		}
	}

	res, err := ag.ExecuteMCPToolForConversation(ctx, holder.Get(), toolName, args)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}
	if res.ExecutionID != "" && record != nil {
		record(res.ExecutionID, compose.GetToolCallID(ctx))
	}
	if res.IsError {
		return ToolErrorPrefix + res.Result, nil
	}
	return res.Result, nil
}

// UnknownToolReminderHandler is used by compose.ToolsNodeConfig.UnknownToolsHandler.
// When the model requests an unregistered tool name, it returns a soft-error tool result (nil error),
// allowing the model to self-correct in the same turn without triggering a full rerun at the run-loop level.
// It does not guess or map names, preventing accidental execution.
func UnknownToolReminderHandler() func(ctx context.Context, name, input string) (string, error) {
	return func(ctx context.Context, name, input string) (string, error) {
		_ = ctx
		_ = input
		requested := strings.TrimSpace(name)
		// Return a soft tool-result error so the graph keeps running and the LLM
		// can correct tool name/arguments within the same run.
		return ToolErrorPrefix + unknownToolReminderText(requested), nil
	}
}

func unknownToolReminderText(requested string) string {
	if requested == "" {
		requested = "(empty)"
	}
	return fmt.Sprintf(`The tool name %q is not registered for this agent.

Please retry using only names that appear in the tool definitions for this turn (exact match, case-sensitive). Do not invent or rename tools; adjust your plan and continue.

(Tool %q is not registered. Use only tool names provided in the current-turn context, matching them exactly; do not rewrite or guess names, and continue with the subsequent steps.)`, requested, requested)
}
