package thinkplan

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func TestThinkPlanHandler(t *testing.T) {
	// 清空内存以确保测试的独立性
	ClearMemory()

	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name: "成功的思考和规划",
			arguments: map[string]interface{}{
				"thought":       "分析当前问题的核心是什么",
				"plan":          "1. 收集信息 2. 分析问题 3. 制定解决方案",
				"action":        "开始收集相关资料和数据",
				"thoughtNumber": "T001",
			},
			expectError: false,
		},
		{
			name: "缺少thought参数",
			arguments: map[string]interface{}{
				"plan":          "制定计划",
				"action":        "执行行动",
				"thoughtNumber": "T002",
			},
			expectError: true,
			errorMsg:    "thought parameter is required",
		},
		{
			name: "缺少plan参数",
			arguments: map[string]interface{}{
				"thought":       "思考内容",
				"action":        "执行行动",
				"thoughtNumber": "T003",
			},
			expectError: true,
			errorMsg:    "plan parameter is required",
		},
		{
			name: "缺少action参数",
			arguments: map[string]interface{}{
				"thought":       "思考内容",
				"plan":          "制定计划",
				"thoughtNumber": "T004",
			},
			expectError: true,
			errorMsg:    "action parameter is required",
		},
		{
			name: "缺少thoughtNumber参数",
			arguments: map[string]interface{}{
				"thought": "思考内容",
				"plan":    "制定计划",
				"action":  "执行行动",
			},
			expectError: true,
			errorMsg:    "thoughtNumber parameter is required",
		},
		{
			name: "空字符串参数",
			arguments: map[string]interface{}{
				"thought":       "",
				"plan":          "制定计划",
				"action":        "执行行动",
				"thoughtNumber": "T005",
			},
			expectError: true,
			errorMsg:    "thought parameter is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: struct {
					Name      string         `json:"name"`
					Arguments map[string]any `json:"arguments,omitempty"`
					Meta      *mcp.Meta      `json:"_meta,omitempty"`
				}{
					Name:      "think_and_plan",
					Arguments: tt.arguments,
				},
			}

			result, err := thinkPlanHandler(context.Background(), request)

			if tt.expectError {
				assert.NotNil(t, result)
				assert.Nil(t, err)
				assert.True(t, result.IsError)
				// 检查第一个内容项是否为TextContent
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(mcp.TextContent); ok {
						assert.Contains(t, textContent.Text, tt.errorMsg)
					}
				}
			} else {
				assert.NotNil(t, result)
				assert.Nil(t, err)
				assert.False(t, result.IsError)
				assert.NotEmpty(t, result.Content)
				// 检查第一个内容项是否为TextContent
				if len(result.Content) > 0 {
					if textContent, ok := result.Content[0].(mcp.TextContent); ok {
						assert.NotEmpty(t, textContent.Text)
					}
				}
			}
		})
	}
}

