# 百灵快传（B0Pass）/ 石头记（StonePass）

局域网文件与内容互传工具：主电脑跑一个 Go 进程，其它设备用浏览器访问。

当前 Web UI 品牌为 **石头记 / StonePass**（`apps/pass/ui/vue/stonePass`）。

## 1. 当前功能

### 1.1 已支持

- [x] 局域网文件共享（浏览 / 上传 / 下载 / 新建文件夹 / 删除）
- [x] 单可执行文件部署（Go `embed` 打包前端；亦可 `Live` 开发模式）
- [x] 传文件：多选上传、拖拽上传、整夹上传（同一「上传」入口；相对路径落盘）
- [x] 列表 **多选 / 全选**、批量删除
- [x] 上传合计超过 **1 GiB** 时前端提示并二次确认
- [x] 文件可选过期时间（默认不过期；到期自动清理）
- [x] 大文件流式上传（服务端按 Content-Length 分流，避免整包进内存）
- [x] 传内容：WebSocket 实时同步；简易 Markdown 展示；最近历史（≤10 条 / 72h，Redis 可选，否则内存）
- [x] 可选共享口令登录（`[gateway] Password` 非空开启 JWT）
- [x] 自定义监听地址 / 域名；自定义共享根目录
- [x] Windows / Linux / macOS

### 1.2 未做 / 未在当前 UI 暴露

- [ ] 细粒度角色（只读 / 仅上传 / 管理）；现为共享口令全权限或完全开放
- [ ] 上传断点续传、限速
- [ ] PDF 在线预览、压缩包在线解压
- [ ] 自动检查更新
- [ ] 二维码主界面、独立安卓 APK、图文模式管理台等历史能力（部分后端 API 仍在，当前 Vue UI 未接）

> 说明：README 以**当前源码与 stonePass UI**为准；旧截图与发行站文案可能仍展示历史界面。

## 2. 快速使用

1. 准备 `config.ini`（可参考下方示例）
2. 运行可执行文件，或源码：`go run main/main.go`（在仓库根目录）
3. 浏览器打开：`http://127.0.0.1:8888/app/pass/`（端口以配置为准）
4. 同局域网其它设备用主机局域网 IP 访问同一地址

### 配置示例

```ini
[gateway]
ListenAddr = ":8888"       # 监听地址
Domain = ""                # 可选访问域名
Password = ""              # 非空则需登录；空则关闭鉴权

[pass]
Path = "files"             # 共享根目录
Live = false               # true：开发时读 dist / 反代 Vite；false：用 embed 的前端
RedisAddr = ""             # 传内容历史；空或连不上则用内存（重启丢失）
RedisPassword = ""
RedisDB = 0
```

注意：

- 勿安装到需管理员权限的系统目录（如 `C:\Program Files`），否则可能无法写共享目录或配置。
- 改 Go 代码或 `Live=false` 下的 embed 前端后，需**重启**进程；`Live=true` 且读磁盘 `ui/dist` 时，前端 `npm run build` 后强制刷新即可。

### 发行版与仓库

- 官网下载：https://4bit.cn/p/b0pass
- 开源中国：https://www.oschina.net/p/b0pass
- GitHub：https://github.com/bitepeng/b0pass
- Gitee：https://gitee.com/b0cloud/b0pass

## 3. 使用场景

- **手机 ↔ 电脑**：电脑启动服务 → 手机浏览器打开主机地址 → 传文件 / 传内容  
- **电脑 ↔ 电脑 / 虚拟机**：对端浏览器访问同一局域网地址即可  
- **办公室 / 家庭临时共享**：走局域网 HTTP，跨 Windows / macOS / Linux / 手机浏览器

## 4. 源码开发

```bash
git clone https://github.com/bitepeng/b0pass.git
cd b0pass
go mod tidy

# 后端（仓库根目录）
go run main/main.go

# 前端开发（可选：config.ini 中 [pass] Live = true，并起 Vite）
cd apps/pass/ui/vue/stonePass
npm install
npm run dev          # 开发
npm run build        # 产物到 apps/pass/ui/dist/
```

### 用 Make 编译 / 部署

```bash
make check     # 检查 go / node / npm 与前后端依赖，缺什么会列出来
make deps      # 安装缺失依赖（go mod + npm install）
make build     # 检查通过后：前端 build + Go 编译 → dist/b0pass
make deploy    # 同上，并附带 config.ini.example、README.txt
```

交叉编译示例：`make deploy GOOS=linux GOARCH=amd64`

Windows 亦可参考 `main/build.bat` 做打包编译。

## 5. 产品入口一览

| 路径 | 说明 |
|---|---|
| `/app/pass/` | StonePass UI（传文件 / 传内容 / 登录） |
| `/pass/*` | 文件与登录等 API |
| `/files/*` | 共享目录静态访问 |
| `/ws` | 传内容 WebSocket |
