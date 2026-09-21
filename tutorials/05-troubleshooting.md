# 5. Troubleshooting

## 5.1 Server không khởi động được

**Triệu chứng:** Lỗi khi chạy `./cyberstrike-ai` hoặc `./run.sh`.

| Lỗi | Nguyên nhân | Xử lý |
|---|---|---|
| `Failed to load config` | `config.yaml` không tồn tại hoặc lỗi YAML. | Kiểm tra file tồn tại. Validate: `python3 -c "import yaml; yaml.safe_load(open('config.yaml'))"` |
| `Failed to open database` | `data/` không ghi được hoặc DB corrupt. | Kiểm tra quyền ghi `data/`. Backup & xóa `data/conversations.db` để reset. |
| `port already in use` | Port 7123/7134 đã bị chiếm. | Đổi port trong `config.yaml` (`server.port`, `mcp.port`). |

---

## 5.2 Lỗi Model/API

**Triệu chứng:** Chat không trả lời, log: `API error`, `401 Unauthorized`.

| Lỗi | Nguyên nhân | Xử lý |
|---|---|---|
| `401 Unauthorized` | API key sai/thiếu. | Kiểm tra `ai.channels.<id>.api_key` trong `config.yaml`. |
| `404 Not Found` | Base URL sai. | Kiểm tra `base_url` (phải có `/v1` path cho OpenAI-compatible). |
| `Rate limit` | Vượt quota API. | Chờ hoặc nâng gói API. |
| `Model not found` | Model name sai. | Kiểm tra model name (vd `MiniMax/MiniMax-M3`, `gpt-4o`). |

**Test kết nối:**
- Vào **System Settings → AI Channel Configuration** → Click **Test connection**.

---

## 5.3 Tool không tìm thấy

**Triệu chứng:** Agent báo `Tool not found` hoặc tool không xuất hiện trong danh sách.

| Lỗi | Nguyên nhân | Xử lý |
|---|---|---|
| Tool chưa cài | Công cụ chưa có trong PATH. | Cài tool: `sudo apt install nmap sqlmap gobuster`. |
| Tool YAML lỗi | File `.yaml` trong `tools/` có lỗi cú pháp. | Kiểm tra YAML syntax. Restart server để reload. |
| Role không có tool | Role bị giới hạn tool scope. | Edit role → Thêm tool vào danh sách. |

---

## 5.4 MCP connection failed

**Triệu chứng:** MCP server không start, tool không hiển thị.

| Lỗi | Nguyên nhân | Xử lý |
|---|---|---|
| `mcp.enabled: false` | MCP chưa bật trong config. | Set `mcp.enabled: true` trong `config.yaml`. |
| Port conflict | Port 7134 đã dùng. | Đổi `mcp.port` trong config. |
| External MCP không kết nối | URL/command sai. | Kiểm tra `external_mcp.servers`. Click **Test connection** trong MCP Management. |

**Xem log MCP:**
- Vào **MCP Monitor** → Xem tool execution logs.

---

## 5.5 Agent/HITL bị stuck

**Triệu chứng:** Chat không phản hồi, agent dừng giữa chừng, hoặc task đợi phê duyệt.

| Lỗi | Nguyên nhân | Xử lý |
|---|---|---|
| `max_iterations` reached | Agent lặp quá nhiều vòng. | Tăng `agent.max_iterations` trong config hoặc kiểm tra prompt. |
| Tool timeout | Tool chạy quá lâu. | Tăng `agent.tool_timeout_minutes` trong config. |
| HITL pending | Đợi phê duyệt. | Vào **Human-in-the-loop** → Approve/Reject. |
| Model error | Model trả về sai định dạng. | Kiểm tra model response trong log. Đổi model/channel. |

**Xem chi tiết:**
- **MCP Monitor** → Tool execution detail.
- **Audit Logs** → Xem thao tác người dùng.
- Terminal log → Xem error message.

---

## Hỗ trợ

Nếu gặp vấn đề không được giải quyết trong tài liệu này:

1. **Kiểm tra logs:**
   - Terminal: xem log khi chạy server.
   - **System Settings → Audit Logs** → Xem thao tác.
   - **MCP Monitor** → Xem tool execution.

2. **Tìm kiếm trong docs:**
   - `docs/en-US/troubleshooting.md` — Troubleshooting chi tiết.
   - `docs/en-US/configuration.md` — Configuration reference.

3. **Community:**
   - GitHub Issues: https://github.com/Ed1s0nZ/CyberStrikeAI/issues

---

**End of User Guide**
