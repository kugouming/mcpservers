package thinkandplan

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const PlanFile = "output/plan.md"

var planFileMutex sync.Mutex

// 确保 plan.md 文件存在
func ensurePlanFile() error {
	planDir := "output"
	if _, err := os.Stat(planDir); os.IsNotExist(err) {
		err = os.MkdirAll(planDir, 0755)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(PlanFile); os.IsNotExist(err) {
		f, err := os.Create(PlanFile)
		if err != nil {
			return err
		}
		defer f.Close()
		f.WriteString("# Task Plan\n\nCreated on: " + time.Now().Format("2006-01-02 15:04:05") + "\n\n## Steps\n\n")
	}
	return nil
}

// 读取 plan.md 内容
func readPlanFile() (string, error) {
	if err := ensurePlanFile(); err != nil {
		return "", err
	}
	b, err := os.ReadFile(PlanFile)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// 写入 plan.md 内容
func writePlanFile(content string) error {
	if err := ensurePlanFile(); err != nil {
		return err
	}
	return os.WriteFile(PlanFile, []byte(content), 0644)
}

// 创建新任务计划
func ThinkAndPlan(taskDesc string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	if strings.Contains(content, "## "+taskDesc+"\n") {
		return fmt.Sprintf("A plan for '%s' already exists in '%s'.", taskDesc, PlanFile), nil
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	newPlan := fmt.Sprintf("\n## %s\n\nCreated: %s\n\n### Planning Notes\n\nThis is a preliminary analysis of the task.\n\n### Steps\n\n[ ] Initialize planning\n[ ] Analyze requirements\n[ ] Design solution\n", taskDesc, timestamp)
	content += newPlan
	if err := writePlanFile(content); err != nil {
		return "", err
	}
	return fmt.Sprintf("Created new plan for '%s' in '%s'. Review and customize the steps as needed.", taskDesc, PlanFile), nil
}

// 向任务添加步骤
func AddStep(stepDesc, taskTitle string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	if taskTitle == "" {
		taskTitle = getLastTaskTitle(content)
		if taskTitle == "" {
			return "No tasks found in the plan. Create a task first using think_and_plan.", nil
		}
	}
	taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `\n([\s\S]+)`)
	taskMatch := taskPattern.FindString(content)
	if taskMatch == "" {
		return fmt.Sprintf("Task '%s' not found in the plan.", taskTitle), nil
	}
	stepsPattern := regexp.MustCompile(`### Steps\n\n([\s\S]*)`)
	stepsMatch := stepsPattern.FindString(taskMatch)
	if stepsMatch == "" {
		// 没有 Steps 段，直接加
		updatedTask := taskMatch + "\n### Steps\n\n[ ] " + stepDesc + "\n"
		content = strings.Replace(content, taskMatch, updatedTask, 1)
	} else {
		updatedSteps := stepsMatch + "[ ] " + stepDesc + "\n"
		updatedTask := strings.Replace(taskMatch, stepsMatch, updatedSteps, 1)
		content = strings.Replace(content, taskMatch, updatedTask, 1)
	}
	if err := writePlanFile(content); err != nil {
		return "", err
	}
	return fmt.Sprintf("Added step '%s' to task '%s'.", stepDesc, taskTitle), nil
}

// 标记步骤完成
func MarkStepComplete(stepText, taskTitle string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	if taskTitle == "" {
		taskTitle = getLastTaskTitle(content)
		if taskTitle == "" {
			return "No tasks found in the plan.", nil
		}
	}
	taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `\n([\s\S]+)`)
	taskMatch := taskPattern.FindString(content)
	if taskMatch == "" {
		return fmt.Sprintf("Task '%s' not found in the plan.", taskTitle), nil
	}
	stepPattern := regexp.MustCompile(`(\[.?\] ` + regexp.QuoteMeta(stepText) + `)([\s\S]*)`)
	if !stepPattern.MatchString(taskMatch) {
		return fmt.Sprintf("Step '%s' not found in task '%s'.", stepText, taskTitle), nil
	}
	updatedTask := stepPattern.ReplaceAllString(taskMatch, "[✅] "+stepText)
	content = strings.Replace(content, taskMatch, updatedTask, 1)
	if err := writePlanFile(content); err != nil {
		return "", err
	}
	return fmt.Sprintf("Marked step '%s' as complete in task '%s'.", stepText, taskTitle), nil
}

// 查看计划
func ReviewPlan(taskTitle string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	if taskTitle == "" {
		return content, nil
	}
	taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `\n([\s\S]+)`)
	taskMatch := taskPattern.FindString(content)
	if taskMatch == "" {
		return fmt.Sprintf("Task '%s' not found in the plan.", taskTitle), nil
	}
	return fmt.Sprintf("# Review of task: '%s'\n\n%s", taskTitle, taskMatch), nil
}

