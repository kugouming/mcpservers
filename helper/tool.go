package helper

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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

// GetConfigDir 获取配置文件目录
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
