package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"cyberstrike-ai/internal/c2"
	"cyberstrike-ai/internal/config"
	"cyberstrike-ai/internal/mcp"
	"cyberstrike-ai/internal/mcp/builtin"
	"cyberstrike-ai/internal/openai"

	"go.uber.org/zap"
)

// Agent is an AI agent.
type Agent struct {
	openAIClient        *openai.Client
	config              *config.OpenAIConfig
	agentConfig         *config.AgentConfig
	mcpServer           *mcp.Server
	externalMCPMgr      *mcp.ExternalMCPManager // External MCP manager
	logger              *zap.Logger
	maxIterations       int
	mu                  sync.RWMutex      // Mutex supporting concurrent updates
	toolNameMapping     map[string]string // Tool-name mapping: OpenAI format -> original format (for external MCP tools)
	promptBaseDir       string            // Base directory for resolving relative system_prompt_path values (normally the config.yaml directory)
	toolDescriptionMode string            // Tool-description mode: "short" | "full"; default: short
}

type agentConversationIDKey struct{}

func withAgentConversationID(ctx context.Context, id string) context.Context {
	id = strings.TrimSpace(id)
	if id == "" || ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, agentConversationIDKey{}, id)
}

func agentConversationIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v, _ := ctx.Value(agentConversationIDKey{}).(string)
	return v
}

// ConversationIDFromContext returns the conversation ID injected into the current Agent request context (for example, for C2 MCP queuing and HITL gating).
func ConversationIDFromContext(ctx context.Context) string {
	return agentConversationIDFromContext(ctx)
}

// NewAgent creates a new Agent.
func NewAgent(cfg *config.OpenAIConfig, agentCfg *config.AgentConfig, mcpServer *mcp.Server, externalMCPMgr *mcp.ExternalMCPManager, logger *zap.Logger, maxIterations int) *Agent {
	// Use the default value of 30 when maxIterations is zero or negative
	if maxIterations <= 0 {
		maxIterations = 30
	}

	// Configure the HTTP transport to optimize connection management and timeouts
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   300 * time.Second,
			KeepAlive: 300 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 60 * time.Minute, // Response-header timeout increased for large responses
		DisableKeepAlives:     false,            // Enable connection reuse
	}

	// Increase the timeout to 30 minutes for long-running AI inference,
	// especially with streaming responses or complex tasks
	httpClient := &http.Client{
		Timeout:   30 * time.Minute, // Increased from 5 to 30 minutes
		Transport: transport,
	}
	llmClient := openai.NewClient(cfg, httpClient, logger)

	return &Agent{
		openAIClient:        llmClient,
		config:              cfg,
		agentConfig:         agentCfg,
		mcpServer:           mcpServer,
		externalMCPMgr:      externalMCPMgr,
		logger:              logger,
		maxIterations:       maxIterations,
		toolNameMapping:     make(map[string]string), // Initialize the tool-name mapping
		toolDescriptionMode: "short",
	}
}

// SetPromptBaseDir sets the base directory for a single agent's relative system_prompt_path (normally the config.yaml directory).
func (a *Agent) SetPromptBaseDir(dir string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.promptBaseDir = strings.TrimSpace(dir)
}

// ChatMessage represents a chat message.
type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	// ToolName applies only to the tool role and is restored from name or tool_name in Eino/trajectory JSON to construct ToolMessage during resume.
	ToolName string `json:"tool_name,omitempty"`
	// ReasoningContent corresponds to OpenAI/DeepSeek reasoning_content and must be sent back when resuming reasoning mode after tool calls (see DeepSeek documentation).
	ReasoningContent string `json:"reasoning_content,omitempty"`
	// ModelFacingTrace is runtime-only metadata: true means Content was already the exact
	// payload seen at the model boundary and must be restored byte-for-byte.
	ModelFacingTrace bool `json:"-"`
}

