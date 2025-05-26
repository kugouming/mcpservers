# README

> 实现参考：[如何让 Agent 规划调用工具](https://mp.weixin.qq.com/s/7XvAcTst9OU_4orDr-dD-w)

## 技术栈
- 编程语言：Golang
- Lib库：github.com/mark3labs/mcp-go


## 代码结构

```
├── cmd         # 存放主程序入口文件
│   ├── client              # 存放客户端程序
│   └── server              # 存放服务端程序
├── go.mod
├── go.sum
├── helper      # 存放一些辅助函数和工具
│   └── tool.go
├── Makefile    # 存放编译和运行脚本
├── output      # 存放编译后的二进制文件
│   ├── client              # 存放客户端二进制文件
│   └── sse                 # 存放SSE二进制文件
├── README.md
└── tools       # 存放各种工具的实现
    ├── elasticsearch       # 存放Elasticsearch工具的实现
    ├── httprequest         # 存放HTTP请求工具的实现
    ├── sequentialthinking  # 存放顺序思考工具的实现
    ├── think               # 存放Think工具的实现
    └── thinkplan           # 存放ThinkPlan工具的实现
```

## 工具定义
```json
{
	"name": "思考和规划",
	"id": "think_and_plan",
	"description": "这是用于系统化思考与规划的工具，支持用户在面对复杂问题或任务时，分阶段梳理思考、规划和行动步骤。工具强调思考（thought）、计划（plan）与实际行动（action）的结合，通过编号（thoughtNumber）追踪过程。该工具不会获取新信息或更改数据库，只会将想法附加到记忆中。当需要复杂推理或某种缓存记忆时，可以使用它。",
	"input_schema": {
		"type": "object",
		"properties": {
			"thought": {
				"type": "string",
				"description": "当前的思考内容，可以是对问题的分析、假设、洞见、反思或对前一步骤的总结。强调深度思考和逻辑推演，是每一步的核心。"
			},
			"plan": {
				"type": "string",
				"description": "针对当前任务拟定的计划或方案，将复杂问题分解为多个可执行步骤。"
			},
			"action": {
				"type": "string",
				"description": "基于当前思考和计划，建议下一步采取的行动步骤，要求具体、可执行、可验证，可以是下一步需要调用的一个或多个工具。"
			},
			"thoughtNumber": {
				"type": "string",
				"description": "当前思考步骤的编号，用于追踪和回溯整个思考与规划过程，便于后续复盘与优化。"
			}
		},
		"required": ["thought", "plan", "action", "thoughtNumber"]
	}
}
```

## 项目需求

请结合上面给出的工具定义，使用 `github.com/mark3labs/mcp-go` Lib，编写一个客户端和服务端程序。可以结合互联网中对“思考和规划”的描述，完善程序，但是务必保证准确的实现工具描述中的功能，同时工具编写完成后，需要给出测试用例和使用建议，确保工具的价值可以在使用过程中的收益最大化。

## 编译和运行
请参考当前项目中已实现的功能和流程来编写程序，确保在当前架构中不做任何调整就可以正常运行。