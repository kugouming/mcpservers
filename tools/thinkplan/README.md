# ThinkPlan 工具 - 系统化思考与规划

## 概述

ThinkPlan 是一个用于系统化思考与规划的 MCP (Model Context Protocol) 工具。它帮助用户在面对复杂问题或任务时，分阶段梳理思考、规划和行动步骤。工具强调思考（thought）、计划（plan）与实际行动（action）的结合，通过编号（thoughtNumber）追踪整个过程。

## 核心特性

- **结构化思考**: 将复杂问题分解为思考、规划、行动三个维度
- **过程追踪**: 通过编号系统追踪思考过程的每个步骤
- **内存管理**: 自动保存和管理所有思考记录
- **并发安全**: 支持多线程并发访问
- **多种接口**: 提供标准输入输出和 HTTP/SSE 两种接口
- **丰富的API**: 提供查询、导出、清理等多种操作

## 工具定义

```json
{
  "name": "think_and_plan",
  "description": "系统化思考与规划工具",
  "parameters": {
    "thought": {
      "type": "string",
      "required": true,
      "description": "当前的思考内容，可以是对问题的分析、假设、洞见、反思或对前一步骤的总结"
    },
    "plan": {
      "type": "string", 
      "required": true,
      "description": "针对当前任务拟定的计划或方案，将复杂问题分解为多个可执行步骤"
    },
    "action": {
      "type": "string",
      "required": true, 
      "description": "基于当前思考和计划，建议下一步采取的行动步骤，要求具体、可执行、可验证"
    },
    "thoughtNumber": {
      "type": "string",
      "required": true,
      "description": "当前思考步骤的编号，用于追踪和回溯整个思考与规划过程"
    }
  }
}
```

## 部署方式

### 1. 标准输入输出模式 (Stdio)

```bash
# 编译
make thinkplan

# 运行
./output/client/thinkplan
```

### 2. HTTP/SSE 服务模式

```bash
# 编译
make sse_thinkplan

# 运行服务器
./output/sse/thinkplan

# 服务将在 http://localhost:8084 启动
```

## 使用示例

### 基础使用

```json
{
  "tool": "think_and_plan",
  "parameters": {
    "thought": "当前项目进度落后，需要分析原因并制定追赶计划",
    "plan": "1. 分析延期原因 2. 重新评估任务优先级 3. 调整资源分配 4. 制定追赶时间表",
    "action": "首先召集团队会议，收集各模块的具体进度和遇到的问题",
    "thoughtNumber": "T001"
  }
}
```

### 连续思考过程

```json
// 第一步：问题分析
{
  "thought": "用户反馈系统响应慢，需要进行性能优化",
  "plan": "1. 性能监控和分析 2. 识别瓶颈 3. 优化方案设计 4. 实施和验证",
  "action": "部署性能监控工具，收集系统运行数据",
  "thoughtNumber": "PERF001"
}

// 第二步：深入分析
{
  "thought": "监控数据显示数据库查询是主要瓶颈，特别是用户查询接口",
  "plan": "1. 分析慢查询日志 2. 优化SQL语句 3. 添加索引 4. 考虑缓存策略",
  "action": "开启数据库慢查询日志，分析最耗时的查询语句",
  "thoughtNumber": "PERF002"
}

// 第三步：解决方案
{
  "thought": "发现用户表的复合查询缺少索引，且存在N+1查询问题",
  "plan": "1. 添加复合索引 2. 重构查询逻辑 3. 实施查询缓存 4. 性能测试验证",
  "action": "创建数据库迁移脚本，添加必要的索引",
  "thoughtNumber": "PERF003"
}
```

## HTTP API 接口

当使用 SSE 模式时，除了标准的 MCP 接口外，还提供以下 HTTP API：

### 获取所有思考记录
```bash
GET /api/thinkplan/memory
```

### 获取思考过程摘要
```bash
GET /api/thinkplan/summary
```

### 根据编号获取特定记录
```bash
GET /api/thinkplan/memory/{thoughtNumber}
```

### 导出所有记录为JSON
```bash
GET /api/thinkplan/export
```

