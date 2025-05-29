package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/kugouming/mcpservers/tools/yapi"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	fmt.Println("=== YAPI 工具示例程序 ===")

	// 检查并加载配置
	config, err := yapi.LoadGlobalConfig()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
		return
	}

	// 打印当前配置
	fmt.Println("\n当前配置信息:")
	cm := yapi.GetConfigManager()
	cm.PrintConfig()

	// 选择运行模式
	mode := getRunMode()
	switch mode {
	case "1":
		runDirectClientExample(config)
	case "2":
		runMCPServerExample()
	case "3":
		runInteractiveExample(config)
	case "4":
		runConfigManagementExample()
	default:
		fmt.Println("无效的选择，退出程序")
	}
}

// getRunMode 获取运行模式
func getRunMode() string {
	fmt.Println("\n请选择运行模式：")
	fmt.Println("1. 直接客户端调用示例")
	fmt.Println("2. MCP 服务器示例")
	fmt.Println("3. 交互式示例")
	fmt.Println("4. 配置管理示例")
	fmt.Print("请输入选择 (1-4): ")

	var choice string
	fmt.Scanln(&choice)
	return choice
}

// runDirectClientExample 运行直接客户端调用示例
func runDirectClientExample(config *yapi.Config) {
	fmt.Println("\n=== 直接客户端调用示例 ===")

	// 使用配置创建客户端
	client := yapi.NewYapiClientFromConfig(config)

	// 示例项目ID (请根据实际情况修改)
	projectID := 123

	fmt.Printf("正在获取项目 %d 的信息...\n", projectID)

	// 1. 获取项目信息
	projectInfo, err := client.GetProjectInfo(projectID)
	if err != nil {
		log.Printf("获取项目信息失败: %v", err)
	} else {
		fmt.Printf("项目名称: %s\n", projectInfo.Data.Name)
		fmt.Printf("项目描述: %s\n", projectInfo.Data.Description)
		fmt.Printf("基础路径: %s\n", projectInfo.Data.BasePath)
	}

	// 2. 获取接口列表
	interfaces, err := client.GetInterfaces(projectID)
	if err != nil {
		log.Printf("获取接口列表失败: %v", err)
		return
	}

	fmt.Printf("\n找到 %d 个接口:\n", interfaces.Data.Count)
	for i, iface := range interfaces.Data.List {
		if i >= 5 { // 只显示前5个接口
			fmt.Printf("... 还有 %d 个接口\n", interfaces.Data.Count-5)
			break
		}
		fmt.Printf("%d. %s - %s %s\n", i+1, iface.Title, iface.Method, iface.Path)
	}

	// 3. 获取第一个接口的详情（如果存在）
	if len(interfaces.Data.List) > 0 {
		firstInterface := interfaces.Data.List[0]
		fmt.Printf("\n正在获取接口 '%s' 的详情...\n", firstInterface.Title)

		detail, err := client.GetInterfaceDetail(firstInterface.ID)
		if err != nil {
			log.Printf("获取接口详情失败: %v", err)
		} else {
			fmt.Printf("接口描述: %s\n", detail.Data.Description)
			fmt.Printf("请求头数量: %d\n", len(detail.Data.ReqHeaders))
			fmt.Printf("查询参数数量: %d\n", len(detail.Data.ReqQuery))
			fmt.Printf("响应体类型: %s\n", detail.Data.ResBodyType)
		}
	}
}

