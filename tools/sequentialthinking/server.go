package sequentialthinking

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SequentialThinkingServer 负责处理和管理思考步骤的核心逻辑
type SequentialThinkingServer struct {
	thoughtHistory []ThoughtData
	branches       map[string][]ThoughtData
}

// NewSequentialThinkingServer 创建一个新的顺序思考服务器实例
func NewSequentialThinkingServer() *SequentialThinkingServer {
	return &SequentialThinkingServer{
		thoughtHistory: make([]ThoughtData, 0),
		branches:       make(map[string][]ThoughtData),
	}
}

// validateThoughtData 验证输入的思考数据是否有效
func (s *SequentialThinkingServer) validateThoughtData(input ThoughtData) error {
	if input.Thought == "" {
		return fmt.Errorf("invalid thought: must not be empty")
	}
	if input.ThoughtNumber < 1 {
		return fmt.Errorf("invalid thoughtNumber: must be greater than 0")
	}
	if input.TotalThoughts < 1 {
		return fmt.Errorf("invalid totalThoughts: must be greater than 0")
	}
	return nil
}

// formatThought 格式化思考步骤，生成美观的控制台输出
func (s *SequentialThinkingServer) formatThought(thoughtData ThoughtData) string {
	var prefix, context string

	// 根据思考类型设置不同的前缀和上下文信息
	if thoughtData.IsRevision {
		prefix = "\033[33m🔄 Revision\033[0m" // 黄色
		context = fmt.Sprintf(" (revising thought %d)", thoughtData.RevisesThought)
	} else if thoughtData.BranchFromThought > 0 {
		prefix = "\033[32m🌿 Branch\033[0m" // 绿色
		context = fmt.Sprintf(" (from thought %d, ID: %s)", thoughtData.BranchFromThought, thoughtData.BranchID)
	} else {
		prefix = "\033[34m💭 Thought\033[0m" // 蓝色
		context = ""
	}

	// 创建思考步骤的标题
	header := fmt.Sprintf("%s %d/%d%s", prefix, thoughtData.ThoughtNumber, thoughtData.TotalThoughts, context)

	// 创建边框，边框长度与标题或思考内容的最大长度相匹配
	borderLength := max(len(header), len(thoughtData.Thought)) + 4
	border := strings.Repeat("─", borderLength)

	// 返回格式化后的思考步骤，包含边框和内容
	return fmt.Sprintf(`
┌%s┐
│ %s │
├%s┤
│ %s │
└%s┘`, border, header, border, padRight(thoughtData.Thought, borderLength-2), border)
}

// ProcessThought 处理思考步骤输入并返回结果
func (s *SequentialThinkingServer) ProcessThought(input ThoughtData) ThoughtResponse {
	// 验证输入数据
	if err := s.validateThoughtData(input); err != nil {
		return s.createErrorResponse(err.Error())
	}

	// 如果当前思考序号超过了估计的总数，则更新总数
	if input.ThoughtNumber > input.TotalThoughts {
		input.TotalThoughts = input.ThoughtNumber
	}

	// 将思考步骤添加到历史记录
	s.thoughtHistory = append(s.thoughtHistory, input)

	// 如果是分支思考，将其添加到相应的分支集合中
	if input.BranchFromThought > 0 && input.BranchID != "" {
		if _, exists := s.branches[input.BranchID]; !exists {
			s.branches[input.BranchID] = make([]ThoughtData, 0)
		}
		s.branches[input.BranchID] = append(s.branches[input.BranchID], input)
	}

	// 格式化思考步骤并输出到控制台
	formattedThought := s.formatThought(input)
	fmt.Fprintln(stderr, formattedThought)

	// 返回处理结果，包含当前思考状态的JSON
	status := ThoughtStatus{
		ThoughtNumber:        input.ThoughtNumber,
		TotalThoughts:        input.TotalThoughts,
		NextThoughtNeeded:    input.NextThoughtNeeded,
		Branches:             make([]string, 0, len(s.branches)),
		ThoughtHistoryLength: len(s.thoughtHistory),
	}

	// 收集所有分支ID
	for branchID := range s.branches {
		status.Branches = append(status.Branches, branchID)
	}

	statusJSON, _ := json.MarshalIndent(status, "", "  ")
	return ThoughtResponse{
		Content: []ContentItem{{
			Type: "text",
			Text: string(statusJSON),
		}},
	}
}

// createErrorResponse 创建错误响应
func (s *SequentialThinkingServer) createErrorResponse(errMsg string) ThoughtResponse {
	errResp := ErrorResponse{
		Error:  errMsg,
		Status: "failed",
	}
	errJSON, _ := json.MarshalIndent(errResp, "", "  ")
	return ThoughtResponse{
		Content: []ContentItem{{
			Type: "text",
			Text: string(errJSON),
		}},
		IsError: true,
	}
}

// padRight 在字符串右侧填充空格到指定长度
func padRight(str string, length int) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(" ", length-len(str))
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 为了测试目的，我们可以模拟标准错误输出
var stderr = &mockWriter{}

type mockWriter struct{}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
