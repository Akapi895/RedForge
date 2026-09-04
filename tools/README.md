# Tool Configuration File Guide

## Overview

Each tool has an independent configuration file stored in the `tools/` directory. This approach makes tool configuration clearer and easier to maintain and manage. The system automatically loads all `.yaml` and `.yml` files under `tools/`.

## Configuration File Format

Each tool configuration file is a YAML file. The following table lists the currently supported top-level fields and whether they are required. Check each field before submitting a configuration:

| Field | Required | Type | Description |
|------|------|------|------|
| `name` | Yes | string | Unique tool identifier; a combination of lowercase letters, numbers, and hyphens is recommended. |
| `command` | Yes | string | Command or script name to execute; it must be available on the system PATH or specified as an absolute path. |
| `enabled` | Yes | bool | Whether to register the tool with MCP. When set to `false`, the tool is ignored. |
| `description` | Yes | string | Detailed description supporting multi-line Markdown, used for deep AI understanding and `resources/read` queries. |
| `short_description` | Optional | string | A 20-50 character summary used in tool lists to reduce token consumption. When omitted, the beginning of `description` is extracted automatically. |
| `args` | Optional | string[] | Fixed arguments prepended to the command line in order, often used to define a default scanning mode. |
| `parameters` | Optional | array | List of parameters configurable at runtime; see the "Parameter Definitions" section. |
| `arg_mapping` | Optional | string | Argument-mapping mode (`auto`/`manual`/`template`); defaults to `auto`. Omit it unless special behavior is required. |

> If a field is invalid or a required field is missing, the system skips that tool during loading and writes a warning to the log without affecting other tools.

## Tool Descriptions

### Short Description (`short_description`)

- **Purpose**: Used in tool lists to reduce the number of tokens sent to the language model
- **Requirement**: One sentence of 20-50 characters describing the tool's primary purpose
- **Example**: `"Network scanning tool for discovering network hosts, open ports, and services"`

### Detailed Description (`description`)

Multi-line text is supported and should include:

1. **Tool functionality**: The tool's primary capabilities
2. **Use cases**: Situations in which the tool should be used
3. **Notes**: Precautions and warnings for use
4. **Examples**: Usage examples (optional)

**Important notes**:

- When the tool list is sent to the language model, `short_description` is used when present
- If `short_description` is absent, the system automatically extracts the first line or first 100 characters from `description`
- The detailed description can be obtained through MCP's `resources/read` endpoint (URI: `tool://tool_name`)

This can substantially reduce token consumption, especially when many tools are available (for example, 100 tools).

## Parameter Definitions

Each parameter may contain the following fields:

- `name`: Parameter name
- `type`: Parameter type (string, int, bool, array)
- `description`: Detailed parameter description (multi-line text is supported)
- `required`: Whether the parameter is required (true/false)
- `default`: Default value
- `flag`: Command-line flag (such as "-u", "--url", or "-p")
- `position`: Positional-argument position (an integer starting at 0)
- `format`: Parameter format ("flag", "positional", "combined", or "template")
- `template`: Template string (used when `format="template"`)
- `options`: List of allowed values (used for enumerated types)

### Parameter Formats

- **`flag`**: Flag argument in the form `--flag value` or `-f value`
  - Example: `flag: "-u"` produces `-u http://example.com`

- **`positional`**: Positional argument added to the command in sequence
  - Example: `position: 0` places it as the first positional argument

- **`combined`**: Combined form written as `--flag=value`
  - Example: `flag: "--level"`, `format: "combined"` produces `--level=3`

- **`template`**: Template form using a custom template string
  - Example: `template: "{flag} {value}"` produces a custom format

### Special Parameters

#### The `additional_args` Parameter

`additional_args` is a special parameter used to pass extra command-line options that are not defined in the parameter list. It is parsed and split on spaces into multiple arguments.

**Use cases:**

- Pass advanced tool options
- Pass arguments not defined in the configuration
- Pass complex combinations of arguments

**Example:**

```yaml
- name: "additional_args"
  type: "string"
  description: "Additional tool arguments, separated by spaces"
  required: false
  format: "positional"
```

**Usage examples:**

- `additional_args: "--script vuln -O"` is parsed as `["--script", "vuln", "-O"]`
- `additional_args: "-T4 --max-retries 3"` is parsed as `["-T4", "--max-retries", "3"]`

**Notes:**

- Arguments are split on spaces, while content inside quotation marks is preserved
- Ensure arguments are formatted correctly to avoid command-injection risks
- This parameter is appended to the end of the command

#### The `scan_type` Parameter (Specific Tools)

Some tools, such as `nmap`, support a `scan_type` parameter that overrides the default scan-type arguments.

**Example (nmap):**

```yaml
- name: "scan_type"
  type: "string"
  description: "Scan-type options that can override the default scan type"
  required: false
  format: "positional"
```

**Usage examples:**

- `scan_type: "-sV -sC"` performs version detection and script scanning
- `scan_type: "-A"` performs a comprehensive scan

**Notes:**

- When `scan_type` is specified, it replaces the default scan-type arguments in the tool configuration
- Separate multiple options with spaces

