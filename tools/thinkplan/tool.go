package thinkplan

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cast"
)

// ThinkPlanEntry 表示一个思考和规划条目
type ThinkPlanEntry struct {
	ThoughtNumber string    `json:"thoughtNumber"`
	Thought       string    `json:"thought"`
	Plan          string    `json:"plan"`
	Action        string    `json:"action"`
	Timestamp     time.Time `json:"timestamp"`
}

// ThinkPlanMemory 存储所有的思考和规划记录
var ThinkPlanMemory []ThinkPlanEntry

// memoryMutex 保护ThinkPlanMemory的并发访问
var memoryMutex sync.RWMutex

// RegisterTool 注册ThinkPlan工具
func RegisterTool(s *server.MCPServer) {
	s.AddTool(thinkPlanTool, thinkPlanHandler)
}

// thinkPlanTool 定义思考和规划工具
var thinkPlanTool = mcp.NewTool("think_and_plan",
	mcp.WithDescription("这是用于系统化思考与规划的工具，支持用户在面对复杂问题或任务时，分阶段梳理思考、规划和行动步骤。工具强调思考（thought）、计划（plan）与实际行动（action）的结合，通过编号（thoughtNumber）追踪过程。该工具不会获取新信息或更改数据库，只会将想法附加到记忆中。当需要复杂推理或某种缓存记忆时，可以使用它。"),
	mcp.WithString("thought",
		mcp.Required(),
		mcp.Description("当前的思考内容，可以是对问题的分析、假设、洞见、反思或对前一步骤的总结。强调深度思考和逻辑推演，是每一步的核心。"),
	),
	mcp.WithString("plan",
		mcp.Required(),
		mcp.Description("针对当前任务拟定的计划或方案，将复杂问题分解为多个可执行步骤。执行步骤以有序列表的形式给出，每个步骤用数字编号。"),
	),
	mcp.WithString("action",
		mcp.Required(),
		mcp.Description("基于当前思考和计划，建议下一步采取的行动步骤，要求具体、可执行、可验证，可以是下一步需要调用的一个或多个工具。执行步骤以有序列表的形式给出，每个步骤用数字编号。"),
	),
	mcp.WithString("thoughtNumber",
		mcp.Required(),
		mcp.Description("当前思考步骤的编号，用于追踪和回溯整个思考与规划过程，便于后续复盘与优化。"),
	),
)

// thinkPlanHandler 处理思考和规划请求
func thinkPlanHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 获取参数
	thought := cast.ToString(request.GetArguments()["thought"])
	plan := cast.ToString(request.GetArguments()["plan"])
	action := cast.ToString(request.GetArguments()["action"])
	thoughtNumber := cast.ToString(request.GetArguments()["thoughtNumber"])

	// 验证必需参数
	if len(thought) == 0 {
		return mcp.NewToolResultError("thought parameter is required"), nil
	}
	if len(plan) == 0 {
		return mcp.NewToolResultError("plan parameter is required"), nil
	}
	if len(action) == 0 {
		return mcp.NewToolResultError("action parameter is required"), nil
	}
	if len(thoughtNumber) == 0 {
		return mcp.NewToolResultError("thoughtNumber parameter is required"), nil
	}

	// 创建新的思考和规划条目
	entry := ThinkPlanEntry{
		ThoughtNumber: thoughtNumber,
		Thought:       thought,
		Plan:          plan,
		Action:        action,
		Timestamp:     time.Now(),
	}

	// 使用互斥锁保护内存访问
	memoryMutex.Lock()
	ThinkPlanMemory = append(ThinkPlanMemory, entry)
	currentCount := len(ThinkPlanMemory)
	memoryMutex.Unlock()

	// 记录到服务器日志
	log.Printf("ThinkPlan Entry [%s]: Thought=%s, Plan=%s, Action=%s",
		thoughtNumber, thought, plan, action)

	// 构建响应内容
	response := fmt.Sprintf(`思考和规划记录 [%s]

🤔 思考内容:
%s

📋 规划方案:
%s

🎯 下一步行动:
%s

⏰ 记录时间: %s

📊 当前已记录 %d 个思考步骤`,
		thoughtNumber, thought, plan, action,
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		currentCount)

	return mcp.NewToolResultText(response), nil
}

// GetMemory 获取所有思考和规划记录
func GetMemory() []ThinkPlanEntry {
	memoryMutex.RLock()
	defer memoryMutex.RUnlock()

	// 返回副本以避免外部修改
	result := make([]ThinkPlanEntry, len(ThinkPlanMemory))
	copy(result, ThinkPlanMemory)
	return result
}

// GetMemoryByNumber 根据编号获取特定的思考和规划记录
func GetMemoryByNumber(thoughtNumber string) *ThinkPlanEntry {
	memoryMutex.RLock()
	defer memoryMutex.RUnlock()

	for _, entry := range ThinkPlanMemory {
		if entry.ThoughtNumber == thoughtNumber {
			// 返回副本以避免外部修改
			entryCopy := entry
			return &entryCopy
		}
	}
	return nil
}

// GetMemoryAsJSON 以JSON格式获取所有记录
func GetMemoryAsJSON() (string, error) {
	memoryMutex.RLock()
	defer memoryMutex.RUnlock()

	data, err := json.MarshalIndent(ThinkPlanMemory, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ClearMemory 清空所有记录（用于测试或重置）
func ClearMemory() {
	memoryMutex.Lock()
	defer memoryMutex.Unlock()

	ThinkPlanMemory = []ThinkPlanEntry{}
	log.Println("ThinkPlan memory cleared")
}

// GetSummary 获取思考和规划过程的摘要
func GetSummary() string {
	memoryMutex.RLock()
	defer memoryMutex.RUnlock()

	if len(ThinkPlanMemory) == 0 {
		return "暂无思考和规划记录"
	}

	summary := fmt.Sprintf("思考和规划过程摘要 (共 %d 个步骤):\n\n", len(ThinkPlanMemory))

	for i, entry := range ThinkPlanMemory {
		summary += fmt.Sprintf("%d. [%s] %s\n",
			i+1, entry.ThoughtNumber, entry.Thought)
		if i < len(ThinkPlanMemory)-1 {
			summary += "\n"
		}
	}

	return summary
}
