# Business Flows

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

规则细节见 [business-rules.md](./business-rules.md)。

---

## 1. 启动

运行 `main` → 加载 config → Apps Run → 监听 `ListenAddr` → 可访问 `/app/pass/`。

## 2. 登录（Password 非空时）

路由守卫 → `#/login` → `POST /pass/login` → JWT（约 24h）存 localStorage + Cookie → API/`/files`/`/ws` 带 token。  
Password 为空：无登录页，中间件放行。

## 3. 传文件

浏览：`file-list`；进目录；下载 `file-download`（及缩略图 `/files`，需 token）。  
上传：选文件 → 弹窗设过期（可空）→ `file-upload` → 刷新。  
新建/删除：`node-add` / `node-delete`；过期：`file-expire` 或上传时带 expire。  
自动清理：启动 + 每分钟 + 列表前 purge（≤72h 规则在过期元数据侧）。

## 4. 传内容

进页拉 `text-history` → 连 `/ws`（可自动重连）→ 同步广播全员（含自己）→ History 落 Redis/内存。  
展示支持简易 Markdown；复制为原文。

## 5. 后端有、UI 未用

`cmd-key`（Win 键鼠）、`cmd-open`、目录树/重命名/内容预览等。勿当当前产品主路径。
