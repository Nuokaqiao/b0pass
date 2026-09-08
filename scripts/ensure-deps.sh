#!/usr/bin/env bash
# 缺什么装什么：前端 node_modules / 后端 go mod
set -e
ROOT="${1:?root dir}"
FRONTEND_DIR="${2:?frontend dir}"

need_npm=0
need_go=0

if [ ! -d "$FRONTEND_DIR/node_modules" ] \
  || [ ! -d "$FRONTEND_DIR/node_modules/vite" ] \
  || [ ! -d "$FRONTEND_DIR/node_modules/vue" ]; then
  need_npm=1
fi

if ! (cd "$ROOT" && go list ./main/ >/dev/null 2>&1); then
  need_go=1
fi

if [ "$need_go" -eq 0 ] && [ "$need_npm" -eq 0 ]; then
  echo "==> 项目依赖已就绪，跳过安装"
  exit 0
fi

if [ "$need_go" -eq 1 ]; then
  echo "==> 安装后端 Go 依赖"
  (cd "$ROOT" && go mod download && go mod tidy)
fi

if [ "$need_npm" -eq 1 ]; then
  echo "==> 安装前端 npm 依赖"
  (cd "$FRONTEND_DIR" && npm install)
fi

echo "依赖安装完成"
