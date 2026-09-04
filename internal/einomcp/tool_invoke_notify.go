package einomcp

import "sync"

// ToolInvokeNotifyHolder is shared by the Eino run loop and MCP/execute bridge; Fire is triggered when the tool returns its raw result.
// The UI tool_result must wait for the ADK schema.Tool event (the body after reduction) and is not pushed by this holder's callback.
type ToolInvokeNotifyHolder struct {
	mu sync.RWMutex
	fn func(toolCallID, toolName, einoAgent string, success bool, content string, invokeErr error)
}

// NewToolInvokeNotifyHolder creates a holder shared between ToolsFromDefinitions and the run loop.
func NewToolInvokeNotifyHolder() *ToolInvokeNotifyHolder {
	return &ToolInvokeNotifyHolder{}
}

// Set is called by runEinoADKAgentLoop before it begins consuming iter; it may be overwritten multiple times (normally only once).
func (h *ToolInvokeNotifyHolder) Set(fn func(toolCallID, toolName, einoAgent string, success bool, content string, invokeErr error)) {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fn = fn
}

// Fire is called by mcpBridgeTool when a tool invocation returns; it is ignored if Set has not been called or toolCallID is empty.
func (h *ToolInvokeNotifyHolder) Fire(toolCallID, toolName, einoAgent string, success bool, content string, invokeErr error) {
	if h == nil {
		return
	}
	h.mu.RLock()
	fn := h.fn
	h.mu.RUnlock()
	if fn == nil {
		return
	}
	fn(toolCallID, toolName, einoAgent, success, content, invokeErr)
}
