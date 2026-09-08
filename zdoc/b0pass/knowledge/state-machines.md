# State Machines

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- apps/pass/lib/chat/**
- apps/pass/ui/vue/stonePass/src/views/TextView.vue

---

## 结论

本项目**没有**订单/发布/审批类的业务状态机。

文件没有 `status` 字段流转；文件系统条目的存在即状态。

以下记录与「连接生命周期」相关的有限状态，避免误以为存在复杂 FSM。

---

## WebSocket 连接生命周期（客户端视角）

状态（UI `status` 文案）：

| 状态 | 触发 |
|---|---|
| 连接中… | 页面挂载，开始 `new WebSocket` |
| 已连接 | `onopen` |
| 连接已断开 | `onclose` |
| 连接异常 | `onerror` |
| 未连接，无法发送 | 发送时 readyState ≠ OPEN |
| 已复制 | 复制操作反馈（非连接状态） |

合法转换（Inferred from UI code）：

```
连接中 → 已连接 → 连接已断开
连接中 → 连接异常
任意发送失败 → 未连接，无法发送（文案）
```

**谁负责**：浏览器 WebSocket + 服务端 Hub register/unregister。

**Retry / Recovery**：前端**无自动重连**逻辑。  
Confidence: Confirmed

---

## Hub 内 Client 生命周期

```
Upgrade 成功
  → register
  → readPump / writePump 运行
  → 读失败或 Hub 踢出
  → unregister + close send + close conn
```

异常：
- write 缓冲区阻塞 → Hub 关闭该 client
- 非预期 close → 打日志后退出 readPump

无「半连接恢复」或消息确认状态。

Confidence: Confirmed

---

## 文件节点

无状态枚举。若需要描述「生命周期」，仅为：

```
不存在 →（upload/add）→ 存在 →（rename）→ 存在(新名) →（delete）→ 不存在
```

无中间态（uploading/ready）持久化。上传过程中文件可能短暂以不完整内容存在于磁盘（流式写入未完成时）。  
Confidence: Confirmed（无中间态元数据）；不完整文件可见性 Inferred
