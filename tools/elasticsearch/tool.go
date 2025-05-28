package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var client IClient

// RegisterTool 注册HTTP请求工具
func RegisterTool(s *server.MCPServer) {
	initClient()
	ListIndicesTool(s)
	GetMappingTool(s)
	SearchTool(s)
	GetShardsTool(s)
}

func initClient() {
	config := &Config{
		URL:       os.Getenv("ES_URL"),
		APIKey:    os.Getenv("ES_API_KEY"),
		Username:  os.Getenv("ES_USERNAME"),
		Password:  os.Getenv("ES_PASSWORD"),
		AuthToken: os.Getenv("ES_AUTHTOKEN"),
		CACert:    os.Getenv("ES_CACERT"),
	}

	var err error
	client, err = NewESClient(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Elasticsearch client: %v", err))
	}
}

// ListIndicesTool 用于列出所有可用的 Elasticsearch 索引
func ListIndicesTool(s *server.MCPServer) {
	tool := mcp.NewTool("list_indices",
		mcp.WithDescription(`List all available Elasticsearch indices`),
		mcp.WithString("indexPattern",
			mcp.Required(),
			mcp.MinLength(1),
			mcp.Description("Index pattern of Elasticsearch indices to list"),
		),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		indices, err := client.ListIndices(request.GetArguments()["indexPattern"].(string))
		if err != nil {
			return mcp.NewToolResultErrorFromErr("list_indices tool failed", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Found %d indices \n\nResult: \n%s", len(indices), mapToText(indices))), nil
	}

	s.AddTool(tool, handler)
}

func GetMappingTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_mappings",
		mcp.WithDescription(`Get field mappings for a specific Elasticsearch index`),
		mcp.WithString("index",
			mcp.Required(),
			mcp.MinLength(1),
			mcp.Description("Name of the Elasticsearch index to get mappings for"),
		),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		index := request.GetArguments()["index"].(string)
		mappings, err := client.GetMapping(index)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("get_mappings tool failed", err), nil
		}

		var result string
		if len(mappings) == 0 {
			result = fmt.Sprintf("No mappings found for index %s", index)
		} else if mapping, has := mappings[index]; has {
			result = fmt.Sprintf("Mappings for index %s: \n%s", index, mapToText(mapping.(map[string]any)["mappings"]))
		} else {
			for k, v := range mappings {
				if vv, ok := v.(map[string]any); ok {
					result += fmt.Sprintf("Mappings for index %s: \n%s\n\n", k, mapToText(vv["mappings"]))
				}
			}
		}
		return mcp.NewToolResultText(result), nil
	}
	s.AddTool(tool, handler)
}

func SearchTool(s *server.MCPServer) {
	tool := mcp.NewTool("search",
		mcp.WithDescription(`Perform an Elasticsearch search with the provided query DSL. Highlights are always enabled.`),
		mcp.WithString("index",
			mcp.Required(),
			mcp.MinLength(1),
			mcp.Description("Name of the Elasticsearch index to search"),
		),
		mcp.WithObject("query",
			mcp.Required(),
			mcp.Description("Complete Elasticsearch query DSL object that can include query, size, from, sort, etc."),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		index := request.GetArguments()["index"].(string)
		query := request.GetArguments()["query"].(map[string]any)
		hits, err := client.Search(index, query)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("search tool failed", err), nil
		}

		total := hits["total"].(float64)
		if total == 0 {
			return mcp.NewToolResultText("No results found"), nil
		}

		hitsArray := hits["hits"].([]any)

		from := 1
		if fromVal, ok := query["from"].(float64); ok {
			from = int(fromVal)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Total results: %.0f, showing %d from position %d \n\nResult: \n%s", total, len(hitsArray), from, mapToText(hitsArray))), nil
	}
	s.AddTool(tool, handler)
}

func GetShardsTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_shards",
		mcp.WithDescription(`Get shard information for all or specific indices`),
		mcp.WithString("index",
			mcp.Description("Optional index name to get shard information for"),
		),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		index := ""
		if indexName, has := request.GetArguments()["index"]; has {
			index = indexName.(string)
		}
		shards, err := client.GetShards(index)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("get_shards tool failed", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Found %d shards \n\nResult: \n%s", len(shards), mapToText(shards))), nil
	}
	s.AddTool(tool, handler)
}

// mapToText 将数据转换为文本格式
func mapToText(indices any) string {
	body, err := json.MarshalIndent(indices, "", "  ")
	if err != nil {
		return fmt.Sprintf("Failed to marshal indices: %v", err)
	}

	return string(body)
}
