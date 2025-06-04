# 📝 任务规划 MCP 服务器

一个 MCP（模型上下文协议）服务器实现，为 Claude 和其他兼容 MCP 的 AI 助手提供任务规划和跟踪工具。


> 参考文档：
> 	1. 英文文档：https://raw.githubusercontent.com/may3rr/think_and_plan_MCP/refs/heads/main/README.md
> 	2. 中文文档：https://raw.githubusercontent.com/may3rr/think_and_plan_MCP/refs/heads/main/README_zh.md

## 📋 概述

这个任务规划 MCP 服务器使 AI 助手能够：

1. 创建结构化的任务计划 ✨
2. 添加并跟踪完成任务的步骤 📊
3. 在完成步骤时将其标记为完成 ✅
4. 记录并解决任务执行过程中出现的问题 🛠️
5. 查看计划的当前状态 👀

所有规划信息都存储在本地的 `plan.md` 文件中，可以由人类查看和编辑。

## ✨ 功能特点

- **结构化规划**：创建带有步骤和规划说明的有组织的任务计划 📑
- **进度跟踪**：在执行步骤时将其标记为完成 ⏱️
- **问题管理**：记录问题及其解决方案 🔍
- **审查能力**：随时查看整个计划或特定任务 👁️
- **基于文件的存储**：所有信息存储在人类可读的 Markdown 文件中 📁

## 🛠️ 提供的工具

此 MCP 服务器提供以下工具：

| 工具 | 描述 |
|------|-------------|
| `think_and_plan` | 为任务创建新的结构化计划 🧠 |
| `add_step` | 向现有任务计划添加新步骤 ➕ |
| `mark_step_complete` | 将步骤标记为已完成 ✓ |
| `review_plan` | 查看当前计划内容 📖 |
| `add_issue` | 记录特定步骤的问题 ⚠️ |
| `resolve_issue` | 标记问题已解决并提供说明 🎯 |
| `update_planning_notes` | 更新任务的规划说明 📝 |
| `check_task_completion` | 检查任务的完成状态 🔄 |
| `delete_step` | 删除用户不需要的步骤 🗑️ |
| `delete_task` | 删除用户不需要的任务 🗑️ |
| `set_priority` | 设置任务的优先级 高🔴中🟠低🟢|

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
	├── thinkplan           # 存放ThinkPlan工具的实现
    └── thinkandplan        # 存放ThinkAndPlan工具的实现
```


## 项目需求

请结合上面给出的工具定义，使用 `github.com/mark3labs/mcp-go` Lib，编写一个客户端和服务端程序。可以结合互联网中对“思考和规划”的描述，完善程序，但是务必保证准确的实现工具描述中的功能，同时工具编写完成后，需要给出测试用例和使用建议，确保工具的价值可以在使用过程中的收益最大化。

## 编译和运行
请参考当前项目中已实现的功能和流程来编写程序，确保在当前架构中不做任何调整就可以正常运行。