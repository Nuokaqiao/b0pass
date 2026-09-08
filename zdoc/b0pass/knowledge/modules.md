# Modules

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

API 明细与 UI 行为分别只维护在 [modules/pass.md](../modules/pass.md)、[modules/stonepass-ui.md](../modules/stonepass-ui.md)，本节不重复罗列。

## 地图

| 模块 | 文档 | 说明 |
|---|---|---|
| pass | [pass.md](../modules/pass.md) | 文件、WS、登录、过期、audit |
| stonePass UI | [stonepass-ui.md](../modules/stonepass-ui.md) | 登录 / 传文件 / 传内容 |
| engine | [engine.md](../modules/engine.md) | 注册、Gin、JWT/CORS、响应 |
| gateway | [gateway.md](../modules/gateway.md) | Listen、Password、开端口 |
| docs | （无独立深文） | `GET /docs/:name` Markdown 站 |
| tools | （无独立深文） | nets / cmd 常用；cache 无引用 |

## UI 实际调用的后端

登录与鉴权：`/pass/auth-status`、`/pass/login`  
文件：`file-list`、`node-add`、`node-delete`、`file-upload`、`file-expire`、`file-download`、`/files/*`  
内容：`/ws`、`/pass/text-history`

后端另有 `node-tree` / `node-rename` / `file-content` / `cmd-*` 等，**当前 UI 未用**。
