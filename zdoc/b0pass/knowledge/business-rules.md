# Business Rules

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- apps/pass/api.go
- apps/pass/lib/files/**
- apps/pass/lib/stream/**
- apps/pass/lib/chat/**
- core/engine/middleware.go

---

## 权限与鉴权

| 规则 | 现状 | Confidence |
|---|---|---|
| 局域网开放访问 | `gateway.Password` 为空时绝大多数 API **无鉴权**；非空时需 JWT | Confirmed |
| JWT | `POST /pass/login` 签发；Header/Query/Cookie `token`；pass API、`/files`、`/ws` 校验 | Confirmed |
| Password 配置 | `[gateway] Password`；空=关闭鉴权 | Confirmed |
| CORS | `Access-Control-Allow-Origin: *` | Confirmed |
| README 规划的 Pread/Pupload/Padmin | 仍未按角色细分；当前为共享口令 | Confirmed（文档/现状） |

**业务含义**：任何能访问端口的客户端可读写共享目录、可发送键鼠命令（若调用 API）、可加入 WS 广播。

---

## 路径与校验

| 规则 | 说明 | Confidence |
|---|---|---|
| 相对路径参数 `f` | 多数 API 直接 `RootPath + f` | Confirmed |
| 空路径拒绝 | rename/delete/content/download 等对空 `f` 返回错误 | Confirmed |
| 路径穿越防护 | **未见** `..` 规范化后必须仍在 Root 内的检查 | Confirmed（缺失） |
| 隐藏文件 | 列表跳过以 `.` 开头的名字 | Confirmed |
| 可执行文件打开 | `cmd-open` 拒绝 `.BAT/.CMD/.EXE` | Confirmed |
| 文本预览大小 | ≥ 2MB 只返回摘要头，不整文件读入 | Confirmed |
| 非文本预览 | `canread=false` 返回 `[[No Preview]]` | Confirmed |

---

## 上传规则

| 规则 | 说明 | Confidence |
|---|---|---|
| 大小分流阈值 | `Content-Length > 4096` → Big 流式；否则 Tiny FormFile | Confirmed |
| 目标目录 | query `f`，默认 `/`；先 `NodeAdd` 确保目录存在 | Confirmed |
| 同名文件 | Big 路径 `os.Create` → 覆盖 | Confirmed |
| 多文件 multipart | Big 路径循环解析多个 boundary 文件 | Confirmed |
| 幂等 | 无 upload-id；重复提交会再次写入/覆盖 | Confirmed |
| 过期（可选） | query `expire` 为 unix 秒；缺省/0 = 不过期 | Confirmed |

---

## 过期删除

| 规则 | 说明 | Confidence |
|---|---|---|
| 默认 | 无过期 | Confirmed |
| 存储 | `{Path}/.b0pass-expire.json`，key 为相对路径 | Confirmed |
| 设置 | `GET /pass/file-expire`；`expire=0` 清除 | Confirmed |
| 清理周期 | 启动一次 + 每分钟；列表前也会 purge | Confirmed |
| 到期行为 | 文件 `Remove`；目录 `RemoveAll` | Confirmed |
| 手动删除 | `node-delete` 同时清元数据 | Confirmed |

---

## 删除 / 重命名

| 规则 | 说明 | Confidence |
|---|---|---|
| 删除 | `os.Remove`：文件可删；非空目录失败 | Confirmed |
| 重命名 | `os.Rename`；空参数拒绝 | Confirmed |
| 前端确认 | 删除前 `confirm` | Confirmed |

---

## WebSocket 规则

| 规则 | 说明 | Confidence |
|---|---|---|
| 最大消息 | `maxMessageSize = 2048000` (~2MB) | Confirmed |
| 心跳 | pingPeriod ≈ 54s，pongWait 60s | Confirmed |
| 广播范围 | 所有已注册 client（含发送者） | Confirmed |
| 慢消费者 | send 阻塞时关闭并踢出该 client | Confirmed |
| 服务端持久化 | 无 | Confirmed |
| 身份 / 房间 | 无用户、无房间；全局一个 Hub | Confirmed |

---

## 唯一性

- 磁盘文件名在同一目录下唯一（OS）
- 无业务层 UUID / 内容去重

---

## 特殊规则

1. 上传分流用 Content-Length，而非文件真实大小字段；multipart 开销可能使「小文件」走 Big 路径。  
   Confidence: Confirmed（阈值逻辑）；实际影响 Inferred
2. Windows 非 Live 模式会尝试释放 `zlib1.dll` 到工作目录。  
   Confidence: Confirmed
3. Gateway `ReadConfig` 把返回中的 Password 设为 `runtime.GOOS`（疑似调试/脱敏行为，意图不完全清楚）。  
   Status: Confirmed 行为；意图 Unknown
