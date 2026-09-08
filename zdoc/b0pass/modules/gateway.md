# Module: gateway

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- core/gateway/**

---

## Purpose

钩子型 App（`APP_HOOK`）：加载监听配置，托管非 pass 应用的静态 UI，并调用 `Gin.Run` 启动 HTTP 服务。

---

## Config

```go
ListenAddr string
Domain     string
Password   string
Live       bool
Debug      bool
```

启动时写入 `engine.Addr` / `engine.Domain`。

---

## Behavior

- `GET /gateway/config`：返回配置（Password 字段被替换为 `runtime.GOOS`）
- Debug 时：`/dev/api` 返回各 App 元信息；`/dev/doc` 静态文档
- 静态 `/app/{name}`：遍历 `engine.App`，**跳过 pass**（pass 自管 UI）
- 最后 `engine.Gin.Run(ListenAddr)`

---

## Relationship

必须与 pass/docs 等同进程启动；真正「开端口」的是 gateway。
