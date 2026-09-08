# Architecture

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

---

## 形态

单体 Go 进程：App 插件式注册（`engine.AppInstall`）+ 单 Gin 端口聚合 UI / API / 静态 / WS。

## 启动

```
main
  ├ 无 config.ini → 写默认并建 files/
  ├ 打印访问地址（不再自动开浏览器）
  └ engine.Run
        ├ 配置灌入各 App
        ├ go App.Run()（pass / docs / gateway…）
        └ gateway 内 Gin.Run 阻塞；main select{} 保活
```

gateway 启动时 `engine.SetAuthPassword(Password)`。

## 调用关系

```
main → engine → pass / gateway / docs
pass → files / stream / chat(+history) / keys / audit
UI   → /pass/* 、/files 、/ws
```

## 数据流

**文件**：`f` 相对 Path；列表 / 上传 / 下载写读磁盘；过期见 `.b0pass-expire.json` + 定时 purge。  
**传内容**：WS 广播全 Hub；广播时写入 History（Redis，失败则内存）；新客户端 `GET /pass/text-history`。  
**鉴权**：Password 非空时 JWT（Header / Query / Cookie `token`）护住 pass API、`/files`、`/ws`。

## 同步性

| 场景 | 模式 |
|---|---|
| HTTP | 请求内同步完成 |
| 大文件上传 | 流式写盘 |
| WS | Hub 单协程 + 每连接 pump |
| 过期清理 | 启动 + 每分钟 + 列表前 |

无 MQ / 独立 Worker。

## 前端托管

| `[pass] Live` | 行为 |
|---|---|
| true | 优先反代 Vite `:5173`，否则 `ui/dist` 磁盘 |
| false | embed `ui/dist` |

gateway **不**托管 pass UI。

## 设计要点

1. 文件系统即主数据模型  
2. 单端口聚合  
3. 大文件自定义 multipart 流式解析（Content-Length > 4096）
