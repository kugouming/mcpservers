# README


## 技术栈
- 编程语言：Golang
- Lib库：github.com/mark3labs/mcp-go


## 代码结构

```
├── cmd         # 存放主程序入口文件
│   ├── client              # 存放客户端程序
│   ├── example             # 存放示例程序,请将生成的工具注册到此处
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
├── config      # 存放配置文件
│   └── switchhosts         # 存放switchhosts工具的配置文件
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
	"name": "本地 Hosts 管理工具",
	"id": "switchhosts",
	"description": "一个本地Hosts管理工具，用于快速切换不同环境下的网络配置。它支持将指定配置中的内容追加到系统Hosts文件中，从而实现快速切换不同环境。Hosts文件中会通过特殊的分隔符（例如：#switchhosts_start和#switchhosts_end）来标识需要追加的内容，确保不会影响到其他手动编辑的Hosts条目。",
	"input_schema": {
		"type": "object",
		"properties": {
			"conf_name": {
				"type": "string",
				"description": "配置名称，用于标识配置目录下不同的Hosts文件。",
			}
		},
		"required": ["conf_name"]
	}
}
```

## 项目需求

请结合上面给出的工具定义，使用 `github.com/mark3labs/mcp-go` Lib，编写一个客户端程序。

已知信息如下：
- 配置文件会存放在`config/switchhosts`目录下，编译后的信息需要分别放在`output/client` `output/config/switchhosts`目录下，
- `config/switchhosts` 目录下有多个配置文件，例如：`hosts_dev.txt`, `hosts_prod.txt`等
- 客户端会调用 mcp tool 工具实现切换Hosts文件的功能，切换Hosts文件时，需要传入配置名称作为参数。
- 客户端切换Hosts文件时，会将`config/switchhosts`目录下对应的配置文件内容追加到系统Hosts文件中。
- 客户端切换Hosts文件时，需要在系统Hosts文件中添加特殊的分隔符（例如：`#switchhosts_start`和`#switchhosts_end`），以确保不会影响到其他手动编辑的Hosts条目。
- 客户端切换Hosts文件时，需要确保系统Hosts文件中原有的内容不被覆盖，可以系统配置内容进行复原（未指定配置文件名称时）。
- 工具可以支持不同操作系统，例如：Windows, MacOS, Linux等。

## 编写要求
- 工具编写完成后，需要给出测试用例和使用建议，确保工具的价值可以在使用过程中的收益最大化。

## 编译和运行
请参考当前项目中已实现的功能和流程来编写程序，确保在当前架构中不做任何调整就可以正常运行。