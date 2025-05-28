package switchhosts

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func TestSwitchhostsHandler(t *testing.T) {
	tempDir := t.TempDir()
	// 模拟 config/switchhosts 目录和配置文件
	configDir := filepath.Join(tempDir, "config", "switchhosts")
	os.MkdirAll(configDir, 0755)
	confName := "dev"
	confFile := filepath.Join(configDir, "hosts_"+confName+".txt")
	hostsContent := "127.0.0.1 dev.local"
	ioutil.WriteFile(confFile, []byte(hostsContent), 0644)

	// 模拟系统 hosts 文件
	hostsPath := filepath.Join(tempDir, "hosts")
	originHosts := "# user hosts\n192.168.1.1 prod.local\n"
	ioutil.WriteFile(hostsPath, []byte(originHosts), 0644)

	// 覆盖 getSystemHostsPath/getConfigFilePath 以便测试
	oldGetSystemHostsPath := getSystemHostsPath
	oldGetConfigFilePath := getConfigFilePath
	getSystemHostsPath = func() string { return hostsPath }
	getConfigFilePath = func(name string) string { return confFile }
	defer func() {
		getSystemHostsPath = oldGetSystemHostsPath
		getConfigFilePath = oldGetConfigFilePath
	}()

	// 构造请求
	request := mcp.CallToolRequest{
		Params: struct {
			Name      string      `json:"name"`
			Arguments interface{} `json:"arguments,omitempty"`
			Meta      *mcp.Meta   `json:"_meta,omitempty"`
		}{
			Name:      "switchhosts",
			Arguments: map[string]interface{}{"conf_name": confName},
		},
	}

	result, err := switchhostsHandler(context.Background(), request)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, confName)

	// 检查 hosts 文件内容
	finalContent, _ := ioutil.ReadFile(hostsPath)
	assert.Contains(t, string(finalContent), hostsContent)
	assert.Contains(t, string(finalContent), "#switchhosts_start")
	assert.Contains(t, string(finalContent), "#switchhosts_end")
}
