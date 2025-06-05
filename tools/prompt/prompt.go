package prompt

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTool 注册所有 prompt 工具到 MCP 服务器
func RegisterTool(s *server.MCPServer) {
	tpls, err := LoadAllPromptTemplates()
	if err != nil {
		panic(fmt.Sprintf("加载prompt模板失败: %v", err))
	}

	for _, tpl := range tpls {
		tpl := tpl // 闭包安全
		tool := buildMcpTool(tpl)
		handler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return runPromptTemplateMcp(tpl, req)
		}
		s.AddTool(*tool, handler)
	}

	// 管理工具: reload_prompts
	s.AddTool(mcp.NewTool("reload_prompts", mcp.WithDescription("重新加载所有 prompt 模板")),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			err := ReloadPromptTemplates()
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText("已重新加载所有 prompt 模板"), nil
		})

	// 管理工具: get_prompt_names
	s.AddTool(mcp.NewTool("get_prompt_names", mcp.WithDescription("获取所有可用的 prompt 名称")),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			names := []string{}
			for _, tpl := range defaultPromptTemplates {
				names = append(names, tpl.Name)
			}
			return mcp.NewToolResultText(strings.Join(names, ", ")), nil
		})
}

// buildMcpTool 构建 mcp.Tool
func buildMcpTool(tpl *PromptTemplate) *mcp.Tool {
	opts := []mcp.ToolOption{
		mcp.WithDescription(tpl.Description),
	}
	for _, a := range tpl.Arguments {
		if a.Required {
			opts = append(opts, mcp.WithString(a.Name, mcp.Required(), mcp.Description(a.Description)))
		} else {
			opts = append(opts, mcp.WithString(a.Name, mcp.Description(a.Description)))
		}
	}
	tool := mcp.NewTool(tpl.Name, opts...)
	return &tool
}

// runPromptTemplateMcp 执行模板，参数替换
func runPromptTemplateMcp(tpl *PromptTemplate, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	// 校验必填参数
	for _, a := range tpl.Arguments {
		if a.Required {
			if _, ok := args[a.Name]; !ok {
				return mcp.NewToolResultError("缺少必需参数: " + a.Name), nil
			}
		}
	}
	// 构建消息
	var messages []string
	for _, msg := range tpl.Messages {
		text := msg.Content.Text
		for _, a := range tpl.Arguments {
			if v, ok := args[a.Name]; ok {
				text = strings.ReplaceAll(text, "{{"+a.Name+"}}", fmt.Sprintf("%v", v))
			}
		}
		messages = append(messages, fmt.Sprintf("[%s] %s", msg.Role, text))
	}
	return mcp.NewToolResultText(strings.Join(messages, "\n")), nil
}
