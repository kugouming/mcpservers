package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kugouming/mcpservers/tools/thinkplan"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// 创建钩子
var hooks = &server.Hooks{}

// 初始化钩子
func init() {
	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		fmt.Printf("开始处理请求: ID=%v, 方法=%s\n", id, method)
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
		"ThinkPlan SSE Server",
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

// WithThinkPlanTool 注册思考和规划工具
func (s *MCPServer) WithThinkPlanTool() *MCPServer {
	thinkplan.RegisterTool(s.server)
	return s
}

// RegisterTool 注册工具的通用方法
func (s *MCPServer) RegisterTool(tool mcp.Tool, handler server.ToolHandlerFunc) *MCPServer {
	s.server.AddTool(tool, handler)
	return s
}

// 添加额外的API端点用于查看思考记录
func setupAdditionalRoutes(r *gin.Engine) {
	// 获取所有思考记录
	r.GET("/api/thinkplan/memory", func(c *gin.Context) {
		memory := thinkplan.GetMemory()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    memory,
			"count":   len(memory),
		})
	})

	// 获取思考记录摘要
	r.GET("/api/thinkplan/summary", func(c *gin.Context) {
		summary := thinkplan.GetSummary()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"summary": summary,
		})
	})

	// 根据编号获取特定记录
	r.GET("/api/thinkplan/memory/:number", func(c *gin.Context) {
		number := c.Param("number")
		entry := thinkplan.GetMemoryByNumber(number)
		if entry == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "思考记录未找到",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    entry,
		})
	})

	// 清空所有记录（用于测试）
	r.DELETE("/api/thinkplan/memory", func(c *gin.Context) {
		thinkplan.ClearMemory()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "所有思考记录已清空",
		})
	})

	// 获取JSON格式的所有记录
	r.GET("/api/thinkplan/export", func(c *gin.Context) {
		jsonData, err := thinkplan.GetMemoryAsJSON()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "导出失败",
				"error":   err.Error(),
			})
			return
		}
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=thinkplan_memory.json")
		c.String(http.StatusOK, jsonData)
	})
}

func main() {
	s := NewMCPServer().WithThinkPlanTool()

	// 设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin路由器
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(LogMiddleware())

	// 设置额外的API路由
	setupAdditionalRoutes(r)

	// 创建SSE服务器
	sseServer := server.NewSSEServer(s.server)

	// 将SSE服务器的SSE端点和处理函数集成到GIN路由中
	r.GET(sseServer.CompleteSsePath(), func(c *gin.Context) {
		sseServer.ServeHTTP(c.Writer, c.Request)
	})

	// 将SSE服务器的消息端点和处理函数集成到GIN路由中
	r.POST(sseServer.CompleteMessagePath(), func(c *gin.Context) {
		sseServer.ServeHTTP(c.Writer, c.Request)
	})

	// 设置根路径，显示服务信息
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":     "ThinkPlan MCP Server",
			"version":     "1.0.0",
			"description": "思考和规划工具服务器",
			"endpoints": gin.H{
				"sse":              sseServer.CompleteSsePath(),
				"message":          sseServer.CompleteMessagePath(),
				"memory":           "/api/thinkplan/memory",
				"summary":          "/api/thinkplan/summary",
				"memory_by_number": "/api/thinkplan/memory/:number",
				"clear_memory":     "DELETE /api/thinkplan/memory",
				"export_memory":    "/api/thinkplan/export",
			},
			"tool": gin.H{
				"name":        "think_and_plan",
				"description": "系统化思考与规划工具",
				"parameters":  []string{"thought", "plan", "action", "thoughtNumber"},
			},
		})
	})

	// 启动服务器
	port := ":8084"
	log.Printf("ThinkPlan MCP Server starting on port %s", port)
	log.Printf("SSE endpoint: http://localhost%s%s", port, sseServer.CompleteSsePath())
	log.Printf("Message endpoint: http://localhost%s%s", port, sseServer.CompleteMessagePath())
	log.Printf("API endpoints: http://localhost%s/api/thinkplan/*", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
