# 共享口令登录（方案 A）

Date: 2026-09-08  
Type: feature

## Summary

`[gateway] Password` 非空时启用整站登录：JWT + 前端路由守卫；空口令关闭鉴权。

## Changes

- `POST /pass/login`、`GET /pass/auth-status`
- pass API / `/files` / `/ws` JWT（Header/Query/Cookie `token`）
- LoginView + 路由守卫；请求自动带 token
- 默认/示例 `config.ini` 含 `Password = "a123"`

## Notes

后期可在 claims 中扩展 Pread/Pupload/Padmin 角色校验。