### Parameter Description Requirements

A parameter description should include:

1. **Purpose**: What the parameter does
2. **Format requirements**: Required value format (such as URL format or port-range format)
3. **Example values**: Specific examples, shown as a list when there are several
4. **Notes**: Relevant considerations such as permission requirements, performance impact, and security warnings

**Recommended description format:**

- Use Markdown to improve readability
- Use `**bold**` to emphasize important information
- Use lists for multiple examples or options
- Use code blocks for complex formats

**Example:**

```yaml
description: |
  Target IP address or domain name. It may be a single IP, IP range, CIDR, or domain name.

  **Example values:**
  - Single IP: "192.168.1.1"
  - IP range: "192.168.1.1-100"
  - CIDR: "192.168.1.0/24"
  - Domain name: "example.com"

  **Notes:**
  - Ensure that the target-address format is correct
  - This parameter is required and cannot be empty
```

## Parameter Types

### Boolean Type (bool)

Boolean parameters receive special handling:

- `true`: Add only the flag, without a value (for example, `--flag`)
- `false`: Do not add any argument
- Multiple input forms are supported: `true`/`false`, `1`/`0`, and `"true"`/`"false"`

**Example:**

```yaml
- name: "verbose"
  type: "bool"
  description: "Verbose output mode"
  required: false
  default: false
  flag: "-v"
  format: "flag"
```

### String Type (string)

The most commonly used parameter type; it supports any string value.

### Integer Type (int/integer)

Used for numeric parameters such as port numbers and levels.

**Example:**

```yaml
- name: "level"
  type: "int"
  description: "Test level, from 1 to 5"
  required: false
  default: 3
  flag: "--level"
  format: "combined"  # --level=3
```

### Array Type (array)

Arrays are automatically converted into comma-delimited strings.

**Example:**

```yaml
- name: "ports"
  type: "array"
  item_type: "number"
  description: "List of ports"
  required: false
  # Input: [80, 443, 8080]
  # Output: "80,443,8080"
```

## Examples

Refer to the existing tool configuration files under `tools/`:

- `nmap.yaml`: Network scanning tool (includes examples of `scan_type` and `additional_args`)
- `sqlmap.yaml`: SQL injection detection tool (includes an `additional_args` example)
- `nikto.yaml`: Web server scanning tool
- `dirb.yaml`: Web directory scanning tool
- `exec.yaml`: System command execution tool

### Complete Example: nmap Tool Configuration

```yaml
name: "nmap"
command: "nmap"
args: ["-sT", "-sV", "-sC"]  # Default scan type
enabled: true

short_description: "Network scanning tool for discovering network hosts, open ports, and services"

description: |
  Network mapping and port scanning tool for discovering hosts, services, and open ports on a network.

  **Main features:**
  - Host discovery: detect active hosts on the network
  - Port scanning: identify open ports on the target host
  - Service identification: detect the types and versions of services running on ports
  - Operating-system detection: identify the target host's operating-system type
  - Vulnerability detection: use NSE scripts to detect common vulnerabilities

parameters:
  - name: "target"
    type: "string"
    description: "Target IP address or domain name"
    required: true
    position: 0
    format: "positional"

  - name: "ports"
    type: "string"
    description: "Port range, for example: 1-1000"
    required: false
    flag: "-p"
    format: "flag"

  - name: "scan_type"
    type: "string"
    description: "Scan-type options, for example: '-sV -sC'"
    required: false
    format: "positional"

  - name: "additional_args"
    type: "string"
    description: "Additional Nmap arguments, for example: '--script vuln -O'"
    required: false
    format: "positional"
```

## Adding a New Tool

To add a new tool, create a YAML file such as `my_tool.yaml` under `tools/`:

```yaml
name: "my_tool"
command: "my-command"
args: ["--default-arg"]  # Fixed arguments (optional)
enabled: true

# Short description (recommended), used in the tool list to reduce token consumption
short_description: "One-sentence description of the tool's purpose"

# Detailed description used for tool documentation and AI understanding
description: |
  Detailed tool description supporting multi-line text and Markdown.

  **Main features:**
  - Feature 1
  - Feature 2

  **Use cases:**
  - Scenario 1
  - Scenario 2

  **Notes:**
  - Usage precautions
  - Permission requirements
  - Performance impact

parameters:
  - name: "target"
    type: "string"
    description: |
      Detailed description of the target parameter.

      **Example values:**
      - "value1"
      - "value2"

      **Notes:**
      - Format requirements
      - Usage restrictions
    required: true
    position: 0  # Positional argument
    format: "positional"

  - name: "option"
    type: "string"
    description: "Option parameter description"
    required: false
    flag: "--option"
    format: "flag"

  - name: "verbose"
    type: "bool"
    description: "Verbose output mode"
    required: false
    default: false
    flag: "-v"
    format: "flag"

  - name: "additional_args"
    type: "string"
    description: "Additional tool arguments, separated by spaces"
    required: false
    format: "positional"
```

After saving the file, restart the service to load the new tool automatically.

