# README

## 技术栈
- 编程语言：Golang 1.23.x
- MCP SDK：github.com/mark3labs/mcp-go


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
    ├── ......              # 存放其他工具的实现
    └── thinkplan           # 存放ThinkPlan工具的实现
```

## 代码使用

- 工具函数：在 `helper/` 目录中，定义了一些辅助函数和工具，在生成代码时需要优先使用已定义的工具函数。
- 工具实现：在 `tools/` 目录中，定义了各种工具的实现。每种工具都是一个独立的文件，文件名即为工具名。
- 第三方库：在 `go.mod` 文件中，定义了项目所使用的第三方库。当文件中不存在时，可以从互联网中获取比较欢迎的库，同时需要考虑 `go.mod` 中的版本依赖。下面为一些常用的库：
* 类型转换库：github.com/spf13/cast
* 配置加载库：github.com/spf13/viper


## 项目需求

请结合上面给出的工具定义，使用 `github.com/mark3labs/mcp-go` Lib，编写一个客户端和服务端程序。

## 编译和运行

请参考当前项目中已实现的功能和流程来编写程序，确保在当前架构中不做任何调整就可以正常运行。



