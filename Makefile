# b0pass / StonePass — 编译与部署
#
# 用法:
#   make          # 一条命令：检查工具 → 装依赖（缺则装）→ 编前端 → 编二进制
#   make deploy   # 同上，并产出 dist/ 部署包
#   make check    # 仅检查（不安装、不编译）
#   make clean    # 清理 dist/
#
# 可选变量:
#   GOOS=linux GOARCH=amd64   # 交叉编译

ROOT         := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
FRONTEND_DIR := $(ROOT)/apps/pass/ui/vue/stonePass
MAIN_DIR     := $(ROOT)/main
OUT_DIR      := $(ROOT)/dist
BIN_NAME     ?= b0pass

GOOS   ?= $(shell go env GOOS 2>/dev/null)
GOARCH ?= $(shell go env GOARCH 2>/dev/null)

ifeq ($(GOOS),windows)
  BIN_FILE := $(OUT_DIR)/$(BIN_NAME).exe
else
  BIN_FILE := $(OUT_DIR)/$(BIN_NAME)
endif

.PHONY: all check tools ensure-deps deps frontend backend build deploy clean help

all: build

help:
	@echo "目标:"
	@echo "  make / make build   一条命令完成：检查工具、安装依赖、编译前后端"
	@echo "  make deploy         同上，并输出到 dist/（含示例配置）"
	@echo "  make check          仅检查依赖是否齐全（不安装、不编译）"
	@echo "  make clean          清理 dist/"
	@echo ""
	@echo "交叉编译: make deploy GOOS=linux GOARCH=amd64"

# ---------- 仅检查（不自动安装）----------
check:
	@bash "$(ROOT)/scripts/check-deps.sh" "$(ROOT)" "$(FRONTEND_DIR)" check

# ---------- 工具链（缺则失败，无法自动装）----------
tools:
	@bash "$(ROOT)/scripts/check-deps.sh" "$(ROOT)" "$(FRONTEND_DIR)" tools

# ---------- 项目依赖：缺则自动安装 ----------
ensure-deps: tools
	@bash "$(ROOT)/scripts/ensure-deps.sh" "$(ROOT)" "$(FRONTEND_DIR)"

deps: ensure-deps

# ---------- 编译 ----------
frontend:
	@echo "==> 构建前端 → apps/pass/ui/dist"
	cd "$(FRONTEND_DIR)" && npm run build
	@test -f "$(ROOT)/apps/pass/ui/dist/index.html" || { echo "[失败] 前端产物缺少 index.html"; exit 1; }
	@echo "前端构建完成"

backend:
	@echo "==> 编译 Go（GOOS=$(GOOS) GOARCH=$(GOARCH)）→ $(BIN_FILE)"
	@mkdir -p "$(OUT_DIR)"
	cd "$(MAIN_DIR)" && \
		CGO_ENABLED=$${CGO_ENABLED:-0} GOOS="$(GOOS)" GOARCH="$(GOARCH)" \
		go build -trimpath -ldflags="-s -w" -o "$(BIN_FILE)" .
	@echo "后端编译完成: $(BIN_FILE)"

build: ensure-deps frontend backend
	@echo "全部完成: $(BIN_FILE)"

# ---------- 部署产出 ----------
deploy: build
	@echo "==> 准备部署目录 $(OUT_DIR)"
	@if [ ! -f "$(OUT_DIR)/config.ini.example" ]; then \
		printf '%s\n' \
			'[gateway]' \
			'ListenAddr = ":8888"' \
			'Domain = ""' \
			'Password = ""' \
			'' \
			'[pass]' \
			'Path = "files"' \
			'Live = false' \
			'RedisAddr = ""' \
			'RedisPassword = ""' \
			'RedisDB = 0' \
			> "$(OUT_DIR)/config.ini.example"; \
		echo "  已写入 config.ini.example"; \
	fi
	@mkdir -p "$(OUT_DIR)/files"
	@printf '%s\n' \
		'# 部署说明' \
		'1. 将本目录中的二进制与 config.ini.example 拷到目标机器' \
		'2. 复制为 config.ini，按需修改 ListenAddr / Password / Path / Redis' \
		'3. 生产请设 [pass] Live = false（前端已 embed 进二进制）' \
		'4. 在含 config.ini 的目录运行: ./$(BIN_NAME)' \
		'5. 浏览器打开: http://127.0.0.1:8888/app/pass/' \
		> "$(OUT_DIR)/README.txt"
	@echo "部署产物已就绪: $(OUT_DIR)/"
	@ls -la "$(OUT_DIR)"

clean:
	@rm -rf "$(OUT_DIR)"
	@echo "已清理 $(OUT_DIR)"
	@echo "提示: 前端静态资源在 apps/pass/ui/dist，如需重建请 make"
