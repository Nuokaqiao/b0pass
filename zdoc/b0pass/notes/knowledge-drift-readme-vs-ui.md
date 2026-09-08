# Knowledge Drift: README / 旧产品描述 vs 当前 stonePass UI

Status: Confirmed（存在漂移）  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

---

## 当前代码是什么（Confirmed）

Vue stonePass 仅提供：

1. 传文件：列表、面包屑、上传、下载、新建目录、删除
2. 传内容：WebSocket + localStorage

后端仍注册更多 API（树、重命名、内容预览、cmd-open、cmd-key、ping…），但 UI 未使用。

---

## Knowledge / 文档原来记录的是什么

根目录 `README.md` 描述并勾选了更广功能，例如：

- 二维码扫码界面
- 共享文件在线管理（图片浏览器、重命名等）
- 图文模式 / 列表模式
- 手机端丰富操作截图
- 键盘鼠标命令（作为已支持能力叙述）

这些内容反映的是**更完整的历史产品形态**或发行版宣传能力，与当前仓库内 Vue UI 不一致。

---

## 差异在哪里

| 能力 | README | 当前 Vue UI | 后端 API |
|---|---|---|---|
| 传文件上传下载 | 有 | 有 | 有 |
| 传内容 WS | 未作为主叙事 | 有 | 有 |
| 二维码 | 有 | 无 | Unknown / 未见 |
| 图片浏览器 | 有 | 仅新窗口/静态图 | 静态 `/files` |
| 重命名 | 有 | 无 | 有 |
| 键鼠控制 | 有 | 无 | 有（Win） |

---

## 哪一个是已确认事实

- **运行时用户界面真相**：以 `apps/pass/ui/vue/stonePass` + 当前 embed dist 为准（Confirmed）
- **后端能力真相**：以 `apps/pass/main.go` / `api.go` 注册路由为准（Confirmed）
- **README**：宣传/历史文档，**不能单独当作当前实现规格**（Confirmed drift）

---

## 建议

1. 后续决定：更新 README 对齐 stonePass，或恢复 UI 能力
2. 在同步 Knowledge 时继续以源码为准，不把 README 勾选框当 Confirmed 实现
