package main

import (
	"log"
	"os"

	"github.com/kugouming/mcpservers/tools/httprequest"
	"github.com/kugouming/mcpservers/tools/sequentialthinking"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 创建MCP服务器实例
	s := server.NewMCPServer(
		"sequential-thinking-server",
		"0.2.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// 注册HTTP请求工具
	httprequest.RegisterTool(s)
	// 注册顺序思考工具
	sequentialthinking.RegisterTool(s)

	// 启动服务器
	log.Println("顺序思考MCP服务器启动中...")
	if err := server.ServeStdio(s); err != nil {
		log.Printf("服务器错误: %v\n", err)
		os.Exit(1)
	}
}
