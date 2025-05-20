package main

import (
	"log"

	"github.com/kugouming/mcpservers/tools/httprequest"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	server *server.MCPServer
}

func NewMCPServer() *MCPServer {
	// 创建一个MCP服务器
	s := server.NewMCPServer(
		"Example Demo",
		"1.0.0",
		server.WithToolCapabilities(true), // 启用工具相关的服务器功能
		server.WithResourceCapabilities(true, true), // 启用资源相关的服务器功能
		server.WithPromptCapabilities(true),         // 启用提示相关的服务器功能
		server.WithRecovery(),                       // 启用恢复机制，在发生错误时能够捕获异常
	)

	return &MCPServer{
		server: s,
	}
}

func (s *MCPServer) WithTools() *MCPServer {
	httprequest.RegisterTool(s.server)

	return s
}

func main() {
	s := NewMCPServer().WithTools()

	// Start the server
	if err := server.ServeStdio(s.server); err != nil { // 同：server.NewStdioServer(s.server).Listen(context.Background(), os.Stdin, os.Stdout)
		log.Printf("Server error: %v\n", err)
	} else {
		log.Printf("Server listening on stdin/stdout")
	}

}
