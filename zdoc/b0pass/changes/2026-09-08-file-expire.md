# 文件可选过期时间

Date: 2026-09-08  
Type: feature

## Summary

传文件支持可选过期时间；默认不过期。到期后由后台清理删除。

## Changes

- 后端：`ExpireStore`（`.b0pass-expire.json`）、`GET /pass/file-expire`、上传 query `expire`、每分钟 cleaner
- 前端：上传下拉选择过期；列表展示剩余时间；「过期」按钮可改/取消
- Knowledge：business-flows / rules / data-model / stonepass-ui / pass 模块已同步

## Notes

需重启 Go 进程以加载 cleaner 与新路由；前端已 `npm run build`。
