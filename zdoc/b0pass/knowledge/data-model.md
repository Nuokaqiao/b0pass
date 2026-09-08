# Data Model

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- apps/pass/lib/files/**
- apps/pass/api.go
- apps/pass/lib/chat/**
- apps/pass/ui/vue/stonePass/src/views/TextView.vue
- core/engine/database.go

---

## 总览

本项目**没有业务数据库表**。

核心「数据」是：

1. 本地文件系统目录树（`[pass] Path`）
2. WebSocket 内存中的连接与广播消息（不持久化）
3. 浏览器 `localStorage` 中的文本消息缓存（仅前端）

`core/engine/database.go` 提供 GORM `Model` / 分页 / BuildWhere 等基础设施，但当前 pass/docs/gateway **未建立 DB 连接、未使用这些模型**。

---

## 核心实体

### 1. 共享根目录（Root Path）

| 字段 | 来源 | 说明 |
|---|---|---|
| Path | `[pass] Path` | 文件操作根路径 |

所有相对路径 `f` 与根路径拼接后访问磁盘。

Confidence: Confirmed

---

### 2. 文件 / 目录节点（运行时视图）

不是持久化 schema，而是 `GetDirTree` 返回的 map：

| 字段 | 含义 |
|---|---|
| indexs | 序号（字符串） |
| name | 文件/目录名 |
| ext | 扩展名大写，目录为 `DIR` |
| size / sizes | 字节数 / 人类可读大小 |
| date | 修改时间 `MM-DD HH:mm` |
| path | 相对 Root 的路径（`/` 风格） |
| type | `file` / `img` / `dir` / `vod` / `htm` / `pdf` |
| canread | 是否像文本可读 |
| header | 文件头采样 |

目录树 API 另有 `FileNode`：`title` / `path` / `spread` / `children`。

隐藏规则：名称以 `.` 开头的条目在列表中跳过。

Confidence: Confirmed  
Source: `apps/pass/lib/files/file_list.go`, `file_node.go`

---

### 3. 文件过期元数据

路径：`{Path}/.b0pass-expire.json`  
结构：`{ "/rel/path": <unix秒>, ... }`  
缺失或值为 0：不过期。

列表项扩展字段：`expire`、`expireAt`、`expireLeft`。

Confidence: Confirmed  
Source: `apps/pass/lib/files/expire.go`, `apps/pass/api.go`

---

### 4. 上传文件

- Tiny 路径：`c.FormFile("file")` → `SaveUploadedFile`
- Big 路径：解析 multipart → `os.Create(RootPath + FileName)` 流式写入
- 可选 query `expire`：写入过期元数据

无独立上传会话表、无分片断点续传元数据、无校验和字段。

Confidence: Confirmed

---

### 5. WebSocket Client / Hub（内存 + 可选 Redis 历史）

| 结构 | 字段 |
|---|---|
| Hub | clients map、broadcast/register/unregister、history store |
| Client | hub、conn、send buffer(256) |
| History | Redis ZSET 或内存；最多 10 条且 72 小时内 |

实时消息：原始文本字节广播。  
历史：`GET /pass/text-history`；写入发生在每次广播时。

Confidence: Confirmed  
Source: `apps/pass/lib/chat/**`

---

### 6. 前端文本消息缓存

localStorage key: `txtdata`  
结构：`[{ key, val, time }, ...]`

仅影响本机浏览器展示，不代表服务端真相。

Confidence: Confirmed

---

### 7. 配置实体

| Section | 关键字段 |
|---|---|
| gateway | ListenAddr, Domain, Password, Live, Debug |
| pass | Live, Path |
| docs | Live |

`Password` 存在于配置，但当前**未见鉴权使用**（ReadConfig 返回时甚至用 `runtime.GOOS` 覆盖 Password 字段）。

Confidence: Confirmed（存在与未用于鉴权）  
Source: `main/config.ini`, `core/gateway/gateway.go`

---

## 实体关系

```
Root Path
  ├── .b0pass-expire.json（过期 map）
  └── File/Dir nodes (树形，磁盘父子关系)
        └── 静态 URL /files{path}
        └── API 相对路径 f
        └── 可选 expire unix

Hub (1) ──< Client (N)
  └── broadcast fan-out（无持久消息实体）
```

---

## 创建 / 修改 / 删除

| 操作 | API | 实现 |
|---|---|---|
| 创建目录 | `node-add` 路径以 `/` 结尾 | `MkdirAll` |
| 创建空文件 | `node-add` 路径无尾 `/` | `WriteFile` 空内容 |
| 重命名 | `node-rename` | `os.Rename` |
| 删除 | `node-delete` | `os.Remove`（非递归）+ 清 expire 元数据 |
| 上传 | `file-upload` | 创建/覆盖同名文件；可选 expire |
| 设过期 | `file-expire` | 写 `.b0pass-expire.json` |
| 过期清理 | 后台 / 列表前 | 到期删除文件或目录 |
| 下载 | `file-download` / `/files` | 读盘 |

---

## 一致性要求

- **无跨服务一致性问题**（单机文件系统）
- 多客户端并发写同一路径：依赖 OS 文件系统语义；**无应用层锁**
- 上传与列表之间：上传完成后客户端主动刷新列表；无服务端推送文件变更事件

Confidence: Confirmed（无锁）；并发行为细节 Inferred 为「最后写入覆盖」

---

## 索引 / 约束

- 无 DB 索引
- 唯一性：同一目录下文件名由 OS 约束；上传同名会覆盖（`os.Create`）
- 路径约束：**缺少**「必须落在 Root Path 内」的规范化校验（见 risks）

---

## 状态字段

文件节点无业务状态机字段（无 draft/published 之类）。  
类型 `type` 由扩展名推断，不是可变状态。
