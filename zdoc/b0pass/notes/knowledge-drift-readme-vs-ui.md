# Drift: README vs stonePass UI

Status: Confirmed  
Last Verified: 2026-09-08  

## 真相源

| 层 | 以谁为准 |
|---|---|
| 用户可见 UI | `apps/pass/ui/vue/stonePass` |
| 后端能力 | `apps/pass` 路由注册 |
| README | 已于 2026-09-08 按当前 UI 改写；发行站旧截图仍可能漂移 |

## 已对齐

根目录 README 功能清单、配置项、开发命令已对应当前登录 / 传文件（夹上传、全选、过期、1G 提示）/ 传内容能力；历史二维码、图文模式等标为「未在当前 UI 暴露」。

## 仍可能漂移

- 官网 / 开源中国页面截图与文案  
- 后端仍注册但 UI 未用的 API（如 cmd-key、部分预览相关）