func TestThinkPlanMemoryFunctions(t *testing.T) {
	// 清空内存
	ClearMemory()

	// 测试初始状态
	assert.Empty(t, GetMemory())
	assert.Equal(t, "暂无思考和规划记录", GetSummary())

	// 添加一些测试数据
	request1 := mcp.CallToolRequest{
		Params: struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments,omitempty"`
			Meta      *mcp.Meta      `json:"_meta,omitempty"`
		}{
			Name: "think_and_plan",
			Arguments: map[string]interface{}{
				"thought":       "第一个思考：分析问题的本质",
				"plan":          "1. 定义问题 2. 收集数据 3. 分析原因",
				"action":        "开始定义问题的边界和范围",
				"thoughtNumber": "T001",
			},
		},
	}

	request2 := mcp.CallToolRequest{
		Params: struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments,omitempty"`
			Meta      *mcp.Meta      `json:"_meta,omitempty"`
		}{
			Name: "think_and_plan",
			Arguments: map[string]interface{}{
				"thought":       "第二个思考：基于收集的数据进行深入分析",
				"plan":          "1. 数据清洗 2. 模式识别 3. 假设验证",
				"action":        "使用统计工具分析数据模式",
				"thoughtNumber": "T002",
			},
		},
	}

	// 执行添加操作
	_, err1 := thinkPlanHandler(context.Background(), request1)
	_, err2 := thinkPlanHandler(context.Background(), request2)

	assert.Nil(t, err1)
	assert.Nil(t, err2)

	// 测试GetMemory
	memory := GetMemory()
	assert.Len(t, memory, 2)
	assert.Equal(t, "T001", memory[0].ThoughtNumber)
	assert.Equal(t, "T002", memory[1].ThoughtNumber)

	// 测试GetMemoryByNumber
	entry1 := GetMemoryByNumber("T001")
	assert.NotNil(t, entry1)
	assert.Equal(t, "第一个思考：分析问题的本质", entry1.Thought)

	entry3 := GetMemoryByNumber("T003")
	assert.Nil(t, entry3)

	// 测试GetSummary
	summary := GetSummary()
	assert.Contains(t, summary, "思考和规划过程摘要 (共 2 个步骤)")
	assert.Contains(t, summary, "[T001] 第一个思考：分析问题的本质")
	assert.Contains(t, summary, "[T002] 第二个思考：基于收集的数据进行深入分析")

	// 测试GetMemoryAsJSON
	jsonData, err := GetMemoryAsJSON()
	assert.Nil(t, err)
	assert.Contains(t, jsonData, "T001")
	assert.Contains(t, jsonData, "T002")
	assert.Contains(t, jsonData, "thoughtNumber")

	// 测试ClearMemory
	ClearMemory()
	assert.Empty(t, GetMemory())
	assert.Equal(t, "暂无思考和规划记录", GetSummary())
}

func TestThinkPlanTool(t *testing.T) {
	// 测试工具定义
	assert.Equal(t, "think_and_plan", thinkPlanTool.Name)
	assert.NotEmpty(t, thinkPlanTool.Description)

	// 检查必需的参数
	inputSchema := thinkPlanTool.InputSchema
	properties := inputSchema.Properties

	// 验证所有必需参数都存在
	requiredParams := []string{"thought", "plan", "action", "thoughtNumber"}
	for _, param := range requiredParams {
		assert.Contains(t, properties, param)
		assert.Contains(t, inputSchema.Required, param)
	}
}

func BenchmarkThinkPlanHandler(b *testing.B) {
	ClearMemory()

	request := mcp.CallToolRequest{
		Params: struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments,omitempty"`
			Meta      *mcp.Meta      `json:"_meta,omitempty"`
		}{
			Name: "think_and_plan",
			Arguments: map[string]interface{}{
				"thought":       "性能测试思考",
				"plan":          "执行性能基准测试",
				"action":        "运行基准测试并收集结果",
				"thoughtNumber": "BENCH001",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := thinkPlanHandler(context.Background(), request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 测试并发安全性
func TestConcurrentAccess(t *testing.T) {
	ClearMemory()

	// 启动多个goroutine同时添加记录
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			request := mcp.CallToolRequest{
				Params: struct {
					Name      string         `json:"name"`
					Arguments map[string]any `json:"arguments,omitempty"`
					Meta      *mcp.Meta      `json:"_meta,omitempty"`
				}{
					Name: "think_and_plan",
					Arguments: map[string]interface{}{
						"thought":       fmt.Sprintf("并发思考 %d", id),
						"plan":          fmt.Sprintf("并发计划 %d", id),
						"action":        fmt.Sprintf("并发行动 %d", id),
						"thoughtNumber": fmt.Sprintf("CONCURRENT_%d", id),
					},
				},
			}
			_, err := thinkPlanHandler(context.Background(), request)
			assert.Nil(t, err)
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有记录都被正确添加
	memory := GetMemory()
	assert.Len(t, memory, 10)
}
