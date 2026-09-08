#!/usr/bin/env bash
# 检查前后端编译依赖
# 用法: check-deps.sh <root> <frontend_dir> [check|tools]
#   tools — 仅检查 go/node/npm 与项目文件（不可自动安装的部分）
#   check — 完整检查含 node_modules / go 模块（默认）
set +e
ROOT="${1:?root dir}"
FRONTEND_DIR="${2:?frontend dir}"
MODE="${3:-check}"
missing=0

echo "==> 检查编译依赖 (${MODE})"

if ! command -v go >/dev/null 2>&1; then
  echo "  [缺] Go 工具链（请安装 Go 1.17+，并确保 go 在 PATH）"
  missing=1
else
  echo "  [ok] go  $(go version | head -1)"
fi

if ! command -v node >/dev/null 2>&1; then
  echo "  [缺] Node.js（前端构建需要，建议 18+）"
  missing=1
else
  echo "  [ok] node $(node -v)"
fi

if ! command -v npm >/dev/null 2>&1; then
  echo "  [缺] npm（通常随 Node.js 安装）"
  missing=1
else
  echo "  [ok] npm  $(npm -v)"
fi

if [ ! -f "$ROOT/go.mod" ]; then
  echo "  [缺] 后端模块文件 go.mod"
  missing=1
else
  echo "  [ok] go.mod"
fi

if [ ! -f "$FRONTEND_DIR/package.json" ]; then
  echo "  [缺] 前端 package.json（路径: apps/pass/ui/vue/stonePass）"
  missing=1
else
  echo "  [ok] 前端 package.json"
fi

if [ "$MODE" = "tools" ]; then
  if [ "$missing" -ne 0 ]; then
    echo ""
    echo "工具链检查未通过，请先安装上述缺失项"
    exit 1
  fi
  echo "工具链检查通过"
  exit 0
fi

if [ ! -d "$FRONTEND_DIR/node_modules" ]; then
  echo "  [缺] 前端依赖 node_modules → 执行 make 会自动安装"
  missing=1
else
  echo "  [ok] 前端 node_modules"
  if [ ! -d "$FRONTEND_DIR/node_modules/vite" ] || [ ! -d "$FRONTEND_DIR/node_modules/vue" ]; then
    echo "  [缺] 前端核心包不完整（vite/vue）→ 执行 make 会自动安装"
    missing=1
  else
    echo "  [ok] 前端核心包 vite / vue"
  fi
fi

if command -v go >/dev/null 2>&1 && [ -f "$ROOT/go.mod" ]; then
  if ! (cd "$ROOT" && go list ./main/ >/dev/null 2>&1); then
    echo "  [缺] Go 依赖无法解析 → 执行 make 会自动安装"
    missing=1
  else
    echo "  [ok] Go 模块可解析"
  fi
fi

if [ "$missing" -ne 0 ]; then
  echo ""
  echo "依赖检查未通过。可直接执行: make   （会自动安装项目依赖后再编译）"
  exit 1
fi

echo "依赖检查通过"
exit 0
