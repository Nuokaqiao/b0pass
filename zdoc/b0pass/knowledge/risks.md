# Risks

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- apps/pass/api.go
- apps/pass/main.go
- apps/pass/lib/files/**
- core/engine/middleware.go
- core/gateway/gateway.go
- README.md

---

## R1. 无有效访问控制（高）

**现象**：文件读写、删除、上传、静态 `/files`、WebSocket、甚至 `cmd-key` 在默认情况下对能访问端口的人均开放。

**证据**：除 `/pass/ping` 外 API 未挂 JWT；Password 未参与鉴权。

**影响**：局域网恶意用户可读写任意共享目录内容，Windows 上可远程模拟键鼠。

README 也将「提升安全性 / JWT」列为未完成。

---

## R2. 路径穿越（高）

**现象**：多数 handler 使用 `config.Path + f`，未确认最终路径仍位于 Root 下。

**影响**：`f=../../...` 可能读写共享根之外的文件（取决于 OS 与 Path 写法）。

`FileContent` 使用了 `path.Clean`，但仍无 Root 前缀校验。

Confidence: Confirmed（缺少校验）；可利用性依赖部署路径（Inferred）

---

## R3. 静态目录整树暴露（高）

`/files` → `http.Dir(config.Path)` 直接挂载共享根。

即使不用 API，也可按 URL 枚举/下载文件（若知道路径；目录列表行为取决于 http.FileServer）。

---

## R4. 硬编码 JWT Secret（中）

`TokenSecret = []byte("01xda2d8f6x9n4x8")` 写死在源码。

当前仅 ping 使用，但若未来扩大鉴权范围，此 secret 不安全。

---

## R5. 并发上传同名文件（中）

无锁、无临时文件+rename 提交。并发或中断可能导致损坏/半截文件。

---

## R6. 删除非递归（低-中）

`os.Remove` 无法删非空目录；用户可能认为「删除文件夹」会递归清理。

---

## R7. cmd-key / cmd-open 攻击面（高，若 API 可达）

当前 Vue UI 未调用，但 API 仍注册。攻击者可直接请求。

非 Windows 上 key 为空操作；open 仍可打开文件。

---

## R8. WebSocket 无来源严格策略（中）

使用默认 Upgrader；跨站场景行为依赖 gorilla 默认 CheckOrigin。  
全局单 Hub，无隔离，任意连接可向所有人广播。

---

## R9. 产品文档与实现漂移（中，维护风险）

README 功能列表 / 截图与当前 stonePass UI 不一致，易误导用户与后续 AI/开发者。  
详见 `notes/knowledge-drift-readme-vs-ui.md`。

---

## R10. 配置中的 Password 明文（低-中）

`config.ini` 可含 Password；且未见加密存储或实际校验用途，易造成「已有密码保护」的错误安全感。

---

## R11. 大目录性能（低-中）

`GetDirTree` / `GetCounts` / `NodeTree` 同步扫盘；超大目录可能拖慢请求。当前主路径 UI 用 file-list（单层 Glob），相对可控。

---

## R12. 共享目录权限位（中）— 部分已修

**现象**：`files` 目录曾为 `drw-r--r--`（无执行位），Unix 下无法创建/写入文件，上传表现为失败或「成功但列表仍空」。

**根因（Confirmed）**：
- 工作区 `files/` 目录 mode 异常
- 旧代码 `NodeAdd` 使用 `os.MkdirAll(..., 0666)`，新建目录也会缺 `+x`

**已做**：本地 `chmod 755 files`；`NodeAdd` 改为 `0755`；`FileUploadTiny` 检查保存错误，避免假成功。

---

## 不存在 / 已排除的风险类型

| 项 | 说明 |
|---|---|
| MQ 丢消息 | 无 MQ |
| 分布式事务 | 单机 |
| 数据迁移脚本 | 无 DB schema |
| 多版本 API 兼容矩阵 | Unknown / 未建立 |
