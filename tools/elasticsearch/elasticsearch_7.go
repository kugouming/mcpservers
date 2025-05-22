package elasticsearch

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// esClient 实现 ESClient 接口的真实客户端
type es7Client struct {
	client *es7.Client
}

func (c *es7Client) ListIndices(pattern string) ([]map[string]any, error) {
	res, err := c.client.Cat.Indices(
		c.client.Cat.Indices.WithIndex(pattern),
		c.client.Cat.Indices.WithFormat("json"),
	)
	if err != nil {
		return nil, fmt.Errorf("get indices failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%v %v", res.StatusCode, string(body))
	}

	var indices []CatIndicesRow
	if err := json.NewDecoder(res.Body).Decode(&indices); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	result := make([]map[string]any, 0, len(indices))
	for _, index := range indices {
		result = append(result, map[string]any{
			"index":     index.Index,
			"health":    index.Health,
			"status":    index.Status,
			"docsCount": index.DocsCount,
			"storeSize": index.StoreSize,
		})
	}

	return result, nil
}

func (c *es7Client) GetMapping(index string) (map[string]any, error) {
	res, err := c.client.Indices.GetMapping(
		c.client.Indices.GetMapping.WithIndex(index),
	)
	if err != nil {
		return nil, fmt.Errorf("get mapping failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%v %v", res.StatusCode, string(body))
	}

	var response map[string]any
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	item := response[index].(map[string]any)
	return item["mappings"].(map[string]any), nil
}

func (c *es7Client) Search(index string, query map[string]any) (map[string]any, error) {
	// 获取映射以识别文本字段
	mappingRes, err := c.GetMapping(index)
	if err != nil {
		return nil, fmt.Errorf("get mapping failed: %w", err)
	}

	// 构建搜索请求
	query["highlight"] = map[string]any{
		"fields":    map[string]any{},
		"pre_tags":  []string{"<em>"},
		"post_tags": []string{"</em>"},
	}

	// 添加文本字段到高亮
	if props, ok := mappingRes["properties"].(map[string]any); ok {
		for field, fieldData := range props {
			if fieldMap, ok := fieldData.(map[string]any); ok {
				if fieldType, ok := fieldMap["type"].(string); ok && fieldType == "text" {
					query["highlight"].(map[string]any)["fields"].(map[string]any)[field] = map[string]any{}
				}
			}
		}
	}

	// 执行搜索
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("params to json failed: %w", err)
	}

	res, err := c.client.Search(
		c.client.Search.WithIndex(index),
		c.client.Search.WithBody(strings.NewReader(string(queryJSON))),
	)
	if err != nil {
		return nil, fmt.Errorf("search run failed: %w", err)
	}
	defer res.Body.Close()

	var searchResponse map[string]any
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("search parse response failed: %w", err)
	}

	if res.IsError() {
		err := searchResponse["error"].(map[string]any)
		return nil, fmt.Errorf("%v %v", err["type"], err["reason"])
	}

	return searchResponse["hits"].(map[string]any), nil
}

func (c *es7Client) GetShards(index string) ([]map[string]any, error) {
	var res *esapi.Response
	var err error

	if index != "" {
		res, err = c.client.Cat.Shards(
			c.client.Cat.Shards.WithIndex(index),
			c.client.Cat.Shards.WithFormat("json"),
		)
	} else {
		res, err = c.client.Cat.Shards(
			c.client.Cat.Shards.WithFormat("json"),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("get shards failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%v %v", res.StatusCode, string(body))
	}

	var shards []map[string]any
	if err := json.NewDecoder(res.Body).Decode(&shards); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	return shards, nil
}