### 清空所有记录
```bash
DELETE /api/thinkplan/memory
```

## 测试用例

### 单元测试

```bash
cd tools/thinkplan
go test -v
```

### 基准测试

```bash
cd tools/thinkplan  
go test -bench=. -benchmem
```

### 集成测试

```bash
# 启动服务器
./output/sse/thinkplan &

# 测试基本功能
curl -X POST http://localhost:8084/api/thinkplan/memory \
  -H "Content-Type: application/json" \
  -d '{
    "thought": "测试思考内容",
    "plan": "测试计划",
    "action": "测试行动",
    "thoughtNumber": "TEST001"
  }'

# 查看结果
curl http://localhost:8084/api/thinkplan/memory
```

## 最佳实践

### 1. 编号规范

建议使用有意义的编号前缀：
- `PROJ001`, `PROJ002` - 项目相关思考
- `BUG001`, `BUG002` - 问题解决过程
- `ARCH001`, `ARCH002` - 架构设计思考
- `PERF001`, `PERF002` - 性能优化过程

### 2. 思考结构

**思考（Thought）**：
- 描述当前的理解和分析
- 包含关键洞察和假设
- 记录重要的发现和疑问

**计划（Plan）**：
- 分解为具体的步骤
- 按优先级排序
- 包含时间估算（可选）

**行动（Action）**：
- 具体可执行的下一步
- 明确的成功标准
- 所需的资源和工具

### 3. 使用场景

**问题解决**：
```json
{
  "thought": "用户报告登录失败，初步怀疑是认证服务问题",
  "plan": "1. 检查认证服务状态 2. 查看错误日志 3. 验证数据库连接 4. 测试认证流程",
  "action": "登录服务器检查认证服务的运行状态和资源使用情况",
  "thoughtNumber": "LOGIN_BUG001"
}
```

**项目规划**：
```json
{
  "thought": "新功能需求复杂，涉及多个系统集成，需要仔细规划",
  "plan": "1. 需求分析和拆解 2. 技术方案设计 3. 接口定义 4. 开发排期 5. 测试策略",
  "action": "组织需求评审会议，邀请产品、开发、测试团队参与",
  "thoughtNumber": "FEATURE_PLAN001"
}
```

**学习研究**：
```json
{
  "thought": "微服务架构在我们的场景下可能带来更好的可扩展性",
  "plan": "1. 研究微服务最佳实践 2. 分析当前系统架构 3. 设计迁移方案 4. 风险评估",
  "action": "阅读《微服务架构设计模式》并整理关键要点",
  "thoughtNumber": "MICROSERVICE_RESEARCH001"
}
```

### 4. 团队协作

- 使用统一的编号规范
- 定期导出思考记录进行团队分享
- 在代码提交信息中引用相关的思考编号
- 建立思考记录的回顾机制

## 性能指标

基于基准测试结果：
- **处理速度**: ~97,912 ns/op (约0.1毫秒每次操作)
- **内存使用**: 923 B/op (每次操作约1KB内存)
- **内存分配**: 14 allocs/op (每次操作14次内存分配)
- **并发安全**: 支持多线程并发访问

## 故障排除

### 常见问题

1. **参数缺失错误**
   ```
   错误: "thought parameter is required"
   解决: 确保所有必需参数都已提供
   ```

2. **服务启动失败**
   ```
   错误: "Failed to start server"
   解决: 检查端口8084是否被占用
   ```

3. **内存使用过高**
   ```
   解决: 定期调用清理API或重启服务
   ```

### 日志查看

服务器日志包含详细的操作记录：
```bash
# 查看服务器日志
./output/sse/thinkplan 2>&1 | tee thinkplan.log
```

## 扩展开发

### 添加新功能

1. 在 `tools/thinkplan/tool.go` 中添加新的函数
2. 在 `cmd/server/sse/thinkplan/main.go` 中添加新的API端点
3. 编写相应的测试用例
4. 更新文档

### 自定义存储

默认使用内存存储，可以扩展为：
- 文件存储
- 数据库存储
- 分布式存储

## 许可证

本项目遵循项目根目录的许可证条款。 