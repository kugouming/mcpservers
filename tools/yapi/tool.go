package yapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kugouming/mcpservers/helper"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// 全局YAPI客户端实例
var yapiClient YapiClient

// RegisterTool 注册YAPI工具到MCP服务器
// 这是主要的注册函数，用于将所有YAPI相关工具注册到MCP服务器
func RegisterTool(s *server.MCPServer) {
	// 初始化YAPI客户端
	initYapiClient()

	// 注册所有工具
	registerGetInterfacesTool(s)
	registerGetInterfaceByPathTool(s)
	registerGetInterfaceByNameTool(s)
	registerGetInterfaceDetailTool(s)
	registerGetProjectInfoTool(s)
}

// initYapiClient 初始化YAPI客户端
func initYapiClient() {
	// 使用新的配置管理系统
	config, err := LoadGlobalConfig()
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	yapiClient = NewYapiClientFromConfig(config)
}

// registerGetInterfacesTool 注册获取接口列表工具
func registerGetInterfacesTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_interfaces",
		mcp.WithDescription("获取YAPI项目的接口列表，返回包含接口ID、标题、路径、方法等信息的列表"),
		mcp.WithNumber("project_id",
			mcp.Required(),
			mcp.Description("YAPI项目ID，用于指定要获取接口列表的项目"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if yapiClient == nil {
			return mcp.NewToolResultError("YAPI客户端未初始化，请检查环境变量配置"), nil
		}

		// 获取项目ID参数
		projectIDFloat, ok := request.GetArguments()["project_id"].(float64)
		if !ok {
			return mcp.NewToolResultError("project_id 参数必须是数字"), nil
		}
		projectID := int(projectIDFloat)

		// 调用YAPI API获取接口列表
		response, err := yapiClient.GetInterfaces(projectID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("获取接口列表失败: %s", err.Error())), nil
		}

		// 格式化输出结果
		var interfaceList []map[string]any
		for _, iface := range response.Data.List {
			interfaceList = append(interfaceList, map[string]any{
				"id":     iface.ID,
				"title":  iface.Title,
				"path":   iface.Path,
				"method": iface.Method,
				"status": iface.Status,
				"tag":    iface.Tag,
			})
		}

		result := map[string]any{
			"total":      response.Data.Total,
			"count":      response.Data.Count,
			"interfaces": interfaceList,
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	s.AddTool(tool, handler)
}

// registerGetInterfaceByPathTool 注册根据接口路径获取接口ID工具
func registerGetInterfaceByPathTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_interface_by_path",
		mcp.WithDescription("根据接口路径获取接口ID，返回包含接口ID、标题、路径、方法等信息的列表"),
		mcp.WithNumber("project_id",
			mcp.Required(),
			mcp.Description("YAPI项目ID，用于指定要获取接口列表的项目"),
		),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("接口路径，用于指定要获取接口的ID"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if yapiClient == nil {
			return mcp.NewToolResultError("YAPI客户端未初始化，请检查环境变量配置"), nil
		}

		// 获取项目ID参数
		projectIDFloat, ok := request.GetArguments()["project_id"].(float64)
		if !ok {
			return mcp.NewToolResultError("project_id 参数必须是数字"), nil
		}
		projectID := int(projectIDFloat)
		path, ok := request.GetArguments()["path"].(string)
		if !ok {
			return mcp.NewToolResultError("path 参数必须是字符串"), nil
		}

		// 调用YAPI API获取接口列表
		response, err := yapiClient.GetInterfaces(projectID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("获取接口列表失败: %s", err.Error())), nil
		}

		// 格式化输出结果
		var interfaceItem, interfaceItem2 []map[string]any
		for _, iface := range response.Data.List {
			if iface.Path == path {
				interfaceItem = []map[string]any{
					{
						"id":     iface.ID,
						"title":  iface.Title,
						"path":   iface.Path,
						"method": iface.Method,
						"status": iface.Status,
						"tag":    iface.Tag,
					},
				}
				break
			} else if strings.HasSuffix(path, iface.Path) {
				interfaceItem2 = append(interfaceItem2, map[string]any{
					"id":     iface.ID,
					"title":  iface.Title,
					"path":   iface.Path,
					"method": iface.Method,
					"status": iface.Status,
					"tag":    iface.Tag,
				})
			}
		}

		if len(interfaceItem) == 0 && len(interfaceItem2) > 0 {
			interfaceItem = interfaceItem2
		}

		if len(interfaceItem) == 0 {
			return mcp.NewToolResultError(fmt.Sprintf("接口 %s 未找到", path)), nil
		}

		return mcp.NewToolResultText(helper.MarshalIndent(interfaceItem)), nil
	}

	s.AddTool(tool, handler)
}

