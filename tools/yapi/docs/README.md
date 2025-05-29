# YAPI 工具说明文档

## 概述

YAPI 工具是一个基于 Golang 和 `github.com/mark3labs/mcp-go` 开发的 MCP（Model Context Protocol）工具，用于与 YAPI（Yet Another API）系统进行交互。通过此工具，您可以轻松获取 YAPI 项目的接口列表、接口详情和项目信息。

## 功能特性

- **🚀 高性能**: 使用 Golang 原生并发特性，支持高并发请求
- **🛡️ 健壮性**: 完善的错误处理和容错机制
- **📝 详细注释**: 每个函数和结构体都有详细的中文注释
- **🧪 完整测试**: 包含单元测试、集成测试和性能测试
- **🔧 易于扩展**: 模块化设计，便于添加新功能

## 支持的工具

### 1. get_interfaces - 获取接口列表
获取指定 YAPI 项目的所有接口列表。

**参数:**
- `project_id` (number, 必需): YAPI 项目 ID

**返回:**
```json
{
  "total": 10,
  "count": 10,
  "interfaces": [
    {
      "id": 1001,
      "title": "获取用户信息",
      "path": "/api/user/info",
      "method": "GET",
      "status": "done",
      "tag": ["用户", "基础"]
    }
  ]
}
```

### 2. get_interface_detail - 获取接口详情
获取指定接口的详细信息，包括请求参数、响应格式等。

**参数:**
- `id` (number, 必需): 接口 ID

**返回:**
```json
{
  "id": 1001,
  "title": "获取用户信息",
  "path": "/api/user/info",
  "method": "GET",
  "status": "done",
  "description": "根据用户ID获取用户详细信息",
  "req_headers": [...],
  "req_query": [...],
  "res_body": "...",
  "add_time": "2022-01-01 00:00:00",
  "up_time": "2022-01-01 00:00:00"
}
```

### 3. get_project_info - 获取项目信息
获取指定项目的基本信息。

**参数:**
- `project_id` (number, 必需): YAPI 项目 ID

**返回:**
```json
{
  "id": 123,
  "name": "测试项目",
  "description": "这是一个测试项目",
  "basepath": "/api",
  "group_id": 1,
  "group_name": "测试分组"
}
```

## 环境配置

在使用 YAPI 工具之前，需要设置以下环境变量：

### 必需环境变量

```bash
# YAPI 服务器地址
export YAPI_BASE_URL="http://your-yapi-server.com"

# YAPI 访问令牌
export YAPI_TOKEN="your_access_token"
```

### 获取 YAPI Token

1. 登录您的 YAPI 系统
2. 进入项目设置页面
3. 在 "Token配置" 部分找到项目 Token
4. 复制 Token 并设置到环境变量中

## 安装和使用

### 1. 克隆项目

```bash
git clone https://github.com/kugouming/mcpservers.git
cd mcpservers
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 设置环境变量

```bash
# 方式一：在终端中设置
export YAPI_BASE_URL="http://your-yapi-server.com"
export YAPI_TOKEN="your_access_token"

# 方式二：在 .env 文件中设置（推荐）
echo "YAPI_BASE_URL=http://your-yapi-server.com" >> .env
echo "YAPI_TOKEN=your_access_token" >> .env
```

### 4. 编译程序

```bash
# 编译服务端程序
make build-sse

# 或者直接运行
go run cmd/server/sse/ginserver/main.go
```

### 5. 测试工具

```bash
# 运行单元测试
go test ./tools/yapi/ -v

# 运行性能测试
go test ./tools/yapi/ -bench=. -benchmem

# 运行覆盖率测试
go test ./tools/yapi/ -cover
```

## 使用示例

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/kugouming/mcpservers/tools/yapi"
)

func main() {
    // 创建 YAPI 客户端
    client := yapi.NewYapiClient("http://your-yapi-server.com", "your_token")
    
    // 获取项目接口列表
    interfaces, err := client.GetInterfaces(123)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("找到 %d 个接口\n", interfaces.Data.Count)
    for _, iface := range interfaces.Data.List {
        fmt.Printf("接口: %s - %s %s\n", iface.Title, iface.Method, iface.Path)
    }
}
```

### 与 MCP 服务器集成

```go
package main

import (
    "github.com/mark3labs/mcp-go/server"
    "github.com/kugouming/mcpservers/tools/yapi"
)

func main() {
    // 创建 MCP 服务器
    mcpServer := server.NewMCPServer("YAPI Server", "1.0.0")
    
    // 注册 YAPI 工具
    yapi.RegisterTool(mcpServer)
    
    // 启动服务器
    mcpServer.Serve()
}
```

## API 参考

### 客户端接口

```go
type YapiClient interface {
    // 获取项目接口列表
    GetInterfaces(projectID int) (*YapiInterfaceListResponse, error)
    
    // 获取接口详细信息
    GetInterfaceDetail(interfaceID int) (*YapiInterfaceDetailResponse, error)
    
    // 获取项目信息
    GetProjectInfo(projectID int) (*YapiProjectResponse, error)
}
```

