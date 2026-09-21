#!/usr/bin/env python3
"""
Pent Claude Agent MCP Server - penetration-testing engineer MCP service

Exposes AI penetration-testing capabilities through MCP: CyberStrikeAI can direct pent_claude_agent to perform penetration-testing tasks.
pent_claude_agent uses the Claude Agent SDK internally and can independently configure MCP servers and tools, running as an independent penetration-testing engineer.

Dependencies: pip install mcp claude-agent-sdk(or use the project venv)
Run: python mcp_pent_claude_agent.py [--config /path/to/config.yaml]
"""

from __future__ import annotations

import argparse
import asyncio
import os
from typing import Any

import yaml
from mcp.server.fastmcp import FastMCP

# Lazy import so a missing package does not affect MCP startup
_claude_sdk_available = False
try:
    from claude_agent_sdk import ClaudeAgentOptions, query

    _claude_sdk_available = True
except ImportError:
    pass

# ---------------------------------------------------------------------------
# Paths and configuration
# ---------------------------------------------------------------------------

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.dirname(os.path.dirname(SCRIPT_DIR))
_DEFAULT_CONFIG_PATH = os.path.join(SCRIPT_DIR, "pent_claude_agent_config.yaml")

# Agent runtime state(simple in-memory state for status)
_last_task: str | None = None
_last_result: str | None = None
_task_count: int = 0


def _load_config(config_path: str | None) -> dict[str, Any]:
    """Load YAML configuration and merge defaults with user configuration."""
    defaults: dict[str, Any] = {
        "cwd": PROJECT_ROOT,
        "allowed_tools": ["Read", "Write", "Bash", "Grep", "Glob"],
        "env": {
            "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
            "DISABLE_TELEMETRY": "1",
            "DISABLE_ERROR_REPORTING": "1",
            "DISABLE_BUG_COMMAND": "1",
        },
        "mcp_servers": {},
        "system_prompt": (
            "You are a professional penetration-testing engineer. Based on the user's task, perform security testing, vulnerability analysis, information gathering, and related work."
            "Execute step by step and output clear, reproducible results. Test only within the authorized scope."
        ),
    }
    path = config_path or os.environ.get("PENT_CLAUDE_AGENT_CONFIG", _DEFAULT_CONFIG_PATH)
    if not os.path.isfile(path):
        return defaults
    try:
        with open(path, "r", encoding="utf-8") as f:
            user = yaml.safe_load(f) or {}
        # Deep merge
        def merge(base: dict, override: dict) -> dict:
            out = dict(base)
            for k, v in override.items():
                if k in out and isinstance(out[k], dict) and isinstance(v, dict):
                    out[k] = merge(out[k], v)
                else:
                    out[k] = v
            return out

        return merge(defaults, user)
    except Exception:
        return defaults


def _resolve_path(s: str) -> str:
    """Resolve path placeholders."""
    return s.replace("${PROJECT_ROOT}", PROJECT_ROOT).replace("${SCRIPT_DIR}", SCRIPT_DIR)


def _build_agent_options(config: dict[str, Any], cwd_override: str | None = None) -> ClaudeAgentOptions:
    """Build ClaudeAgentOptions from configuration."""
    raw_cwd = cwd_override or config.get("cwd", PROJECT_ROOT)
    cwd = _resolve_path(str(raw_cwd)) if isinstance(raw_cwd, str) else str(raw_cwd)
    env = dict(os.environ)
    env.update(config.get("env", {}))
    mcp_servers = config.get("mcp_servers") or {}
    # Resolve path placeholders
    for name, cfg in list(mcp_servers.items()):
        if isinstance(cfg, dict):
            args = cfg.get("args") or []
            cfg = dict(cfg)
            cfg["args"] = [_resolve_path(str(a)) for a in args]
            mcp_servers[name] = cfg

    return ClaudeAgentOptions(
        cwd=cwd,
        allowed_tools=config.get("allowed_tools", ["Read", "Write", "Bash", "Grep", "Glob"]),
        disallowed_tools=config.get("disallowed_tools", []),
        mcp_servers=mcp_servers,
        env=env,
        system_prompt=config.get("system_prompt"),
        setting_sources=config.get("setting_sources", ["user", "project"]),
    )