// MarshalJSON customizes JSON serialization by converting arguments in tool_calls to JSON strings.
func (cm ChatMessage) MarshalJSON() ([]byte, error) {
	// Build the serialized structure
	aux := map[string]interface{}{
		"role": cm.Role,
	}

	// Add content when present
	if cm.Content != "" {
		aux["content"] = cm.Content
	}
	if cm.ReasoningContent != "" {
		aux["reasoning_content"] = cm.ReasoningContent
	}

	// Add tool_call_id when present
	if cm.ToolCallID != "" {
		aux["tool_call_id"] = cm.ToolCallID
	}
	if cm.ToolName != "" {
		aux["tool_name"] = cm.ToolName
	}

	// Convert tool_calls arguments to JSON strings
	if len(cm.ToolCalls) > 0 {
		toolCallsJSON := make([]map[string]interface{}, len(cm.ToolCalls))
		for i, tc := range cm.ToolCalls {
			// Convert arguments to a JSON string
			argsJSON := ""
			if tc.Function.Arguments != nil {
				argsBytes, err := json.Marshal(tc.Function.Arguments)
				if err != nil {
					return nil, err
				}
				argsJSON = string(argsBytes)
			}

			toolCallsJSON[i] = map[string]interface{}{
				"id":   tc.ID,
				"type": tc.Type,
				"function": map[string]interface{}{
					"name":      tc.Function.Name,
					"arguments": argsJSON,
				},
			}
		}
		aux["tool_calls"] = toolCallsJSON
	}

	return json.Marshal(aux)
}

// OpenAIRequest represents an OpenAI API request.
type OpenAIRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Tools    []Tool        `json:"tools,omitempty"`
	Stream   bool          `json:"stream,omitempty"`
}

// OpenAIResponse represents an OpenAI API response.
type OpenAIResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Error   *Error   `json:"error,omitempty"`
}

// Choice represents a response choice.
type Choice struct {
	Message      MessageWithTools `json:"message"`
	FinishReason string           `json:"finish_reason"`
}

// MessageWithTools represents a message containing tool calls.
type MessageWithTools struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Tool represents an OpenAI tool definition.
type Tool struct {
	Type     string             `json:"type"`
	Function FunctionDefinition `json:"function"`
}

