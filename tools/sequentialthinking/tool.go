package sequentialthinking

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cast"
)

// RegisterTool 注册顺序思考工具
func RegisterTool(s *server.MCPServer) {
	s.AddTool(sequentialThinkingTool, toolHandler)
}

// sequentialThinkingTool 定义了顺序思考工具的配置
var sequentialThinkingTool = mcp.NewTool("sequentialthinking",
	mcp.WithDescription(`一个用于动态和反思性解决问题的详细工具。
这个工具通过灵活的思考过程帮助分析问题，可以适应和进化。
每个思考都可以基于、质疑或修改之前的见解，随着理解的深入而发展。

使用场景：
- 将复杂问题分解为步骤
- 需要修改的计划和设计
- 可能需要修正的分析
- 初始范围不明确的问题
- 需要多步骤解决方案的问题
- 需要在多个步骤中保持上下文的任务
- 需要过滤无关信息的情况

主要特点：
- 可以随着进展调整总思考步骤数
- 可以质疑或修改之前的思考
- 即使在看似结束时也可以添加更多思考
- 可以表达不确定性并探索替代方案
- 思考不必线性构建 - 可以分支或回溯
- 生成解决方案假设
- 基于思考链验证假设
- 重复过程直到满意
- 提供正确答案

参数说明：
- thought(思考): 当前思考步骤, 可以包括：
* 常规分析步骤
* 对之前思考的修改
* 对之前决策的质疑
* 需要更多分析的认知
* 方法的改变
* 假设的生成
* 假设的验证
- next_thought_needed(需要下一步思考): 如果需要更多思考则为true, 即使在看似结束时
- thought_number(思考编号): 序列中的当前编号(如果需要可以超过初始总数)
- total_thoughts(总思考数): 当前估计需要的思考数(可以上下调整)
- is_revision(是否修改): 布尔值, 表示此思考是否修改之前的思考
- revises_thought(修改的思考): 如果is_revision为true, 表示正在重新考虑哪个思考编号
- branch_from_thought(分支来源): 如果分支, 表示分支点的思考编号
- branch_id(分支标识): 当前分支的标识符(如果有)
- needs_more_thoughts(需要更多思考): 如果到达结束但意识到需要更多思考

你应该：
1. 从初始估计的思考数开始, 但准备好调整
2. 随时质疑或修改之前的思考
3. 如果需要, 即使在"结束"时也不要犹豫添加更多思考
4. 在存在不确定性时表达出来
5. 标记修改之前思考或分支到新路径的思考
6. 忽略与当前步骤无关的信息
7. 在适当时生成解决方案假设
8. 基于思考链验证假设
9. 重复过程直到对解决方案满意
10. 提供单一、理想的正确答案作为最终输出
11. 只有在真正完成并达到满意答案时才将next_thought_needed设置为false`),
	mcp.WithString("thought",
		mcp.Required(),
		mcp.Description("当前思考步骤"),
	),
	mcp.WithBoolean("nextThoughtNeeded",
		mcp.Required(),
		mcp.Description("是否需要另一个思考步骤"),
	),
	mcp.WithNumber("thoughtNumber",
		mcp.Required(),
		mcp.Description("当前思考编号"),
	),
	mcp.WithNumber("totalThoughts",
		mcp.Required(),
		mcp.Description("估计需要的总思考数"),
	),
	mcp.WithBoolean("isRevision",
		mcp.Description("是否修改之前的思考"),
	),
	mcp.WithNumber("revisesThought",
		mcp.Description("正在重新考虑的思考编号"),
	),
	mcp.WithNumber("branchFromThought",
		mcp.Description("分支点的思考编号"),
	),
	mcp.WithString("branchId",
		mcp.Description("分支标识符"),
	),
	mcp.WithBoolean("needsMoreThoughts",
		mcp.Description("是否需要更多思考"),
	),
)

func toolHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析请求参数并转换为ThoughtData结构体
	thoughtData := ThoughtData{
		Thought:           cast.ToString(req.Params.Arguments["thought"]),
		ThoughtNumber:     cast.ToInt(req.Params.Arguments["thoughtNumber"]),
		TotalThoughts:     cast.ToInt(req.Params.Arguments["totalThoughts"]),
		NextThoughtNeeded: cast.ToBool(req.Params.Arguments["nextThoughtNeeded"]),
		IsRevision:        cast.ToBool(req.Params.Arguments["isRevision"]),
		RevisesThought:    cast.ToInt(req.Params.Arguments["revisesThought"]),
		BranchFromThought: cast.ToInt(req.Params.Arguments["branchFromThought"]),
		BranchID:          cast.ToString(req.Params.Arguments["branchId"]),
		NeedsMoreThoughts: cast.ToBool(req.Params.Arguments["needsMoreThoughts"]),
	}

	// 处理思考步骤
	response := NewSequentialThinkingServer().ProcessThought(thoughtData)

	// 将响应转换为MCP工具结果
	if response.IsError {
		return mcp.NewToolResultError(response.Content[0].Text), nil
	}

	return mcp.NewToolResultText(response.Content[0].Text), nil
}
