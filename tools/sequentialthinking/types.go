package sequentialthinking

// ThoughtData 定义了思考步骤的所有属性
type ThoughtData struct {
	Thought           string `json:"thought"`                     // 当前思考的内容
	ThoughtNumber     int    `json:"thoughtNumber"`               // 当前思考的序号
	TotalThoughts     int    `json:"totalThoughts"`               // 估计的总思考步骤数
	IsRevision        bool   `json:"isRevision,omitempty"`        // 是否修改之前的思考
	RevisesThought    int    `json:"revisesThought,omitempty"`    // 修改的是哪个思考
	BranchFromThought int    `json:"branchFromThought,omitempty"` // 从哪个思考分支出来的
	BranchID          string `json:"branchId,omitempty"`          // 分支标识符
	NeedsMoreThoughts bool   `json:"needsMoreThoughts,omitempty"` // 是否需要更多思考
	NextThoughtNeeded bool   `json:"nextThoughtNeeded"`           // 是否需要下一个思考步骤
}

// ThoughtResponse 定义了思考服务器的响应结构
type ThoughtResponse struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ContentItem 定义了响应内容的单个项目
type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ThoughtStatus 定义了思考状态的JSON响应结构
type ThoughtStatus struct {
	ThoughtNumber        int      `json:"thoughtNumber"`
	TotalThoughts        int      `json:"totalThoughts"`
	NextThoughtNeeded    bool     `json:"nextThoughtNeeded"`
	Branches             []string `json:"branches"`
	ThoughtHistoryLength int      `json:"thoughtHistoryLength"`
}

// ErrorResponse 定义了错误响应的结构
type ErrorResponse struct {
	Error  string `json:"error"`
	Status string `json:"status"`
}
