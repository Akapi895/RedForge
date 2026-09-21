# CyberStrikeAI User Guide (Hướng dẫn sử dụng)

> **Phiên bản:** 2026-09-21  
> **Mục tiêu:** Tài liệu hướng dẫn người dùng mới làm quen với CyberStrikeAI — từ cài đặt, cấu hình, đến sử dụng các tính năng chính.  
> **Lưu ý:** Tài liệu này tập trung vào quy trình pentest an toàn, có kiểm soát. Không đề cập sâu C2/WebShell (xem tài liệu riêng cho các tính năng nâng cao).

---

## 1. Tổng quan

### 1.1 CyberStrikeAI là gì?

CyberStrikeAI là **nền tảng pentest tự động hóa AI-native** — kết hợp lập kế hoạch, thực thi, giám sát con người, bằng chứng (evidence), và phân tích trong một workspace có kiểm soát. Hệ thống được xây dựng bằng Go, sử dụng các agent AI để thực hiện các tác vụ bảo mật theo ngữ cảnh (reconnaissance, vulnerability analysis, exploitation guidance, reporting).

**Đặc điểm chính:**
- **Agent-based execution**: AI phân tích yêu cầu, lập kế hoạch, gọi công cụ (tools), và tổng hợp kết quả.
- **Human-in-the-loop (HITL)**: Hỗ trợ phê duyệt trước khi thực thi các công cụ nhạy cảm.
- **Evidence-based**: Mọi kết quả được lưu trữ và liên kết (fact board, assets, vulnerabilities).
- **Modular & extensible**: Hỗ trợ MCP (Model Context Protocol), skills, roles, workflows.
- **RBAC**: Quản lý người dùng, vai trò, và phân quyền chi tiết.

### 1.2 Các chức năng chính

| Chức năng | Mô tả |
|---|---|
| **Chat / Agent** | Giao tiếp với AI, yêu cầu thực hiện tác vụ (scan, phân tích, tìm kiếm thông tin). |
| **Project & Attack Chain** | Tổ chức engagement theo dự án, liên kết các sự kiện thành chuỗi tấn công. |
| **Asset Management** | Quản lý tài sản (host, domain, IP, service), import/export, tìm kiếm qua FOFA/ZoomEye/Shodan. |
| **Vulnerability Management** | Theo dõi lỗ hổng, phân loại mức độ, liên kết với asset/project. |
| **Task Management** | Tạo batch tasks (nhiều tác vụ hàng loạt), theo dõi tiến độ. |
| **MCP Integration** | Kết nối công cụ bên ngoài qua Model Context Protocol (HTTP/stdio/SSE). |
| **Knowledge Base** | Lưu trữ tri thức (tài liệu, CVE, best practices) để AI truy xuất khi cần. |
| **Skills & Roles** | Skills = gói kỹ năng (SKILL.md + script); Roles = persona (prompt + công cụ). |
| **Audit & Monitor** | Ghi log thao tác, theo dõi tool execution, cảnh báo. |

### 1.3 Phạm vi tài liệu

Tài liệu này hướng dẫn:
- Cài đặt và khởi chạy (local, Docker).
- Cấu hình cơ bản (AI, tools, MCP, HITL).
- Quy trình sử dụng từ tạo project đến quản lý kết quả.
- Troubleshooting các lỗi thường gặp.

**Không đề cập sâu:**
- C2 (Command & Control) — xem `docs/en-US/c2.md`.
- WebShell management — xem `docs/en-US/webshell.md`.
- Plugin development (Burp Suite, browser extension) — xem `docs/en-US/plugin-development.md`.

---

## 2. Cài đặt và khởi chạy

### 2.1 Yêu cầu môi trường

| Thành phần | Yêu cầu |
|---|---|
| **Go** | 1.21+ (cho build từ source) |
| **Python** | 3.10+ (cho một số tool/scripts) |
| **Docker** | 20.10+ (nếu dùng Docker — không bắt buộc) |
| **Security tools** | nmap, sqlmap, nuclei, gobuster... (cài riêng vào PATH) |
| **AI provider** | API key từ OpenAI/DeepSeek/Qwen/Claude... |

### 2.2 Clone repo

```bash
git clone https://github.com/Ed1s0nZ/CyberStrikeAI.git
cd CyberStrikeAI
```

### 2.3 Cài dependencies

