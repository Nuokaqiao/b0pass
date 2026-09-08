# Risks

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

只记**仍成立**的风险。已缓解项不堆历史故事。

## 高

**R1 共享口令即全权限**  
Password 为空则完全开放；非空则凡持 JWT 者可读写、删、WS、调用 cmd API。无角色隔离。

**R2 路径穿越**  
多数 handler `Path + f`，缺「必须落在 Root 内」校验。

**R3 `/files` 整树暴露**  
静态挂载共享根；有 token（或鉴权关闭）即可按路径取文件。

**R4 cmd-key / cmd-open**  
UI 未用但 API 仍在；鉴权关闭或口令泄露时可被直接调用。

## 中

**R5 JWT Secret 硬编码**于源码。  
**R6 并发同名上传**无锁 / 无临时文件提交。  
**R7 WS** 全局单 Hub；默认 Upgrader 来源策略需留意。  
**R8** 官网 / 发行页截图文案可能仍展示旧 UI；仓库 README 已对齐当前实现。  
**R9 Redis 未启动时历史进内存**，重启丢失（易被当成「没存 Redis」）。

## 低–中

**R11** 大目录同步扫盘性能。  
**R12** `config.ini` Password 明文。
