# Skills Guide

Skills provide Agents with on-demand domain capabilities, procedures, templates, and reference material. They are suitable for stable methodology, not one-off task input.

## Structure

Default directory:

```yaml
skills_dir: skills
```

Recommended structure:

```text
skills/
  api-security-testing/
    SKILL.md
  ssrf-testing/
    SKILL.md
  cyberstrike-eino-demo/
    SKILL.md
    REFERENCE.md
    assets/
```

Every Skill must contain at least `SKILL.md`.

## `SKILL.md`

`SKILL.md` uses YAML front matter:

```markdown
---
name: ssrf-testing
description: SSRF identification, validation, bypass, and remediation workflow
---

# SSRF Testing

Use this Skill when the task involves server-side request forgery, URL callbacks, cloud metadata access, or internal-network probing.
```

The `description` is important because the Agent uses it to decide when to load the Skill.

## Progressive Disclosure

Eino Skills support on-demand loading. Configure them with:

```yaml
multi_agent:
  eino_skills:
    disable: false
    filesystem_tools: true
    skill_tool_name: skill
```

The Agent initially sees only the Skill name and description. It calls `skill` to read details when needed, reducing context use.

## Content Suitable for a Skill

Good candidates include:

- a vulnerability-testing procedure;
- a security-audit checklist;
- a report template;
- a tool-combination method;
- an internal standard;
- common false-positive decisions.

Do not use a Skill for temporary target information, API keys, passwords, cookies, frequently changing scan results, or large unstructured raw logs.

## Supporting Files

A Skill may include `REFERENCE.md`, templates, dictionaries, or examples. `SKILL.md` should explain when to read these files.

- Keep the main file short and clear.
- Split references by topic.
- Read large files only when necessary.

## Binding to Roles

A role can prompt the Agent to use a class of Skills; Skills can also be related through page management and roles.

- Keep general Skills unbound so descriptions can trigger them automatically.
- Bind high-risk Skills to specialized roles.
- Avoid overlapping descriptions for Skills in the same category.

## Development Recommendations

Skill content should cover:

1. Trigger conditions.
2. Objectives and boundaries.
3. Operating steps.
4. Tool recommendations.
5. Output format.
6. Risks and prohibited actions.
7. Reference material.

Write so the Agent can execute it, not merely so a person can read it.

## Troubleshooting

If a Skill is not used, check whether its `description` is too narrow or vague, the task lacks trigger keywords, `multi_agent.eino_skills.disable: true` is set, or the front matter is invalid.

If too much Skill content is loaded, split supporting files, state “read Y only when X is needed” in `SKILL.md`, and remove duplication.

## Skill Design in Depth

The core value of a Skill is not teaching a concept; it is giving the Agent an executable procedure at the right time. Pay particular attention to trigger and stop conditions.

Recommended structure:

```markdown
## When to use
Describe the trigger conditions.

## Preconditions
State required user input and target requirements.

## Procedure
Give ordered steps with tools, inputs, and decision criteria.

## Stop conditions
State when to stop, request approval, or hand off to a person.

## Output
Define the final result format.
```

## Anti-Patterns

| Anti-pattern | Result | Fix |
| --- | --- | --- |
| Description too broad: `for security testing` | triggers almost every task | specify vulnerability, scenario, and signal |
| Encyclopedia-style content | Agent does not know the next step | rewrite as a procedure and decision tree |
| Sensitive configuration in a Skill | leakage and misuse | use runtime configuration or user input |
| One Skill containing everything | high read cost and noisy retrieval | split by vulnerability or task |
| No stop condition | Agent may keep expanding scope | state when to stop and request approval |

## Skill vs Knowledge Base

- Skill: tells the Agent how to act and emphasizes procedures.
- Knowledge base: provides facts, cases, and references and emphasizes retrieval.

For SSRF, the Skill describes how to test, decide, and stop; the knowledge base stores cloud-provider metadata addresses, historical bypasses, and remediation guidance.

## Local Tool Risk

`filesystem_tools: true` exposes local read, write, and execute capability. It is useful for development and automation but is also a security boundary. In production:

- constrain the workspace with `workspace_root_dir`;
- use HITL for write and execute actions;
- do not add `execute` to the global allowlist;
- explicitly forbid out-of-scope file access in the Skill.

## Source Anchors

- Skill package validation: `internal/skillpackage/validate.go`
- Skill service: `internal/skillpackage/service.go`
- Eino Skills integration: `internal/multiagent/eino_skills.go`
- Skills Handler: `internal/handler/skills.go`
