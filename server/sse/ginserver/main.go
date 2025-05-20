package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cast"
)

// 创建钩子
var hooks = &server.Hooks{}

// 初始化钩子
func init() {
	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		fmt.Printf("开始处理请求: ID=%v, 方法=%s\n", id, method)
	})

	// 添加 prompt 获取前钩子
	hooks.AddBeforeGetPrompt(func(ctx context.Context, id any, message *mcp.GetPromptRequest) {
		fmt.Printf("开始获取 prompt: ID=%v, Prompt名称=%s\n", id, message.Params.Name)
	})

	// 添加 prompt 获取后钩子
	hooks.AddAfterGetPrompt(func(ctx context.Context, id any, message *mcp.GetPromptRequest, result *mcp.GetPromptResult) {
		fmt.Printf("完成获取 prompt: ID=%v, Prompt名称=%s\n", id, message.Params.Name)
	})

	// 添加工具调用前钩子
	hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		fmt.Printf("开始调用工具: ID=%v, 工具名称=%s\n", id, message.Params.Name)
	})

	// 添加工具调用后钩子
	hooks.AddAfterCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		fmt.Printf("完成调用工具: ID=%v, 工具名称=%s\n", id, message.Params.Name)
	})

	// 添加成功钩子
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		fmt.Printf("请求成功: ID=%v, 方法=%s\n", id, method)
	})

	// 添加错误钩子
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		fmt.Printf("请求错误: ID=%v, 方法=%s, 错误=%v\n", id, method, err)
	})
}

type MCPServer struct {
	server *server.MCPServer
}

// LogWriter 用于捕获响应内容
type LogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应并同时记录
func (w LogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewMCPServer() *MCPServer {
	// 创建一个MCP服务器
	s := server.NewMCPServer(
		"Example Demo",
		"1.0.0",
		server.WithToolCapabilities(true),           // 启用工具相关的服务器功能
		server.WithPromptCapabilities(true),         // 启用提示相关的服务器功能
		server.WithResourceCapabilities(true, true), // 启用资源相关的服务器功能
		server.WithRecovery(),                       // 启用恢复机制，在发生错误时能够捕获异常
		server.WithLogging(),                        // 启用日志记录
		server.WithHooks(hooks),                     // 启用钩子
	)

	return &MCPServer{
		server: s,
	}
}

// LogMiddleware 记录请求和响应
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重置请求体，以便后续中间件和处理程序可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 记录请求信息
		log.Printf("\n\n[请求] %s %s\n请求头: %v\n请求体: %s\n\n",
			c.Request.Method, c.Request.URL.Path,
			c.Request.Header, string(requestBody))

		// 捕获响应
		responseBody := &bytes.Buffer{}
		logWriter := &LogWriter{ResponseWriter: c.Writer, body: responseBody}
		c.Writer = logWriter

		// 处理请求
		c.Next()

		// 记录响应信息
		log.Printf("\n\n[响应] %s %s\n状态码: %d\n响应头: %v\n响应体: %s\n处理时间: %v\n\n",
			c.Request.Method, c.Request.URL.Path,
			c.Writer.Status(), c.Writer.Header(),
			responseBody.String(), time.Since(startTime))
	}
}

// WithCalculatorTool 注册计算器工具
func (s *MCPServer) WithCalculatorTool() *MCPServer {
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
	return s.RegisterTool(calculatorTool, calculatorHandler)

}

// WithHttpRequestTool 注册HTTP请求工具
func (s *MCPServer) WithHttpRequestTool() *MCPServer {
	httpTool := mcp.NewTool("http_request",
		mcp.WithDescription("Make HTTP requests to external APIs"),
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
			mcp.Description("Request body (for POST/PUT)"),
		),
	)

	httpHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		method := cast.ToString(request.Params.Arguments["method"])
		url := cast.ToString(request.Params.Arguments["url"])

		log.Printf("request: [%T: %v]\n", request.Params.Arguments["headers"], request.Params.Arguments["headers"])
		// 尝试解析headers参数
		// headers := make(map[string]interface{})
		// if h, ok := request.Params.Arguments["headers"].(map[string]interface{}); ok {
		// 	headers = h
		// }
		// 尝试解析headers参数
		headers := make(map[string]string)
		headersStr := cast.ToString(request.Params.Arguments["headers"])
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
		if b, ok := request.Params.Arguments["body"].(string); ok {
			body = b
		}

		log.Printf("Tool input: \n%s %s\n%v\n\n %v\n\n", strings.ToUpper(method), url, headers, body)

		statusCode, responseBody, err := HttpRequest(ctx, method, url, headers, body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("执行请求失败", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Status: %d\nBody: %s", statusCode, string(responseBody))), nil
	}

	return s.RegisterTool(httpTool, httpHandler)
}

func (s *MCPServer) RegisterTool(tool mcp.Tool, handler server.ToolHandlerFunc) *MCPServer {
	s.server.AddTool(tool, handler)
	return s
}

func HttpRequest(ctx context.Context, method, url string, headers map[string]string, body string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
	if err != nil {
		return 0, nil, err
	}

	// 添加请求头
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, responseBody, nil
}

func main() {
	s := NewMCPServer()
	s.WithHttpRequestTool()
	s.WithCalculatorTool()

	// Start the SSE server with GIN
	r := gin.Default()

	// 添加日志中间件
	r.Use(LogMiddleware())

	sseServer := server.NewSSEServer(s.server)

	// 将 SSESever 的 SSE 端点和处理函数集成到 GIN 路由中
	r.GET(sseServer.CompleteSsePath(), func(c *gin.Context) {
		sseServer.ServeHTTP(c.Writer, c.Request)
	})

	// 将 SSESever 的消息端点和处理函数集成到 GIN 路由中
	r.POST(sseServer.CompleteMessagePath(), func(c *gin.Context) {
		sseServer.ServeHTTP(c.Writer, c.Request)
	})

	// 启动 GIN 服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start GIN server: %v", err)
	}
}