#### Option A: Dùng `run.sh` (khuyến nghị cho local)

```bash
chmod +x run.sh
./run.sh
```

Script tự động:
- Kiểm tra Go/Python.
- Tạo Python venv và cài dependencies (`requirements.txt`).
- Tải Go dependencies.
- Build và start server.

#### Option B: Build thủ công

```bash
# Go dependencies
go mod download

# Build
go build -o cyberstrike-ai ./cmd/server

# Python venv (cho tool scripts)
python3 -m venv venv
source venv/bin/activate  # Linux/macOS
# hoặc venv\Scripts\activate  # Windows
pip install -r requirements.txt
```

### 2.4 Chạy server

```bash
# HTTPS (self-signed cert — khuyến nghị cho local)
./cyberstrike-ai --https

# Hoặc HTTP (không khuyến khích)
./cyberstrike-ai --http
```

Hoặc dùng `run.sh` (đã cấu hình HTTPS mặc định).

### 2.5 Truy cập Web UI

- Mở trình duyệt: `https://localhost:8080/` (hoặc `http://localhost:8080/` nếu dùng HTTP).
- Chấp nhận cảnh báo chứng chỉ tự-signed (nếu dùng HTTPS).
- Đăng nhập lần đầu:
  - Username: `admin`
  - Password: **tự sinh ngẫu nhiên** — xem log terminal để lấy password, hoặc chạy `./cyberstrike-ai --reset-admin-password` để đặt lại.

---

## 3. Cấu hình

### 3.1 `config.yaml`

File cấu hình chính, nằm ở root repo. Nhiều trường có thể chỉnh qua **Web Settings** (System Settings → Basic Settings), nhưng một số thay đổi cần restart server.

**Cấu trúc chính:**
```yaml
server:
  host: 0.0.0.0
  port: 8080
  tls_enabled: true
  tls_auto_self_sign: true

auth:
  session_duration_hours: 12

ai:
  default_channel: <channel-id>
  channels:
    <channel-id>:
      name: Display name
      provider: openai_compatible  # hoặc claude
      base_url: https://...
      api_key: sk-...
      model: gpt-4o
      max_total_tokens: 120000
      max_completion_tokens: 16384

agent:
  max_iterations: 12000
  tool_timeout_minutes: 60

mcp:
  enabled: true  # bật MCP server HTTP
  port: 8081

hitl:
  default_mode: off  # off | approval | review_edit
  tool_whitelist: [...]
```

### 3.2 AI provider / model

**Cách 1: Web UI** (khuyến nghị)

1. Vào **System Settings → Basic Settings → AI Channel Configuration**.
2. Click `+` để thêm channel mới.
3. Điền:
   - **Name**: Tên hiển thị.
   - **Provider**: `openai_compatible` (OpenAI/DeepSeek/Qwen...) hoặc `claude` (Anthropic).
   - **Base URL**: Ví dụ `https://api.openai.com/v1`, `https://api.deepseek.com/v1`.
   - **API Key**: Khóa API (không commit vào git!).
   - **Model**: Tên model (`gpt-4o`, `deepseek-chat`, `qwen3-max`...).
   - **Max tokens**: `max_total_tokens` (context), `max_completion_tokens` (output).
4. Click **Save changes**.
5. Set **default channel**: chọn channel → **Set as default** → Save.

**Cách 2: Chỉnh `config.yaml` trực tiếp**

```yaml
ai:
  default_channel: netmind
  channels:
    netmind:
      name: Viettel MiniMax M3
      provider: openai_compatible
      base_url: https://stream-netmind.viettel.vn/aigw/ai/v1
      api_key: "${CYBERSTRIKE_AI_API_KEY}"  # hoặc khóa thật
      model: MiniMax/MiniMax-M3
      max_total_tokens: 120000
      max_completion_tokens: 16384
```

Sau đó restart server hoặc dùng **System Settings → Apply config**.

### 3.3 Tools

System tự động load tools từ thư mục `tools/` (file `.yaml`). Mỗi tool có:
- `name`, `description`
- `command` (lệnh thực thi)
- `args` (tham số)
- `timeout`, `security` (nếu cần)

**Kiểm tra tools:**
- Vào **Tools Management** (nếu có) hoặc **Chat → @mention tool**.
- Tools không được cài vào PATH sẽ bị bỏ qua (không lỗi server).

