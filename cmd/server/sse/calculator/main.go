package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cast"
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

// RegisterCalculatorTool 注册计算器工具
func (s *MCPServer) RegisterCalculatorTool() *MCPServer {
	// 计算器参数声明
	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"), // 设置工具描述：执行基本算术运算
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"), // 设置操作参数：操作类型（加、减、乘、除）
			mcp.Enum("add", "subtract", "multiply", "divide"),                             // 设置操作类型枚举
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"), // 设置第一个数字参数：第一个数字
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"), // 设置第二个数字参数：第二个数字
		),
	)

	// 计算器处理程序
	calculatorHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := cast.ToString(request.Params.Arguments["operation"])
		x := cast.ToFloat64(request.Params.Arguments["x"])
		y := cast.ToFloat64(request.Params.Arguments["y"])

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return mcp.NewToolResultError("cannot divide by zero"), nil
			}
			result = x / y
		}

		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	}

	// 添加计算器处理程序
	s.server.AddTool(calculatorTool, calculatorHandler)

	return s
}

func main() {
	s := NewMCPServer()
	s = s.RegisterCalculatorTool()

	// Start the SSE server
	if err := server.NewSSEServer(s.server).Start(":8081"); err != nil {
		log.Fatalf("Failed to start SSE server: %v", err)
	} else {
		log.Printf("SSE server listening on :8080")
	}
}
