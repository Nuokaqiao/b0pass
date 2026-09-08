# Overview

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- README.md
- main/main.go
- apps/pass/**
- core/**
- go.mod

---

## 项目是什么

**百灵快传（B0Pass）** 是一个局域网大文件传输 / 共享工具。

运行方式：在一台「主电脑」上启动单个 Go 可执行程序，同一局域网内的手机或其它电脑通过浏览器访问 HTTP 服务，进行文件互传与文本同步。

产品对外名称：百灵快传 / B0Pass。  
当前前端品牌文案：**石头记 / StonePass**（Vue 前端）。

Source:
- README.md
- apps/pass/ui/vue/stonePass/src/views/HomeView.vue

---

## 解决的业务问题

- 手机与电脑之间传大文件（无需网盘、数据线）
- 电脑与电脑、虚拟机与宿主机之间共享文件
- 局域网内快速同步文字 / 链接（WebSocket「传内容」）
- 把主电脑某个目录暴露为可浏览、可上传、可下载的共享空间

---

## 技术栈

| 层 | 技术 | Confidence |
|---|---|---|
| 语言 / 运行时 | Go 1.17（module: `b0go`） | Confirmed |
| HTTP 框架 | Gin | Confirmed |
| 前端 | Vue 3 + Vue Router + Vite | Confirmed |
| 实时通信 | gorilla/websocket | Confirmed |
| 配置解析 | BurntSushi/toml（读取 `config.ini` 风格分段配置） | Confirmed |
| 压缩 | gin-contrib/gzip | Confirmed |
| Windows 键鼠控制 | robotgo（仅 Windows build tag） | Confirmed |
| JWT 库 | golang-jwt/jwt（仅 `/pass/ping` 使用中间件） | Confirmed |
| GORM | 在 `core/engine/database.go` 提供工具类型；**当前业务未连接数据库** | Confirmed |

---

## 项目入口

1. 进程入口：`main/main.go`
2. 启动引擎：`engine.Run("config.ini")`
3. 通过 blank import 注册应用：
   - `_ "b0go/apps/pass"`
   - `_ "b0go/apps/docs"`
   - `_ "b0go/core/gateway"`
4. 各 App 在 `init()` 中 `engine.AppInstall`，在 `Run` 中按配置调用各自 `Run()`
5. Gateway 调用 `engine.Gin.Run(ListenAddr)` 真正监听端口

默认监听：`:8888`  
默认文件根目录：`files`（可由 `[pass] Path` 覆盖）

---

## 主要模块（高层）

| 模块 | 路径 | 职责 |
|---|---|---|
| main | `main/` | 进程入口、缺省配置生成、启动后打开浏览器 |
| engine | `core/engine/` | 应用安装、路由注册、Gin 引擎、JWT/CORS、统一响应 |
| gateway | `core/gateway/` | 监听地址、Debug API 文档、其它 App 静态资源托管 |
| pass | `apps/pass/` | 核心业务：文件管理、上传下载、WebSocket 传内容、键鼠命令 |
| docs | `apps/docs/` | Markdown 文档渲染应用 |
| stonePass UI | `apps/pass/ui/vue/stonePass/` | 当前用户界面（传文件 / 传内容） |
| tools | `core/tools/` | 文件、网络、命令打开、LRU cache 等通用工具 |

---

## 外部依赖 / 基础设施

| 类型 | 现状 | Confidence |
|---|---|---|
| 数据库 | **未使用**（无业务路径无 DB Open） | Confirmed |
| Cache | LRU 实现存在于 `core/tools/cache`，**业务未引用** | Confirmed |
| Message Queue | **无** | Confirmed |
| 第三方云服务 | **无**（纯局域网本地服务） | Confirmed |
| 静态文件存储 | 本地文件系统目录 `[pass] Path` | Confirmed |

---

## 服务关系

单进程、单 HTTP 服务：

```
Browser / Phone
      │
      ▼
┌─────────────────────────────┐
│  Gin (gateway ListenAddr)   │
│  ├─ /app/pass  → Vue UI     │
│  ├─ /pass/*    → pass API   │
│  ├─ /files/*   → 静态共享目录 │
│  ├─ /ws        → 文本广播   │
│  └─ /docs/*    → docs app   │
└─────────────────────────────┘
      │
      ▼
 本地目录 config.Path
```

无微服务拆分；无独立 Worker / Cron / Consumer 进程。

---

## 产品现状 vs README 描述

README 仍描述较完整的旧产品能力（二维码、图片浏览器、列表/图文模式、手机端操作截图等）。

当前 Vue 前端（最近一次重构）**仅保留**：
- 传文件（列表、上传、下载、新建目录、删除、可选过期自动删除）
- 传内容（WebSocket 文本同步）

后端仍保留部分旧 API（目录树、重命名、文件内容、cmd-open、cmd-key 等），但当前 stonePass UI **未调用**它们。

详见：`notes/knowledge-drift-readme-vs-ui.md`
