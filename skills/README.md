# Skills Directory (Agent Skills / Eino)

- Each skill is a **subdirectory** whose root must contain **`SKILL.md`** (YAML front matter with `name`, `description`, and Markdown body); see [agentskills.io](https://agentskills.io/specification.md).
- **The directory name must match `name`.**
- **Runtime loading**: in an **Eino DeepAgent (multi-agent)** session, the ADK **`skill` middleware** progressively discloses skills (the system prompt lists each skill's name/description, and the model calls the **`skill`** tool to load the full `SKILL.md`). You may enable **`multi_agent.eino_skills.filesystem_tools`** to access package scripts and resources through tools such as `read_file` / `execute`, just as on the local machine.
- **Web management**: HTTP `/api/skills/*` continues to list, edit, and upload package files (implemented by `internal/skillpackage`, not MCP).
- **Runtime**: in multi-agent (DeepAgent) sessions, the ADK **`skill`** tool loads skills progressively; the single-agent MCP loop does not include Skills, so enable multi-agent or use a later single-agent Eino path.
