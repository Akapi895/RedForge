package projectprompt

// ShellExecExecuteGuidanceSection returns the shared shell-tool guidance.
func ShellExecExecuteGuidanceSection() string {
	return `Shell (exec/execute): prefer a dedicated MCP tool when one exists. Use exec for system commands, pipelines, working directories, and background processes. Use execute for scripts under skills/ together with read_file and skill. Split multi-step scans into separate calls; do not chain several scanners in one shell command. Write long scripts, request bodies, and payloads to the session workspace with write_file before executing a short command. Do not embed long content directly in command arguments. Store downloads and temporary files in the session workspace specified by the system prompt; do not use /tmp.`
}

// ShellExecExecuteGuidanceReconSuffix is optional guidance for reconnaissance
// sub-agents.
func ShellExecExecuteGuidanceReconSuffix() string {
	return `For enumeration, prefer dedicated MCP capabilities such as subfinder or amass; do not build long scanner chains with exec or execute.`
}