// 添加问题
func AddIssue(issueDesc, stepText, taskTitle string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	if taskTitle == "" {
		taskTitle = getLastTaskTitle(content)
		if taskTitle == "" {
			return "No tasks found in the plan.", nil
		}
	}
	taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `\n([\s\S]+)`)
	taskMatch := taskPattern.FindString(content)
	if taskMatch == "" {
		return fmt.Sprintf("Task '%s' not found in the plan.", taskTitle), nil
	}
	stepPattern := regexp.MustCompile(`(\[.?\] ` + regexp.QuoteMeta(stepText) + `)([\s\S]*)`)
	stepMatch := stepPattern.FindString(taskMatch)
	if stepMatch == "" {
		return fmt.Sprintf("Step '%s' not found in task '%s'.", stepText, taskTitle), nil
	}
	issueNote := "\n    - ⚠️ ISSUE: " + issueDesc
	updatedStep := stepMatch + issueNote
	updatedTask := strings.Replace(taskMatch, stepMatch, updatedStep, 1)
	content = strings.Replace(content, taskMatch, updatedTask, 1)
	if err := writePlanFile(content); err != nil {
		return "", err
	}
	return fmt.Sprintf("Added issue '%s' to step '%s' in task '%s'.", issueDesc, stepText, taskTitle), nil
}

// 解决问题
func ResolveIssue(stepText, resolution, taskTitle string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	if taskTitle == "" {
		taskTitle = getLastTaskTitle(content)
		if taskTitle == "" {
			return "No tasks found in the plan.", nil
		}
	}
	taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `\n([\s\S]+)`)
	taskMatch := taskPattern.FindString(content)
	if taskMatch == "" {
		return fmt.Sprintf("Task '%s' not found in the plan.", taskTitle), nil
	}
	stepPattern := regexp.MustCompile(`(\[.?\] ` + regexp.QuoteMeta(stepText) + `)([\s\S]*)`)
	stepMatch := stepPattern.FindString(taskMatch)
	if stepMatch == "" {
		return fmt.Sprintf("Step '%s' not found in task '%s'.", stepText, taskTitle), nil
	}
	if !strings.Contains(stepMatch, "⚠️ ISSUE:") {
		return fmt.Sprintf("No issues found for step '%s' in task '%s'.", stepText, taskTitle), nil
	}
	issuePattern := regexp.MustCompile(`(    - ⚠️ ISSUE: [^\n]+)`)
	resolvedStep := issuePattern.ReplaceAllString(stepMatch, "$1 (✓ RESOLVED: "+resolution+")")
	updatedTask := strings.Replace(taskMatch, stepMatch, resolvedStep, 1)
	content = strings.Replace(content, taskMatch, updatedTask, 1)
	if err := writePlanFile(content); err != nil {
		return "", err
	}
	return fmt.Sprintf("Marked issues as resolved in step '%s' for task '%s'.", stepText, taskTitle), nil
}

