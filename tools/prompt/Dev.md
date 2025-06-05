# README

> 实现参考：[MCP Prompt架构优化：从Tools到原生Prompt的性能重构](https://mp.weixin.qq.com/s/43aqvs7uG1rhAlsAt4x7IQ)
> Javascript版本参考：[mcp-prompt-server](https://github.com/gdli6177/mcp-prompt-server)

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

## 代码实现参考

```js
import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import fs from 'fs-extra';
import path from 'path';
import { fileURLToPath } from 'url';
import YAML from 'yaml';
import { z } from 'zod';

// 获取当前文件的目录路径
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// 预设prompts的目录路径
const PROMPTS_DIR = path.join(__dirname, 'prompts');

// 存储所有加载的prompts
let loadedPrompts = [];

/**
 * 从prompts目录加载所有预设的prompt
 */
async function loadPrompts() {
  try {
    // 确保prompts目录存在
    await fs.ensureDir(PROMPTS_DIR);
    
    // 读取prompts目录中的所有文件
    const files = await fs.readdir(PROMPTS_DIR);
    
    // 过滤出YAML和JSON文件
    const promptFiles = files.filter(file => 
      file.endsWith('.yaml') || file.endsWith('.yml') || file.endsWith('.json')
    );
    
    // 加载每个prompt文件
    const prompts = [];
    for (const file of promptFiles) {
      const filePath = path.join(PROMPTS_DIR, file);
      const content = await fs.readFile(filePath, 'utf8');
      
      let prompt;
      if (file.endsWith('.json')) {
        prompt = JSON.parse(content);
      } else {
        // 假设其他文件是YAML格式
        prompt = YAML.parse(content);
      }
      
      // 确保prompt有name字段
      if (!prompt.name) {
        console.warn(`Warning: Prompt in ${file} is missing a name field. Skipping.`);
        continue;
      }
      
      prompts.push(prompt);
    }
    
    loadedPrompts = prompts;
    console.log(`Loaded ${prompts.length} prompts from ${PROMPTS_DIR}`);
    return prompts;
  } catch (error) {
    console.error('Error loading prompts:', error);
    return [];
  }
}

/**
 * 启动MCP服务器
 */
async function startServer() {
  // 加载所有预设的prompts
  await loadPrompts();
  
  // 创建MCP服务器
  const server = new McpServer({
    name: "mcp-prompt-server",
    version: "1.0.0"
  });
  
  // 为每个预设的prompt创建一个工具
  loadedPrompts.forEach(prompt => {
    // 构建工具的输入schema
    const schemaObj = {};
    
    if (prompt.arguments && Array.isArray(prompt.arguments)) {
      prompt.arguments.forEach(arg => {
        // 默认所有参数都是字符串类型
        schemaObj[arg.name] = z.string().describe(arg.description || `参数: ${arg.name}`);
      });
    }
    
    // 注册工具
    server.tool(
      prompt.name,
      schemaObj,
      async (args) => {
        // 处理prompt内容
        let promptText = '';
        
        if (prompt.messages && Array.isArray(prompt.messages)) {
          // 只处理用户消息
          const userMessages = prompt.messages.filter(msg => msg.role === 'user');
          
          for (const message of userMessages) {
            if (message.content && typeof message.content.text === 'string') {
              let text = message.content.text;
              
              // 替换所有 {{arg}} 格式的参数
              for (const [key, value] of Object.entries(args)) {
                text = text.replace(new RegExp(`{{${key}}}`, 'g'), value);
              }
              
              promptText += text + '\n\n';
            }
          }
        }
        
        // 返回处理后的prompt内容
        return {
          content: [
            {
              type: "text",
              text: promptText.trim()
            }
          ]
        };
      },
      {
        description: prompt.description || `Prompt: ${prompt.name}`
      }
    );
  });
  
  // 添加管理工具 - 重新加载prompts
  server.tool(
    "reload_prompts",
    {},
    async () => {
      await loadPrompts();
      return {
        content: [
          {
            type: "text",
            text: `成功重新加载了 ${loadedPrompts.length} 个prompts。`
          }
        ]
      };
    },
    {
      description: "重新加载所有预设的prompts"
    }
  );
  
  // 添加管理工具 - 获取prompt名称列表
  server.tool(
    "get_prompt_names",
    {},
    async () => {
      const promptNames = loadedPrompts.map(p => p.name);
      return {
        content: [
          {
            type: "text",
            text: `可用的prompts (${promptNames.length}):\n${promptNames.join('\n')}`
          }
        ]
      };
    },
    {
      description: "获取所有可用的prompt名称"
    }
  );
  
  // 创建stdio传输层
  const transport = new StdioServerTransport();
  
  // 连接服务器
  await server.connect(transport);
  console.log('MCP Prompt Server is running...');
}

// 启动服务器
startServer().catch(error => {
  console.error('Failed to start server:', error);
  process.exit(1);
});
```