# Data Model

Status: Confirmed  
Last Verified: 2026-09-08  
Verified Commit: 9000df3  

无业务 DB。真相在磁盘 + 可选 Redis + 浏览器缓存。

## 1. 共享根目录

`[pass] Path`。相对路径参数 `f` 与之拼接（路径穿越风险见 risks）。

## 2. 目录列表项（运行时）

`GetDirTree` 返回 map：`name` / `path` / `type` / `size(s)` / `date` / `canread`…  
`.` 开头名称列表跳过。可选扩展：`expire` / `expireAt` / `expireLeft`。

## 3. 过期元数据

`{Path}/.b0pass-expire.json`：`{ "/rel": unix秒 }`；0 或缺失 = 不过期。

## 4. 上传

Tiny：`FormFile`；Big：流式 multipart。可选 `expire`。无分片会话 / 断点元数据。

## 5. 传内容

| 层 | 内容 |
|---|---|
| Hub/Client | 内存连接与广播 |
| History | Redis ZSET `stonepass:text:history`，≤10 条且 ≤72h；连不上则进程内存 |
| 前端 | `localStorage.txtdata` 仅缓存兜底 |

## 6. 配置

| Section | 字段 |
|---|---|
| gateway | ListenAddr, Domain, **Password**（空=关鉴权）, Live, Debug |
| pass | Live, Path, RedisAddr, RedisPassword, RedisDB |
| docs | Live |

`GET /gateway/config` 返回时会把 Password 换成 `runtime.GOOS`（脱敏/怪异行为，意图 Unknown）。

## 关系

```
Path/
  .b0pass-expire.json
  …files…
Hub → Clients
History(Redis|memory)
```

## 写操作一览

| 操作 | 效果 |
|---|---|
| node-add `/` 结尾 | MkdirAll 0755 |
| node-delete | Remove（非递归）+ 清 expire |
| upload | 覆盖同名；可选 expire |
| file-expire | 写元数据 |
| purge | 文件 Remove / 目录 RemoveAll |
| text-sync | 广播 + History.Add |
