package thinkandplan

import (
	"context"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func cleanPlanFile() {
	os.Remove(PlanFile)
}

func TestThinkAndPlan_AllFeatures(t *testing.T) {
	cleanPlanFile()
	defer cleanPlanFile()

	// 1. 创建新任务
	msg, err := ThinkAndPlan("测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "已为 '测试任务1' 创建新任务计划")

	// 2. 添加步骤
	msg, err = AddStep("步骤A", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "已为任务 '测试任务1' 添加步骤")

	// 3. 标记步骤完成
	msg, err = MarkStepComplete("步骤A", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "已将步骤 '步骤A' 标记为完成")

	// 4. 添加问题
	msg, err = AddIssue("有bug", "步骤A", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "已为步骤 '步骤A' 添加问题")

	// 5. 解决问题
	msg, err = ResolveIssue("步骤A", "已修复", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "已为步骤 '步骤A' 的问题添加解决说明")

	// 6. 更新规划说明
	msg, err = UpdatePlanningNotes("新规划说明", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "已更新任务 '测试任务1' 的规划说明")

	// 7. 检查任务完成度
	msg, err = CheckTaskCompletion("测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "任务 '测试任务1' 完成度：")

	// 8. 设置优先级
	msg, err = SetPriority("high", "步骤A", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "已为步骤 '步骤A' 设置优先级")

	msg, err = SetPriority("medium", "", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "已为任务 '测试任务1' 设置优先级")

	// 9. 删除步骤
	msg, err = DeleteStep("步骤A", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "已删除步骤 '步骤A'")

	// 10. 删除任务
	msg, err = DeleteTask("测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "已删除任务 '测试任务1'")
}

func TestMCPHandlers(t *testing.T) {
	cleanPlanFile()
	defer cleanPlanFile()

	ctx := context.Background()

	// think_and_plan
	result, err := thinkAndPlanHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "think_and_plan",
			Arguments: map[string]any{"task_description": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)

	// add_step
	result, err = addStepHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "add_step",
			Arguments: map[string]any{"step_description": "MCP步骤A", "task_title": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)

	// mark_step_complete
	result, err = markStepCompleteHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "mark_step_complete",
			Arguments: map[string]any{"step_text": "MCP步骤A", "task_title": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)

	// review_plan
	result, err = reviewPlanHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "review_plan",
			Arguments: map[string]any{"task_title": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)

	// add_issue
	result, err = addIssueHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "add_issue",
			Arguments: map[string]any{"issue_description": "MCP问题", "step_text": "MCP步骤A", "task_title": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)

	// resolve_issue
	result, err = resolveIssueHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "resolve_issue",
			Arguments: map[string]any{"step_text": "MCP步骤A", "resolution_text": "MCP已解决", "task_title": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)

	// update_planning_notes
	result, err = updatePlanningNotesHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "update_planning_notes",
			Arguments: map[string]any{"notes": "MCP新说明", "task_title": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)

	// check_task_completion
	result, err = checkTaskCompletionHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "check_task_completion",
			Arguments: map[string]any{"task_title": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)

	// set_priority
	result, err = setPriorityHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "set_priority",
			Arguments: map[string]any{"priority": "high", "step_text": "MCP步骤A", "task_title": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)

	// delete_step
	result, err = deleteStepHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "delete_step",
			Arguments: map[string]any{"step_text": "MCP步骤A", "task_title": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)

	// delete_task
	result, err = deleteTaskHandler(ctx, mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name:      "delete_task",
			Arguments: map[string]any{"task_title": "MCP任务1"},
		},
	})
	assert.Nil(t, err)
	assert.False(t, result.IsError)
}
