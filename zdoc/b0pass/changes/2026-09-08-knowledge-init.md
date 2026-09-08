# CHG-2026-001: 首次 Project Knowledge 初始化

Change ID: CHG-2026-001  
Status: Documented  
Date: 2026-09-08  

Related Decision:
- decisions/0001-vue-frontend-simplify.md

Related Knowledge:
- knowledge/*
- modules/*
- index.md

---

## What changed

在 `zdoc/b0pass` 建立第一版长期 Project Knowledge（非业务代码变更）。

对照源码 commit `de3640e9dd7cf65ba64d606ae6a8021807b4abe3` 完成整体地图、数据模型、流程、规则、可靠性与风险记录。

---

## Why

支持后续人类与 AI 增量维护，避免每次从零理解仓库。

---

## Previous / New

- Previous：仅有 protocol、空 index、空目录
- New：overview / architecture / modules / data-model / flows / rules / state-machines / integrations / reliability / risks + 核心 module 文档

---

## Impact

- 业务代码：无
- 认知资产：有

---

## Risks

文档可能随代码演进出漂移；需按 protocol 做 Sync。
