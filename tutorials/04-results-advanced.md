# 4. Kết quả & Tính năng Nâng cao

## 4.1 Quản lý Kết quả

| Tab | Mục đích |
|---|---|
| **Facts** | Tri thức được trích xuất từ chat (target, finding, chain). Liên kết với vulnerability. |
| **Assets** | Tài sản (host, domain, IP, service). Filter theo status, project, risk level. Import/export CSV/XLSX. |
| **Vulnerabilities** | Lỗ hổng phát hiện. Filter theo severity (critical/high/medium/low/info), status (open/confirmed/fixed/false_positive/ignored). |
| **Tasks** | Tác vụ agent đang chạy hoặc đã hoàn thành. Xem timeline, result, error. |

### Facts

Facts là đơn vị tri thức (knowledge unit) trong project blackboard.

- **Key**: Mã định danh (ví dụ: `finding:subdomain`).
- **Category**: target, finding, exploit, auth, infra, chain, poc, business, note.
- **Confidence**: tentative / confirmed / deprecated.

![Facts tab](./assets/results/01-facts.png)

### Assets

Quản lý tài sản phát hiện:

- Filter theo: status (new/active/deprecated), project, risk level, scan state.
- Batch actions: gán project, bulk edit, merge duplicates, export.
- Import: Drag-drop CSV/XLSX.

![Asset Library](./assets/results/02-assets.png)

### Vulnerabilities

Theo dõi lỗ hổng:

- Card hiển thị: title, severity, status, description, reproduction steps, evidence, impact, recommendation.
- Tạo từ fact: **Link to existing vulnerability** hoặc **Create vulnerability from fact**.
- Export: Markdown (.md cho từng vulnerability).

![Vulnerabilities list](./assets/results/03-vulnerabilities.png)

### Tasks

Xem và theo dõi tác vụ:

- Filter: status (running/completed/failed/cancelled), time (last 15m/1h/24h/7d).
- Detail: ID, session, type, command, timeline, result.

---

## 4.2 Tính năng Nâng cao

### Workflows

Workflow là đồ thị các node (start, tool, agent, condition, hitl, output, end) để tự động hóa tác vụ phức tạp.

- Kéo-thả node, kết nối, cấu hình properties.
- **AI generate**: Nhập mô tả tự nhiên → System tạo draft graph.
- **Dry-run**: Chạy thử, xem trace.
- Export/Import: File `.csapkg.zip`.

![Workflow canvas](./assets/advanced/01-workflow.png)

### MCP (Model Context Protocol)

Theo dõi và quản lý MCP servers đã cấu hình:

- **MCP Monitor**: Xem tool execution real-time, terminate/cancel.
- **MCP Management**: Thêm/xóa external servers, test connection.
- Tool MCP xuất hiện tự động trong Chat khi server kết nối thành công.

![MCP Monitor](./assets/advanced/02-mcp.png)

### Knowledge Base

Kho tri thức (tài liệu, CVE, best practices) để AI truy xuất khi cần (RAG).

- Thêm item: Category, title, content (markdown/text).
- **Scan** → Phát hiện item mới.
- **Build Index** → Tạo vector embeddings.
- AI tự động retrieval khi cần. Xem **Retrieval Logs** để theo dõi.

![Knowledge Base](./assets/advanced/03-knowledge.png)

### Skills

Skill là gói kỹ năng (SKILL.md + script/file) cho agent.

- Xem: **Skills Monitor** (thống kê calls) hoặc **Management** (CRUD).
- Tạo: Name (lowercase, hyphenated), description, content (SKILL.md).
- Load trong Chat: Dùng tool `skill` → Chọn skill → Load.

![Skills Monitor](./assets/advanced/04-skills.png)

### Agents

Agents là sub-agents (multi-agent) định nghĩa bằng file Markdown trong `agents/`.

- **Orchestrator**: Agent chủ, điều phối.
- **Sub**: Agent con, thực thi tác vụ cụ thể.
- Cấu hình: tools, max iterations, instruction (system prompt).

![Agents Management](./assets/advanced/05-agents.png)

### Roles

Role là persona (prompt + tool scope) cho agent.

- Chọn trong Chat → Tool list tự cập nhật.
- Tạo: Name, icon (emoji), description, user prompt, tools.
- Enable/disable role.

![Roles list](./assets/advanced/06-roles.png)

---

## Lưu ý

**C2 (Command & Control)** và **WebShell Management** là tính năng out-of-scope trong tài liệu này — xem `docs/en-US/c2.md` và `docs/en-US/webshell.md` để biết chi tiết.

---

**Tiếp theo:** [5. Troubleshooting →](05-troubleshooting.md)
