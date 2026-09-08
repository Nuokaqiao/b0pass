# Reliability

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

## 幂等

上传同名非幂等（覆盖）。`MkdirAll` 可重复。delete 第二次失败。WS 每发一条。

## 并发

同名并发上传无锁 → 可能损坏/覆盖。Hub clients 仅 Hub 协程访问。列表与删除有常规 FS 竞态。

## 失败与恢复

- 崩溃：已落盘保留；进行中上传可能半截文件；**无**上传会话恢复  
- 多文件 multipart 中途失败：已成功文件保留（部分成功）  
- WS：慢客户端被踢，队列丢；前端可自动重连，**在途消息不补**（可再拉 history）  
- Redis 不可用：历史回退内存（重启丢）

## 大数据

流式上传避整包内存。预览/WS 约 2MB 上限。超大目录扫盘可能慢（`file-count`/`node-tree` UI 未用）。

## 下载

标准文件服务，支持 Range；无限速。
