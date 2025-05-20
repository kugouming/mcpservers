package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 创建 MCP 服务器
	s := server.NewMCPServer(
		"提示词工具示例",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// 创建一个只返回提示词的工具
	promptTool := mcp.NewTool("get_prompt",
		mcp.WithDescription("获取特定场景的提示词"),
		mcp.WithString("scenario",
			mcp.Required(),
			mcp.Description("需要提示词的场景"),
			mcp.Enum("code_review", "bug_fix", "feature_request"),
		),
	)

	// 添加工具处理函数
	s.AddTool(promptTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scenario := request.Params.Arguments["scenario"].(string)

		var promptText string
		switch scenario {
		case "code_review":
			promptText = `作为代码审查者，请按照以下步骤审查代码：  
1. 检查代码的可读性和清晰度  
2. 识别潜在的性能问题  
3. 查找安全漏洞  
4. 评估代码是否符合项目的编码标准  
5. 提供具体的改进建议  
  
请提供详细的反馈，包括代码中的优点和需要改进的地方。`

		case "bug_fix":
			promptText = `作为调试专家，请按照以下步骤分析和修复bug：  
1. 描述问题的症状  
2. 分析可能的原因  
3. 提出解决方案  
4. 说明如何验证修复是否成功  
5. 建议如何防止类似问题再次发生  
  
请尽可能详细地解释每个步骤，以便其他开发人员能够理解问题和解决方案。`

		case "feature_request":
			promptText = `作为产品设计师，请按照以下步骤评估新功能请求：  
1. 分析功能的目标和价值  
2. 评估实现的复杂性和可行性  
3. 考虑与现有功能的集成  
4. 提出可能的实现方案  
5. 讨论潜在的用户体验影响  
  
请提供全面的分析，包括功能的优点、挑战和建议的实施路径。`

		default:
			return mcp.NewToolResultError("不支持的场景: " + scenario), nil
		}

		return mcp.NewToolResultText(promptText), nil
	})

	// 启动服务器
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("服务器错误: %v\n", err)
	}
}
