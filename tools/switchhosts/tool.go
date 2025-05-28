package switchhosts

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cast"
)

const (
	switchhostsStart = "#switchhosts_start"
	switchhostsEnd   = "#switchhosts_end"
)

// RegisterTool 注册本地 Hosts 管理工具
func RegisterTool(s *server.MCPServer) {
	s.AddTool(switchhostsTool, switchhostsHandler)
	s.AddTool(listhostsTool, listhostsHandler)
	s.AddTool(viewhostsTool, viewhostsHandler)
}

var switchhostsTool = mcp.NewTool(
	"switchhosts",
	mcp.WithDescription("本地 Hosts 管理工具，用于快速切换不同环境下的网络配置。"),
	mcp.WithString("conf_name",
		mcp.Required(),
		mcp.Description("配置名称，用于标识配置目录下不同的Hosts文件。"),
	),
)

var viewhostsTool = mcp.NewTool(
	"switchhosts_view",
	mcp.WithDescription("查看当前系统 hosts 文件内容。"),
	mcp.WithString("conf_name",
		mcp.Required(),
		mcp.Description("配置名称，用于标识配置目录下不同的Hosts文件。"),
	),
)

var listhostsTool = mcp.NewTool(
	"switchhosts_list",
	mcp.WithDescription("列出所有可用的配置名称。"),
)

// getHostFileName 获取 hosts 文件名
var getHostFileName = func(confName string) string {
	return fmt.Sprintf("hosts_%s.txt", confName)
}

// getSystemHostsPath 获取当前系统的 hosts 文件路径
var getSystemHostsPath = func() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("SystemRoot"), "System32", "drivers", "etc", "hosts")
	case "darwin", "linux":
		return "/etc/hosts"
	default:
		return "/etc/hosts"
	}
}

// getConfigDir 获取配置文件目录
var getConfigDir = func() string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Join(filepath.Dir(filepath.Dir(exePath)), "config", "switchhosts")
}

// getConfigFilePath 获取指定配置名称的 hosts 配置文件路径
var getConfigFilePath = func(confName string) string {
	confDir := getConfigDir()
	return filepath.Join(confDir, getHostFileName(confName))
}

var listConfigNames = func() []string {
	confDir := getConfigDir()
	files, err := os.ReadDir(confDir)
	if err != nil {
		return []string{}
	}
	confNames := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), "hosts_") && strings.HasSuffix(file.Name(), ".txt") {
			confNames = append(confNames, strings.TrimPrefix(strings.TrimSuffix(file.Name(), ".txt"), "hosts_"))
		}
	}
	return confNames
}

// removeSwitchhostsBlock 移除 hosts 文件中的 switchhosts 区块
func removeSwitchhostsBlock(content string) string {
	startIdx := strings.Index(content, switchhostsStart)
	endIdx := strings.Index(content, switchhostsEnd)
	if startIdx == -1 || endIdx == -1 || endIdx < startIdx {
		return content
	}
	// 保留 end 后的内容
	return content[:startIdx] + content[endIdx+len(switchhostsEnd):]
}

// switchhostsHandler 处理 hosts 切换请求
func switchhostsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	confName := cast.ToString(request.GetArguments()["conf_name"])

	var (
		err         error
		confContent []byte
	)
	// 读取配置文件内容
	if confName != "" {
		configPath := getConfigFilePath(confName)
		confContent, err = os.ReadFile(configPath)
		if err != nil && !os.IsNotExist(err) {
			return mcp.NewToolResultError(fmt.Sprintf("读取配置文件失败: %v", err)), nil
		}
	}

	debug := fmt.Sprintf("ConfContent: \n%s\n\n", string(confContent))

	// 读取系统 hosts 文件
	hostsPath := getSystemHostsPath()
	hostsContent, err := os.ReadFile(hostsPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("读取系统 hosts 文件失败: %v", err)), nil
	}

	debug += fmt.Sprintf("HostsContent: \n%s\n\n", string(hostsContent))

	// 移除系统 hosts 文件中的 switchhosts 区块
	newContent := removeSwitchhostsBlock(string(hostsContent))

	// 追加分隔符和配置内容
	if len(confContent) > 0 {
		block := "\n" + switchhostsStart + "\n" + string(confContent) + "\n" + switchhostsEnd + "\n"
		newContent += block
	}

	debug += fmt.Sprintf("NewContent: \n%s\n\n", newContent)

	// 写回系统 hosts 文件
	err = os.WriteFile(hostsPath, []byte(newContent), 0644)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("写入系统 hosts 文件失败: %v", err)), nil
	}

	log.Printf("Debug message: %v", debug)

	return mcp.NewToolResultText(fmt.Sprintf("已成功切换 hosts 配置为: %s\n\n%v", confName, debug)), nil
}

func listhostsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	confNames := listConfigNames()
	if len(confNames) == 0 {
		return mcp.NewToolResultText("暂无可用的配置名称"), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("可用的配置名称: %v", strings.Join(confNames, ", "))), nil
}

func viewhostsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	confName := cast.ToString(request.GetArguments()["conf_name"])
	if confName == "" {
		confName = "default"
	}

	hostsPath := getSystemHostsPath()
	if confName != "default" {
		hostsPath = getConfigFilePath(confName)
	}

	hostsContent, err := os.ReadFile(hostsPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("读取配置文件失败: %v", err)), nil
	}

	result := fmt.Sprintf("系统 hosts 文件内容: \n%s", string(hostsContent))
	if confName != "default" {
		result = fmt.Sprintf("Host %s 的内容: \n%s", confName, string(hostsContent))
	}

	return mcp.NewToolResultText(result), nil
}
