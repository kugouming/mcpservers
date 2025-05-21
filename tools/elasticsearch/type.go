package elasticsearch

// IClient 定义 Elasticsearch 客户端接口
type IClient interface {
	ListIndices(pattern string) ([]map[string]any, error)
	GetMapping(index string) (map[string]any, error)
	Search(index string, query map[string]any) (map[string]any, error)
	GetShards(index string) ([]map[string]any, error)
}

// Config 定义 Elasticsearch 配置
type Config struct {
	URL      string `json:"url"`
	APIKey   string `json:"apiKey,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	CACert   string `json:"caCert,omitempty"`
}

type CatIndicesRow struct {
	Health       string `json:"health"`              // "green", "yellow", or "red"
	Status       string `json:"status"`              // "open" or "closed"
	Index        string `json:"index"`               // index name
	UUID         string `json:"uuid"`                // index uuid
	Pri          int    `json:"pri,string"`          // number of primary shards
	Rep          int    `json:"rep,string"`          // number of replica shards
	DocsCount    int    `json:"docs.count,string"`   // number of available documents
	DocsDeleted  int    `json:"docs.deleted,string"` // number of deleted documents
	StoreSize    string `json:"store.size"`          // store size of primaries & replicas, e.g. "4.6kb"
	PriStoreSize string `json:"pri.store.size"`      // store size of primaries, e.g. "230b"
	DatasetSize  string `json:"dataset.size"`        // store size of primaries, e.g. "230b"
}
