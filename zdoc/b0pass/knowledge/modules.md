# Modules

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- apps/**
- core/**

---

## 模块地图

| 模块 | 文档 | 类型 | 说明 |
|---|---|---|---|
| pass | [modules/pass.md](../modules/pass.md) | 核心业务 App | 文件共享、上传、WS、键鼠 |
| stonePass UI | [modules/stonepass-ui.md](../modules/stonepass-ui.md) | 前端 | 传文件 / 传内容 |
| engine | [modules/engine.md](../modules/engine.md) | 框架 | App 注册、Gin、中间件 |
| gateway | [modules/gateway.md](../modules/gateway.md) | 钩子 App | 监听与静态托管 |
| docs | （本节简述） | 附属 App | Markdown 文档站 |
| tools | （本节简述） | 工具库 | nets / files / cmd / cache |

---

## docs（简述）

路径：`apps/docs/`

- 将 `docs/*.md` 渲染为 HTML（`show.html` 模板）
- 路由：`GET /docs/:name`（`-` 替换为路径分隔）
- 静态：`/docs_root` → `docs/`
- 与 pass 核心传文件业务**弱相关**

Confidence: Confirmed

---

## tools（简述）

路径：`core/tools/`

| 包 | 用途 | 业务使用情况 |
|---|---|---|
| `nets` | 出网 IP、HTTP client | pass `ReadIP`、main 打印地址 |
| `files`（core） | PathExists 等 | main 检查 config |
| `cmd` | 系统打开 URL/文件 | main 开浏览器、pass `CmdOpen` |
| `cache` | LRU | **未被业务引用** |
| `encode.go` / `datastruct.go` | 通用工具 | Unknown（未深入） |

---

## 当前前端实际依赖的后端能力

Confirmed（来自 `apps/pass/ui/vue/stonePass/src/api/pass.js` + TextView）：

- `GET /pass/file-list`
- `GET /pass/node-delete`
- `GET /pass/node-add`
- `POST /pass/file-upload`
- `GET /pass/file-download`
- `GET /files/*`
- `WS /ws`

后端存在但当前 UI 未调用：

- `/pass/ping`（JWT）
- `/pass/read-config`、`/pass/read-ip`
- `/pass/cmd-open`、`/pass/cmd-key`
- `/pass/node-tree`、`/pass/node-rename`
- `/pass/file-count`、`/pass/file-content`
