# CyberStrikeAI — Hướng dẫn Docker hóa & Sử dụng (portable)

> Hướng dẫn **tổng quát, chạy được trên bất kỳ máy nào**. Chỉ cần clone repo về,
> chạy `docker compose`, không phụ thuộc đường dẫn/người dùng. Không cần đụng đến
> Go, Python, hoặc bất kỳ cài đặt thủ công nào khác.
>
> Lộ trình nhanh: **clone → prepare → compose up** (3 bước, xem mục 3).

- [Authorization](#important-authorization)
- [1. Điều kiện tiên quyết](#1-điều-kiện-tiên-quyết)
- [2. Cấu trúc các file Docker](#2-cấu-trúc-các-file-docker)
- [3. Bật hệ thống (3 bước)](#3-bật-hệ-thống-3-bước)
- [4. Thiết lập sau khi bật](#4-thiết-lập-sau-khi-bật)
- [5. Cách dùng các tính năng chính](#5-cách-dùng-các-tính-năng-chính)
- [6. Cấu hình nâng cao](#6-cấu-hình-nâng-cao)
- [7. Gỡ lỗi & khắc phục sự cố](#7-gỡ-lỗi--khắc-phục-sự-cố)
- [8. Bảo mật & Backup](#8-bảo-mật--backup)

---

## Important Authorization

> [!IMPORTANT]
> **CyberStrikeAI là hệ thống offensive-security quyền cao.** Nó có thể thực thi lệnh,
> gọi MCP tool, quản lý WebShell, và chạy C2 listener.
> **Chỉ sử dụng trên hệ thống bạn sở hữu hoặc được cấp phép rõ ràng.** Không expose
> port ra Internet công cộng. Xem `docs/en-US/security-model.md` và `docs/en-US/security-hardening.md`.

---

## 1. Điều kiện tiên quyết

| Thành phần | Yêu cầu |
|---|---|
| Docker Engine | 20.10+ (hoặc Docker Desktop; đã kiểm tra với Docker 28) |
| Docker Compose | v2 trở lên (`docker compose`) |
| Git | Để `git clone` và chạy `./docker/prepare.sh` |
| Kết nối Internet | Pull base image + tải Go modules (qua `goproxy.cn`) + pip packages |

Không cần cài Go, Python, SQLite riêng — tất cả nằm trong image. Không bắt buộc WSL: trên
Windows nhìn chung dùng Docker được cả từ PowerShell/CMD lẫn từ WSL (xem mục 3b).

> **Lưu ý về antivirus (Windows Defender/EDR):**
> Repo chứa mã C2/payload hợp lệ bị một số phần mềm diệt virus gắn cờ nhầm là malware và
> **tự động xóa** (vd `internal/c2/payload_oneliner.go`, `knowledge_base/SQL Injection/*.md`).
> Điều này khiến build lỗi `undefined: c2.OnelinerKind`. Bước "prepare" (mục 3a/3b) sẽ khôi
> phục các file đó từ git bằng một lệnh tuỳ nền tảng. Nếu file vẫn bị xóa, thêm repo vào danh
> sách loại trừ của antivirus rồi chạy lại bước prepare.

---

## 2. Cấu trúc các file Docker

| File | Vai trò |
|---|---|
| `Dockerfile` | Multi-stage build: biên dịch Go binary (cÓ cgo/SQLite), dựng Python venv, copy `web/ tools/ skills/ roles/ agents/ knowledge_base/` vào image chạy ở `/app`. |
| `docker/entrypoint.sh` | Khi container khởi động: tạo `config.yaml` từ example (nếu chưa có), set password admin từ env `CYBERSTRIKE_ADMIN_PASSWORD`, rồi start server. |
| `docker/prepare.sh` | **Portable helper (Linux/macOS/WSL)**: khôi phục các file bị antivirus xóa (false positive) để build luôn chạy được. Chạy trước khi build. |
| `docker/prepare.ps1` | **Portable helper (Windows thuần, không WSL)**: tương đương `prepare.sh`, chạy bằng PowerShell. |
| `docker/prepare.bat` | **Bộ khởi chạy Windows**: gọi `prepare.ps1` với quyền bypass execution-policy — chạy/double-click trực tiếp, không cần cài đặt gì. |
| `docker-compose.yml` | Khai báo build, port 7123/7134, volume `./data`, env, healthcheck. |
| `.dockerignore` | Loại bỏ khỏi build context: `data/`, `venv/`, `.git`, `docs/`, secret files,... |
| `cmd/set-admin-password/` | Helper Go set password admin **không tương tác** (lệnh `--reset-admin-password` gốc chỉ chạy tương tác, không dùng được trong container). |

**Port:**
- `7123` — Web UI chính (HTTPS self-signed theo mặc định).
- `7134` — MCP server (đã bật mặc định trong `config.example.yaml`; bind `0.0.0.0` để truy cập được từ host qua port map).

---

## 3. Bật hệ thống (3 bước)

Docker chạy ổn định trên **Linux/macOS, WSL, và Windows thuần** (không WSL). Chỉ khác nhau ở
**bước "prepare"** (khôi phục file bị antivirus xóa) — chọn đúng theo hệ điều hành bên dưới.

Chung: `--build` build image lần đầu (vài phút do tải Go deps + pip); `CYBERSTRIKE_ADMIN_PASSWORD`
là password admin mặc định để **khởi tạo tài khoản `admin`** lần đầu và **reset** mỗi lần bật
(nếu env vẫn set); dữ liệu lưu vào thư mục **`./data`** trên host, persist qua restart.

### 3a. Linux / macOS / WSL (dùng bash)

```bash
# 1) Clone repo
git clone https://github.com/Ed1s0nZ/CyberStrikeAI.git
cd CyberStrikeAI

# 2) Khôi phục file bị antivirus gắn cờ nhầm (chạy lại nhiều lần được)
./docker/prepare.sh

# 3) Build & bật với password admin mặc định
CYBERSTRIKE_ADMIN_PASSWORD='admin123' docker compose up -d --build
```

### 3b. Windows thuần — không WSL (dùng PowerShell / CMD)

Trên Windows thuần, chỉ cần **Docker Desktop** (bản thân nó dùng WSL2 bên dưới nhưng **người
dùng không phải gõ lệnh trong WSL**). Mọi lệnh dưới đây gõ trực tiếp trong **PowerShell** hoặc
**CMD**; không cần mở WSL, không cần touch file `.sh`.

```powershell
# 1) Clone repo (PowerShell)
git clone https://github.com/Ed1s0nZ/CyberStrikeAI.git
cd CyberStrikeAI

# 2) Khôi phục file bị antivirus gắn cờ nhầm —
#    dùng prepare.bat (không lo execution-policy, double-click cũng được)
.\docker\prepare.bat
#    hoặc gõ trực tiếp bằng PowerShell:
#    powershell -ExecutionPolicy Bypass -File .\docker\prepare.ps1

# 3) Build & bật với password admin mặc định (PowerShell)
$env:CYBERSTRIKE_ADMIN_PASSWORD='admin123'
docker compose up -d --build
```

> Git Bash dùng được luôn cả `./docker/prepare.sh` — nhưng trên Windows thuần, `prepare.bat`
> là cách đơn giản và đáng tin nhất (không phụ thuộc shell/execution policy).

Kiểm tra:

```bash
docker compose ps
docker compose logs -f cyberstrike-ai
```

Thấy dòng `● ONLINE   https://127.0.0.1:7123/` là đã bật xong.
Mở **`https://localhost:7123/`** — chấp nhận cảnh báo chứng chỉ self-signed 1 lần.

### Không đặt password (dùng cơ chế tự sinh)

Chung cho cả 3a và 3b — chỉ cần bỏ env `CYBERSTRIKE_ADMIN_PASSWORD`:

```bash
docker compose up -d --build
```

Trường hợp này container giữ nguyên mặc định của hệ thống: lần đầu **tự sinh password**
cho `admin` và in ra log container:

```bash
# Linux/macOS/WSL
docker compose logs cyberstrike-ai | grep -i password
# Windows (PowerShell)
docker compose logs cyberstrike-ai | findstr /i password
```

### Các lệnh hữu ích khác

Các lệnh `docker compose` này giống hệt trên mọi nền tảng (PowerShell/CMD/WSL):

```bash
docker compose logs -f cyberstrike-ai   # log real-time
docker compose down                     # dừng, giữ dữ liệu trong ./data
docker compose down -v                  # xóa cả volume (CẨN THẬN: mất dữ liệu)
docker compose up -d                    # tái dùng image đã build, không build lại
```

---

## 4. Thiết lập sau khi bật

### 4a. Đăng nhập
Mở `https://localhost:7123/` → đăng nhập `admin` / mật khẩu bạn đã đặt. **Đổi password ngay**.
Nếu quên: set lại `CYBERSTRIKE_ADMIN_PASSWORD` rồi restart container (entrypoint sẽ reset).

### 4b. Cấu hình AI Channel (bắt buộc trước khi dùng AI)
**System Settings → Basic Settings → AI Channel Configuration**:
- **Provider**: `openai_compatible` (OpenAI/DeepSeek/Qwen...) hoặc `claude` (Anthropic).
- **Base URL**: vd `https://api.openai.com/v1`, `https://api.deepseek.com/v1`, `https://dashscope.aliyuncs.com/compatible-mode/v1`.
- **API key**, **Model** (`gpt-4o`, `deepseek-chat`, `qwen3-max`...), token limits.
- **Save changes**, set channel mặc định, test probe.

> Cách 2: sửa trực tiếp `config.yaml` (mặc định có placeholder `qwen-max` với `api_key: sk-xxxxxxx` —
> **nhớ thay key thật**). Mẫu xem `config.example.yaml`.

### 4c. Cài công cụ bảo mật (tùy chọn)
Hệ thống **không tự cài** `nmap`, `sqlmap`... — `tools/*.yaml` chỉ mô tả lệnh. Muốn dùng tool phải cài:

```bash
docker compose exec cyberstrike-ai bash
apt-get update && apt-get install -y nmap masscan sqlmap nikto gobuster hydra hashcat john
exit
```

> Re-create container sẽ mất gói cài thủ công. Cách bền vững: mở rộng `Dockerfile`
> (thêm `RUN apt-get install ...`), hoặc mount binaries từ host.

---

## 5. Cách dùng các tính năng chính

### 5a. Chat / Agent AI
Nhập yêu cầu, ví dụ:
- `Scan open ports on 192.168.1.1`
- `Check if https://example.com/page?id=1 is vulnerable to SQL injection`
- `Enumerate subdomains for example.com, then run nuclei against the results`

Chọn **Agent mode**: `eino_single`, `deep`, `plan_execute`, `supervisor`. Xem `agents/` cho các
role con (`recon`, `penetration`, `post-exploitation`, `cleanup-rollback`...).

### 5b. Projects & Attack Chains
- **Project**: gom nhiều hội thoại/phiên về cùng mục tiêu.
- **Attack Chain**: nối sự kiện liên phiên, chấm rủi ro, xem đồ thị, **replay từng bước**.

### 5c. Asset Management
Chuẩn hóa & khử trùng domain/IP/port/service; import/export XLSX/CSV; bộ lọc nâng cao;
gán chủ sở hữu; tích hợp **FOFA / ZoomEye / Quake / Shodan** (backend proxy che key).

### 5d. Vulnerability Management
Phân loại nghiêm trọng, vòng đời, lọc, thống kê; liên kết lỗ hổng với asset/project.

### 5e. Tools & Skills
- **Tools**: 100+ recipe trong `tools/` (nmap, sqlmap, nuclei, metasploit, hashcat, bloodhound, impacket...), truy cập theo role.
- **Skills**: chuẩn Agent Skills, thư mục `skills/` (recon, post-exploitation, active-directory-attack...), load theo nhu cầu.
- **Roles**: prompt + chính sách tool theo tình huống trong `roles/` (CTF, Web, AD, Cloud...).

### 5f. MCP Integration
Hỗ trợ MCP **HTTP, stdio, SSE, federation**, dynamic tool discovery. MCP server (HTTP) **đã bật
mặc định** và lắng nghe port `7134`; cấu hình thêm qua `external_mcp.servers`.

### 5g. Quản lý người dùng & RBAC
**Platform permissions → User management**: tạo user, gán role (`管理员`, `操作员`, `审计员`, `只读用户`).
Mọi thay đổi ghi **audit log**.

### 5h. Human-in-the-Loop (HITL)
Bật phê duyệt (`approval`) hoặc `review_edit` để con người duyệt tool call trước khi thực thi;
whitelist tool; **Call blocking** trong **Security → Call blocking** (mặc định chặn domain chính phủ).

### 5i. WebShell & C2 (nâng cao)
> Chỉ dùng cho hệ thống bạn sở hữu / được ủy quyền.
- **WebShell**: quản lý kết nối, terminal ảo, thao tác file, workflow AI.
- **C2**: listener, beacon mã hóa, session, task queue, payload helper, sự kiện real-time.
- **Không expose ra Internet.**

---

## 6. Cấu hình nâng cao

### Biến môi trường
| Biến | Ý nghĩa |
|---|---|
| `CYBERSTRIKE_ADMIN_PASSWORD` | Set/reset password admin khi bật container. |
| `CYBERSTRIKE_HTTPS` | `1`/`true`/`yes` → bật HTTPS. |
| `GOPROXY` | Proxy Go khi **build** (mặc định `goproxy.cn,direct`). |

### Command tùy chỉnh
Sửa `command` trong `docker-compose.yml`, vd `command: ["--http"]` để dùng HTTP thay vì HTTPS.

### Persist config qua Web UI
Thêm mount cho `config.yaml` (sau khi đã có file `config.yaml` ở host):

```yaml
    volumes:
      - ./data:/app/data
      - ./config.yaml:/app/config.yaml:ro
```

---

## 7. Gỡ lỗi & khắc phục sự cố

| Triệu chứng | Nguyên nhân & xử lý |
|---|---|
| Build lỗi `undefined: c2.OnelinerKind / GenerateOneliner` | File bị antivirus xóa. Chạy `./docker/prepare.sh` (Linux/WSL) hoặc `.\docker\prepare.bat` (Windows), rồi build lại. Nếu vẫn mất → thêm repo vào danh sách loại trừ antivirus. |
| Build lỗi tải Go deps (`stream error`) | Mạng tới `proxy.golang.org` chập chờn. Đã mặc định `goproxy.cn`; override: `docker compose build --build-arg GOPROXY=https://goproxy.cn,direct` |
| Trình duyệt báo "không an toàn" | Chứng chỉ tự ký — **Advanced → Proceed**, hoặc dùng `--http`. |
| Log `client sent an HTTP request to an HTTPS server` | Đang mở `http://` trên cổng TLS. Dùng `https://`. |
| Login ra `密码不能为空` | Body JSON lỗi (thiếu `Content-Type`). Server nhận password là OK: đúng → token, sai → 401 (đã test). |
| Không thấy dòng `ONLINE` | Chờ thêm hoặc xem `docker compose logs cyberstrike-ai` để tìm lỗi cấu hình/process isolation. |
| Dữ liệu biến mất sau `down` | Dữ liệu nằm trong `./data` (mount) — chỉ mất nếu `down -v` hoặc xóa thư mục. |

---

## 8. Bảo mật & Backup

- **Không** expose port 7123/7134 ra Internet công cộng; nếu cần truy cập xa, đặt sau reverse
  proxy TLS + auth, hoặc dùng VPN/SSH tunnel.
- Backup trước khi nâng cấp: `config.yaml`, `data/`, `tools/`, `skills/`, `roles/`, `agents/`, `knowledge_base/`.
- Không commit secret; dùng biến môi trường cho API key.
- Password admin tối thiểu 8 ký tự (do hệ thống yêu cầu).
- Xem thêm: `docs/en-US/security-hardening.md`, `docs/en-US/security-model.md`, `docs/en-US/runbooks.md`, `docs/en-US/troubleshooting.md`.