// FunctionDefinition represents a function definition.
type FunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Error represents an OpenAI error.
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// ToolCall represents a tool call.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call.
type FunctionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// UnmarshalJSON customizes JSON parsing to handle arguments represented as either a string or an object.
func (fc *FunctionCall) UnmarshalJSON(data []byte) error {
	type Alias FunctionCall
	aux := &struct {
		Name      string      `json:"name"`
		Arguments interface{} `json:"arguments"`
		*Alias
	}{
		Alias: (*Alias)(fc),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	fc.Name = aux.Name

	// Handle arguments represented as either a string or an object
	switch v := aux.Arguments.(type) {
	case map[string]interface{}:
		fc.Arguments = v
	case string:
		// If it is a string, try parsing it as JSON
		if err := json.Unmarshal([]byte(v), &fc.Arguments); err != nil {
			// If parsing fails, create a map containing the original string
			fc.Arguments = map[string]interface{}{
				"raw": v,
			}
		}
	case nil:
		fc.Arguments = make(map[string]interface{})
	default:
		// For other types, try converting to a map
		fc.Arguments = map[string]interface{}{
			"value": v,
		}
	}

	return nil
}

// ProgressCallback is the progress callback function type.
type ProgressCallback func(eventType, message string, data interface{})

// EinoSingleAgentSystemInstruction is used by Eino adk.ChatModelAgent.Instruction and includes system_prompt_path.
func (a *Agent) EinoSingleAgentSystemInstruction() string {
	systemPrompt := DefaultSingleAgentSystemPrompt()
	if a.agentConfig != nil {
		if p := strings.TrimSpace(a.agentConfig.SystemPromptPath); p != "" {
			path := p
			a.mu.RLock()
			base := a.promptBaseDir
			a.mu.RUnlock()
			if !filepath.IsAbs(path) && base != "" {
				path = filepath.Join(base, path)
			}
			if b, err := os.ReadFile(path); err != nil {
				a.logger.Warn("Failed to read single-agent system_prompt_path; using the built-in prompt", zap.String("path", path), zap.Error(err))
			} else if s := strings.TrimSpace(string(b)); s != "" {
				systemPrompt = s
			}
		}
	}
	return systemPrompt
}

// getAvailableTools obtains available tools.
// It retrieves the tool list dynamically from the MCP server; tool_description_mode controls the description mode.
// roleTools is the role-configured list in toolKey format; when empty or nil, all tools are used (default role).
func (a *Agent) getAvailableTools(roleTools []string) []Tool {
	// Build the role's tool set for fast lookup
	roleToolSet := make(map[string]bool)
	if len(roleTools) > 0 {
		for _, toolKey := range roleTools {
			roleToolSet[toolKey] = true
		}
	}

	// Obtain all registered internal tools from the MCP server
	mcpTools := a.mcpServer.GetAllTools()

	// Convert to OpenAI-format tool definitions
	tools := make([]Tool, 0, len(mcpTools))
	for _, mcpTool := range mcpTools {
		// If a role tool list is specified, add only listed tools
		if len(roleToolSet) > 0 {
			toolKey := mcpTool.Name // Built-in tools use the tool name as the key
			if !roleToolSet[toolKey] {
				continue // Skip tools not present in the role tool list
			}
		}
		description := a.pickToolDescription(mcpTool.ShortDescription, mcpTool.Description)

		// Convert schema types to OpenAI-standard types
		convertedSchema := a.convertSchemaTypes(mcpTool.InputSchema)

		tools = append(tools, Tool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        mcpTool.Name,
				Description: description, // Use a short description to reduce token consumption
				Parameters:  convertedSchema,
			},
		})
	}

	// Obtain external MCP tools
	if a.externalMCPMgr != nil {
		// Increase the timeout to 30 seconds because connecting to a remote server through a proxy may take longer
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		externalTools, err := a.externalMCPMgr.GetAllTools(ctx)
		extMap := make(map[string]string)
		if err != nil {
			a.logger.Warn("Failed to obtain external MCP tools", zap.Error(err))
		} else {
			// Obtain external MCP configuration to check tool enablement
			externalMCPConfigs := a.externalMCPMgr.GetConfigs()

			// Add enabled external MCP tools to the tool list
			for _, externalTool := range externalTools {
				// External tools use "mcpName::toolName" as toolKey
				externalToolKey := externalTool.Name

				// If a role tool list is specified, add only listed tools
				if len(roleToolSet) > 0 {
					if !roleToolSet[externalToolKey] {
						continue // Skip tools not present in the role tool list
					}
				}

				// Parse the tool name: mcpName::toolName
				var mcpName, actualToolName string
				if idx := strings.Index(externalTool.Name, "::"); idx > 0 {
					mcpName = externalTool.Name[:idx]
					actualToolName = externalTool.Name[idx+2:]
				} else {
					continue // Skip incorrectly formatted tools
				}

				// Check whether the tool is enabled
				enabled := false
				if cfg, exists := externalMCPConfigs[mcpName]; exists {
					// First check whether the external MCP is enabled
					if !cfg.ExternalMCPEnable {
						enabled = false // All tools are disabled when the MCP is disabled
					} else {
						// The MCP is enabled; check the individual tool's status
						// Default to enabled when ToolEnabled is empty or does not contain the tool (backward compatibility)
						if cfg.ToolEnabled == nil {
							enabled = true // No tool status configured; default to enabled
						} else if toolEnabled, exists := cfg.ToolEnabled[actualToolName]; exists {
							enabled = toolEnabled // Use the configured tool status
						} else {
							enabled = true // Tool absent from configuration; default to enabled
						}
					}
				}

				// Add enabled tools only
				if !enabled {
					continue
				}

				description := a.pickToolDescription(externalTool.ShortDescription, externalTool.Description)

				// Convert schema types to OpenAI-standard types
				convertedSchema := a.convertSchemaTypes(externalTool.InputSchema)

				// Replace "::" with "__" in tool names to satisfy OpenAI naming rules
				// OpenAI tool names may contain only [a-zA-Z0-9_-]
				openAIName := strings.ReplaceAll(externalTool.Name, "::", "__")

				// Save the name mapping (OpenAI format -> original format)
				extMap[openAIName] = externalTool.Name

				tools = append(tools, Tool{
					Type: "function",
					Function: FunctionDefinition{
						Name:        openAIName, // Use an OpenAI-compliant name
						Description: description,
						Parameters:  convertedSchema,
					},
				})
			}
		}
		a.mu.Lock()
		a.toolNameMapping = extMap
		a.mu.Unlock()
	}

	a.logger.Debug("Obtained available tool list",
		zap.Int("internalTools", len(mcpTools)),
		zap.Int("totalTools", len(tools)),
	)

	return tools
}

