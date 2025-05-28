package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
)

type SSEServer struct {
	server *server.SSEServer
}

func NewSSEServer(sev *server.MCPServer) *SSEServer {
	sse := server.NewSSEServer(sev,
		server.WithStaticBasePath("/mcp"),
		server.WithKeepAlive(true),
		// 此选项对 HandleSSE2 方法无效
		server.WithDynamicBasePath(func(r *http.Request, sessionID string) string {
			// 这里可以根据请求头、环境变量、租户信息等动态拼接前缀
			if getEnv() != "dev" {
				gatewayPrefix := "/supplygoods"
				staticBase := "/mcp"
				return gatewayPrefix + staticBase
			}
			return "/mcp"
		}),
		// server.WithUseFullURLForMessageEndpoint(true),
		// server.WithBaseURL("http://localhost:8080/supplygoods/"),
	)
	return &SSEServer{
		server: sse,
	}
}

func (sse *SSEServer) SsePath() string {
	return sse.server.CompleteSsePath()
}

func (sse *SSEServer) MessagePath() string {
	return sse.server.CompleteMessagePath()
}

func (sse *SSEServer) HandleSSE(c *gin.Context) {
	func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/sse") {
			sse.server.SSEHandler().ServeHTTP(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/message") {
			sse.server.MessageHandler().ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	}(c.Writer, c.Request)
}

func (sse *SSEServer) HandleSSE2(c *gin.Context) {
	sse.server.ServeHTTP(c.Writer, c.Request)
}

func getEnv() string {
	return os.Getenv("ENV")
}
