# Vue 前端（stonePass）

## 开发

```bash
# 终端 1：后端（config.ini 中 [pass] Live = true）
# 终端 2：
cd apps/pass/ui/vue/stonePass
npm install
npm run dev
```

访问后端地址（如 `http://127.0.0.1:8888/app/pass/`），Live 模式下会代理到 Vite。

也可直接打开 Vite：`http://127.0.0.1:5173/`（已代理 `/pass` `/files` `/ws`）。

## 生产构建

```bash
cd apps/pass/ui/vue/stonePass
npm run build
```

产物输出到 `apps/pass/ui/dist/`，由 Go `embed` 打包。将 `[pass] Live = false`。

## 页面

- `/` 首页：大按钮「传文件」「传内容」
- `/files` 传文件
- `/text` 传内容（WebSocket）