**Cài security tools (ví dụ):**
```bash
# Kali/Ubuntu
sudo apt install nmap sqlmap nikto gobuster hydra hashcat nuclei

# macOS (Homebrew)
brew install nmap sqlmap gobuster
```

### 3.4 MCP (Model Context Protocol)

MCP cho phép kết nối công cụ bên ngoài (HTTP/stdio/SSE).

**Bật MCP:**
```yaml
mcp:
  enabled: true
  port: 8081
```

**Cấu hình external MCP servers:**
```yaml
external_mcp:
  servers:
    my-mcp-server:
      transport: http  # hoặc stdio, sse
      url: https://mcp-server.example.com/mcp
      # hoặc command + args cho stdio
      description: "Mô tả server"
```

**Test kết nối:**
- Vào **MCP Management** (Web UI).
- Click **Test connection** hoặc xem **MCP Monitor** để theo dõi tool execution.

### 3.5 Human-in-the-Loop (HITL)

HITL cho phép phê duyệt trước khi thực thi công cụ nhạy cảm.

**Cấu hình:**
```yaml
hitl:
  default_mode: approval  # off | approval | review_edit
  default_reviewer: human  # human | audit_agent
  tool_whitelist: [read_file, ls, glob, grep]  # tool không cần phê duyệt
```

**Web UI:**
- Vào **Security → Human-in-the-Loop**.
- Chọn **mode** (approval/review_edit).
- Chọn **reviewer** (human/audit agent).
- Edit **tool whitelist**.

**Trong Chat:**
- Mỗi session có thể cấu hình HITL riêng (settings popover).
- Khi HITL bật, tool call sẽ chờ phê duyệt → xem **Security → HITL Pending**.

### 3.6 Các config quan trọng khác

| Config | Mục đích |
|---|---|
| `agent.max_iterations` | Số vòng lặp tối đa cho agent (mặc định 12000). |
| `agent.tool_timeout_minutes` | Thời gian chờ tool tối đa (mặc định 60 phút). |
| `knowledge.enabled` | Bật/tắt Knowledge Base. |
| `c2.enabled` | Bật/tắt C2 (chỉ bật khi cần, trong môi trường được phép). |
| `audit.enabled` | Ghi log audit (mặc định true). |
| `monitor.retention_days` | Số ngày giữ tool execution logs (mặc định 90). |

---

## 4. Quy trình sử dụng cơ bản

### 4.1 Tạo Project

1. Vào **Projects**.
2. Click **New Project**.
3. Điền:
   - **Name**: Tên project.
   - **Description**: Mô tả.
   - **Scope** (optional): JSON mô tả phạm vi (target, constraints).
4. Click **Create**.

**Project** là container để gom nhiều conversation, fact, asset, vulnerability liên quan đến cùng một engagement.

### 4.2 Tạo Conversation

1. Vào **Chat**.
2. Click **New Conversation** (hoặc chọn project từ dropdown).
3. Nhập yêu cầu trong composer (khung nhập).

**Conversation** là một phiên hội thoại với AI, có thể liên kết với project.

### 4.3 Chọn Role

1. Trong Chat, click **Role selector** (side panel hoặc dropdown).
2. Chọn role phù hợp (ví dụ: `渗透测试`, `Web 安全`, `CTF`).
3. Role định nghĩa:
   - **Persona** (system prompt).
   - **Tool scope** (công cụ được phép dùng).

**Roles** nằm trong `roles/`, có thể tùy chỉnh.

### 4.4 Chọn Agent mode

1. Trong Chat, click **Agent mode selector**.
2. Chọn mode:
   - **eino_single**: Single agent (đơn giản).
   - **deep**: Multi-agent với deep reasoning.
   - **plan_execute**: Planner → Executor → Replanner loop.
   - **supervisor**: Supervisor điều phối sub-agents.

**Multi-agent modes** chỉ hiện nếu `multi_agent.enabled: true`.

### 4.5 Gửi task

Nhập yêu cầu trong composer, ví dụ:
```
Scan open ports on 192.168.1.1
Check if https://example.com/page?id=1 is vulnerable to SQL injection
Enumerate subdomains for example.com, then run nuclei against the results
```

Click **Send** (hoặc Enter).

### 4.6 Theo dõi quá trình thực thi

