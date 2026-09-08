# Business Rules

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

---

## 鉴权

| 规则 | 说明 |
|---|---|
| 开关 | `[gateway] Password` 空 = 关闭；非空 = 需登录 |
| 登录 | `POST /pass/login` 校验共享口令，发 JWT |
| 携带 | Header / Query / Cookie 名均为 `token` |
| 范围 | pass 受保护路由、`/files`、`/ws`；`auth-status`/`login` 公开 |
| CORS | `*` |
| 角色 | 无；非 README 规划的 Pread/Pupload/Padmin |
| 审计 | `[audit]` 日志含 ip / device / ua |

Password 开启时，能登录者仍拥有完整读写与 cmd API 能力（共享口令模型）。

## 路径

| 规则 | 说明 |
|---|---|
| 拼接 | 多为 `Path + f` |
| 空 f | 若干写/读接口拒绝 |
| 穿越 | **缺**「Clean 后必须仍在 Root」硬校验（高风险） |
| 隐藏 | 列表跳过 `.` 开头名 |
| cmd-open | 拒 `.BAT/.CMD/.EXE` |

## 上传 / 下载

| 规则 | 说明 |
|---|---|
| 分流 | Content-Length > 4096 → Big 流式，否则 Tiny |
| 同名 | 覆盖 |
| 断点上传 | **无** |
| 断点下载 | HTTP Range（`c.File` / StaticFS）可用 |
| 限速 | **无**应用层限速 |
| 过期 | query `expire` unix；0/缺省不过期 |

## 过期删除

默认不过期。元数据 `.b0pass-expire.json`。清理：启动 + 1m + 列表前。到期删文件或 `RemoveAll` 目录。手动删同步清元数据。

## 删除 / 重命名

`Remove` 非递归（非空目录失败）。前端删前 `confirm`。

## 传内容 / WS

| 规则 | 说明 |
|---|---|
| 上限 | 单消息约 2MB |
| 心跳 | ping≈54s / pongWait 60s |
| 范围 | 全局单 Hub，含回显发送者 |
| 历史 | ≤10 条且 ≤72h（Redis，否则内存） |
| 房间/用户 | 无 |

## 其它

- 上传分流看 Content-Length，multipart 开销可能让「小文件」走 Big。  
- Windows 非 Live 可能释放 `zlib1.dll`。