// 更新规划说明
func UpdatePlanningNotes(notes, taskTitle string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	if taskTitle == "" {
		taskTitle = getLastTaskTitle(content)
		if taskTitle == "" {
			return "No tasks found in the plan.", nil
		}
	}
	taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `\n([\s\S]+)`)
	taskMatch := taskPattern.FindString(content)
	if taskMatch == "" {
		return fmt.Sprintf("Task '%s' not found in the plan.", taskTitle), nil
	}
	notesPattern := regexp.MustCompile(`### Planning Notes\n\n([\s\S]*)`)
	notesMatch := notesPattern.FindString(taskMatch)
	if notesMatch == "" {
		updatedTask := taskMatch + "\n### Planning Notes\n\n" + notes + "\n"
		content = strings.Replace(content, taskMatch, updatedTask, 1)
	} else {
		updatedNotes := "### Planning Notes\n\n" + notes + "\n"
		updatedTask := strings.Replace(taskMatch, notesMatch, updatedNotes, 1)
		content = strings.Replace(content, taskMatch, updatedTask, 1)
	}
	if err := writePlanFile(content); err != nil {
		return "", err
	}
	return fmt.Sprintf("Updated planning notes for task '%s'.", taskTitle), nil
}

// 检查任务完成度
func CheckTaskCompletion(taskTitle string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	if taskTitle == "" {
		taskTitle = getLastTaskTitle(content)
		if taskTitle == "" {
			return "No tasks found in the plan.", nil
		}
	}
	taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `\n([\s\S]+)`)
	taskMatch := taskPattern.FindString(content)
	if taskMatch == "" {
		return fmt.Sprintf("Task '%s' not found in the plan.", taskTitle), nil
	}
	incompleteSteps := regexp.MustCompile(`\[ \] (.+)`).FindAllStringSubmatch(taskMatch, -1)
	completedSteps := regexp.MustCompile(`\[✅\] (.+)`).FindAllStringSubmatch(taskMatch, -1)
	total := len(incompleteSteps) + len(completedSteps)
	if total == 0 {
		return fmt.Sprintf("No steps found for task '%s'.", taskTitle), nil
	}
	percent := float64(len(completedSteps)) / float64(total) * 100
	result := fmt.Sprintf("Task '%s' completion status: \n- %d of %d steps completed (%.1f%%)\n", taskTitle, len(completedSteps), total, percent)
	if len(incompleteSteps) > 0 {
		result += "\nRemaining steps:\n"
		for _, s := range incompleteSteps {
			result += "- " + s[1] + "\n"
		}
	}
	if len(completedSteps) == total {
		result += "\n🎉 All steps completed!"
	}
	return result, nil
}

// 删除步骤
func DeleteStep(stepText, taskTitle string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	if taskTitle == "" {
		taskTitle = getLastTaskTitle(content)
		if taskTitle == "" {
			return "No tasks found in the plan.", nil
		}
	}
	taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `\n([\s\S]+)`)
	taskMatch := taskPattern.FindString(content)
	if taskMatch == "" {
		return fmt.Sprintf("Task '%s' not found in the plan.", taskTitle), nil
	}
	stepPattern := regexp.MustCompile(`(\[.?\] ` + regexp.QuoteMeta(stepText) + `.*)`)
	stepMatch := stepPattern.FindString(taskMatch)
	if stepMatch == "" {
		return fmt.Sprintf("Step '%s' not found in task '%s'.", stepText, taskTitle), nil
	}
	updatedTask := strings.Replace(taskMatch, stepMatch+"\n", "", 1)
	if updatedTask == taskMatch {
		updatedTask = strings.Replace(taskMatch, stepMatch, "", 1)
	}
	content = strings.Replace(content, taskMatch, updatedTask, 1)
	if err := writePlanFile(content); err != nil {
		return "", err
	}
	return fmt.Sprintf("Deleted step '%s' from task '%s'.", stepText, taskTitle), nil
}

// 删除任务
func DeleteTask(taskTitle string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `\n([\s\S]+)`)
	taskMatch := taskPattern.FindString(content)
	if taskMatch == "" {
		return fmt.Sprintf("Task '%s' not found in the plan.", taskTitle), nil
	}
	updatedContent := strings.Replace(content, taskMatch, "", 1)
	updatedContent = regexp.MustCompile(`\n{3,}`).ReplaceAllString(updatedContent, "\n\n")
	if err := writePlanFile(updatedContent); err != nil {
		return "", err
	}
	return fmt.Sprintf("Deleted task '%s' from the plan.", taskTitle), nil
}