// runMCPServerExample 运行 MCP 服务器示例
func runMCPServerExample() {
	fmt.Println("\n=== MCP 服务器示例 ===")

	// 创建 MCP 服务器
	mcpServer := server.NewMCPServer(
		"YAPI MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// 注册 YAPI 工具
	yapi.RegisterTool(mcpServer)

	fmt.Println("MCP 服务器已启动，YAPI 工具已注册")
	fmt.Println("可以通过 MCP 协议调用以下工具:")
	fmt.Println("- get_interfaces: 获取接口列表")
	fmt.Println("- get_interface_detail: 获取接口详情")
	fmt.Println("- get_project_info: 获取项目信息")

	fmt.Println("\n服务器正在运行，按 Ctrl+C 退出...")

	// 在实际应用中，这里会启动 MCP 服务器的网络监听
	select {} // 无限等待
}

// runInteractiveExample 运行交互式示例
func runInteractiveExample(config *yapi.Config) {
	fmt.Println("\n=== 交互式示例 ===")

	// 使用配置创建客户端
	client := yapi.NewYapiClientFromConfig(config)

	for {
		fmt.Println("\n请选择操作：")
		fmt.Println("1. 获取项目信息")
		fmt.Println("2. 获取接口列表")
		fmt.Println("3. 获取接口详情")
		fmt.Println("4. 查看当前配置")
		fmt.Println("5. 退出")
		fmt.Print("请输入选择 (1-5): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			handleGetProjectInfo(client)
		case "2":
			handleGetInterfaces(client)
		case "3":
			handleGetInterfaceDetail(client)
		case "4":
			yapi.GetConfigManager().PrintConfig()
		case "5":
			fmt.Println("退出程序")
			return
		default:
			fmt.Println("无效的选择，请重新输入")
		}
	}
}

// runConfigManagementExample 运行配置管理示例
func runConfigManagementExample() {
	fmt.Println("\n=== 配置管理示例 ===")

	cm := yapi.GetConfigManager()

	for {
		fmt.Println("\n请选择配置操作：")
		fmt.Println("1. 查看当前配置")
		fmt.Println("2. 生成示例配置文件")
		fmt.Println("3. 保存当前配置到文件")
		fmt.Println("4. 查看配置来源")
		fmt.Println("5. 验证配置")
		fmt.Println("6. 返回主菜单")
		fmt.Print("请输入选择 (1-6): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			cm.PrintConfig()
		case "2":
			handleGenerateExampleConfig(cm)
		case "3":
			handleSaveConfig(cm)
		case "4":
			handleShowConfigSource(cm)
		case "5":
			handleValidateConfig()
		case "6":
			return
		default:
			fmt.Println("无效的选择，请重新输入")
		}
	}
}

// handleGenerateExampleConfig 处理生成示例配置文件
func handleGenerateExampleConfig(cm *yapi.ConfigManager) {
	fmt.Print("请输入文件名 (默认: yapi.example.yaml): ")
	var filename string
	fmt.Scanln(&filename)
	if filename == "" {
		filename = "yapi.example.yaml"
	}

	err := cm.GenerateExampleConfig(filename)
	if err != nil {
		fmt.Printf("生成示例配置文件失败: %v\n", err)
	} else {
		fmt.Printf("示例配置文件已生成: %s\n", filename)
	}
}

// handleSaveConfig 处理保存配置
func handleSaveConfig(cm *yapi.ConfigManager) {
	fmt.Print("请输入文件名 (默认: yapi.yaml): ")
	var filename string
	fmt.Scanln(&filename)
	if filename == "" {
		filename = "yapi.yaml"
	}

	err := cm.SaveConfigFile(filename)
	if err != nil {
		fmt.Printf("保存配置文件失败: %v\n", err)
	} else {
		fmt.Printf("配置文件已保存: %s\n", filename)
	}
}

// handleShowConfigSource 处理显示配置来源
func handleShowConfigSource(cm *yapi.ConfigManager) {
	configKeys := []string{
		"base_url", "token", "timeout", "retry_count",
		"enable_cache", "cache_ttl", "cache_max_size",
		"log_level", "log_format", "enable_metrics",
	}

	fmt.Println("\n配置项来源:")
	for _, key := range configKeys {
		source := cm.GetConfigSource(key)
		fmt.Printf("  %s: %s\n", key, source)
	}
}

// handleValidateConfig 处理验证配置
func handleValidateConfig() {
	err := yapi.ValidateEnvironment()
	if err != nil {
		fmt.Printf("配置验证失败: %v\n", err)
	} else {
		fmt.Println("配置验证通过")
	}
}

// handleGetProjectInfo 处理获取项目信息
func handleGetProjectInfo(client yapi.YapiClient) {
	fmt.Print("请输入项目ID: ")
	var projectIDStr string
	fmt.Scanln(&projectIDStr)

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		fmt.Printf("无效的项目ID: %v\n", err)
		return
	}

	fmt.Printf("正在获取项目 %d 的信息...\n", projectID)
	response, err := client.GetProjectInfo(projectID)
	if err != nil {
		fmt.Printf("获取项目信息失败: %v\n", err)
		return
	}

	project := response.Data
	fmt.Printf("\n项目信息:\n")
	fmt.Printf("ID: %d\n", project.ID)
	fmt.Printf("名称: %s\n", project.Name)
	fmt.Printf("描述: %s\n", project.Description)
	fmt.Printf("基础路径: %s\n", project.BasePath)
	fmt.Printf("分组: %s (ID: %d)\n", project.GroupName, project.GroupID)
}

