package elasticsearch

import (
	"fmt"
	"os"
	"testing"

	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/assert"
)

var config = &Config{
	URL:      "http://localhost:9200",
	Username: "elastic",
	Password: "fyOGTrjf",
}

func init() {
	var err error
	client, err = NewESClient(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Elasticsearch client: %v", err))
	}
}

func TestGetVersion(t *testing.T) {
	cfg7 := es7.Config{
		Addresses: []string{config.URL},
	}
	if config.APIKey != "" {
		cfg7.APIKey = config.APIKey
	} else if config.Username != "" && config.Password != "" {
		cfg7.Username = config.Username
		cfg7.Password = config.Password
	}

	if config.CACert != "" {
		caCert, _ := os.ReadFile(config.CACert)
		cfg7.CACert = caCert
	}

	client7, _ := es7.NewClient(cfg7)
	version := GetVersion(client7)
	assert.Equal(t, 9, version)
}

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
