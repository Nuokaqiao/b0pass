# Module: stonePass UI

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

路径：`apps/pass/ui/vue/stonePass/` → build 到 `apps/pass/ui/dist/`（`base: '/app/pass/'`）。

## 路由

| Hash | 视图 |
|---|---|
| `/login` | 共享口令（Password 开启时） |
| `/` | 传文件 / 传内容入口；可退出 |
| `/files` | 传文件 |
| `/text` | 传内容 |

鉴权开启时整站守卫；token 走 Header/Cookie/Query。

## 传文件

面包屑；上传文件 / 上传文件夹（`webkitdirectory`）；拖拽文件或目录（`webkitGetAsEntry` 递归）。
相对路径落到当前目录下子路径；上传前弹窗设过期；所选合计 **>1 GiB** 时警告并可二次确认。
列表支持单选勾选、**全选**、取消选择、批量删除。
进度条、下载/复制链接/过期/删除；点文件名不下载。空目录不会单独创建（仅随文件建父目录）。

## 传内容

历史在上、输入在下；新消息在下。进页拉 `text-history`；WS 自动重连；简易 Markdown 展示；复制原文。

## 调用的 API

`auth-status` / `login` / `file-list` / `node-add` / `node-delete` / `file-upload` / `file-expire` / `file-download` / `text-history` / `/files` / `/ws`

## Dev / Prod

Live+Vite 反代，或 `npm run build` + embed（`Live=false` 改 UI 后需重启进程）。

相对旧产品：无二维码、无键鼠 UI。README 可能漂移，见 notes。
