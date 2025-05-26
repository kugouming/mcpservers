package think

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cast"
)

// RegisterTool 注册Think工具
func RegisterTool(s *server.MCPServer) {
	s.AddTool(thinkTool, thinkHandler)
}

// Add the "think" tool
var thinkTool = mcp.NewTool("think",
	mcp.WithDescription("Use the tool to think about something. It will not obtain new information or change the database, but just append the thought to the log. Use it when complex reasoning or some cache memory is needed."),
	mcp.WithString("thought",
		mcp.Required(),
		mcp.Description("A thought to think about."),
	),
)

func thinkHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 使用 RequireString 方法获取必需的 thought 参数
	thought := cast.ToString(request.Params.Arguments["thought"])

	if len(thought) == 0 {
		return mcp.NewToolResultError("thought parameter is required"), nil
	}

	// 记录思考过程（这将在服务器日志中可见，但不会发送给用户）
	log.Printf("Thinking process: %s", thought)

	// 简单返回思考内容本身，正如 Anthropic 博客文章中提到的
	return mcp.NewToolResultText(thought), nil
}
