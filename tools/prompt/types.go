package prompt

// PromptTemplate 表示一个 prompt 模板
type PromptTemplate struct {
	Name        string           `yaml:"name"        json:"name"`
	Description string           `yaml:"description" json:"description"`
	Arguments   []PromptArgument `yaml:"arguments"   json:"arguments"`
	Messages    []PromptMessage  `yaml:"messages"    json:"messages"`
}

// PromptArgument 表示模板参数
type PromptArgument struct {
	Name        string `yaml:"name"        json:"name"`
	Description string `yaml:"description" json:"description"`
	Required    bool   `yaml:"required"    json:"required"`
}

// PromptMessage 表示一条消息
type PromptMessage struct {
	Role    string        `yaml:"role"    json:"role"`
	Content PromptContent `yaml:"content" json:"content"`
}

type PromptContent struct {
	Type string `yaml:"type" json:"type"`
	Text string `yaml:"text" json:"text"`
}
