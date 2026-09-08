# Drift: README vs stonePass UI

Status: Confirmed  
Last Verified: 2026-09-08  

## 真相源

| 层 | 以谁为准 |
|---|---|
| 用户可见 UI | `apps/pass/ui/vue/stonePass` |
| 后端能力 | `apps/pass` 路由注册 |
| README | 宣传/历史，**不能单独当规格** |

## 主要差

当前 UI：**登录（可选）+ 传文件（含过期）+ 传内容（Markdown/历史）**。  
README 仍强调二维码、图文模式、完整管理台、键鼠等；后端 API 仍可能保留对应能力但 UI 未接。

## 建议

改 README 对齐，或恢复 UI；Knowledge 继续以源码为准。