async def _run_claude_agent(prompt: str, config_path: str | None = None, cwd: str | None = None) -> str:
    """Run the Claude Agent internally and return the last round's text result."""
    global _last_task, _last_result, _task_count
    _last_task = prompt
    _task_count += 1

    if not _claude_sdk_available:
        _last_result = "Error: claude-agent-sdk is not installed; run pip install claude-agent-sdk"
        return _last_result

    config = _load_config(config_path)
    options = _build_agent_options(config, cwd_override=cwd)

    messages: list[Any] = []
    try:
        async for message in query(prompt=prompt, options=options):
            messages.append(message)
    except Exception as e:
        _last_result = f"Agent execution error: {e}"
        return _last_result

    if not messages:
        _last_result = "(no output)"
        return _last_result

    # During multi-turn iteration, take the last ResultMessage(the latest result)
    result_msgs = [m for m in messages if hasattr(m, "result") and getattr(m, "result", None) is not None]
    last = result_msgs[-1] if result_msgs else messages[-1]
    # Extract text, preferring ResultMessage.result to avoid outputting metadata
    if hasattr(last, "result") and last.result is not None:
        text = last.result
    elif hasattr(last, "content") and last.content:
        parts = []
        for block in last.content:
            if hasattr(block, "text") and block.text:
                parts.append(block.text)
        text = "\n".join(parts) if parts else "(no output)"
    else:
        text = "(no output)"
    _last_result = text
    return _last_result


# ---------------------------------------------------------------------------
# MCP service and tools
# ---------------------------------------------------------------------------

app = FastMCP(
    name="pent-claude-agent",
    instructions="Penetration-testing engineer MCP: after receiving a task, start a Claude Agent internally to independently perform penetration testing, vulnerability analysis, and related work, then return the results.",
)


@app.tool(
    description="Execute a penetration-testing task. After receiving a task description, pent_claude_agent acts as an independent penetration-testing engineer, uses Claude Agent to execute the task, and returns the results. Supports port scanning, vulnerability probing, Web security testing, information gathering, and more.",
)
async def pent_claude_run_pentest_task(task: str) -> str:
    """Run a penetration testing task. The agent executes independently and returns results."""
    return await _run_claude_agent(task)


@app.tool(
    description="Analyze vulnerability information. Pass a vulnerability description, PoC, impact scope, and related details for professional Agent analysis and remediation recommendations.",
)
async def pent_claude_analyze_vulnerability(vuln_info: str) -> str:
    """Analyze vulnerability information and provide remediation suggestions."""
    prompt = f"Professionally analyze the following vulnerability information, including risk level, impact scope, exploitation method, and remediation recommendations.\n\n{vuln_info}"
    return await _run_claude_agent(prompt)


@app.tool(
    description="Execute a specified task. This is a general task entry point; the Agent automatically selects suitable tools and methods based on the task.",
)
async def pent_agent_execute(task: str) -> str:
    """Execute a task. The agent chooses appropriate tools and methods."""
    return await _run_claude_agent(task)


@app.tool(
    description="Perform security diagnostics on a target. Pass a URL, IP, domain, or similar target; the Agent performs an initial security assessment and diagnosis.",
)
async def pent_agent_diagnose(target: str) -> str:
    """Diagnose a target (URL, IP, domain) for security assessment."""
    prompt = f"Perform security diagnostics and an initial assessment of the following target: {target}\n\nInclude reachability, open services, common vulnerability surfaces, and related details."
    return await _run_claude_agent(prompt)


@app.tool(
    description="Get the current status of pent_claude_agent: latest task, result summary, execution count, and more.",
)
def pent_claude_status() -> str:
    """Get the current status of pent_claude_agent."""
    global _last_task, _last_result, _task_count
    lines = [
        f"Task executions: {_task_count}",
        f"Latest task: {_last_task or '-'}",
        f"Latest result summary: {(str(_last_result or '-')[:200] + '...') if _last_result and len(str(_last_result)) > 200 else (_last_result or '-')}",
        f"Claude SDK available: {_claude_sdk_available}",
    ]
    return "\n".join(lines)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Pent Claude Agent MCP Server")
    parser.add_argument(
        "--config",
        default=None,
        help="Path to pent_claude_agent config YAML (env: PENT_CLAUDE_AGENT_CONFIG)",
    )
    args, _ = parser.parse_known_args()
    # Store the config path in the environment for tool calls
    if args.config:
        os.environ["PENT_CLAUDE_AGENT_CONFIG"] = args.config
    app.run(transport="stdio")
