# State Machines

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

无订单类业务 FSM。文件「存在与否」即状态。

## WS（前端）

`connecting` → `open` →（断线）`reconnect` 自动重试 → 再 `open`；或失败文案。  
发送时非 OPEN 则提示无法同步。

## Hub Client

Upgrade → register → pumps → 读失败/踢出 → unregister。

## 文件

`不存在 → upload/add → 存在 → rename? → delete → 不存在`。  
无 uploading 持久态；流式上传中途可能见到半截文件。
