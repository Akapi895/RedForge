# Citations and External Links (Example)

This file verifies that the **`references/`** directory inside a skill package is correctly recognized by the listing API, HTTP `resource_path`, and multi-agent local filesystem tools.

## Test Procedure (Authorized Environment)

1. The `package_files` in the response to `GET /api/skills/cyberstrike-eino-demo` should include `references/citations.md`.
2. `GET /api/skills/cyberstrike-eino-demo?resource_path=references/citations.md` should return this file's content.
3. With multi-agent mode and `eino_skills.filesystem_tools` enabled, this file should be readable by relative path.

## Placeholder Citation

- [OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/) (link-format example only)
