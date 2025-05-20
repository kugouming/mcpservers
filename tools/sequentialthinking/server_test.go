package sequentialthinking

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSequentialThinkingServer_ProcessThought(t *testing.T) {
	tests := []struct {
		name     string
		input    ThoughtData
		wantErr  bool
		validate func(t *testing.T, response ThoughtResponse)
	}{
		{
			name: "正常思考步骤",
			input: ThoughtData{
				Thought:           "这是一个测试思考",
				ThoughtNumber:     1,
				TotalThoughts:     3,
				NextThoughtNeeded: true,
			},
			wantErr: false,
			validate: func(t *testing.T, response ThoughtResponse) {
				if response.IsError {
					t.Error("期望非错误响应")
				}
				var status ThoughtStatus
				if err := json.Unmarshal([]byte(response.Content[0].Text), &status); err != nil {
					t.Errorf("无法解析响应JSON: %v", err)
				}
				if status.ThoughtNumber != 1 {
					t.Errorf("期望ThoughtNumber为1，得到 %d", status.ThoughtNumber)
				}
				if status.TotalThoughts != 3 {
					t.Errorf("期望TotalThoughts为3，得到 %d", status.TotalThoughts)
				}
				if !status.NextThoughtNeeded {
					t.Error("期望NextThoughtNeeded为true")
				}
			},
		},
		{
			name: "修改之前的思考",
			input: ThoughtData{
				Thought:           "修改之前的思考",
				ThoughtNumber:     2,
				TotalThoughts:     3,
				IsRevision:        true,
				RevisesThought:    1,
				NextThoughtNeeded: true,
			},
			wantErr: false,
			validate: func(t *testing.T, response ThoughtResponse) {
				if response.IsError {
					t.Error("期望非错误响应")
				}
				var status ThoughtStatus
				if err := json.Unmarshal([]byte(response.Content[0].Text), &status); err != nil {
					t.Errorf("无法解析响应JSON: %v", err)
				}
				if status.ThoughtNumber != 2 {
					t.Errorf("期望ThoughtNumber为2，得到 %d", status.ThoughtNumber)
				}
			},
		},
		{
			name: "分支思考",
			input: ThoughtData{
				Thought:           "分支思考",
				ThoughtNumber:     3,
				TotalThoughts:     3,
				BranchFromThought: 1,
				BranchID:          "branch-1",
				NextThoughtNeeded: true,
			},
			wantErr: false,
			validate: func(t *testing.T, response ThoughtResponse) {
				if response.IsError {
					t.Error("期望非错误响应")
				}
				var status ThoughtStatus
				if err := json.Unmarshal([]byte(response.Content[0].Text), &status); err != nil {
					t.Errorf("无法解析响应JSON: %v", err)
				}
				if len(status.Branches) != 1 {
					t.Errorf("期望有1个分支，得到 %d", len(status.Branches))
				}
				if status.Branches[0] != "branch-1" {
					t.Errorf("期望分支ID为branch-1，得到 %s", status.Branches[0])
				}
			},
		},
		{
			name: "无效的思考步骤",
			input: ThoughtData{
				Thought:           "",
				ThoughtNumber:     0,
				TotalThoughts:     3,
				NextThoughtNeeded: true,
			},
			wantErr: true,
			validate: func(t *testing.T, response ThoughtResponse) {
				if !response.IsError {
					t.Error("期望错误响应")
				}
				var errResp ErrorResponse
				if err := json.Unmarshal([]byte(response.Content[0].Text), &errResp); err != nil {
					t.Errorf("无法解析错误响应JSON: %v", err)
				}
				if errResp.Status != "failed" {
					t.Errorf("期望状态为failed，得到 %s", errResp.Status)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewSequentialThinkingServer()
			response := server.ProcessThought(tt.input)

			if (response.IsError) != tt.wantErr {
				t.Errorf("ProcessThought() 错误 = %v, 期望错误 %v", response.IsError, tt.wantErr)
			}

			if tt.validate != nil {
				tt.validate(t, response)
			}
		})
	}
}

func TestSequentialThinkingServer_FormatThought(t *testing.T) {
	server := NewSequentialThinkingServer()

	tests := []struct {
		name     string
		input    ThoughtData
		contains []string
	}{
		{
			name: "普通思考格式化",
			input: ThoughtData{
				Thought:           "测试思考",
				ThoughtNumber:     1,
				TotalThoughts:     3,
				NextThoughtNeeded: true,
			},
			contains: []string{"💭 Thought", "1/3", "测试思考"},
		},
		{
			name: "修改思考格式化",
			input: ThoughtData{
				Thought:           "修改思考",
				ThoughtNumber:     2,
				TotalThoughts:     3,
				IsRevision:        true,
				RevisesThought:    1,
				NextThoughtNeeded: true,
			},
			contains: []string{"🔄 Revision", "2/3", "修改思考", "revising thought 1"},
		},
		{
			name: "分支思考格式化",
			input: ThoughtData{
				Thought:           "分支思考",
				ThoughtNumber:     3,
				TotalThoughts:     3,
				BranchFromThought: 1,
				BranchID:          "branch-1",
				NextThoughtNeeded: true,
			},
			contains: []string{"🌿 Branch", "3/3", "分支思考", "from thought 1", "ID: branch-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := server.formatThought(tt.input)
			for _, str := range tt.contains {
				if !strings.Contains(formatted, str) {
					t.Errorf("格式化输出中未找到期望的字符串 %q", str)
				}
			}
		})
	}
}
