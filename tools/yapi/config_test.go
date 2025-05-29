package yapi

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "", config.BaseURL)
	assert.Equal(t, "", config.Token)
	assert.Equal(t, 30, config.Timeout)
	assert.Equal(t, 3, config.RetryCount)
	assert.Equal(t, true, config.EnableCache)
	assert.Equal(t, 300, config.CacheTTL)
	assert.Equal(t, 1000, config.CacheMaxSize)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "text", config.LogFormat)
	assert.Equal(t, false, config.EnableMetrics)
	assert.Equal(t, []string{}, config.TrustedHosts)
}

func TestConfigManager_LoadConfig(t *testing.T) {
	// 备份原始环境变量
	originalBaseURL := os.Getenv("YAPI_BASE_URL")
	originalToken := os.Getenv("YAPI_TOKEN")
	defer func() {
		if originalBaseURL != "" {
			os.Setenv("YAPI_BASE_URL", originalBaseURL)
		} else {
			os.Unsetenv("YAPI_BASE_URL")
		}
		if originalToken != "" {
			os.Setenv("YAPI_TOKEN", originalToken)
		} else {
			os.Unsetenv("YAPI_TOKEN")
		}
	}()

	t.Run("通过环境变量配置", func(t *testing.T) {
		os.Setenv("YAPI_BASE_URL", "http://test-env.com")
		os.Setenv("YAPI_TOKEN", "test_env_token")
		os.Setenv("YAPI_TIMEOUT", "60")

		cm := NewConfigManager()
		err := cm.LoadConfig()
		require.NoError(t, err)

		config := cm.GetConfig()
		assert.Equal(t, "http://test-env.com", config.BaseURL)
		assert.Equal(t, "test_env_token", config.Token)
		assert.Equal(t, 60, config.Timeout)
	})

	t.Run("缺少必需配置项", func(t *testing.T) {
		os.Unsetenv("YAPI_BASE_URL")
		os.Unsetenv("YAPI_TOKEN")

		cm := NewConfigManager()
		err := cm.LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "base_url 配置项不能为空")
	})
}

