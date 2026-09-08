# Open Questions / Unknowns

Status: Note（非权威）  
Date: 2026-09-08  

---

## Unknown

1. **gateway.Password 的预期用途**：配置存在，但无校验逻辑；ReadConfig 用 GOOS 覆盖返回值的意图不明。
2. **JWT 完整登录/发牌流程**：有 Parse/Middleware，未见签发 API；是否半成品。
3. **GORM / database.go 是否为框架遗留**：是否计划后续接入，或仅从 B0Boot 模板残留。
4. **LRU cache 是否有外部调用方**：仓库内无引用。
5. **发行版是否仍带旧 UI**：本 Knowledge 仅验证当前 git 树；下载站二进制可能不同。
6. **安卓 Uniapp 客户端**与当前后端的协议是否仍兼容（README 提及，本仓库未见对应源码）。
7. **上传 Content-Length 4096 阈值**是否为历史经验值，有无文档说明。
8. **Domain 配置**在反向代理/HTTPS 下的完整推荐部署方式（代码仅用于展示 URL）。
9. **多文件 multipart 与浏览器 XHR 单文件 FormData**在 Big/Tiny 分流下的边界用例覆盖（测试不足）。
10. **jobs/cron**：确认无；若外部用 systemd timer 包装，属部署层 Unknown。

---

## 建议下一步深入

见 `index.md` 的 Recommended Next Deep-Dives。
