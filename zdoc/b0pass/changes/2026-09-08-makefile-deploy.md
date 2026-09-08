# Change: Make 编译部署

Date: 2026-09-08  
Type: tooling  

## Summary

新增根目录 `Makefile` + `scripts/check-deps.sh`：先检查 go/node/npm 与前后端依赖，齐全再编译；缺失项汇总提示。`make deploy` 产出到 `dist/`。
