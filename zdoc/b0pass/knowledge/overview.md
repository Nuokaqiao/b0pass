# Overview

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

---

## 是什么

局域网单机文件共享 / 内容同步工具。主电脑跑一个 Go 进程，其它设备用浏览器访问。

- 对外名：百灵快传 / B0Pass  
- 当前 UI 品牌：**石头记 / StonePass**

## 解决什么

- 手机 ↔ 电脑大文件互传（无网盘）
- 局域网共享目录浏览 / 上传 / 下载（可选过期）
- 跨设备同步文字 / 链接（WebSocket + 可选 Redis 最近历史）

## 技术栈

| 层 | 技术 |
|---|---|
| Go 1.17 / Gin / embed 或 Vite Live | 后端与托管 |
| Vue 3 + Router + Vite | 前端 |
| gorilla/websocket | 传内容 |
| go-redis（可选） | 传内容历史 |
| golang-jwt | 共享口令登录 |
| BurntSushi/toml | `config.ini` |
| robotgo | 仅 Windows 键鼠 API |

无业务数据库。GORM / LRU 代码残留，业务未用。

## 入口

`main/main.go` → `engine.Run("config.ini")` → blank import 注册 pass / docs / gateway → gateway `Gin.Run(ListenAddr)`。

默认 `:8888`，共享根目录 `[pass] Path`（默认 `files`）。

## 模块一览

| 模块 | 职责 |
|---|---|
| pass | 文件 API、静态 `/files`、WS、登录、过期清理 |
| stonePass UI | 传文件 / 传内容 / 登录 |
| engine | App 注册、路由、JWT/CORS、响应 |
| gateway | 监听、Password、非 pass 静态 |
| docs | Markdown 文档站（弱相关） |

## 拓扑

```
Browser → Gin(:ListenAddr)
           ├ /app/pass  UI
           ├ /pass/*    API（Password 非空则 JWT）
           ├ /files/*   静态共享根
           ├ /ws        传内容
           └ /docs/*    docs
         → 本地 Path；可选 Redis
```

## 文档漂移

根目录 README 已按当前 Vue UI 改写（2026-09-08）。官网截图与未接 UI 的后端 API 仍可能不一致，见 [notes/knowledge-drift-readme-vs-ui.md](../notes/knowledge-drift-readme-vs-ui.md)。