func TestConfigManager_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "有效配置",
			config: &Config{
				BaseURL:      "http://valid.com",
				Token:        "valid_token",
				Timeout:      30,
				RetryCount:   3,
				CacheTTL:     300,
				CacheMaxSize: 1000,
				LogLevel:     "info",
				LogFormat:    "text",
			},
			wantErr: false,
		},
		{
			name: "缺少BaseURL",
			config: &Config{
				Token:        "valid_token",
				Timeout:      30,
				RetryCount:   3,
				CacheTTL:     300,
				CacheMaxSize: 1000,
				LogLevel:     "info",
				LogFormat:    "text",
			},
			wantErr: true,
			errMsg:  "base_url 配置项不能为空",
		},
		{
			name: "缺少Token",
			config: &Config{
				BaseURL:      "http://valid.com",
				Timeout:      30,
				RetryCount:   3,
				CacheTTL:     300,
				CacheMaxSize: 1000,
				LogLevel:     "info",
				LogFormat:    "text",
			},
			wantErr: true,
			errMsg:  "token 配置项不能为空",
		},
		{
			name: "无效的日志级别",
			config: &Config{
				BaseURL:      "http://valid.com",
				Token:        "valid_token",
				Timeout:      30,
				RetryCount:   3,
				CacheTTL:     300,
				CacheMaxSize: 1000,
				LogLevel:     "invalid",
				LogFormat:    "text",
			},
			wantErr: true,
			errMsg:  "无效的日志级别",
		},
		{
			name: "无效的超时时间",
			config: &Config{
				BaseURL:      "http://valid.com",
				Token:        "valid_token",
				Timeout:      0,
				RetryCount:   3,
				CacheTTL:     300,
				CacheMaxSize: 1000,
				LogLevel:     "info",
				LogFormat:    "text",
			},
			wantErr: true,
			errMsg:  "timeout 必须大于 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ConfigManager{config: tt.config}
			err := cm.validateConfig()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfigManager_GenerateExampleConfig(t *testing.T) {
	cm := NewConfigManager()

	// 生成示例配置文件
	filename := "test_example.yaml"
	err := cm.GenerateExampleConfig(filename)
	require.NoError(t, err)

	// 检查文件是否存在
	_, err = os.Stat(filename)
	assert.NoError(t, err)

	// 清理测试文件
	defer os.Remove(filename)
}

func TestConfigManager_SaveConfigFile(t *testing.T) {
	// 备份原始环境变量
	originalBaseURL := os.Getenv("YAPI_BASE_URL")
	originalToken := os.Getenv("YAPI_TOKEN")
	defer func() {
		if originalBaseURL != "" {
			os.Setenv("YAPI_BASE_URL", originalBaseURL)
		} else {
			os.Unsetenv("YAPI_BASE_URL")
		}
		if originalToken != "" {
			os.Setenv("YAPI_TOKEN", originalToken)
		} else {
			os.Unsetenv("YAPI_TOKEN")
		}
	}()

	os.Setenv("YAPI_BASE_URL", "http://test.com")
	os.Setenv("YAPI_TOKEN", "test_token")

	cm := NewConfigManager()
	err := cm.LoadConfig()
	require.NoError(t, err)

	// 保存配置文件
	filename := "test_config.yaml"
	err = cm.SaveConfigFile(filename)
	require.NoError(t, err)

	// 检查文件是否存在
	_, err = os.Stat(filename)
	assert.NoError(t, err)

	// 清理测试文件
	defer os.Remove(filename)
}

func TestConfigManager_PrintConfig(t *testing.T) {
	// 备份原始环境变量
	originalBaseURL := os.Getenv("YAPI_BASE_URL")
	originalToken := os.Getenv("YAPI_TOKEN")
	defer func() {
		if originalBaseURL != "" {
			os.Setenv("YAPI_BASE_URL", originalBaseURL)
		} else {
			os.Unsetenv("YAPI_BASE_URL")
		}
		if originalToken != "" {
			os.Setenv("YAPI_TOKEN", originalToken)
		} else {
			os.Unsetenv("YAPI_TOKEN")
		}
	}()

	os.Setenv("YAPI_BASE_URL", "http://test.com")
	os.Setenv("YAPI_TOKEN", "test_token_123456789")

	cm := NewConfigManager()
	err := cm.LoadConfig()
	require.NoError(t, err)

	// 测试打印配置（不会真正输出，只检查不会panic）
	assert.NotPanics(t, func() {
		cm.PrintConfig()
	})
}

func TestLoadGlobalConfig(t *testing.T) {
	// 备份原始环境变量
	originalBaseURL := os.Getenv("YAPI_BASE_URL")
	originalToken := os.Getenv("YAPI_TOKEN")
	defer func() {
		if originalBaseURL != "" {
			os.Setenv("YAPI_BASE_URL", originalBaseURL)
		} else {
			os.Unsetenv("YAPI_BASE_URL")
		}
		if originalToken != "" {
			os.Setenv("YAPI_TOKEN", originalToken)
		} else {
			os.Unsetenv("YAPI_TOKEN")
		}
	}()

	os.Setenv("YAPI_BASE_URL", "http://global-test.com")
	os.Setenv("YAPI_TOKEN", "global_test_token")

	config, err := LoadGlobalConfig()
	require.NoError(t, err)
	assert.Equal(t, "http://global-test.com", config.BaseURL)
	assert.Equal(t, "global_test_token", config.Token)
}

func TestValidateEnvironment(t *testing.T) {
	// 备份原始环境变量
	originalBaseURL := os.Getenv("YAPI_BASE_URL")
	originalToken := os.Getenv("YAPI_TOKEN")
	defer func() {
		if originalBaseURL != "" {
			os.Setenv("YAPI_BASE_URL", originalBaseURL)
		} else {
			os.Unsetenv("YAPI_BASE_URL")
		}
		if originalToken != "" {
			os.Setenv("YAPI_TOKEN", originalToken)
		} else {
			os.Unsetenv("YAPI_TOKEN")
		}
	}()

	t.Run("环境变量完整", func(t *testing.T) {
		os.Setenv("YAPI_BASE_URL", "http://test.com")
		os.Setenv("YAPI_TOKEN", "test_token")

		err := ValidateEnvironment()
		assert.NoError(t, err)
	})

	t.Run("缺少环境变量", func(t *testing.T) {
		os.Unsetenv("YAPI_BASE_URL")
		os.Unsetenv("YAPI_TOKEN")

		err := ValidateEnvironment()
		assert.Error(t, err)
	})
}

func TestContains(t *testing.T) {
	slice := []string{"debug", "info", "warn", "error"}

	assert.True(t, contains(slice, "info"))
	assert.False(t, contains(slice, "invalid"))
	assert.False(t, contains([]string{}, "test"))
}
