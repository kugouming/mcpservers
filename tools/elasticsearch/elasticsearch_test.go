package elasticsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var client, _ = NewESClient(&Config{
	URL:      "http://localhost:9200",
	Username: "elastic",
	Password: "fyOGTrjf",
})

func TestListIndices(t *testing.T) {
	indices, err := client.ListIndices(".internal.alerts-transform.health.alerts-default-000001")
	assert.NoError(t, err)
	t.Logf("所有索引: %v", indices)
}

func TestGetMapping(t *testing.T) {
	mapping, err := client.GetMapping(".internal.alerts-transform.health.alerts-default-000001")
	assert.NoError(t, err)
	t.Logf("映射: %v", mapping)
}

func TestGetShards(t *testing.T) {
	shards, err := client.GetShards(".internal.alerts-transform.health.alerts-default-000001")
	assert.NoError(t, err)
	t.Logf("分片: %v", shards)
}

func TestSearch(t *testing.T) {
	query := map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"message": "error",
			},
		},
	}
	results, err := client.Search(".internal.alerts-transform.health.alerts-default-000001", query)
	assert.NoError(t, err)
	t.Logf("搜索结果: %v", results)
}
