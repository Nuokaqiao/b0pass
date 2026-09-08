# DEC-0001: Vue 重建前端并精简为传文件 / 传内容

Decision ID: DEC-0001  
Status: Accepted（已体现在当前代码）  
Date: 2026-09-08  

Related Knowledge:
- knowledge/overview.md
- modules/stonepass-ui.md
- notes/knowledge-drift-readme-vs-ui.md

---

## 问题

旧前端功能面较广（README 可见图文模式、二维码、图片浏览等），需要可维护的现代前端并收敛产品能力。

---

## 上下文

Commit `de3640e`：`refactor(pass): 用 Vue 重建石头记前端，精简为传文件与传内容`

---

## 选择

- 使用 Vue 3 + Vite + Vue Router
- 首页仅保留「传文件」「传内容」
- 构建产物进入 `apps/pass/ui/dist`，由 Go embed / Live 代理托管

---

## 为何选择（Confirmed vs Inferred）

| 点 | Confidence |
|---|---|
| 使用 Vue 重建并精简两个功能 | Confirmed（commit + 代码） |
| 精简的产品/商业动机细节 | Unknown（无更多设计文档） |
| 旧 UI 是否完全删除 | Inferred：当前源码树以 stonePass 为主；历史 UI 细节未完整比对 |

---

## 后果

- 后端若干 API 暂时无官方 UI 入口，但仍暴露
- README / 截图与真实 UI 不一致（Knowledge Drift）
- 开发体验依赖 Vite Live 代理或直连

---

## 被拒绝的方案

Unknown（无记录）