func (a *Agent) pickToolDescription(shortDesc, fullDesc string) string {
	a.mu.RLock()
	mode := strings.TrimSpace(strings.ToLower(a.toolDescriptionMode))
	a.mu.RUnlock()
	if mode == "full" {
		return fullDesc
	}
	if shortDesc != "" {
		return shortDesc
	}
	return fullDesc
}

// convertSchemaTypes recursively converts schema types to OpenAI-standard types.
func (a *Agent) convertSchemaTypes(schema map[string]interface{}) map[string]interface{} {
	if schema == nil {
		return schema
	}

	// Create a new schema copy
	converted := make(map[string]interface{})
	for k, v := range schema {
		converted[k] = v
	}

	// Convert types in properties
	if properties, ok := converted["properties"].(map[string]interface{}); ok {
		convertedProperties := make(map[string]interface{})
		for propName, propValue := range properties {
			if prop, ok := propValue.(map[string]interface{}); ok {
				convertedProp := make(map[string]interface{})
				for pk, pv := range prop {
					if pk == "type" {
						// Convert the type
						if typeStr, ok := pv.(string); ok {
							convertedProp[pk] = a.convertToOpenAIType(typeStr)
						} else {
							convertedProp[pk] = pv
						}
					} else {
						convertedProp[pk] = pv
					}
				}
				convertedProperties[propName] = convertedProp
			} else {
				convertedProperties[propName] = propValue
			}
		}
		converted["properties"] = convertedProperties
	}

	return converted
}

// convertToOpenAIType converts configuration types to OpenAI/JSON Schema standard types.
func (a *Agent) convertToOpenAIType(configType string) string {
	switch configType {
	case "bool":
		return "boolean"
	case "int", "integer":
		return "number"
	case "float", "double":
		return "number"
	case "string", "array", "object":
		return configType
	default:
		// Return the original type by default
		return configType
	}
}

// ToolExecutionResult is an MCP tool-execution result used by the Eino bridge and monitoring persistence.
type ToolExecutionResult struct {
	Result      string
	ExecutionID string
	IsError     bool
}

func buildToolFailureMessage(toolName, detail string, err error) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Tool invocation failed\n\n")
	fmt.Fprintf(&b, "Tool name: %s\n", toolName)
	fmt.Fprintf(&b, "Error details: %s", detail)
	return strings.TrimRight(b.String(), "\n")
}

