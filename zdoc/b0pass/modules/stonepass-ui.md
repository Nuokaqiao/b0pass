# Module: stonePass UI

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: working-tree（相对 de3640e 有未提交 UI 优化）  

Related Paths:
- apps/pass/ui/vue/stonePass/**
- apps/pass/ui/dist/**

---

## Purpose

当前用户可见前端「石头记 / StonePass」：精简为 **传文件** 与 **传内容** 两个入口。

源自提交：`refactor(pass): 用 Vue 重建石头记前端，精简为传文件与传内容`。

---

## Stack

- Vue 3 + Vue Router（hash history）
- Vite 5
- 无状态管理库；API 用 fetch / XHR

---

## Routes

| Hash 路由 | 视图 | 功能 |
|---|---|---|
| `/login` | LoginView | 共享口令登录（`gateway.Password` 非空时） |
| `/` | HomeView | 两大入口：传文件 / 传内容；可退出 |
| `/files` | FilesView | 目录浏览、拖拽/选择上传、下载、复制链接、新建文件夹、删除、过期 |
| `/text` | TextView | WebSocket 内容同步（剪贴板语义） |

未登录（鉴权开启时）整站路由守卫拦截至 `/login`。请求头 / Cookie / Query 携带 `token`。

Base：生产构建 `base: '/app/pass/'`，产物输出到 `apps/pass/ui/dist/`。

---

## FilesView（传文件）

Confirmed（当前源码）：

- 面包屑导航；目录优先排序
- 拖拽到页面上传；多文件顺序上传 + 进度条
- 上传：先选文件 → 再设过期（数字 + 分钟/小时/天，空=不过期）→ 开始上传；拖拽同流程
- 卡片列表：下载（主操作）、过期设置（同单位）、复制下载链接、删除；有过期时显示剩余时间
- 点文件名不下载；仅文件夹可点进目录
- 新建文件夹用页内表单（非 `prompt`）
- 图片缩略图走 `/files...`

---

## TextView（传内容）

Confirmed（当前源码）：

- 布局：最近内容在上，输入区在下；消息时间正序（新在下）
- 文案：同步 / 最近内容（非聊天语义）
- 简易 Markdown 渲染展示（`marked` + `DOMPurify`）；复制为原文
- 卡片：点击复制原文、URL「打开链接」
- WebSocket 自动重连 + 状态胶囊；进入时拉服务端历史
- 本地 `localStorage.txtdata` 作缓存兜底

---

## API 使用面

仅：

- `/pass/file-list`
- `/pass/node-add`
- `/pass/node-delete`
- `/pass/file-upload`（可选 query `expire=unix秒`）
- `/pass/file-expire`（设置/清除过期；`expire=0` 不过期）
- `/pass/file-download`
- `/files...`
- `/ws`

---

## Dev vs Prod

- Dev：`npm run dev` + 后端 `[pass] Live=true` 反代，或直连 Vite（自带 proxy）
- Prod：`npm run build` → embed；`Live=false`

---

## Notes

- 无登录页、无二维码页、无键鼠控制 UI（相对旧产品）
- README / 旧截图与当前 UI 仍可能漂移（见 notes/knowledge-drift-readme-vs-ui.md）