### 创建客户端

```go
func NewYapiClient(baseURL, token string) YapiClient
```

### 注册 MCP 工具

```go
func RegisterTool(s *server.MCPServer)
```

### 验证环境配置

```go
func ValidateEnvironment() error
```

## 错误处理

工具提供了完善的错误处理机制：

### 常见错误类型

1. **环境变量未设置**
   ```
   错误: YAPI_BASE_URL 环境变量未设置
   解决: 设置正确的 YAPI 服务器地址
   ```

2. **Token 无效**
   ```
   错误: HTTP请求失败，状态码: 401
   解决: 检查 YAPI_TOKEN 是否正确
   ```

3. **项目或接口不存在**
   ```
   错误: HTTP请求失败，状态码: 404
   解决: 检查项目ID或接口ID是否正确
   ```

4. **网络连接问题**
   ```
   错误: 请求失败: dial tcp: no such host
   解决: 检查网络连接和 YAPI 服务器地址
   ```

### 错误处理最佳实践

```go
response, err := client.GetInterfaces(projectID)
if err != nil {
    // 记录错误日志
    log.Printf("获取接口列表失败: %v", err)
    
    // 返回友好的错误信息
    return fmt.Errorf("无法获取项目 %d 的接口列表: %w", projectID, err)
}

// 检查 YAPI 返回的错误
if response.ErrCode != 0 {
    return fmt.Errorf("YAPI 错误: %s (错误码: %d)", response.ErrMsg, response.ErrCode)
}
```

## 性能优化

### 1. HTTP 客户端配置
```go
httpClient: &http.Client{
    Timeout: 30 * time.Second, // 设置合理的超时时间
}
```

### 2. 并发请求示例
```go
// 并发获取多个接口详情
var wg sync.WaitGroup
results := make(chan *YapiInterfaceDetailResponse, len(interfaceIDs))

for _, id := range interfaceIDs {
    wg.Add(1)
    go func(interfaceID int) {
        defer wg.Done()
        detail, err := client.GetInterfaceDetail(interfaceID)
        if err == nil {
            results <- detail
        }
    }(id)
}

wg.Wait()
close(results)
```

### 3. 缓存机制（推荐）
```go
// 可以使用 Redis 或内存缓存来缓存接口信息
// 减少对 YAPI 服务器的请求频率
```

## 测试

项目包含完整的测试套件：

### 运行测试

```bash
# 运行所有测试
go test ./tools/yapi/ -v

# 运行特定测试
go test ./tools/yapi/ -run TestYapiClient_GetInterfaces -v

# 运行性能测试
go test ./tools/yapi/ -bench=BenchmarkYapiClient -benchmem

# 生成测试覆盖率报告
go test ./tools/yapi/ -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 测试覆盖的场景

- ✅ 正常接口调用
- ✅ 错误处理（网络错误、认证失败等）
- ✅ 参数验证
- ✅ JSON 序列化/反序列化
- ✅ 环境变量验证
- ✅ MCP 工具注册
- ✅ 性能基准测试

## 故障排查

### 1. 检查环境变量
```bash
echo $YAPI_BASE_URL
echo $YAPI_TOKEN
```

### 2. 验证网络连接
```bash
curl "$YAPI_BASE_URL/api/interface/list?project_id=123&token=$YAPI_TOKEN"
```

### 3. 查看日志
```bash
# 启用详细日志
export LOG_LEVEL=debug
go run cmd/server/sse/ginserver/main.go
```

### 4. 使用测试工具验证
```bash
# 运行环境验证测试
go test ./tools/yapi/ -run TestValidateEnvironment -v
```

## 贡献指南

欢迎为项目做出贡献！请遵循以下步骤：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/new-feature`)
3. 提交更改 (`git commit -am 'Add some feature'`)
4. 推送到分支 (`git push origin feature/new-feature`)
5. 创建 Pull Request

### 代码规范

- 所有公共函数必须有详细的中文注释
- 新功能必须包含相应的测试用例
- 遵循 Go 语言的标准代码格式 (`go fmt`)
- 通过所有测试和 linter 检查

## 许可证

本项目采用 MIT 许可证。详细信息请参阅 [LICENSE](../../LICENSE) 文件。

## 联系方式

如有问题或建议，请通过以下方式联系：

- GitHub Issues: [提交问题](https://github.com/kugouming/mcpservers/issues)
- Email: your-email@example.com

## 更新日志

### v1.0.0 (2024-01-XX)
- ✨ 初始版本发布
- ✨ 支持获取接口列表、接口详情、项目信息
- ✨ 完整的测试覆盖
- ✨ 详细的中文文档

### 计划功能
- 🔄 支持接口数据缓存
- 🔄 支持批量操作
- 🔄 支持 WebSocket 实时更新
- �� 支持更多 YAPI API 端点 