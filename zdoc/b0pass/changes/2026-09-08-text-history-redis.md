# 传内容 Redis 历史

Date: 2026-09-08  
Type: feature

## Summary

传内容同步时写入 Redis（失败回退内存），新用户进入可拉取最近历史：最多 10 条且不超过 72 小时。

## Config

```ini
[pass]
RedisAddr = "127.0.0.1:6379"
RedisPassword = ""
RedisDB = 0
```

`RedisAddr` 为空或 Ping 失败 → 内存存储（仅本进程有效）。

## API

`GET /pass/text-history`（需登录，若开启鉴权）
