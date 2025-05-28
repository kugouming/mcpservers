# Go 相关变量
GO=go
GOFLAGS=-v
LDFLAGS=-w -s

# 目录相关变量
WORK_DIR=$(shell pwd)
CLIENT_DIR=cmd/client
OUTPUT_DIR=output
CLIENT_OUTPUT_DIR=output/client
SSE_DIR=cmd/server/sse
SSE_OUTPUT_DIR=output/sse

# 获取所有 client 和 sse 子目录
CMDS := $(shell find $(CLIENT_DIR) -maxdepth 1 -mindepth 1 -type d -exec basename {} \;)
SSE_CMDS := $(shell find $(SSE_DIR) -maxdepth 1 -mindepth 1 -type d -exec basename {} \;)
SSE_TARGETS := $(addprefix sse_,$(SSE_CMDS))

# 默认目标
.PHONY: all
all: $(CMDS) $(SSE_TARGETS)
	@echo "Sync config..."
	@cp -r ./config $(OUTPUT_DIR)/
	@echo ""
	@echo "Build done."
	@echo ""
	@echo "Compile output:"
	@echo ""
	@echo "================================"
	@echo ""
	@tree ./$(OUTPUT_DIR)

# 创建输出目录
$(CLIENT_OUTPUT_DIR) $(SSE_OUTPUT_DIR):
	@mkdir -p $@

# 为每个 cmd 创建编译目标
.PHONY: $(CMDS)
$(CMDS): $(CLIENT_OUTPUT_DIR)
	@echo "Building $@..."
	@cd $(CLIENT_DIR)/$@ && $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o ../../../$(CLIENT_OUTPUT_DIR)/$@

# 为每个 sse cmd 创建编译目标
.PHONY: $(SSE_TARGETS)
$(SSE_TARGETS): $(SSE_OUTPUT_DIR)
	@cmd=$$(echo $@ | sed 's/^sse_//'); \
	echo "Building SSE $$cmd..."; \
	cd $(SSE_DIR)/$$cmd && $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o ../../../../$(SSE_OUTPUT_DIR)/$$cmd

# 输出所有命令
.PHONY: list
list:
	@echo "Available commands to build:"
	@echo "Standard commands:"
	@for cmd in $(CMDS); do \
		echo "  $(WORK_DIR)/$(CLIENT_OUTPUT_DIR)/$$cmd"; \
	done
	@echo "SSE commands:"
	@for cmd in $(SSE_CMDS); do \
		echo "  $(WORK_DIR)/$(SSE_OUTPUT_DIR)/$$cmd"; \
	done

# 输出MCPSever配置
.PHONY: config
config:
	@echo "{"
	@echo "  \"mcpServers\": {"
	@# 输出标准命令配置
	@for cmd in $(CMDS); do \
		echo "    \"$$cmd\": {"; \
		echo "      \"command\": \"$(WORK_DIR)/$(CLIENT_OUTPUT_DIR)/$$cmd\","; \
		echo "      \"args\": [],"; \
		echo "      \"env\": {}"; \
		if [ "$$cmd" != "$$(echo $(CMDS) | rev | cut -d' ' -f1 | rev)" ] || [ -n "$(SSE_CMDS)" ]; then \
			echo "    },"; \
		else \
			echo "    }"; \
		fi; \
	done
	@# 输出 SSE 命令配置
	@for cmd in $(SSE_CMDS); do \
		port=$$(grep -aoE ":[0-9]{4}" $(SSE_DIR)/$$cmd/main.go | head -n1 || echo ":8080"); \
		echo "    \"sse_$$cmd\": {"; \
		echo "      \"transport\": \"sse\","; \
		echo "      \"url\": \"http://127.0.0.1$$port/sse\""; \
		if [ "$$cmd" != "$$(echo $(SSE_CMDS) | rev | cut -d' ' -f1 | rev)" ]; then \
			echo "    },"; \
		else \
			echo "    }"; \
		fi; \
	done
	@echo "  }"
	@echo "}"

# 清理目标
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(OUTPUT_DIR)

# 帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Build all commands (including SSE commands)"
	@echo "  list       - List all commands"
	@echo "  config     - Output MCPSever config"
	@echo "  clean      - Remove all built binaries"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Available commands to build:"
	@echo "Standard commands:"
	@for cmd in $(CMDS); do \
		echo "  $$cmd"; \
	done
	@echo "SSE commands:"
	@for cmd in $(SSE_CMDS); do \
		echo "  sse_$$cmd"; \
	done 