// handleGetInterfaces 处理获取接口列表
func handleGetInterfaces(client yapi.YapiClient) {
	fmt.Print("请输入项目ID: ")
	var projectIDStr string
	fmt.Scanln(&projectIDStr)

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		fmt.Printf("无效的项目ID: %v\n", err)
		return
	}

	fmt.Printf("正在获取项目 %d 的接口列表...\n", projectID)
	response, err := client.GetInterfaces(projectID)
	if err != nil {
		fmt.Printf("获取接口列表失败: %v\n", err)
		return
	}

	fmt.Printf("\n找到 %d 个接口:\n", response.Data.Count)
	for i, iface := range response.Data.List {
		fmt.Printf("%d. [%d] %s - %s %s (状态: %s)\n",
			i+1, iface.ID, iface.Title, iface.Method, iface.Path, iface.Status)
	}
}

// handleGetInterfaceDetail 处理获取接口详情
func handleGetInterfaceDetail(client yapi.YapiClient) {
	fmt.Print("请输入接口ID: ")
	var interfaceIDStr string
	fmt.Scanln(&interfaceIDStr)

	interfaceID, err := strconv.Atoi(interfaceIDStr)
	if err != nil {
		fmt.Printf("无效的接口ID: %v\n", err)
		return
	}

	fmt.Printf("正在获取接口 %d 的详情...\n", interfaceID)
	response, err := client.GetInterfaceDetail(interfaceID)
	if err != nil {
		fmt.Printf("获取接口详情失败: %v\n", err)
		return
	}

	detail := response.Data
	fmt.Printf("\n接口详情:\n")
	fmt.Printf("ID: %d\n", detail.ID)
	fmt.Printf("标题: %s\n", detail.Title)
	fmt.Printf("路径: %s\n", detail.Path)
	fmt.Printf("方法: %s\n", detail.Method)
	fmt.Printf("状态: %s\n", detail.Status)
	fmt.Printf("描述: %s\n", detail.Description)

	if len(detail.ReqHeaders) > 0 {
		fmt.Printf("\n请求头 (%d个):\n", len(detail.ReqHeaders))
		for i, header := range detail.ReqHeaders {
			if i >= 3 { // 只显示前3个
				fmt.Printf("  ... 还有 %d 个请求头\n", len(detail.ReqHeaders)-3)
				break
			}
			fmt.Printf("  %s: %s (%s)\n", header.Name, header.Value, header.Description)
		}
	}

	if len(detail.ReqQuery) > 0 {
		fmt.Printf("\n查询参数 (%d个):\n", len(detail.ReqQuery))
		for i, param := range detail.ReqQuery {
			if i >= 3 { // 只显示前3个
				fmt.Printf("  ... 还有 %d 个查询参数\n", len(detail.ReqQuery)-3)
				break
			}
			required := "可选"
			if param.Required == "1" {
				required = "必需"
			}
			fmt.Printf("  %s (%s): %s - %s\n", param.Name, param.Type, param.Description, required)
		}
	}

	if detail.ResBodyType != "" {
		fmt.Printf("\n响应体类型: %s\n", detail.ResBodyType)
		if len(detail.ResBody) > 100 {
			fmt.Printf("响应体示例: %s...\n", detail.ResBody[:100])
		} else {
			fmt.Printf("响应体示例: %s\n", detail.ResBody)
		}
	}
}

// init 初始化函数，检查配置
func init() {
	// 检查配置是否完整
	err := yapi.ValidateEnvironment()
	if err != nil {
		fmt.Printf("警告: %v\n", err)
		fmt.Println("请设置环境变量或配置文件，例如:")
		fmt.Println("环境变量:")
		fmt.Println("  export YAPI_BASE_URL=\"http://your-yapi-server.com\"")
		fmt.Println("  export YAPI_TOKEN=\"your_access_token\"")
		fmt.Println("配置文件 (yapi.yaml):")
		fmt.Println("  base_url: \"http://your-yapi-server.com\"")
		fmt.Println("  token: \"your_access_token\"")
		fmt.Println()
	}
}
