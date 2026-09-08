# Module: gateway

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

路径：`core/gateway/**`。`APP_HOOK`：开端口、注入 Password、托管非 pass 静态。

## Config

`ListenAddr`, `Domain`, `Password`, `Live`, `Debug`  
Run 开头：`engine.SetAuthPassword(Password)`，并设置 `engine.Addr` / `Domain`。

## 行为

- `GET /gateway/config`：返回配置（Password 字段被换成 `GOOS`）
- Debug：`/dev/api`、`/dev/doc`
- 静态 `/app/{name}`：**跳过 pass**
- `Gin.Run(ListenAddr)`