- **Timeline** hiển thị:
  - Tool calls (công cụ được gọi).
  - Tool results (kết quả).
  - Reasoning chain (suy luận).
  - Iterations (vòng lặp).
- **Progress indicator**:
  - Đang chạy: hiển thị thời gian elapsed.
  - Đợi phê duyệt (HITL): hiển thị countdown.
  - Hoàn thành: hiển thị summary.
- **Stop task**: Click nút **Stop** để hủy tác vụ đang chạy.

---

## 5. Quản lý kết quả

### 5.1 Fact Board

**Fact** là đơn vị tri thức (knowledge unit) được lưu trong project blackboard.

**Xem facts:**
1. Vào **Projects → [project name] → Facts tab**.
2. Mỗi fact có:
   - **Key**: Mã định danh (ví dụ: `finding:subdomain`).
   - **Category**: target, finding, exploit, auth, infra, chain, poc, business, note.
   - **Summary**: Tóm tắt.
   - **Body**: Chi tiết (có thể là attack chain).
   - **Confidence**: tentative / confirmed / deprecated.

**Liên kết fact với vulnerability:**
- Trong fact detail, chọn **Link to existing vulnerability** hoặc **Create vulnerability from fact**.

### 5.2 Assets

**Asset** là tài sản (host, domain, IP, service).

**Xem assets:**
1. Vào **Assets → Library**.
2. Filter theo:
   - **Status**: new / active / deprecated.
   - **Project**: liên kết với project.
   - **Risk level**: severe / high / medium / low / info.
   - **Scan state**: never / scanned / overdue.
3. **Batch actions**:
   - **Gán project**.
   - **Bulk edit** (status, responsible person, department, tags...).
   - **Merge duplicates**.
   - **Export** (CSV/XLSX).

**Import assets:**
- Click **Import** → Drag-drop file XLSX/CSV.
- Template có sẵn: **Download template**.

### 5.3 Vulnerabilities

**Vulnerability** là lỗ hổng được phát hiện.

**Xem vulnerabilities:**
1. Vào **Vulnerabilities**.
2. Filter theo:
   - **Severity**: critical / high / medium / low / info.
   - **Status**: open / confirmed / fixed / false_positive / ignored.
   - **Project / Asset / Conversation / Task**.
3. **Card hiển thị**:
   - Title, severity badge, status.
   - Description, reproduction steps.
   - Evidence, impact, recommendation.
   - Linked facts.

**Tạo vulnerability:**
- Click **New Vulnerability**.
- Điền:
  - **Conversation ID** (bắt buộc — từ đâu phát hiện).
  - **Title**, **Description**.
  - **Severity**.
  - **Type** (SQLi, XSS, RCE...).
  - **Target** (URL, host...).
  - **Reproduction steps**, **Evidence**, **Impact**, **Recommendation**.
- Click **Save**.

**Export:**
- Click **Export** → Chọn **Markdown** (tạo file .md cho từng vulnerability).

### 5.4 Tasks

**Task** là một tác vụ agent-loop đang chạy hoặc đã hoàn thành.

**Xem tasks:**
1. Vào **Tasks**.
2. Filter theo:
   - **Status**: running / completed / failed / cancelled.
   - **Time**: last 15m / 1h / 24h / 7d.
3. **Task detail**:
   - ID, session, type, command.
   - Timeline (created → sent → completed).
   - Result (output, error).

**Batch tasks:**
- Vào **Tasks → Batch queues**.
- Click **New Queue**.
- Nhập danh sách task (mỗi dòng 1 task).
- Cấu hình:
  - **Role**, **Agent mode**.
  - **HITL policy**.
  - **Schedule**: manual hoặc cron.
  - **Concurrency**.
- Click **Create** → **Start**.

### 5.5 Evidence / trạng thái

**Evidence** được lưu qua:
- **Chat uploads** (file đính kèm): Vào **Chat Files**.
- **Tool execution outputs**: Vào **MCP Monitor**.
- **C2 events/sessions**: Vào **C2 → Events/Sessions** (nếu bật C2).

**Trạng thái:**
- **Running**: Đang thực thi.
- **Pending approval**: Đợi phê duyệt (HITL).
- **Completed**: Hoàn thành.
- **Failed**: Lỗi (xem error message).
- **Cancelled**: Đã hủy.

---

## 6. Các chức năng mở rộng

### 6.1 Workflows

