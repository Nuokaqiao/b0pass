# Architecture

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- main/main.go
- core/engine/**
- core/gateway/**
- apps/pass/**
- apps/docs/**

---

## 总体形态

**单体 Go 进程 + 嵌入式 / 开发代理前端**。

基于自研应用框架风格（代码注释称 B0Boot-Go）：各业务以 App 形式注册到 `engine.App`，由引擎按 `config.ini` 分段配置启动。

---

## 启动时序

```
main.main()
  ├─ 若不存在 config.ini → 写入默认配置，创建 files/
  ├─ goroutine: 打印访问地址，cmd.Open 打开浏览器
  └─ engine.Run(config.ini)
        ├─ toml 解码配置到各 App.Config
        ├─ go App.Run()  （pass / docs / gateway 等）
        └─ gateway.Run() 内 Gin.Run(ListenAddr) 阻塞监听
```

Confidence: Confirmed  
Source: `main/main.go`, `core/engine/engine.go`, `core/gateway/gateway.go`

注意：各 App 的 `Run` 用 `go config.Run()` 启动，gateway 内的 `Gin.Run` 在其中一个 goroutine 中阻塞。`main` 末尾 `select {}` 防止主 goroutine 退出。

---

## 模块调用关系

```
                    ┌──────────────┐
                    │    main      │
                    └──────┬───────┘
                           │ engine.Run
           ┌───────────────┼───────────────┐
           ▼               ▼               ▼
      ┌────────┐     ┌──────────┐    ┌─────────┐
      │  pass  │     │ gateway  │    │  docs   │
      └───┬────┘     └────┬─────┘    └────┬────┘
          │               │               │
          │ 注册路由       │ Gin.Run       │ 注册 /docs
          ▼               ▼               ▼
      ┌─────────────────────────────────────┐
      │            engine.Gin               │
      └─────────────────────────────────────┘
          │
   ┌──────┼──────────┬────────────┐
   ▼      ▼          ▼            ▼
 files  stream     chat         keys
 节点/列表 流式上传  WS Hub     键鼠(Win)
```

---

## 数据流动

### 文件路径

1. UI 请求 `/pass/file-list?f=相对路径`
2. Handler 将 `config.Path + f` 交给 `files.GetDirTree`
3. 返回相对 `Path` 的条目列表（name/path/type/size…）
4. 下载：`/pass/file-download` 或静态 `/files{path}`
5. 上传：`POST /pass/file-upload?f=目录` → 写入 `Path + f + filename`

持久化介质：**本地文件系统**，无 DB 元数据表。

### 文本同步

1. 客户端连接 `ws://host/ws`
2. `chat.Hub` 登记 Client
3. 任意 Client 发消息 → Hub.broadcast → 所有 Client（含发送者）

消息不落服务端磁盘；前端用 `localStorage` key `txtdata` 缓存展示记录。

---

## 同步 / 异步

| 场景 | 模式 | 说明 |
|---|---|---|
| HTTP API | 同步请求/响应 | 文件读写在请求 goroutine 内完成 |
| 大文件上传 | 同步流式写盘 | 边读 Body 边写文件，避免整文件进内存 |
| WebSocket | 异步广播 | Hub 单 goroutine + 每连接 read/write pump |
| App 启动 | 异步 | 各 App `go Run()` |
| 打开浏览器 | 异步 | main 中独立 goroutine |

无消息队列、无后台任务调度器。

---

## 入口类型

| 入口 | 路径 / 机制 | 用途 |
|---|---|---|
| HTTP API | `/pass/*` | 文件与命令 API |
| 静态 UI | `/app/pass/` | Vue 前端（Live 时反代 Vite） |
| 静态文件 | `/files/*` | 直接暴露共享目录 |
| WebSocket | `/ws` | 传内容广播 |
| Docs | `/docs/:name` | Markdown → HTML |
| Dev API | `/dev/api`（Debug 时） | 应用元信息 JSON |
| 根路径 | `/` | 跳转到 `/app/pass/` |

**不存在**：独立 Worker、Cron、MQ Consumer。

---

## 前端托管策略

`[pass] Live = true`：
- 尝试反代 `http://127.0.0.1:5173`（Vite）
- 失败则回退 `ui/dist` 磁盘静态目录

`Live = false`：
- `embed` 的 `ui/dist` 通过 `StaticFS` 提供

Source: `apps/pass/main.go` `routeStatic`

Gateway 对 pass 的 UI **跳过**托管（由 pass 自行处理）。

---

## 核心依赖关系

- `pass` → `engine`（路由、响应、CORS）
- `pass` → `files` / `stream` / `chat` / `keys`
- `pass` → `core/tools/cmd`、`nets`
- `gateway` → `engine`（监听、静态路由）
- `docs` → `engine`
- UI → `/pass/*`、`/files`、`/ws`（相对当前 host）

---

## 设计特征（Confirmed）

1. **文件系统即数据模型**：无独立元数据库
2. **App 插件式注册**：blank import + init 安装
3. **单端口聚合**：UI、API、静态文件、WS 同进程
4. **大文件流式上传**：自定义 multipart boundary 解析（非标准 `c.FormFile` 大文件路径）
