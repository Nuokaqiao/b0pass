# b0pass Project Knowledge

## Project

- Name: b0pass（百灵快传 / 石头记 StonePass）
- Source: `/Users/lmc10304/go/src/my_project/b0pass`
- Protocol: [protocol.md](./protocol.md)
- Config: [project.yaml](./project.yaml)

## Verification

- Last Sync: 2026-09-08
- Verified Commit: `de3640e9dd7cf65ba64d606ae6a8021807b4abe3`
- Note: 首次初始化；对照当前 git HEAD 源码建立。

## Knowledge Health

| Area | Status | Verified Commit |
|---|---|---|
| Overview | Fresh | de3640e |
| Architecture | Fresh | de3640e |
| Modules | Fresh | de3640e |
| Data Model | Fresh | de3640e |
| Business Flows | Fresh | de3640e |
| Business Rules | Fresh | de3640e |
| State Machines | Fresh | de3640e |
| Integrations | Fresh | de3640e |
| Jobs and Workers | N/A（无独立 Job/Worker） | de3640e |
| Reliability | Fresh | de3640e |
| Risks | Fresh | de3640e |

## Knowledge Documents

| Doc | Description |
|---|---|
| [knowledge/overview.md](./knowledge/overview.md) | 项目目的、技术栈、入口、模块与基础设施 |
| [knowledge/architecture.md](./knowledge/architecture.md) | 启动时序、调用关系、数据流、入口类型 |
| [knowledge/modules.md](./knowledge/modules.md) | 模块地图与 API 使用面 |
| [knowledge/data-model.md](./knowledge/data-model.md) | 文件系统实体、WS/前端缓存、无 DB |
| [knowledge/business-flows.md](./knowledge/business-flows.md) | 启动、传文件、传内容等业务流程 |
| [knowledge/business-rules.md](./knowledge/business-rules.md) | 鉴权、路径、上传、WS 等规则 |
| [knowledge/state-machines.md](./knowledge/state-machines.md) | 无业务 FSM；WS/文件存在性生命周期 |
| [knowledge/integrations.md](./knowledge/integrations.md) | OS/开发集成；无云依赖 |
| [knowledge/reliability.md](./knowledge/reliability.md) | 幂等、并发、失败恢复 |
| [knowledge/risks.md](./knowledge/risks.md) | 安全与可靠性风险清单 |

未创建 `jobs-and-workers.md`：当前无 Cron/Worker/Consumer，避免空文档。

## Modules

| Module | Doc |
|---|---|
| pass | [modules/pass.md](./modules/pass.md) |
| stonePass UI | [modules/stonepass-ui.md](./modules/stonepass-ui.md) |
| engine | [modules/engine.md](./modules/engine.md) |
| gateway | [modules/gateway.md](./modules/gateway.md) |

## Decisions

| ID | Doc |
|---|---|
| DEC-0001 | [decisions/0001-vue-frontend-simplify.md](./decisions/0001-vue-frontend-simplify.md) |

## Changes

| ID | Doc |
|---|---|
| CHG-2026-001 | [changes/2026-09-08-knowledge-init.md](./changes/2026-09-08-knowledge-init.md) |

## Notes

| Note | Doc |
|---|---|
| README vs UI 漂移 | [notes/knowledge-drift-readme-vs-ui.md](./notes/knowledge-drift-readme-vs-ui.md) |
| 未决问题 | [notes/open-questions.md](./notes/open-questions.md) |

## Known Gaps

1. Password / JWT 完整安全模型未实现，意图不完全清楚
2. 发行版二进制与当前 git UI 是否一致：Unknown
3. Uniapp 手机端源码不在本仓库：Unknown
4. stream 上传与异常半截文件的清理策略未设计
5. 路径穿越加固方案未落地
6. docs App、tools 细节仅作浅层记录
7. 无自动化测试覆盖核心上传/安全路径（file_test 存在但未在本次深挖）

## Recommended Next Deep-Dives

1. **上传路径加固**：Big/Tiny 分流、覆盖语义、半截文件、路径穿越修复设计
2. **鉴权方案**：Password/JWT 与 README 规划权限模型如何落地
3. **cmd-key / cmd-open**：是否保留、如何限制仅本机或鉴权后可用
4. **WS 传内容**：重连、房间/配对、是否回显自身、历史是否要服务端存储
5. **测试现状**：`apps/pass/lib/files/file_test.go` 与其它测试缺口
6. **对齐 README**：更新文档或恢复 UI 能力，消除对外描述漂移

## How to use

复杂问题请按 protocol：先读本 index → 相关 Knowledge → 必要时回源码校验 → 增量更新 Knowledge。
