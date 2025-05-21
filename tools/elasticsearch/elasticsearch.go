package elasticsearch

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	es7 "github.com/elastic/go-elasticsearch/v7"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cast"
)

func NewESClient(config *Config) (IClient, error) {
	if config == nil || config.URL == "" {
		return nil, errors.New("配置为空")
	}

	cfg7 := es7.Config{
		Addresses: []string{config.URL},
	}
	cfg8 := es8.Config{
		Addresses: []string{config.URL},
	}

	if config.APIKey != "" {
		cfg7.APIKey = config.APIKey
		cfg8.APIKey = config.APIKey
	} else if config.Username != "" && config.Password != "" {
		cfg7.Username = config.Username
		cfg7.Password = config.Password
		cfg8.Username = config.Username
		cfg8.Password = config.Password
	}

	if config.CACert != "" {
		caCert, err := os.ReadFile(config.CACert)
		if err != nil {
			return nil, err
		}
		cfg7.CACert = caCert
		cfg8.CACert = caCert
	}

	client7, err := es7.NewClient(cfg7)
	if err != nil {
		return nil, err
	}

	version := GetVersion(client7)
	if version < 8 {
		return &es7Client{client: client7}, nil
	}

	client8, err := es8.NewClient(cfg8)
	if err != nil {
		return nil, err
	}
	return &es8Client{client: client8}, nil
}

func GetVersion(client *es7.Client) int {
	info, err := client.Info()
	if err != nil {
		return 7
	}
	defer info.Body.Close()

	var response map[string]any
	if err := json.NewDecoder(info.Body).Decode(&response); err != nil {
		return 7
	}

	version := response["version"].(map[string]any)
	number := version["number"].(string)
	parts := strings.Split(number, ".")

	return cast.ToInt(parts[0])
}
