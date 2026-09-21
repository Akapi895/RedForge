# CyberStrikeAI User Guide (Hướng dẫn sử dụng)

> **Phiên bản:** 2026-09-21  
> **Mục tiêu:** Tài liệu hướng dẫn người dùng mới làm quen với CyberStrikeAI — từ cài đặt, cấu hình, đến sử dụng các tính năng chính.  
> **Lưu ý:** Tài liệu này tập trung vào quy trình pentest an toàn, có kiểm soát. Không đề cập sâu C2/WebShell.

---

## 📚 Mục lục

### 1. [Tổng quan](01-overview.md)
- CyberStrikeAI là gì
- Các chức năng chính
- Phạm vi tài liệu

### 2. [Cài đặt và khởi chạy](02-installation.md)
- Yêu cầu môi trường
- Clone repo
- Cài dependencies
- Chạy `run.sh`
- Truy cập Web UI

### 3. [Cấu hình](03-configuration.md)
- `config.yaml`
- AI provider / model
- Tools
- MCP
- Human-in-the-loop
- Các config quan trọng khác

### 4. [Quy trình sử dụng cơ bản](04-basic-usage.md)
- Tạo Project
- Tạo Conversation
- Chọn Role
- Chọn Agent mode
- Gửi task
- Theo dõi quá trình thực thi

### 5. [Quản lý kết quả](05-results-management.md)
- Fact Board
- Assets
- Vulnerabilities
- Tasks
- Evidence / trạng thái

### 6. [Các chức năng mở rộng](06-advanced-features.md)
- Workflows
- MCP
- Knowledge
- Skills
- Agents
- Roles

### 7. [Troubleshooting](07-troubleshooting.md)
- Không chạy được server
- Không gọi được model
- Tool not found
- MCP lỗi
- Agent không chạy / không gọi tool

---

## 📸 Hướng dẫn chèn ảnh minh họa

Để thêm ảnh vào tài liệu:

1. **Chụp ảnh màn hình** từ Web UI hoặc terminal.
2. **Lưu vào folder phù hợp** trong `assets/`:
   - `assets/overview/` — Ảnh tổng quan hệ thống, dashboard
   - `assets/installation/` — Ảnh quá trình cài đặt, terminal
   - `assets/configuration/` — Ảnh các trang cấu hình, settings
   - `assets/basic-usage/` — Ảnh chat, project, conversation
   - `assets/results/` — Ảnh fact board, assets, vulnerabilities
   - `assets/advanced/` — Ảnh workflows, MCP, knowledge, skills
   - `assets/troubleshooting/` — Ảnh lỗi, log terminal

3. **Thêm ảnh vào file markdown tương ứng** theo cú pháp:
   ```markdown
   ![Mô tả ảnh](./assets/<folder>/<ten-anh>.png)
   ```

4. **Đặt tên ảnh** theo thứ tự và mô tả:
   - `01-dashboard.png` — Dashboard overview
   - `02-run-sh-output.png` — Output khi chạy run.sh
   - `03-ai-channel-config.png` — Cấu hình AI channel
   - `04-new-project-form.png` — Form tạo project mới
   - ...

5. **Kích thước ảnh khuyến nghị:**
   - Ảnh UI: tối đa 1920px chiều rộng
   - Ảnh terminal: giữ nguyên tỷ lệ
   - Nén ảnh trước khi thêm (sử dụng TinyPNG hoặc tương tự)

---

## 📖 Tài liệu tham khảo

- **Deployment**: `docs/en-US/deployment.md`
- **Configuration**: `docs/en-US/configuration.md`
- **Security**: `docs/en-US/security-model.md`, `docs/en-US/security-hardening.md`
- **Asset Management**: `docs/en-US/asset-management.md`
- **Knowledge Base**: `docs/en-US/knowledge-base.md`
- **Skills**: `docs/en-US/skills-guide.md`
- **Workflows**: `docs/en-US/workflow-graph.md`
- **MCP**: `docs/en-US/mcp-federation.md`
- **HITL**: `docs/en-US/hitl-best-practices.md`
- **Troubleshooting**: `docs/en-US/troubleshooting.md`
