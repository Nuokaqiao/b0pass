# Change: 目录上传与超 1G 提示

Date: 2026-09-08  
Type: feature  

## Summary

传文件支持选择/拖拽文件夹，按相对路径上传到当前目录下对应子路径；所选文件合计超过 1 GiB 时弹窗警告并二次确认。

## Scope

- Frontend: `FilesView.vue`、`pass.js`（`saveAsName`）
- Backend: 无改动（沿用 `file-upload?f=` + `NodeAdd` 建目录）

## Notes

- 空文件夹不会单独落盘
- 相对路径拒绝 `..`
