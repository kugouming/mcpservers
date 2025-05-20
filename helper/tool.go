package helper

import (
	"context"
	"io"
	"log"
	"net/http"
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