**Workflow** là đồ thị (graph) các node (start, tool, agent, condition, hitl, output, end) để tự động hóa tác vụ phức tạp.

**Tạo workflow:**
1. Vào **Workflows**.
2. Click **New Workflow**.
3. Kéo-thả node từ palette.
4. Kết nối node (connect mode).
5. Cấu hình từng node (properties panel).
6. Click **Save** → **Validate**.

**AI generate workflow:**
- Click **AI workflow** → Nhập mô tả tự nhiên → System tạo draft graph → **Apply**.

**Dry-run:**
- Nhập message → System chạy thử → Xem trace.

**Export/Import:**
- **Export**: Tải file `.csapkg.zip`.
- **Import**: Upload `.csapkg.zip` → Resolve conflict (keep/overwrite/rename).

### 6.2 MCP

**MCP (Model Context Protocol)** cho phép kết nối công cụ bên ngoài.

**Cấu hình:**
- Xem mục 3.4.

**Sử dụng:**
- Trong Chat, tool MCP xuất hiện trong danh sách.
- **MCP Monitor** theo dõi tool execution, terminate/cancel.

**External MCP servers:**
- HTTP: `url: https://.../mcp`.
- Stdio: `command` + `args`.
- SSE: `url` + `transport: sse`.

### 6.3 Knowledge Base

**Knowledge Base** là kho tri thức (tài liệu, CVE, best practices) để AI truy xuất khi cần.

**Bật Knowledge:**
```yaml
knowledge:
  enabled: true
  base_path: knowledge_base
```

**Thêm knowledge:**
1. Vào **Knowledge → Management**.
2. Click **Add Item**.
3. Điền:
   - **Category**.
   - **Title**.
   - **Content** (markdown/text).
4. Click **Save**.

**Index:**
- Click **Scan** → Phát hiện item mới.
- Click **Build Index** → Tạo vector embeddings.

**Retrieval:**
- AI tự động truy xuất knowledge khi cần (RAG).
- Xem **Knowledge → Retrieval Logs** để biết query nào đã lấy gì.

### 6.4 Skills

**Skill** là gói kỹ năng (SKILL.md + script/file) cho agent.

**Xem skills:**
1. Vào **Skills → Monitor** (thống kê calls) hoặc **Management** (CRUD).
2. Mỗi skill có:
   - **Name**, **Version**, **Description**.
   - **Files** (SKILL.md + scripts).
   - **Tags**.

**Tạo skill:**
1. Click **New Skill**.
2. Điền:
   - **Name** (lowercase, hyphenated: `my-skill`).
   - **Description**.
   - **Content** (SKILL.md).
3. Click **Save**.

**Load skill trong Chat:**
- Dùng tool `skill` (nếu được enable) → Chọn skill → Load.

### 6.5 Agents

**Agents** là sub-agents (multi-agent) định nghĩa bằng file Markdown trong `agents/`.

**Xem agents:**
1. Vào **Agents Management**.
2. Mỗi agent có:
   - **Name**, **ID**, **Description**.
   - **Role**: orchestrator (chủ) hoặc sub (con).
   - **Tools**, **Max iterations**.
   - **Instruction** (system prompt).

**Tạo agent:**
1. Click **New Agent**.
2. Điền:
   - **Filename** (`.md`).
   - **Role**: orchestrator/sub.
   - **Name**, **Description**.
   - **Tools** (danh sách).
   - **Instruction**.
3. Click **Save**.

### 6.6 Roles

**Role** là persona (prompt + tool scope) cho agent.

**Xem roles:**
1. Vào **Roles**.
2. Mỗi role có:
   - **Name**, **Icon**, **Description**.
   - **User prompt** (gợi ý cho user).
   - **Enabled** (bật/tắt).
   - **Tools** (được phép dùng).

**Tạo role:**
1. Click **New Role**.
2. Điền:
   - **Name**.
   - **Icon** (emoji).
   - **Description**.
   - **User prompt**.
   - **Tools** (chọn từ danh sách).
3. Click **Save**.

**Chọn role trong Chat:**
- Click **Role selector** → Chọn role → Tool list tự cập nhật.

---

## 7. Troubleshooting

### 7.1 Không chạy được server

**Triệu chứng:**
- Lỗi khi chạy `./cyberstrike-ai` hoặc `./run.sh`.
- Log: `Failed to load config`, `Failed to open database`.

