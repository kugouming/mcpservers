package thinkandplan

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cast"
)

// RegisterTool 注册 thinkandplan 工具到 MCP 服务器
func RegisterTool(s *server.MCPServer) {
	s.AddTool(thinkAndPlanTool, thinkAndPlanHandler)
	s.AddTool(addStepTool, addStepHandler)
	s.AddTool(markStepCompleteTool, markStepCompleteHandler)
	s.AddTool(reviewPlanTool, reviewPlanHandler)
	s.AddTool(addIssueTool, addIssueHandler)
	s.AddTool(resolveIssueTool, resolveIssueHandler)
	s.AddTool(updatePlanningNotesTool, updatePlanningNotesHandler)
	s.AddTool(checkTaskCompletionTool, checkTaskCompletionHandler)
	s.AddTool(deleteStepTool, deleteStepHandler)
	s.AddTool(deleteTaskTool, deleteTaskHandler)
	s.AddTool(setPriorityTool, setPriorityHandler)
	s.AddResource(getPlanResource, getPlanHandler)
}

var thinkAndPlanTool = mcp.NewTool("think_and_plan",
	mcp.WithDescription("🧠 Think through a task and create a structured plan."),
	mcp.WithString("task_description", mcp.Required(), mcp.Description("A description of the task to be planned")),
)

func thinkAndPlanHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskDesc := cast.ToString(req.GetArguments()["task_description"])
	if taskDesc == "" {
		return mcp.NewToolResultError("task_description not provided"), nil
	}
	msg, err := ThinkAndPlan(taskDesc)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

var addStepTool = mcp.NewTool("add_step",
	mcp.WithDescription("➕ Add a new step to the plan."),
	mcp.WithString("step_description", mcp.Required(), mcp.Description("Description of the step to add")),
	mcp.WithString("task_title", mcp.Description("The task title to add the step to (uses the latest task if not specified)")),
)

func addStepHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stepDesc := cast.ToString(req.GetArguments()["step_description"])
	taskTitle := cast.ToString(req.GetArguments()["task_title"])
	if stepDesc == "" {
		return mcp.NewToolResultError("step_description not provided"), nil
	}
	msg, err := AddStep(stepDesc, taskTitle)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

var markStepCompleteTool = mcp.NewTool("mark_step_complete",
	mcp.WithDescription("✓ Mark a step as completed in the plan."),
	mcp.WithString("step_text", mcp.Required(), mcp.Description("Text of the step to mark as complete")),
	mcp.WithString("task_title", mcp.Description("The task title containing the step (uses the latest task if not specified)")),
)

func markStepCompleteHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stepText := cast.ToString(req.GetArguments()["step_text"])
	taskTitle := cast.ToString(req.GetArguments()["task_title"])
	if stepText == "" {
		return mcp.NewToolResultError("step_text not provided"), nil
	}
	msg, err := MarkStepComplete(stepText, taskTitle)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

var reviewPlanTool = mcp.NewTool("review_plan",
	mcp.WithDescription("📖 Review the current plan and return its contents."),
	mcp.WithString("task_title", mcp.Description("Specific task to review (reviews the entire plan if not specified)")),
)

func reviewPlanHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskTitle := cast.ToString(req.GetArguments()["task_title"])
	msg, err := ReviewPlan(taskTitle)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

var addIssueTool = mcp.NewTool("add_issue",
	mcp.WithDescription("⚠️ Add an issue note to a specific step in the plan."),
	mcp.WithString("issue_description", mcp.Required(), mcp.Description("Description of the issue")),
	mcp.WithString("step_text", mcp.Required(), mcp.Description("The step text to add the issue to")),
	mcp.WithString("task_title", mcp.Description("The task title (uses the latest task if not specified)")),
)

func addIssueHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	issueDesc := cast.ToString(req.GetArguments()["issue_description"])
	stepText := cast.ToString(req.GetArguments()["step_text"])
	taskTitle := cast.ToString(req.GetArguments()["task_title"])
	if issueDesc == "" || stepText == "" {
		return mcp.NewToolResultError("issue_description & step_text not provided"), nil
	}
	msg, err := AddIssue(issueDesc, stepText, taskTitle)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

