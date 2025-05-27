package httprequest

import (
	"context"
	"fmt"
	"strings"

	"github.com/kugouming/mcpservers/helper"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cast"
)

// RegisterTool 注册HTTP请求工具
func RegisterTool(s *server.MCPServer) {
	s.AddTool(httpTool, httpHandler)
}

// httpTool 定义了HTTP请求工具的配置
var httpTool = mcp.NewTool("http_request",
	mcp.WithDescription(`HTTP 请求解析工具，对用户输入的内容首先进行参数解析，再对解析后的参数进行 HTTP 请求。
- 支持从用户输入内容中按照不同格式解析出请求所需参数;
  * 根据用户的输入先识别出输入的数据格式, 如: CURL 命令、REST 风格请求等数据;
  * 再根据数据格式找到与下面提供的格式说明中对应的格式, 提取出请求所需参数;
- 对解析出的HTTP请求参数发起HTTP请求;
- 禁止对请求结果做解析, 直接按照原样(json格式)输出。

## CURL 命令示例

示例输入：
curl -X POST https://api.example.com/data -H "Content-Type: application/json" -d '{"key": "value"}'
	  
输出参数：
{
	"method": "POST",
	"url": "https://api.example.com/data",
	"headers": "Content-Type: application/json",
	"body": "{\"key\": \"value\"}"
}

## REST 风格请求示例 1:

示例输入：
POST https://api.example.com/data
Content-Type: application/json

{"key": "value"}

示例输出：
{
	"method": "POST",
	"url": "https://api.example.com/data",
	"headers": "Content-Type: application/json",
	"body": "{\"key\": \"value\"}"
}

## REST 风格请求示例 2:

示例输入：
POST https://api.example.com/data
Content-Type: application/json

{
	"key": "value"
}

示例输出：
{
	"method": "POST",
	"url": "https://api.example.com/data",
	"headers": "Content-Type: application/json",
	"body": "{\"key\": \"value\"}"
}

错误输出(重点规避: headers数据丢失; body中末尾大括号丢失)：
{
	"method": "POST",
	"url": "https://api.example.com/data",
	"headers": "",
	"body": "{\"key\": \"value\""
}
`),
	mcp.WithString("method",
		mcp.Required(),
		mcp.Description("HTTP method to use"),
		mcp.Enum("GET", "POST", "PUT", "DELETE"),
	),
	mcp.WithString("url",
		mcp.Required(),
		mcp.Description("URL to send the request to"),
		mcp.Pattern("^https?://.*"),
	),
	mcp.WithString("headers",
		mcp.Description("Request headers"),
	),
	mcp.WithString("body",
		mcp.Description("Request body"),
	),
)

func httpHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := cast.ToString(request.GetArguments()["method"])
	url := cast.ToString(request.GetArguments()["url"])

	// 尝试解析headers参数
	headers := make(map[string]string)
	headersStr := cast.ToString(request.GetArguments()["headers"])
	if headersStr != "" {
		// 这里可以添加更复杂的header解析逻辑
		// 简单处理：假设格式为 "Key1: Value1\nKey2: Value2"
		headerLines := strings.Split(headersStr, "\n")
		for _, line := range headerLines {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if key != "" && value != "" {
					headers[key] = value
				}
			}
		}
	}

	// 尝试解析body参数
	body := ""
	if b, ok := request.GetArguments()["body"].(string); ok {
		body = strings.ReplaceAll(b, "\"\"", "\"")
	}

	statusCode, responseBody, err := helper.HttpRequest(ctx, method, url, headers, body)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("执行请求失败", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Status: %d\nBody: %s", statusCode, string(responseBody))), nil
}
