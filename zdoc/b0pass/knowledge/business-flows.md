# Business Flows

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- main/main.go
- apps/pass/**
- apps/pass/ui/vue/stonePass/**

---

## Flow 1：启动主电脑服务

**参与者**：用户（主电脑）

1. 用户运行可执行文件 / `go run main/main.go`
2. 若无 `config.ini`，写入默认配置并创建 `files/`
3. 引擎加载 gateway / pass / docs
4. 监听 `ListenAddr`（默认 `:8888`）
5. 控制台打印局域网访问地址；尝试用系统浏览器打开

**结果**：局域网内可通过 `http://{ip}:{port}/app/pass/` 访问。

Confidence: Confirmed

---

## Flow 2：传文件 — 浏览与下载

**参与者**：浏览器用户（局域网或公网部署均可）

1. 打开首页 → 点「传文件」
2. 前端 `GET /pass/file-list?f=/` 加载根目录
3. 进入子目录：再次 file-list
4. 点击文件：`window.open(/pass/file-download?f=...)`  
   或图片预览使用 `/files{path}`
5. 「复制链接」复制同源 download URL（Confirmed：当前 FilesView）

**数据变化**：无服务端写操作；仅读盘。

Confidence: Confirmed

---

## Flow 3：传文件 — 上传

1. 用户在当前目录选择文件，或将文件拖到页面（可多选，顺序上传）
2. `POST /pass/file-upload?f={当前目录}/[&expire={unix秒}]`，FormData 字段 `file`（默认无 expire = 不过期）
3. 服务端确保目录存在（`NodeAdd`）
4. Content-Length > 4096 → 流式大文件路径；否则 Tiny `FormFile` 路径
5. 文件写入 `Path + 目录 + 文件名`；若带 expire 则写入 `{Path}/.b0pass-expire.json`
6. 前端刷新列表

**结果**：共享目录出现新文件，其它设备可通过列表/静态 URL 访问。

Confidence: Confirmed

---

## Flow 4：传文件 — 新建目录 / 删除

- 新建：`GET /pass/node-add?f={dir}/{name}/`
- 删除：确认后 `GET /pass/node-delete?f={path}` → `os.Remove`；并清除过期元数据
- 改过期：`GET /pass/file-expire?f={path}&expire={unix|0}`

**边界**：非空目录手动删除会失败（`os.Remove` 非递归）。过期清理对目录使用 `RemoveAll`。

Confidence: Confirmed

---

## Flow 5：传文件 — 过期自动删除

1. 上传时选过期，或列表点「过期」设置 unix 截止时间（0 = 清除/不过期）
2. 元数据存于共享根目录 `.b0pass-expire.json`（隐藏，列表不展示）
3. 进程启动时清理一次，之后每分钟 `PurgeExpired`
4. 到期则删除文件（目录则 `RemoveAll`）并去掉元数据条目
5. 列表接口在返回前也会先 purge，并附带 `expire` / `expireAt` / `expireLeft`

Confidence: Confirmed

---

## Flow 6：传内容 — 跨设备内容同步

1. 用户打开「传内容」（剪贴板语义：同步 / 最近内容）
2. 浏览器连接 `ws(s)://{host}/ws`；断开后前端自动重连
3. 输入内容点「同步」→ Hub 广播给所有在线连接
4. 各端 `onmessage` 追加到「最近内容」列表（新在上），并写入 `localStorage.txtdata`（约 100 条上限）
5. 用户可复制内容；若含 URL 可「打开链接」

**结果**：在线设备实时看到同一条内容；离线后服务端不保留历史。

**注意**：发送端也会收到自己的广播（Hub 向所有 client 发送）。

Confidence: Confirmed

---

## Flow 7：主电脑远程键鼠（后端能力，当前 UI 未用）

1. 调用 `GET /pass/cmd-key?k=...`
2. Windows：`robotgo` 模拟按键/鼠标
3. 非 Windows：空实现

Confidence: Confirmed（API 存在）；产品是否仍对外宣传：见 README drift

---

## Flow 8：在主电脑打开文件（后端能力，当前 UI 未用）

1. `GET /pass/cmd-open?f=相对路径`
2. 禁止 `.BAT/.CMD/.EXE`
3. `cmd.Open(RootPath + f)` 调系统打开

Confidence: Confirmed
