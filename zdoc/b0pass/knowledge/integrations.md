# Integrations

Status: Confirmed  
Freshness: Fresh  

Last Verified: 2026-09-08  
Verified Commit: de3640e9dd7cf65ba64d606ae6a8021807b4abe3  

Related Paths:
- apps/pass/lib/keys/**
- core/tools/cmd/**
- apps/pass/main.go
- go.mod

---

## 外部系统

本项目面向局域网本地使用，**无云 API、支付、账号体系、对象存储、MQ**。

---

## 操作系统集成

| 集成 | 用途 | 平台 | Confidence |
|---|---|---|---|
| 系统打开 URL/文件 | `cmd.Open`（open / xdg-open / start） | Win/Mac/Linux | Confirmed |
| 键鼠模拟 | `robotgo` via `keys.SendKey` | **仅 Windows**；其它平台空函数 | Confirmed |
| 出网网卡 IP | `nets.GetOutBoundIP` 展示访问地址 | 通用 | Confirmed |
| 释放 zlib1.dll | Windows 且非 Live 时从 embed 写出 | Windows | Confirmed |

---

## 开发期集成

| 集成 | 说明 |
|---|---|
| Vite Dev Server | Live 模式下反代 `127.0.0.1:5173` |
| Vite proxy | `/pass` `/files` `/ws` `/gateway` → `:8888` |

---

## 发行 / 官网（文档提及，非代码依赖）

- 官网下载：https://4bit.cn/p/b0pass
- GitHub / Gitee 仓库链接见 README

这些不是运行时集成。

---

## 未使用但存在的库能力

| 库 | 代码位置 | 业务使用 |
|---|---|---|
| gorm | `core/engine/database.go` | 未连接 DB |
| jwt | `middleware.go` | 仅 ping 校验；无登录发牌流程 |
| LRU cache | `core/tools/cache` | 无引用 |