// executeToolViaMCP executes a tool through MCP.
// Even when execution fails, it returns a result rather than an error so the AI can handle the failure.
func (a *Agent) executeToolViaMCP(ctx context.Context, toolName string, args map[string]interface{}) (*ToolExecutionResult, error) {
	a.logger.Info("Executing tool through MCP",
		zap.String("tool", toolName),
		zap.Any("args", args),
	)

	// Automatically add conversation_id for the record_vulnerability tool
	if toolName == builtin.ToolRecordVulnerability {
		conversationID := agentConversationIDFromContext(ctx)
		if conversationID != "" {
			args["conversation_id"] = conversationID
			a.logger.Debug("Automatically added conversation_id to record_vulnerability tool",
				zap.String("conversation_id", conversationID),
			)
		} else {
			a.logger.Warn("conversation_id is empty during record_vulnerability tool invocation")
		}
	}

	var result *mcp.ToolResult
	var executionID string
	var err error

	// Per-tool execution timeout prevents an individual tool from hanging for a long time (for example, still running after 30 minutes)
	toolCtx := ctx
	var toolCancel context.CancelFunc
	if a.agentConfig != nil && a.agentConfig.ToolTimeoutMinutes > 0 {
		toolCtx, toolCancel = context.WithTimeout(ctx, time.Duration(a.agentConfig.ToolTimeoutMinutes)*time.Minute)
		defer func() {
			if toolCancel != nil {
				toolCancel()
			}
		}()
	}
	// Asynchronous HITL waiting for dangerous C2 tasks must bind to the entire Agent-run ctx, not the per-tool child ctx (which is canceled on return)
	toolCtx = c2.WithHITLRunContext(toolCtx, ctx)

	// Check whether this is an external MCP tool using the tool-name mapping
	a.mu.RLock()
	originalToolName, isExternalTool := a.toolNameMapping[toolName]
	a.mu.RUnlock()

	if isExternalTool && a.externalMCPMgr != nil {
		// Invoke the external MCP tool using its original name
		a.logger.Debug("Invoking external MCP tool",
			zap.String("openAIName", toolName),
			zap.String("originalName", originalToolName),
		)
		result, executionID, err = a.externalMCPMgr.CallTool(toolCtx, originalToolName, args)
	} else {
		// Invoke an internal MCP tool
		result, executionID, err = a.mcpServer.CallTool(toolCtx, toolName, args)
	}

	// If invocation fails (for example, missing tool or timeout), return a friendly error message instead of raising an exception
	if err != nil {
		detail := err.Error()
		timeoutMinutes := 10
		if a.agentConfig != nil && a.agentConfig.ToolTimeoutMinutes > 0 {
			timeoutMinutes = a.agentConfig.ToolTimeoutMinutes
		}
		if errors.Is(err, context.Canceled) {
			detail = "The tool invocation was terminated manually from the MCP monitoring page. The agent will carry this result into subsequent steps; the overall task is not stopped."
		} else if errors.Is(err, context.DeadlineExceeded) {
			detail = fmt.Sprintf("Tool execution was terminated automatically after exceeding %d minutes (adjust agent.tool_timeout_minutes in config.yaml)", timeoutMinutes)
		}
		errorMsg := buildToolFailureMessage(toolName, detail, err)

		return &ToolExecutionResult{
			Result:      errorMsg,
			ExecutionID: executionID,
			IsError:     true,
		}, nil // Return a nil error so the caller can process the result
	}

	// Format the result
	var resultText strings.Builder
	for _, content := range result.Content {
		resultText.WriteString(content.Text)
		resultText.WriteString("\n")
	}

	resultStr := resultText.String()

	return &ToolExecutionResult{
		Result:      resultStr,
		ExecutionID: executionID,
		IsError:     result != nil && result.IsError,
	}, nil
}

// UpdateConfig updates the OpenAI configuration.
func (a *Agent) UpdateConfig(cfg *config.OpenAIConfig) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.config = cfg

	a.logger.Info("Agent configuration updated",
		zap.String("base_url", cfg.BaseURL),
		zap.String("model", cfg.Model),
	)
}

// UpdateMaxIterations updates the maximum iteration count.
func (a *Agent) UpdateMaxIterations(maxIterations int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if maxIterations > 0 {
		a.maxIterations = maxIterations
		a.logger.Info("Agent maximum iteration count updated", zap.Int("max_iterations", maxIterations))
	}
}

// UpdateToolDescriptionMode updates the tool-description mode (short/full).
func (a *Agent) UpdateToolDescriptionMode(mode string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	mode = strings.TrimSpace(strings.ToLower(mode))
	if mode != "full" {
		mode = "short"
	}
	a.toolDescriptionMode = mode
	a.logger.Debug("Agent tool-description mode updated", zap.String("tool_description_mode", mode))
}

// RepairOrphanToolMessages removes unpaired tool messages and incomplete tool_calls to prevent OpenAI errors.
// It also ensures that tool_calls in historical messages serve only as contextual memory and do not trigger re-execution.
// This exported method can be called when restoring message history.
func (a *Agent) RepairOrphanToolMessages(messages *[]ChatMessage) bool {
	return a.repairOrphanToolMessages(messages)
}

