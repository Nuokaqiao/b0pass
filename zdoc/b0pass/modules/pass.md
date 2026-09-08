# Module: pass

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

路径：`apps/pass/**`。appId=`pass`。

## 职责

UI 托管、`/files`、文件 CRUD/上传下载/过期、WS+历史、登录鉴权挂载、audit、cmd API。

## 主要 API（`/pass`）

公开：`GET /auth-status`、`POST /login`  
JWT*：`file-*`、`node-*`、`text-history`、`read-*`、`cmd-*`、`ping` 等  

\* Password 空则中间件放行。Token：Header/Query/Cookie。  
另：`/files/*`、`/ws` 同样鉴权；`/` → UI。

完整表不必在多处复制；以 `main.go` `routeApi` 为准。

## Config

`Live`, `Path`, `RedisAddr`, `RedisPassword`, `RedisDB`

历史 key：`stonepass:text:history`。

## 子包

`files`（含 expire）、`stream`、`chat`（hub/history）、`keys`、`audit`

流程 / 规则 / 风险 → knowledge 对应文档。
