# Reliability

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- apps/pass/api.go
- apps/pass/lib/stream/upload.go
- apps/pass/lib/chat/**
- core/engine/middleware.go
- core/engine/engine.go

---

## Idempotency（幂等）

| 操作 | 是否幂等 | 说明 |
|---|---|---|
| 上传同名文件 | 否（覆盖写） | 重复请求会再次 Create/写入 |
| node-add 已存在目录 | 基本可重复 | `MkdirAll` |
| node-add 已存在文件 | 覆盖为空文件风险 | `WriteFile` 空内容 |
| node-delete | 第二次失败 | 文件已不存在 |
| WS 发送 | 非幂等 | 每次广播一条 |

无 upload token / 请求去重。

---

## Concurrency / Race

- 多客户端同时上传同名文件：无文件锁 → 内容可能交错损坏或互相覆盖。  
  Confidence: Inferred（无锁是 Confirmed；损坏模式 Inferred）
- Hub 的 clients map 仅在 Hub 单 goroutine 内访问 → 注册表本身安全。  
  Confidence: Confirmed
- 列表与删除并发：可能列出已删文件或删失败，属常规 FS 竞态。

---

## Transaction Boundary

- 无数据库事务
- 单文件上传是单次流式写；**多文件 multipart** 中途失败时：已写完的文件保留，后续失败返回错误（部分成功）。  
  Confidence: Confirmed

---

## Retry

- HTTP API：无服务端自动重试
- 前端上传：失败则 toast，无自动重试
- WebSocket：无自动重连

---

## Message Loss / Ordering（WS）

- 广播是 best-effort；慢客户端被踢掉时，其未发送队列丢失
- 无消息 ID、无 ACK、无持久化 → 断线即丢在途消息
- 同一连接上 writePump 会合并队列消息（换行拼接），顺序保持发送顺序

Confidence: Confirmed

---

## Large Data

- 大文件：流式解析 multipart，避免整包进内存（设计目标 Confirmed，见 stream 包注释）
- 文本预览限制 2MB
- WS 单消息上限约 2MB
- `GetCounts` 递归遍历目录：大目录可能阻塞/耗 CPU（当前 UI 未调用）

---

## Failure Recovery

- 进程崩溃：已写入磁盘的文件保留；进行中的上传可能残留不完整文件
- 无上传会话恢复 / 断点续传
- 无崩溃后清理临时文件的逻辑被发现

---

## Backward Compatibility

- API 路径以 `/pass/...` 为主，长期风格稳定（Inferred）
- 前端从旧 UI 迁到 Vue stonePass：功能表面积缩小（Confirmed by commit message + 代码）
- embed UI 与 Live 代理双模式并存

---

## Exactly-once

不提供 exactly-once。上传与 WS 均为 at-most/at-least 的朴素语义，需调用方自行约束。