// registerGetInterfaceByNameTool 注册根据接口名称获取接口ID工具
func registerGetInterfaceByNameTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_interface_by_name",
		mcp.WithDescription("根据接口名称获取接口ID，返回包含接口ID、标题、路径、方法等信息的列表"),
		mcp.WithNumber("project_id",
			mcp.Required(),
			mcp.Description("YAPI项目ID，用于指定要获取接口列表的项目"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("接口名称，用于指定要获取接口的ID"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if yapiClient == nil {
			return mcp.NewToolResultError("YAPI客户端未初始化，请检查环境变量配置"), nil
		}

		// 获取项目ID参数
		projectIDFloat, ok := request.GetArguments()["project_id"].(float64)
		if !ok {
			return mcp.NewToolResultError("project_id 参数必须是数字"), nil
		}
		projectID := int(projectIDFloat)
		name, ok := request.GetArguments()["name"].(string)
		if !ok {
			return mcp.NewToolResultError("name 参数必须是字符串"), nil
		}

		// 调用YAPI API获取接口列表
		response, err := yapiClient.GetInterfaces(projectID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("获取接口列表失败: %s", err.Error())), nil
		}

		// 格式化输出结果
		var interfaceItem map[string]any
		for _, iface := range response.Data.List {
			if iface.Title == name {
				interfaceItem = map[string]any{
					"id":     iface.ID,
					"title":  iface.Title,
					"path":   iface.Path,
					"method": iface.Method,
					"status": iface.Status,
					"tag":    iface.Tag,
				}
				break
			}
		}

		if interfaceItem == nil {
			return mcp.NewToolResultError(fmt.Sprintf("接口名称 %s 未找到", name)), nil
		}

		return mcp.NewToolResultText(helper.MarshalIndent(interfaceItem)), nil
	}

	s.AddTool(tool, handler)
}

// registerGetInterfaceDetailTool 注册获取接口详情工具
func registerGetInterfaceDetailTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_interface_detail",
		mcp.WithDescription("获取YAPI接口的详细信息，包括请求参数、响应格式、文档说明等完整信息"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("接口ID，用于指定要获取详细信息的接口"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if yapiClient == nil {
			return mcp.NewToolResultError("YAPI客户端未初始化，请检查环境变量配置"), nil
		}

		// 获取接口ID参数
		interfaceIDFloat, ok := request.GetArguments()["id"].(float64)
		if !ok {
			return mcp.NewToolResultError("id 参数必须是数字"), nil
		}
		interfaceID := int(interfaceIDFloat)

		// 调用YAPI API获取接口详情
		response, err := yapiClient.GetInterfaceDetail(interfaceID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("获取接口详情失败: %s", err.Error())), nil
		}

		// 格式化输出结果
		detail := response.Data
		result := map[string]any{
			"id":             detail.ID,
			"title":          detail.Title,
			"path":           detail.Path,
			"method":         detail.Method,
			"status":         detail.Status,
			"description":    detail.Description,
			"markdown":       detail.Markdown,
			"req_headers":    detail.ReqHeaders,
			"req_query":      detail.ReqQuery,
			"req_body_form":  detail.ReqBodyForm,
			"req_body_other": detail.ReqBodyOther,
			"res_body":       detail.ResBody,
			"res_body_type":  detail.ResBodyType,
			"tag":            detail.Tag,
			"project_id":     detail.ProjectID,
			"add_time":       time.Unix(detail.AddTime, 0).Format("2006-01-02 15:04:05"),
			"up_time":        time.Unix(detail.UpTime, 0).Format("2006-01-02 15:04:05"),
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	s.AddTool(tool, handler)
}

// registerGetProjectInfoTool 注册获取项目信息工具
func registerGetProjectInfoTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_project_info",
		mcp.WithDescription("获取YAPI项目的基本信息，包括项目名称、描述、基础路径等"),
		mcp.WithNumber("project_id",
			mcp.Required(),
			mcp.Description("YAPI项目ID，用于指定要获取信息的项目"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if yapiClient == nil {
			return mcp.NewToolResultError("YAPI客户端未初始化，请检查环境变量配置"), nil
		}

		// 获取项目ID参数
		projectIDFloat, ok := request.GetArguments()["project_id"].(float64)
		if !ok {
			return mcp.NewToolResultError("project_id 参数必须是数字"), nil
		}
		projectID := int(projectIDFloat)

		// 调用YAPI API获取项目信息
		response, err := yapiClient.GetProjectInfo(projectID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("获取项目信息失败: %s", err.Error())), nil
		}

		// 格式化输出结果
		project := response.Data
		result := map[string]any{
			"id":          project.ID,
			"name":        project.Name,
			"description": project.Description,
			"basepath":    project.BasePath,
			"group_id":    project.GroupID,
			"group_name":  project.GroupName,
			"color":       project.Color,
			"icon":        project.Icon,
			"add_time":    time.Unix(project.AddTime, 0).Format("2006-01-02 15:04:05"),
			"up_time":     time.Unix(project.UpTime, 0).Format("2006-01-02 15:04:05"),
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	s.AddTool(tool, handler)
}

// FormatJSONResponse 格式化JSON响应，用于美化输出
func FormatJSONResponse(data any) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("格式化JSON失败: %v", err)
	}
	return string(jsonData)
}
