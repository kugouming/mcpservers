# YAPI 工具快速开始

## 快速安装

```bash
# 1. 克隆项目
git clone https://github.com/kugouming/mcpservers.git
cd mcpservers

# 2. 安装依赖
go mod tidy
```

## 配置设置

### 方式一：环境变量（推荐）

```bash
export YAPI_BASE_URL="http://your-yapi-server.com"
export YAPI_TOKEN="your_access_token"
export YAPI_TIMEOUT="30"
```

### 方式二：配置文件

```bash
# 创建配置文件
cat > yapi.yaml << EOF
base_url: "http://your-yapi-server.com"
token: "your_access_token" 
timeout: 30
log_level: "info"
EOF
```

## 快速测试

```bash
# 验证配置
go test ./tools/yapi/ -v

# 运行示例程序
go run tools/yapi/example/main.go
```

## 基本使用

### 1. 使用配置文件/环境变量

```go
package main

import (
    "fmt"
    "github.com/kugouming/mcpservers/tools/yapi"
)

func main() {
    // 自动加载配置（环境变量优先级高于配置文件）
    config, err := yapi.LoadGlobalConfig()
    if err != nil {
        panic(err)
    }
    
    // 创建客户端
    client := yapi.NewYapiClientFromConfig(config)
    
    // 获取项目接口列表
    interfaces, err := client.GetInterfaces(123)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("找到 %d 个接口\n", interfaces.Data.Count)
}
```

### 2. 直接指定参数

```go
package main

import (
    "fmt"
    "github.com/kugouming/mcpservers/tools/yapi"
)

func main() {
    // 直接创建客户端
    client := yapi.NewYapiClient("http://your-yapi.com", "your_token")
    
    // 获取项目接口列表
    interfaces, err := client.GetInterfaces(123)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("找到 %d 个接口\n", interfaces.Data.Count)
}
```

### 3. 集成到 MCP 服务器

```go
package main

import (
    "github.com/mark3labs/mcp-go/server"
    "github.com/kugouming/mcpservers/tools/yapi"
)

func main() {
    // 创建 MCP 服务器
    s := server.NewMCPServer("YAPI Server", "1.0.0")
    
    // 注册 YAPI 工具（自动使用配置文件/环境变量）
    yapi.RegisterTool(s)
    
    // 启动服务器
    s.Serve()
}
```

## 可用工具

| 工具名称 | 功能描述 | 参数 | 示例 |
|---------|---------|------|------|
| `get_interfaces` | 获取项目接口列表 | `project_id` | `{"project_id": 123}` |
| `get_interface_detail` | 获取接口详情 | `id` | `{"id": 1001}` |
| `get_project_info` | 获取项目信息 | `project_id` | `{"project_id": 123}` |

## 配置管理

### 配置优先级
1. **环境变量** (最高优先级)
2. **配置文件**
3. **默认值** (最低优先级)

### 支持的环境变量

| 配置项 | 环境变量 | 默认值 | 说明 |
|-------|----------|--------|------|
| base_url | `YAPI_BASE_URL` | - | YAPI 服务器地址 |
| token | `YAPI_TOKEN` | - | 访问令牌 |
| timeout | `YAPI_TIMEOUT` | 30 | 请求超时（秒） |
| retry_count | `YAPI_RETRY_COUNT` | 3 | 重试次数 |
| log_level | `YAPI_LOG_LEVEL` | info | 日志级别 |

### 配置文件位置
系统按以下顺序查找 `yapi.yaml`：
1. `./yapi.yaml` (当前目录)
2. `./config/yapi.yaml`
3. `./configs/yapi.yaml` 
4. `~/.yapi/yapi.yaml` (用户目录)
5. `/etc/yapi/yapi.yaml` (系统目录)

## 配置管理示例

```bash
# 运行示例程序的配置管理模式
go run tools/yapi/example/main.go
# 选择 "4. 配置管理示例"

# 可以进行以下操作：
# 1. 查看当前配置
# 2. 生成示例配置文件
# 3. 保存当前配置到文件
# 4. 查看配置来源
# 5. 验证配置
```

## 常见问题

**Q: 如何获取 YAPI Token？**
A: 登录 YAPI → 进入项目 → 设置 → Token 配置

**Q: 环境变量和配置文件哪个优先级高？**
A: 环境变量优先级更高，会覆盖配置文件中的设置

**Q: 如何检查当前使用的配置？**
A: 使用配置管理器：`yapi.GetConfigManager().PrintConfig()`

**Q: 支持哪些配置文件格式？**
A: 支持 YAML、JSON、TOML 等格式，推荐使用 YAML

**Q: 如何验证配置是否正确？**
A: 运行 `yapi.ValidateEnvironment()` 或测试 `go test ./tools/yapi/ -run TestValidateEnvironment -v`

## 故障排查

```bash
# 1. 检查环境变量
echo $YAPI_BASE_URL
echo $YAPI_TOKEN

# 2. 验证配置
go test ./tools/yapi/ -run TestValidateEnvironment -v

# 3. 测试网络连接
curl "$YAPI_BASE_URL/api/interface/list?project_id=123&token=$YAPI_TOKEN"

# 4. 查看详细日志
export YAPI_LOG_LEVEL=debug
go run tools/yapi/example/main.go
```

## 性能表现

- ✅ 接口列表查询：平均 121ms
- ✅ 接口详情查询：平均 110ms  
- ✅ 内存使用：约 10-12KB per request
- ✅ 并发安全：支持多 goroutine 并发调用
- ✅ 配置缓存：支持热重载，无需重启

## 新特性 (v1.1.0)

- 🆕 支持配置文件和环境变量
- 🆕 多种配置文件格式支持
- 🆕 配置热重载功能
- 🆕 详细的配置验证
- 🆕 配置来源追踪
- 🆕 示例配置文件生成

## 更多帮助

- 📖 [详细文档](README.md)
- 🧪 [测试示例](tool_test.go)
- 💡 [示例程序](example/main.go)
- ⚙️ [配置示例](yapi.example.yaml) 