// repairOrphanToolMessages removes unpaired tool messages and incomplete tool_calls to prevent OpenAI errors.
// It also ensures that tool_calls in historical messages serve only as contextual memory and do not trigger re-execution.
func (a *Agent) repairOrphanToolMessages(messages *[]ChatMessage) bool {
	if messages == nil {
		return false
	}

	msgs := *messages
	if len(msgs) == 0 {
		return false
	}

	pending := make(map[string]int)
	cleaned := make([]ChatMessage, 0, len(msgs))
	removed := false

	for _, msg := range msgs {
		switch strings.ToLower(msg.Role) {
		case "assistant":
			if len(msg.ToolCalls) > 0 {
				// Record all tool_call IDs
				for _, tc := range msg.ToolCalls {
					if tc.ID != "" {
						pending[tc.ID]++
					}
				}
			}
			cleaned = append(cleaned, msg)
		case "tool":
			callID := msg.ToolCallID
			if callID == "" {
				removed = true
				continue
			}
			if count, exists := pending[callID]; exists && count > 0 {
				if count == 1 {
					delete(pending, callID)
				} else {
					pending[callID] = count - 1
				}
				cleaned = append(cleaned, msg)
			} else {
				removed = true
				continue
			}
		default:
			cleaned = append(cleaned, msg)
		}
	}

	// If unmatched tool_calls remain (an assistant message has tool_calls without corresponding tool responses),
	// remove them from the final assistant message so the AI does not execute them again
	if len(pending) > 0 {
		// Search backward for the final assistant message
		for i := len(cleaned) - 1; i >= 0; i-- {
			if strings.ToLower(cleaned[i].Role) == "assistant" && len(cleaned[i].ToolCalls) > 0 {
				// Remove unmatched tool_calls
				originalCount := len(cleaned[i].ToolCalls)
				validToolCalls := make([]ToolCall, 0)
				for _, tc := range cleaned[i].ToolCalls {
					if tc.ID != "" && pending[tc.ID] > 0 {
						// Remove this tool_call because it has no corresponding tool response
						removed = true
						delete(pending, tc.ID)
					} else {
						validToolCalls = append(validToolCalls, tc)
					}
				}
				// Update the message's ToolCalls
				if len(validToolCalls) != originalCount {
					cleaned[i].ToolCalls = validToolCalls
					a.logger.Info("Removed incomplete tool_calls to prevent re-execution",
						zap.Int("removed_count", originalCount-len(validToolCalls)),
					)
				}
				break
			}
		}
	}

	if removed {
		a.logger.Warn("Repaired tool messages and tool_calls in conversation history",
			zap.Int("original_messages", len(msgs)),
			zap.Int("cleaned_messages", len(cleaned)),
		)
		*messages = cleaned
	}

	return removed
}

// ToolsForRole returns OpenAI-function tool definitions consistent with the single-Agent loop for orchestration layers such as Eino DeepAgent to bind MCP tools.
func (a *Agent) ToolsForRole(roleTools []string) []Tool {
	return a.getAvailableTools(roleTools)
}

// ExecuteMCPToolForConversation executes an MCP tool in a specified conversation context, matching main-Agent-loop behavior such as automatic conversation_id injection.
func (a *Agent) ExecuteMCPToolForConversation(ctx context.Context, conversationID, toolName string, args map[string]interface{}) (*ToolExecutionResult, error) {
	ctx = withAgentConversationID(ctx, conversationID)
	ctx = mcp.WithMCPConversationID(ctx, conversationID)
	return a.executeToolViaMCP(ctx, toolName, args)
}

// BeginLocalToolExecution writes running status when a tool outside the CallTool path starts, allowing the MCP monitoring page to display it as running.
func (a *Agent) BeginLocalToolExecution(ctx context.Context, toolName string, args map[string]interface{}) string {
	if a == nil || a.mcpServer == nil {
		return ""
	}
	return a.mcpServer.BeginToolExecution(ctx, toolName, args)
}

// FinishLocalToolExecution completes a record created by BeginLocalToolExecution; when executionID is empty, it writes a completed record in one operation.
func (a *Agent) FinishLocalToolExecution(ctx context.Context, executionID, toolName string, args map[string]interface{}, resultText string, invokeErr error) string {
	if a == nil || a.mcpServer == nil {
		return ""
	}
	return a.mcpServer.FinishToolExecution(ctx, executionID, toolName, args, resultText, invokeErr)
}

// AppendLocalToolExecutionPartialOutput records a bounded live-output preview for a running local tool.
func (a *Agent) AppendLocalToolExecutionPartialOutput(executionID, chunk string) {
	if a == nil || a.mcpServer == nil {
		return
	}
	a.mcpServer.AppendToolExecutionPartialOutput(executionID, chunk)
}

func (a *Agent) RegisterLocalToolExecutionCancel(executionID string, cancel context.CancelFunc) {
	if a == nil || a.mcpServer == nil {
		return
	}
	a.mcpServer.RegisterToolExecutionCancel(executionID, cancel)
}

func (a *Agent) UnregisterLocalToolExecutionCancel(executionID string) {
	if a == nil || a.mcpServer == nil {
		return
	}
	a.mcpServer.UnregisterToolExecutionCancel(executionID)
}