**Nguyên nhân & xử lý:**
| Lỗi | Nguyên nhân | Xử lý |
|---|---|---|
| `Failed to load config` | `config.yaml` không tồn tại hoặc lỗi YAML. | Kiểm tra `config.yaml` có tồn tại. Dùng `python3 -c "import yaml; yaml.safe_load(open('config.yaml'))"` để validate YAML. |
| `Failed to open database` | `data/` không ghi được hoặc DB corrupt. | Kiểm tra quyền ghi `data/`. Backup & xóa `data/conversations.db` để reset (mất dữ liệu). |
| `port already in use` | Port 8080/8081 đã bị占用. | Đổi port trong `config.yaml` (`server.port`, `mcp.port`). |
| `Go version too old` | Go < 1.21. | Cài Go mới từ https://go.dev/dl/. |

### 7.2 Không gọi được model

**Triệu chứng:**
- Chat không trả lời, hoặc log: `API error`, `401 Unauthorized`.

**Nguyên nhân & xử lý:**
| Lỗi | Nguyên nhân | Xử lý |
|---|---|---|
| `401 Unauthorized` | API key sai/thiếu. | Kiểm tra `ai.channels.<id>.api_key`. Set đúng key (không commit vào git!). |
| `404 Not Found` | Base URL sai. | Kiểm tra `base_url` (phải có `/v1` path cho OpenAI-compatible). |
| `Rate limit` | Vượt quota API. | Chờ hoặc nâng gói API. |
| `Model not found` | Model name sai. | Kiểm tra model name (vd `gpt-4o`, `deepseek-chat`). |

**Test kết nối:**
- Vào **System Settings → AI Channel Configuration** → Click **Test connection** hoặc **Bulk probe**.

### 7.3 Tool not found

**Triệu chứng:**
- Agent báo `Tool not found` hoặc tool không xuất hiện trong danh sách.

**Nguyên nhân & xử lý:**
| Lỗi | Nguyên nhân | Xử lý |
|---|---|---|
| Tool không trong PATH | Công cụ chưa cài. | Cài tool vào PATH (vd `sudo apt install nmap`). |
| Tool YAML lỗi | File `.yaml` trong `tools/` có lỗi. | Kiểm tra cú pháp YAML. Restart server để reload. |
| Role không có tool | Role bị giới hạn tool. | Edit role → Thêm tool vào danh sách. |

### 7.4 MCP lỗi

**Triệu chứng:**
- MCP server không start, hoặc tool không hiển thị.

**Nguyên nhân & xử lý:**
| Lỗi | Nguyên nhân | Xử lý |
|---|---|---|
| `mcp.enabled: false` | MCP chưa bật. | Set `mcp.enabled: true` trong config. |
| Port conflict | Port 8081 đã dùng. | Đổi `mcp.port`. |
| External MCP không kết nối | URL/command sai. | Kiểm tra `external_mcp.servers`. Click **Test connection**. |
| Tool not registered | Tool chưa được MCP server expose. | Kiểm tra MCP server log. Restart MCP server. |

**Xem log MCP:**
- Vào **MCP Monitor** → Xem tool execution logs.
- Terminal: log server có dòng `启动 MCP server`.

### 7.5 Agent không chạy / không gọi tool

**Triệu chứng:**
- Chat không phản hồi, hoặc agent dừng giữa chừng.

**Nguyên nhân & xử lý:**
| Lỗi | Nguyên nhân | Xử lý |
|---|---|---|
| `max_iterations` reached | Agent lặp quá nhiều vòng. | Tăng `agent.max_iterations` hoặc kiểm tra prompt. |
| Tool timeout | Tool chạy quá lâu. | Tăng `agent.tool_timeout_minutes` hoặc kiểm tra tool. |
| HITL pending | Đợi phê duyệt. | Vào **Security → HITL Pending** → Approve/Reject. |
| Model error | Model lỗi/trả về sai định dạng. | Kiểm tra model response trong log. Đổi model/channel. |
| Tool guard blocked | Tool bị chặn bởi rule. | Vào **Security → Tool Guard** → Kiểm tra rule. |

**Xem chi tiết:**
- **MCP Monitor** → Tool execution detail.
- **Audit Logs** → Xem thao tác người dùng.
- Terminal log → Xem error message.

---

## Tài liệu tham khảo

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

---

**End of Guide**
