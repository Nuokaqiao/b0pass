# Module: engine

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- core/engine/**

---

## Purpose

应用运行时框架：安装 App、加载配置、共享 Gin 实例、统一 JSON 响应与中间件。

---

## Key Concepts

- `App map[string]*AppConfig`：已安装应用
- App 类型常量：`APP_HOOK` / `APP_APP` 等
- `AppInstall`：记录 Name/Config/UIFS/Run/Dir
- `Router`：注册 `/ {appId}{url}` 并记入 App.Router 元数据
- `Run(configFile)`：读配置、toml 解码到各 App.Config、`go App.Run()`

---

## Response Convention

```json
{ "code": 0|400|401, "msg": "...", "data": ... }
```

`OK` / `ERR` / `JSON` / `PAGE` 辅助函数。

---

## Middleware

- Gzip + Recovery（引擎 init）
- `CorsMiddleware`（pass run 时挂载）
- `JWTMiddleware`（按路由选择性挂载）

---

## Database helpers

`database.go` 提供 GORM Model / Paginate / BuildWhere。  
**当前无业务调用链连接到真实 DB。**

---

## Risks for consumers

- JWT secret 硬编码
- CORS 允许任意 Origin