// RecordLocalToolExecution writes a completed tool invocation outside the CallTool path to the MCP monitoring store (consistent with CallTool persistence) and returns executionId.
// It is used for cases such as Eino filesystem execute so the assistant bubble's penetration-testing details link to monitoring consistently with ordinary MCP calls.
func (a *Agent) RecordLocalToolExecution(ctx context.Context, toolName string, args map[string]interface{}, resultText string, invokeErr error) string {
	return a.FinishLocalToolExecution(ctx, "", toolName, args, resultText, invokeErr)
}

// UpdateMCPExecutionDisplayResult updates a tool result in the monitoring store to the displayed body sent to the model after reduction.
func (a *Agent) UpdateMCPExecutionDisplayResult(executionID, resultText string) {
	if a == nil || strings.TrimSpace(executionID) == "" {
		return
	}
	text := resultText
	if strings.TrimSpace(text) == "" {
		text = "(No output)"
	}
	tr := &mcp.ToolResult{
		Content: []mcp.Content{{Type: "text", Text: text}},
	}
	if a.mcpServer != nil {
		_ = a.mcpServer.UpdateToolExecutionResult(executionID, tr)
	}
}

// MCPExecutionResultText returns the monitor-facing result text after storage
// guards such as large-output spilling have been applied.
func (a *Agent) MCPExecutionResultText(executionID string) string {
	if a == nil || a.mcpServer == nil || strings.TrimSpace(executionID) == "" {
		return ""
	}
	exec, ok := a.mcpServer.GetExecution(executionID)
	if !ok || exec == nil || exec.Result == nil {
		return ""
	}
	return mcp.ToolResultPlainText(exec.Result)
}

// CancelMCPToolExecutionWithNote cancels an in-progress MCP tool (internal first, then external), matching "Terminate Tool" on the monitoring page; a non-empty note is merged into the text returned to the model.
func (a *Agent) CancelMCPToolExecutionWithNote(executionID, note string) bool {
	executionID = strings.TrimSpace(executionID)
	note = strings.TrimSpace(note)
	if executionID == "" {
		return false
	}
	if a.mcpServer != nil && a.mcpServer.CancelToolExecutionWithNote(executionID, note) {
		return true
	}
	if a.externalMCPMgr != nil && a.externalMCPMgr.CancelToolExecutionWithNote(executionID, note) {
		return true
	}
	return false
}

// CancelRunningMCPToolsForConversation cancels all currently running internal/external MCP executions
// owned by the conversation. It is used when a session ends or the user stops a task.
func (a *Agent) CancelRunningMCPToolsForConversation(conversationID, note string) int {
	conversationID = strings.TrimSpace(conversationID)
	if a == nil || conversationID == "" {
		return 0
	}
	note = strings.TrimSpace(note)
	seen := make(map[string]struct{})
	cancelled := 0
	cancelIfConversationMatches := func(execID string, get func(string) (*mcp.ToolExecution, bool), cancel func(string, string) bool) {
		execID = strings.TrimSpace(execID)
		if execID == "" {
			return
		}
		if _, ok := seen[execID]; ok {
			return
		}
		seen[execID] = struct{}{}
		exec, ok := get(execID)
		if !ok || exec == nil || strings.TrimSpace(exec.ConversationID) != conversationID {
			return
		}
		if cancel(execID, note) {
			cancelled++
		}
	}
	if a.mcpServer != nil {
		for execID := range a.mcpServer.ActiveRunningExecutionIDs() {
			cancelIfConversationMatches(execID, a.mcpServer.GetExecution, a.mcpServer.CancelToolExecutionWithNote)
		}
	}
	if a.externalMCPMgr != nil {
		for execID := range a.externalMCPMgr.ActiveRunningExecutionIDs() {
			cancelIfConversationMatches(execID, a.externalMCPMgr.GetExecution, a.externalMCPMgr.CancelToolExecutionWithNote)
		}
	}
	return cancelled
}

// extractQuotedToolName attempts to extract a quoted tool name from an error message.
func extractQuotedToolName(errMsg string) string {
	start := strings.Index(errMsg, "\"")
	if start == -1 {
		return ""
	}
	rest := errMsg[start+1:]
	end := strings.Index(rest, "\"")
	if end == -1 {
		return ""
	}
	return rest[:end]
}
