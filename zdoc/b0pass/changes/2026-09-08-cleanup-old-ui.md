# Change: 清理旧 pass UI 残留

Date: 2026-09-08  
Type: cleanup  

## Summary

删除空目录 `apps/pass/ui/dist_old/`（旧 LayUI 迁出后的空壳），以及未引用的测试页 `apps/pass/lib/stream/upload.html`。

当前 pass 前端仅保留：`apps/pass/ui/vue/stonePass/`（源码）与 `apps/pass/ui/dist/`（构建产物 / embed）。
