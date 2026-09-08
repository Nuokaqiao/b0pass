# b0pass Project Knowledge

## Project

- Name: b0pass（百灵快传 / 石头记 StonePass）
- Source: `/Users/lmc10304/go/src/my_project/b0pass`
- Protocol: [protocol.md](./protocol.md)
- Config: [project.yaml](./project.yaml)

## Verification

- Last Sync: 2026-09-08
- Verified Commit: `9000df3`（对照当前源码整理；修正鉴权 / 传内容历史等漂移）

## Knowledge Health

以源码为准。各文档头部 `Verified Commit` 可能略旧；冲突时信任当前实现并更新 Knowledge。

| Area | Note |
|---|---|
| Overview / Architecture / Flows / Rules | 已含登录 JWT、文件过期、Redis 文本历史、audit |
| State machines | 刻意极简（无业务 FSM） |
| Jobs and Workers | N/A，不建空文档 |

## Knowledge Documents

| Doc | 写什么（避免重复） |
|---|---|
| [overview](./knowledge/overview.md) | 产品定位、栈、入口、模块一览 |
| [architecture](./knowledge/architecture.md) | 启动、调用关系、数据流、托管模式 |
| [modules](./knowledge/modules.md) | 模块地图；细节见 `modules/*` |
| [data-model](./knowledge/data-model.md) | 文件系统 / 过期元数据 / WS·Redis 历史 / 配置 |
| [business-flows](./knowledge/business-flows.md) | 用户可感知主流程 |
| [business-rules](./knowledge/business-rules.md) | 鉴权、路径、上传、过期、WS 约束 |
| [state-machines](./knowledge/state-machines.md) | 仅连接/文件存在生命周期 |
| [integrations](./knowledge/integrations.md) | OS / Vite / Redis；无云依赖 |
| [reliability](./knowledge/reliability.md) | 幂等、并发、失败与续传边界 |
| [risks](./knowledge/risks.md) | 仍存在的安全与可靠性风险 |

## Modules

| Module | Doc |
|---|---|
| pass | [modules/pass.md](./modules/pass.md) |
| stonePass UI | [modules/stonepass-ui.md](./modules/stonepass-ui.md) |
| engine | [modules/engine.md](./modules/engine.md) |
| gateway | [modules/gateway.md](./modules/gateway.md) |

## Decisions / Changes / Notes

| 类型 | 入口 |
|---|---|
| Decisions | [DEC-0001 Vue 精简](./decisions/0001-vue-frontend-simplify.md) |
| Changes | [changes/](./changes/)（expire / login / redis-history / knowledge-init…） |
| Notes | [README 漂移](./notes/knowledge-drift-readme-vs-ui.md)、[未决问题](./notes/open-questions.md) |

## Known Gaps（仍开放）

1. 路径穿越：多数 API 仍缺「规范化后必须在 Root 内」校验
2. 共享口令无角色细分（README 的 Pread/Pupload/Padmin 未落地）
3. 上传无断点续传；中断可能留半截文件
4. README / 发行版宣传与当前 Vue UI 仍可能不一致
5. Uniapp 手机端源码不在本仓库
6. 核心路径自动化测试不足

## Recommended Next Deep-Dives

1. 路径穿越加固与上传半截文件清理
2. 角色权限（若需要超出共享口令）
3. 上传断点续传（若产品需要）
4. 对齐或改写 README

## How to use

复杂任务：本 index → 相关 Knowledge / modules → 回源码校验 → 增量更新。不要复制整份 Knowledge 到别处。
