# Change: 目录删除改为 RemoveAll

Date: 2026-09-08  
Type: bugfix  

## Summary

删文件夹时改用 `os.RemoveAll`，避免列表不展示的隐藏文件（如 macOS `.DS_Store`）导致「目录非空无法删除」。同时清理该目录及子路径的过期元数据。

## Root cause

`GetDirTree` 跳过 `.` 开头名；`NodeRemove` 原用 `os.Remove`，非空目录失败。
