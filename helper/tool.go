package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// HttpRequest 函数用于发起HTTP请求
//
// 参数:
//
//	ctx:     context.Context   用于控制请求的上下文
//	method:  string            HTTP请求的方法（如GET, POST等）
//	url:     string            请求的URL
//	headers: map[string]string 请求的HTTP头信息
//	body:    string            请求体内容
//
// 返回值:
//
//	int    HTTP响应的状态码
//	[]byte HTTP响应的内容
//	error  可能发生的错误
func HttpRequest(ctx context.Context, method, url string, headers map[string]string, body string) (int, []byte, error) {
	log.Printf("\n\n\tMethod: %s \n\tUrl: %s \n\tHeaders: %v \n\tBody: %s\n\n", method, url, headers, body)
	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
	if err != nil {
		return 0, nil, err
	}

	// 添加请求头
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	log.Printf("Code: %d Body: %s\n", resp.StatusCode, string(responseBody))

	return resp.StatusCode, responseBody, nil
}

// GetConfigDir 根据工具名称获取配置目录路径
//
// 参数:
//
//	toolName: 工具名称，若为空则不添加工具名称到配置目录路径中
//
// 返回值:
//
//	返回配置目录路径，若获取可执行文件路径失败则返回空字符串
func GetConfigDir(toolName string) string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	dirList := []string{filepath.Dir(filepath.Dir(exePath)), "config"}
	if toolName != "" {
		dirList = append(dirList, toolName)
	}
	return filepath.Join(dirList...)
}

// MarshalIndent 将数据转换为文本格式
func MarshalIndent(indices any) string {
	body, err := json.MarshalIndent(indices, "", "  ")
	if err != nil {
		return fmt.Sprintf("Failed to marshal indices: %v", err)
	}

	return string(body)
}
