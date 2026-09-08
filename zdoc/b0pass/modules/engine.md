# Module: engine

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

路径：`core/engine/**`

## 职责

App 安装、读 config、共享 Gin、统一 JSON、`SetAuthPassword` / JWT / CORS / PathAuth。

## 要点

- `AppInstall` + `Run(config)` → `go App.Run()`
- 路由：`/{appId}{url}`
- 响应：`{ code, msg, data }`（0 / 400 / 401）
- JWT：Password 空则 `EnsureAuth` 放行；否则校验 token
- `PathAuthMiddleware`：给 `/files`、`/ws` 等前缀用
- `database.go`：GORM 辅助，**业务未连库**

## 风险

JWT Secret 硬编码；CORS `*`。
