# 3. Quy trình Sử dụng

## Quy trình hoàn chỉnh

### Bước 1: Tạo Project

1. Vào **Projects** → Click **New Project**.
2. Điền:
   - **Name**: Tên project.
   - **Description**: Mô tả engagement.
   - **Scope** (optional): JSON mô tả phạm vi (target, constraints).
3. Click **Create**.

![New Project form](./assets/basic-usage/01-new-project.png)

Project là container để gom nhiều conversation, fact, asset, vulnerability liên quan đến cùng một engagement.

---

### Bước 2: Mở Chat

1. Vào **Chat**.
2. Click **New Conversation** hoặc chọn project từ dropdown.
3. Nhập yêu cầu trong composer (khung nhập).

![Chat interface](./assets/basic-usage/02-chat-interface.png)

---

### Bước 3: Chọn Role

1. Click **Role selector** (dropdown hoặc side panel).
2. Chọn role phù hợp (ví dụ: `Penetration Testing`, `API Security Assessment`, `Controlled Validation`).
3. Role định nghĩa:
   - **Persona** (system prompt).
   - **Tool scope** (công cụ được phép dùng).

![Role selector](./assets/basic-usage/03-role-selector.png)

---

### Bước 4: Chọn Agent Mode

1. Click **Agent mode selector**.
2. Chọn mode:
   - **eino_single**: Single agent (đơn giản).
   - **deep**: DeepAgent mode cho các task phức tạp.
   - **plan_execute**: Planner → Executor → Replanner loop.
   - **supervisor**: Supervisor điều phối sub-agents.

![Agent mode selector](./assets/basic-usage/04-agent-mode.png)

---

### Bước 5: Gửi Task

Nhập yêu cầu trong composer, ví dụ:

```
Scan open ports on 192.168.1.1
Check if https://example.com/page?id=1 is vulnerable to SQL injection
Enumerate subdomains for example.com, then run nuclei against the results
```

Click **Send** (hoặc Enter).

---

### Bước 6: Theo dõi và Review Kết quả

- **Timeline** hiển thị:
  - Tool calls (công cụ được gọi).
  - Tool results (kết quả).
  - Execution trace / Agent steps (các bước thực thi).
  - Iterations (vòng lặp).
- **Progress indicator**:
  - Đang chạy: hiển thị thời gian elapsed.
  - Đợi phê duyệt (HITL): hiển thị countdown.
  - Hoàn thành: hiển thị summary.
- **Stop task**: Click nút **Stop** để hủy tác vụ đang chạy.

![Timeline tool calls](./assets/basic-usage/05-timeline.png)

---

**Tiếp theo:** [4. Kết quả & Tính năng Nâng cao →](04-results-advanced.md)
