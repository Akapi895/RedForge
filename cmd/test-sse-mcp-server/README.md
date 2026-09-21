# SSE MCP Test Server

This is a test server for validating external MCP functionality in SSE mode.

## Usage

### 1. Start the test server

```bash
cd cmd/test-sse-mcp-server
go run main.go
```

The server starts at `http://127.0.0.1:8082` and provides these endpoints:
- `GET /sse` - SSE event stream endpoint
- `POST /message` - Message receiver endpoint

### 2. Add the configuration in CyberStrikeAI

Add an external MCP configuration in the web interface using this JSON:

```json
{
  "test-sse-mcp": {
    "transport": "sse",
    "url": "http://127.0.0.1:8082/sse",
    "description": "SSE MCP test server",
    "timeout": 30
  }
}
```

### 3. Test the functionality

The test server provides two test tools:

1. **test_echo** - Echoes input text
   - Parameter: `text` (string) - Text to echo

2. **test_add** - Adds two numbers
   - Parameter: `a` (number) - First number
   - Parameter: `b` (number) - Second number

## How It Works

1. The client establishes an SSE connection through `GET /sse` and receives server-sent events.
2. The client sends MCP protocol messages through `POST /message`.
3. After processing a message, the server pushes the response over the SSE connection.

## Logs

The server outputs these logs:
- SSE client connections and disconnections
- Received requests (method name and ID)
- Tool call details
