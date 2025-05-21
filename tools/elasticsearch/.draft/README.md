# Elasticsearch MCP Server

这是一个使用 Go 语言实现的 Elasticsearch MCP 服务器，它提供了与 Elasticsearch 交互的工具和资源。

## 功能特性

- 列出所有 Elasticsearch 索引
- 获取特定索引的详细信息
- 批量写入文档到 Elasticsearch
- 获取 Elasticsearch 容器日志
- 访问配置文件
- 提供索引分析提示模板

## 环境要求

- Go 1.21 或更高版本
- Elasticsearch 7.x
- Docker（用于运行 Elasticsearch）

## 配置

在项目根目录创建 `.env` 文件，包含以下配置：

```env
ELASTIC_HOST=http://localhost:9200
ELASTIC_USERNAME=elastic
ELASTIC_PASSWORD=your_password
```

## 安装

```bash
# 克隆仓库
git clone https://github.com/kugouming/mcpservers.git
cd mcpservers/tools/elasticsearch

# 安装依赖
go mod download

# 编译
go build -o elasticsearch-mcp-server
```

## 运行

```bash
./elasticsearch-mcp-server
```

## 使用方法

服务器启动后，可以通过 MCP 客户端访问以下功能：

### 工具函数

- `list_indices`: 列出所有 Elasticsearch 索引
- `get_index`: 获取特定索引的详细信息
- `write_documents`: 批量写入文档到 Elasticsearch

### 资源

- `es://logs`: 获取 Elasticsearch 容器日志
- `file://docker-compose.yaml`: 获取 docker-compose.yaml 文件内容
- `file://movies.csv`: 获取 movies.csv 文件内容

### 提示模板

- `es_prompt`: 创建索引分析提示

## 示例

```go
// 列出所有索引
indices, err := mcpClient.CallTool("list_indices", nil)

// 获取特定索引信息
indexInfo, err := mcpClient.CallTool("get_index", map[string]interface{}{
    "index": "my_index",
})

// 写入文档
docs := []map[string]interface{}{
    {"title": "示例文档1", "content": "内容1"},
    {"title": "示例文档2", "content": "内容2"},
}
result, err := mcpClient.CallTool("write_documents", map[string]interface{}{
    "index": "my_index",
    "documents": docs,
})
```

## 许可证

MIT 