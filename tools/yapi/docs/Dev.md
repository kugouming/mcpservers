# README


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
    ├── yapi                # 存放YAPI工具的实现
	├── think               # 存放Think工具的实现
    └── thinkplan           # 存放ThinkPlan工具的实现
```


## 项目需求
请结合当前目录下的 `index.ts` 文件中的工具实现，使用 `github.com/mark3labs/mcp-go` Lib，编写一个客户端程序。

## 编写要求
- 编写工具时，需要考虑到工具的通用性和可扩展性。
- 编写工具时，需要考虑到工具的健壮性和容错性。
- 需要对编写的每个工具进行详细的注释说明，确保其他开发者可以快速理解和使用这些工具。
- 工具编写完成后，需要给出测试用例和使用建议，确保工具的价值可以在使用过程中的收益最大化。


## 编译和运行
请参考当前项目中已实现的功能和流程来编写程序，确保在当前架构中不做任何调整就可以正常运行。