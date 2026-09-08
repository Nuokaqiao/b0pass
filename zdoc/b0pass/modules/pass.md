# Module: pass

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- apps/pass/**

---

## Purpose

局域网文件共享与内容同步的核心业务 App（appId = `pass`）。

---

## Responsibilities

- 托管 / 代理前端 UI（`/app/pass`）
- 暴露共享目录静态资源（`/files`）
- 文件列表、增删改名、上传下载 API
- WebSocket 文本广播（`/ws`）
- 主电脑打开文件、键鼠命令（平台相关）

---

## Entry Points

| 入口 | 位置 |
|---|---|
| `init` / `run` | `apps/pass/main.go` |
| HTTP handlers | `apps/pass/api.go` |
| WS | `routeWs` → `lib/chat` |
| 流式上传 | `lib/stream/upload.go` |
| 文件节点 | `lib/files/*` |
| 键鼠 | `lib/keys/*` |

---

## Important APIs

前缀：`/pass`

| Method | Path | Auth | Handler |
|---|---|---|---|
| GET | /auth-status | 无 | AuthStatus（`enabled`） |
| POST | /login | 无 | Login（共享口令 → JWT） |
| GET | /ping | JWT* | Ping |
| GET | /read-config | JWT* | ReadConfig |
| GET | /read-ip | JWT* | ReadIP |
| GET | /cmd-open | JWT* | CmdOpen |
| GET | /cmd-key | JWT* | CmdKey |
| GET | /node-tree | JWT* | NodeTree |
| GET | /node-add | JWT* | NodeAdd |
| GET | /node-rename | JWT* | NodeRename |
| GET | /node-delete | JWT* | NodeRemove |
| GET | /file-count | JWT* | FileCount |
| GET | /file-list | JWT* | FileList（含 expire / expireAt / expireLeft） |
| GET | /file-expire | JWT* | FileExpire（`f` + `expire` unix；0=清除） |
| GET | /file-content | JWT* | FileContent |
| GET | /file-download | JWT* | FileDownload |
| GET | /text-history | JWT* | TextHistory（最近≤10 条且≤72h） |
| POST | /file-upload | JWT* | FileUpload（可选 query `expire=unix秒`） |
| GET | /ws | JWT* | ServeWs |

\* `gateway.Password` 为空时 JWT 中间件自动放行。Token 可来自 Header / Query / Cookie `token`。

另：`GET /files/*` 静态（同样鉴权）；`GET /` 跳转 UI。

---

## Config

```go
type AppConfig struct {
  Live bool
  Path string // 文件根目录
  RedisAddr, RedisPassword string
  RedisDB int
}
```

传内容历史：Redis ZSET `stonepass:text:history`；`RedisAddr` 为空或连不上时回退进程内存。
---

## Dependencies

- `core/engine`
- `lib/files`, `lib/stream`, `lib/chat`, `lib/keys`
- `core/tools/cmd`, `core/tools/nets`
- embed `ui/dist`

---

## Business Flows

见 `knowledge/business-flows.md` Flow 2–7。

---

## Reliability / Risks

见 `knowledge/reliability.md`、`knowledge/risks.md`（鉴权、路径穿越、并发上传尤为关键）。

---

## Subpackages

| 包 | 职责 |
|---|---|
| files | 列表、树、读写、节点 CRUD |
| stream | multipart 流式解析与写盘 |
| chat | Hub/Client WebSocket 广播 |
| keys | Windows robotgo；其它平台 noop |