var resolveIssueTool = mcp.NewTool("resolve_issue",
	mcp.WithDescription("🎯 Mark an issue as resolved for a specific step."),
	mcp.WithString("step_text", mcp.Required(), mcp.Description("The step text containing the issue")),
	mcp.WithString("resolution_text", mcp.Required(), mcp.Description("Description of how the issue was resolved")),
	mcp.WithString("task_title", mcp.Description("The task title (uses the latest task if not specified)")),
)

func resolveIssueHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stepText := cast.ToString(req.GetArguments()["step_text"])
	resolution := cast.ToString(req.GetArguments()["resolution_text"])
	taskTitle := cast.ToString(req.GetArguments()["task_title"])
	if stepText == "" || resolution == "" {
		return mcp.NewToolResultError("step_text & resolution_text not provided"), nil
	}
	msg, err := ResolveIssue(stepText, resolution, taskTitle)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

var updatePlanningNotesTool = mcp.NewTool("update_planning_notes",
	mcp.WithDescription("📝 Update the planning notes for a task."),
	mcp.WithString("notes", mcp.Required(), mcp.Description("The new planning notes")),
	mcp.WithString("task_title", mcp.Description("The task title (uses the latest task if not specified)")),
)

func updatePlanningNotesHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	notes := cast.ToString(req.GetArguments()["notes"])
	taskTitle := cast.ToString(req.GetArguments()["task_title"])
	if notes == "" {
		return mcp.NewToolResultError("notes not provided"), nil
	}
	msg, err := UpdatePlanningNotes(notes, taskTitle)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

var checkTaskCompletionTool = mcp.NewTool("check_task_completion",
	mcp.WithDescription("🔄 Check if all steps in a task are marked as complete."),
	mcp.WithString("task_title", mcp.Description("The task title (uses the latest task if not specified)")),
)

func checkTaskCompletionHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskTitle := cast.ToString(req.GetArguments()["task_title"])
	msg, err := CheckTaskCompletion(taskTitle)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

var deleteStepTool = mcp.NewTool("delete_step",
	mcp.WithDescription("🗑️ Delete a step from the plan."),
	mcp.WithString("step_text", mcp.Required(), mcp.Description("Text of the step to delete")),
	mcp.WithString("task_title", mcp.Description("The task title (uses the latest task if not specified)")),
)

func deleteStepHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stepText := cast.ToString(req.GetArguments()["step_text"])
	taskTitle := cast.ToString(req.GetArguments()["task_title"])
	if stepText == "" {
		return mcp.NewToolResultError("step_text not provided"), nil
	}
	msg, err := DeleteStep(stepText, taskTitle)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

var deleteTaskTool = mcp.NewTool("delete_task",
	mcp.WithDescription("🗑️ Delete an entire task from the plan."),
	mcp.WithString("task_title", mcp.Required(), mcp.Description("The title of the task to delete")),
)

func deleteTaskHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	taskTitle := cast.ToString(req.GetArguments()["task_title"])
	if taskTitle == "" {
		return mcp.NewToolResultError("task_title not provided"), nil
	}
	msg, err := DeleteTask(taskTitle)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

var setPriorityTool = mcp.NewTool("set_priority",
	mcp.WithDescription("🔴 Set priority for a task or step."),
	mcp.WithString("priority", mcp.Required(), mcp.Description("Priority level (high, medium, low)")),
	mcp.WithString("step_text", mcp.Description("Text of the step to prioritize (if None, sets priority for the task)")),
	mcp.WithString("task_title", mcp.Description("The task title (uses the latest task if not specified)")),
)

func setPriorityHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	priority := cast.ToString(req.GetArguments()["priority"])
	stepText := cast.ToString(req.GetArguments()["step_text"])
	taskTitle := cast.ToString(req.GetArguments()["task_title"])
	if priority == "" {
		return mcp.NewToolResultError("priority not provided"), nil
	}
	msg, err := SetPriority(priority, stepText, taskTitle)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(msg), nil
}

// 资源 handler
var getPlanResource = mcp.NewResource("plan://{task_title}", "Get the plan contents as a resource.")

func getPlanHandler(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	taskTitle := ""
	if v, ok := req.Params.Arguments["task_title"]; ok {
		taskTitle = cast.ToString(v)
	}
	plan, err := GetPlanResource(taskTitle)
	if err != nil {
		return nil, err
	}
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      "plan://" + taskTitle,
			MIMEType: "text/markdown",
			Text:     plan,
		},
	}, nil
}
