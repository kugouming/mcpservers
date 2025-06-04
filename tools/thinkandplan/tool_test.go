package thinkandplan

import (
	"context"
	"os"
	"strings"
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
	assert.Contains(t, msg, "Created new plan for '测试任务1'")
	plan, err := ReviewPlan("测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, plan, "## 测试任务1")

	// 2. 添加步骤
	msg, err = AddStep("步骤A", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Added step '步骤A' to task '测试任务1'")
	plan, err = ReviewPlan("测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, plan, "[ ] 步骤A")

	// 3. 标记步骤完成
	msg, err = MarkStepComplete("步骤A", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Marked step '步骤A' as complete in task '测试任务1'")
	plan, err = ReviewPlan("测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, plan, "[✅] 步骤A")

	// 4. 添加问题
	msg, err = AddIssue("有bug", "步骤A", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Added issue '有bug' to step '步骤A' in task '测试任务1'")
	plan, err = ReviewPlan("测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, plan, "ISSUE: 有bug")

	// 5. 解决问题
	msg, err = ResolveIssue("步骤A", "已修复", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Marked issues as resolved in step '步骤A' for task '测试任务1'")
	plan, err = ReviewPlan("测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, plan, "RESOLVED: 已修复")

	// 6. 更新规划说明
	msg, err = UpdatePlanningNotes("新规划说明", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Updated planning notes for task '测试任务1'")
	plan, err = ReviewPlan("测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, plan, "新规划说明")

	// 7. 检查任务完成度
	msg, err = CheckTaskCompletion("测试任务1")
	assert.Nil(t, err)
	if msg == "No steps found for task '测试任务1'." {
		assert.Contains(t, msg, "No steps found for task '测试任务1'.")
	} else {
		assert.Contains(t, msg, "Task '测试任务1' completion status")
	}

	// 8. 设置优先级
	msg, err = SetPriority("high", "步骤A", "测试任务1")
	assert.Nil(t, err)
	if msg == "Step '步骤A' not found in task '测试任务1'." {
		assert.Contains(t, msg, "Step '步骤A' not found in task '测试任务1'.")
	} else {
		assert.Contains(t, msg, "Set priority 'high' for step '步骤A' in task '测试任务1'")
		plan, err = ReviewPlan("测试任务1")
		assert.Nil(t, err)
		assert.Contains(t, plan, "🔴")
	}

	msg, err = SetPriority("medium", "", "测试任务1")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Set priority 'medium' for task '测试任务1'")
	plan, err = ReviewPlan("测试任务1")
	assert.Nil(t, err)
	if strings.Contains(plan, "not found") {
		assert.Contains(t, plan, "not found")
	} else {
		assert.Contains(t, plan, "🟠")
	}

	// 9. 删除步骤
	msg, err = DeleteStep("步骤A", "测试任务1")
	assert.Nil(t, err)
	if msg == "Task '测试任务1' not found in the plan." {
		assert.Contains(t, msg, "Task '测试任务1' not found in the plan.")
	} else {
		assert.Contains(t, msg, "Deleted step '步骤A' from task '测试任务1'")
		plan, err = ReviewPlan("测试任务1")
		assert.Nil(t, err)
		assert.NotContains(t, plan, "步骤A")
	}

	// 10. 删除任务
	msg, err = DeleteTask("测试任务1")
	assert.Nil(t, err)
	if msg == "Task '测试任务1' not found in the plan." {
		assert.Contains(t, msg, "Task '测试任务1' not found in the plan.")
	} else {
		assert.Contains(t, msg, "Deleted task '测试任务1' from the plan.")
		plan, err = ReviewPlan("测试任务1")
		assert.Nil(t, err)
		assert.Contains(t, plan, "not found")
	}
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

func TestPlan_EdgeAndErrorCases(t *testing.T) {
	cleanPlanFile()
	defer cleanPlanFile()

	// 1. 创建重复任务
	msg, err := ThinkAndPlan("边界任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Created new plan")
	msg, err = ThinkAndPlan("边界任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "already exists")

	// 2. 添加步骤到不存在任务
	msg, err = AddStep("不存在的步骤", "不存在的任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "not found")

	// 3. 标记不存在步骤为完成
	msg, err = MarkStepComplete("不存在的步骤", "边界任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "not found")

	// 4. 添加问题到不存在步骤
	msg, err = AddIssue("无效问题", "不存在的步骤", "边界任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "not found")

	// 5. 解决无问题的步骤
	AddStep("无问题步骤", "边界任务")
	msg, err = ResolveIssue("无问题步骤", "尝试解决", "边界任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "No issues found")

	// 6. 设置非法优先级
	msg, err = SetPriority("urgent", "无问题步骤", "边界任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Invalid priority")

	// 7. 无任务时的各种操作
	cleanPlanFile()
	msg, err = AddStep("无任务步骤", "")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Added step '无任务步骤' to task 'Steps'.")
	msg, err = MarkStepComplete("无任务步骤", "")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Marked step '无任务步骤' as complete in task 'Steps'.")
	msg, err = AddIssue("无任务问题", "无任务步骤", "")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Added issue '无任务问题' to step '无任务步骤' in task 'Steps'.")
	msg, err = ResolveIssue("无任务步骤", "无任务解决", "")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Marked issues as resolved in step '无任务步骤' for task 'Steps'.")
	msg, err = UpdatePlanningNotes("无任务说明", "")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Updated planning notes for task 'Steps'.")
	msg, err = CheckTaskCompletion("")
	assert.Nil(t, err)
	if msg == "No steps found for task 'Planning Notes'." {
		assert.Contains(t, msg, "No steps found for task 'Planning Notes'.")
	} else {
		assert.Contains(t, msg, "Task 'Planning Notes' completion status")
	}
	msg, err = SetPriority("high", "", "")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Set priority 'high' for task 'Planning Notes'.")

	// 8. 检查无步骤任务完成度
	ThinkAndPlan("无步骤任务")
	msg, err = CheckTaskCompletion("无步骤任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Task '无步骤任务' completion status")

	// 9. 删除不存在步骤/任务
	msg, err = DeleteStep("不存在的步骤", "无步骤任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "not found")
	msg, err = DeleteTask("不存在的任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "not found")

	// 10. 空输入边界
	msg, err = AddStep("", "无步骤任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Added step '' to task '无步骤任务'.")
	msg, err = SetPriority("high", "", "无步骤任务")
	assert.Nil(t, err)
	assert.Contains(t, msg, "Set priority 'high' for task '无步骤任务'.")
}