// 设置优先级
func SetPriority(priority, stepText, taskTitle string) (string, error) {
	planFileMutex.Lock()
	defer planFileMutex.Unlock()
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	priority = strings.ToLower(priority)
	valid := map[string]string{"high": "🔴", "medium": "🟠", "low": "🟢"}
	if _, ok := valid[priority]; !ok {
		return fmt.Sprintf("Invalid priority '%s'. Please use one of: high,medium,low.", priority), nil
	}
	if taskTitle == "" {
		taskTitle = getLastTaskTitle(content)
		if taskTitle == "" {
			return "No tasks found in the plan.", nil
		}
	}
	taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `(?: \[🔴🟠🟢]\])?\n([\s\S]+)`)
	taskMatch := taskPattern.FindString(content)
	if taskMatch == "" {
		return fmt.Sprintf("Task '%s' not found in the plan.", taskTitle), nil
	}
	if stepText != "" {
		stepPattern := regexp.MustCompile(`(\[.?\] ` + regexp.QuoteMeta(stepText) + `)([\s\S]*)`)
		stepMatch := stepPattern.FindString(taskMatch)
		if stepMatch == "" {
			return fmt.Sprintf("Step '%s' not found in task '%s'.", stepText, taskTitle), nil
		}
		cleanStep := regexp.MustCompile(` [🔴🟠🟢] `).ReplaceAllString(stepMatch, " ")
		prioritizedStep := strings.Replace(cleanStep, "] ", "] "+valid[priority]+" ", 1)
		updatedTask := strings.Replace(taskMatch, stepMatch, prioritizedStep, 1)
		content = strings.Replace(content, taskMatch, updatedTask, 1)
		if err := writePlanFile(content); err != nil {
			return "", err
		}
		return fmt.Sprintf("Set priority '%s' for step '%s' in task '%s'.", priority, stepText, taskTitle), nil
	} else {
		taskHeading := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `( \[🔴🟠🟢]\])?\n`)
		taskHeadingMatch := taskHeading.FindString(taskMatch)
		if taskHeadingMatch != "" {
			cleanHeading := regexp.MustCompile(` \[🔴🟠🟢\]`).ReplaceAllString(taskHeadingMatch, "")
			prioritizedHeading := strings.Replace(cleanHeading, "\n", " ["+valid[priority]+"]\n", 1)
			updatedTask := strings.Replace(taskMatch, taskHeadingMatch, prioritizedHeading, 1)
			content = strings.Replace(content, taskMatch, updatedTask, 1)
			if err := writePlanFile(content); err != nil {
				return "", err
			}
		}
		return fmt.Sprintf("Set priority '%s' for task '%s'.", priority, taskTitle), nil
	}
}

// 获取最后一个任务标题
func getLastTaskTitle(content string) string {
	taskSections := regexp.MustCompile(`## (.+?)\n`).FindAllStringSubmatch(content, -1)
	if len(taskSections) == 0 {
		return ""
	}
	return taskSections[len(taskSections)-1][1]
}

// Plan 资源获取
func GetPlanResource(taskTitle string) (string, error) {
	if err := ensurePlanFile(); err != nil {
		return "", err
	}
	content, err := readPlanFile()
	if err != nil {
		return "", err
	}
	if taskTitle != "" && taskTitle != "all" {
		taskPattern := regexp.MustCompile(`## ` + regexp.QuoteMeta(taskTitle) + `\n([\s\S]+)`)
		taskMatch := taskPattern.FindString(content)
		if taskMatch == "" {
			return "Task '" + taskTitle + "' not found in the plan.", nil
		}
		return "# Task: " + taskTitle + "\n\n" + taskMatch, nil
	}
	return content, nil
}