### Tool Configuration Best Practices

1. **Parameter design**
   - Define commonly used parameters individually so that the AI can understand and use them easily
   - Use `additional_args` for flexibility and advanced usage
   - Provide clear descriptions and examples for parameters

2. **Description optimization**
   - Use `short_description` to reduce token consumption
   - Make `description` detailed enough to help the AI understand the tool's purpose
   - Use Markdown to improve readability

3. **Default values**
   - Set reasonable defaults for commonly used parameters
   - Boolean defaults are usually `false`
   - Set numeric defaults according to the tool's characteristics

4. **Parameter validation**
   - Clearly state parameter format requirements in descriptions
   - Provide several example values
   - Explain parameter restrictions and precautions

5. **Security**
   - Add warnings to descriptions for dangerous operations
   - State permission requirements
   - Remind users to operate only in authorized environments

6. **Single-execution duration and timeout (best practice)**
   - If a tool frequently runs for a long time (for example, it still shows "Running" after 10-30 minutes), treat this as an abnormally long hang and consider the following:
     - Set the maximum duration of one tool execution in `agent.tool_timeout_minutes` in **config.yaml** (default: 10 minutes). After the timeout, the process is terminated automatically and its resources are released.
     - Increase the value appropriately for longer scans (for example, to 20 or 30). Setting it to 0 (unlimited) is not recommended.
     - On the task-monitoring page, use "Stop Task" for the entire task to interrupt the current conversation and subsequent tool calls.
     - Where possible, implement tools so that they can be interrupted or have a built-in timeout (for example, a timeout inside the script), allowing them to cooperate with the system timeout.

## Disabling a Tool

To disable a tool, set the `enabled` field in its configuration file to `false`, or delete/rename the configuration file.

After a tool is disabled, it does not appear in the tool list and cannot be called by the AI.

## Tool Configuration Validation

The system performs basic validation while loading tool configurations:

- Check required fields (`name`, `command`, and `enabled`)
- Validate the parameter-definition format
- Check whether parameter types are supported

If a configuration is invalid, the system displays a warning in the startup log but does not prevent the server from starting. The invalid tool configuration is skipped, while other tools continue to work normally.

## Frequently Asked Questions

### Q: How do I pass multiple parameter values?

A: For array parameters, the system automatically converts the value into a comma-delimited string. To pass several independent arguments, use `additional_args`.

### Q: How do I override a tool's default arguments?

A: Some tools, such as `nmap`, support `scan_type` for overriding the default scan type. In other cases, use `additional_args`.

### Q: What should I do if a tool has been showing "Running" for more than 30 minutes?

A: Treat this as an abnormally long hang and take the following steps:

1. Configure `agent.tool_timeout_minutes` in **config.yaml** (default: 10); a single tool execution that exceeds this duration is terminated automatically.
2. Use "Stop Task" for that task on the monitoring page to interrupt it immediately.
3. If the tool genuinely requires more time, increase `tool_timeout_minutes` appropriately, but do not set it to 0.

### Q: What should I do if tool execution fails?

A: Check the following:

1. Whether the tool is installed and available on the system PATH
2. Whether the tool configuration is correct
3. Whether the argument format meets the requirements
4. The server logs for detailed error information

### Q: How do I test a tool configuration?

A: Use `cmd/test-config/main.go` to test configuration loading:

```bash
go run cmd/test-config/main.go
```

### Q: How is argument order controlled?

A: Use the `position` field to control positional-argument order. **The argument at position 0 (such as Gobuster's `dir` subcommand) immediately follows the command name and precedes all flag arguments**, supporting CLIs that require a "subcommand + options" form. Remaining flag arguments are added in their order in the `parameters` list, followed by the remaining positional arguments in positions 1, 2, and so on. `additional_args` is appended to the end of the command.

## Tool Configuration Templates

### Basic Tool Template

```yaml
name: "tool_name"
command: "command"
enabled: true

short_description: "Short description (20-50 characters)"

description: |
  Detailed description of the tool's functionality, use cases, and precautions.

parameters:
  - name: "target"
    type: "string"
    description: "Target parameter description"
    required: true
    position: 0
    format: "positional"

  - name: "additional_args"
    type: "string"
    description: "Additional tool arguments"
    required: false
    format: "positional"
```

### Tool Template with Flag Arguments

```yaml
name: "tool_name"
command: "command"
enabled: true

short_description: "Short description"

description: |
  Detailed description.

parameters:
  - name: "target"
    type: "string"
    description: "Target"
    required: true
    flag: "-t"
    format: "flag"

  - name: "option"
    type: "bool"
    description: "Option"
    required: false
    default: false
    flag: "--option"
    format: "flag"

  - name: "level"
    type: "int"
    description: "Level"
    required: false
    default: 3
    flag: "--level"
    format: "combined"

  - name: "additional_args"
    type: "string"
    description: "Additional arguments"
    required: false
    format: "positional"
```

## Related Documentation

- Main project README: see `README.md` for complete project documentation
- Tool list: see all tool configuration files under `tools/`
- API documentation: see the API endpoint documentation in the main README
