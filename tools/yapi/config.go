package yapi

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/kugouming/mcpservers/helper"
	"github.com/spf13/viper"
)

// Config YAPI工具配置结构体
type Config struct {
	// 基本配置
	BaseURL string `mapstructure:"base_url" json:"base_url" yaml:"base_url"` // YAPI服务器地址
	Token   string `mapstructure:"token" json:"token" yaml:"token"`          // 访问令牌

	// HTTP客户端配置
	Timeout     int  `mapstructure:"timeout" json:"timeout" yaml:"timeout"`                // 请求超时时间（秒）
	RetryCount  int  `mapstructure:"retry_count" json:"retry_count" yaml:"retry_count"`    // 重试次数
	EnableCache bool `mapstructure:"enable_cache" json:"enable_cache" yaml:"enable_cache"` // 是否启用缓存

	// 缓存配置
	CacheTTL     int `mapstructure:"cache_ttl" json:"cache_ttl" yaml:"cache_ttl"`                // 缓存存活时间（秒）
	CacheMaxSize int `mapstructure:"cache_max_size" json:"cache_max_size" yaml:"cache_max_size"` // 缓存最大条目数

	// 日志配置
	LogLevel  string `mapstructure:"log_level" json:"log_level" yaml:"log_level"`    // 日志级别
	LogFormat string `mapstructure:"log_format" json:"log_format" yaml:"log_format"` // 日志格式

	// 高级配置
	EnableMetrics bool     `mapstructure:"enable_metrics" json:"enable_metrics" yaml:"enable_metrics"` // 是否启用指标收集
	TrustedHosts  []string `mapstructure:"trusted_hosts" json:"trusted_hosts" yaml:"trusted_hosts"`    // 信任的主机列表
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		BaseURL:       "",
		Token:         "",
		Timeout:       30,
		RetryCount:    3,
		EnableCache:   true,
		CacheTTL:      300,  // 5分钟
		CacheMaxSize:  1000, // 1000个条目
		LogLevel:      "info",
		LogFormat:     "text",
		EnableMetrics: false,
		TrustedHosts:  []string{},
	}
}

// ConfigManager 配置管理器
type ConfigManager struct {
	viper  *viper.Viper
	config *Config
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager() *ConfigManager {
	v := viper.New()

	// 设置配置文件名和路径
	v.SetConfigName("yapi")                      // 配置文件名（不包含扩展名）
	v.SetConfigType("yaml")                      // 配置文件类型
	v.AddConfigPath(".")                         // 当前目录
	v.AddConfigPath(helper.GetConfigDir(""))     // config 目录
	v.AddConfigPath(helper.GetConfigDir("yapi")) // config/yapi 目录
	v.AddConfigPath("$HOME/.yapi")               // 用户主目录下的 .yapi 目录
	v.AddConfigPath("/etc/yapi/")                // 系统配置目录

	// 设置环境变量前缀
	v.SetEnvPrefix("YAPI")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv() // 自动绑定环境变量

	// 设置默认值
	setDefaults(v)

	return &ConfigManager{
		viper:  v,
		config: DefaultConfig(),
	}
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	defaults := DefaultConfig()

	v.SetDefault("base_url", defaults.BaseURL)
	v.SetDefault("token", defaults.Token)
	v.SetDefault("timeout", defaults.Timeout)
	v.SetDefault("retry_count", defaults.RetryCount)
	v.SetDefault("enable_cache", defaults.EnableCache)
	v.SetDefault("cache_ttl", defaults.CacheTTL)
	v.SetDefault("cache_max_size", defaults.CacheMaxSize)
	v.SetDefault("log_level", defaults.LogLevel)
	v.SetDefault("log_format", defaults.LogFormat)
	v.SetDefault("enable_metrics", defaults.EnableMetrics)
	v.SetDefault("trusted_hosts", defaults.TrustedHosts)
}

// LoadConfig 加载配置
// 优先级：环境变量 > 配置文件 > 默认值
func (cm *ConfigManager) LoadConfig() error {
	// 尝试读取配置文件
	if err := cm.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到，使用默认配置和环境变量
			fmt.Printf("配置文件未找到，使用默认配置和环境变量\n")
		} else {
			// 配置文件存在但读取出错
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	} else {
		fmt.Printf("使用配置文件: %s\n", cm.viper.ConfigFileUsed())
	}

	// 将配置解析到结构体
	if err := cm.viper.Unmarshal(cm.config); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	// 验证必需的配置项
	if err := cm.validateConfig(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	return nil
}

// validateConfig 验证配置
func (cm *ConfigManager) validateConfig() error {
	if cm.config.BaseURL == "" {
		return fmt.Errorf("base_url 配置项不能为空，请设置 YAPI_BASE_URL 环境变量或在配置文件中指定")
	} else {
		cm.config.BaseURL = strings.TrimRight(cm.config.BaseURL, "/")
	}

	if cm.config.Token == "" {
		return fmt.Errorf("token 配置项不能为空，请设置 YAPI_TOKEN 环境变量或在配置文件中指定")
	}

	if cm.config.Timeout <= 0 {
		return fmt.Errorf("timeout 必须大于 0")
	}

	if cm.config.RetryCount < 0 {
		return fmt.Errorf("retry_count 不能为负数")
	}

	if cm.config.CacheTTL <= 0 {
		return fmt.Errorf("cache_ttl 必须大于 0")
	}

	if cm.config.CacheMaxSize <= 0 {
		return fmt.Errorf("cache_max_size 必须大于 0")
	}

	// 验证日志级别
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, cm.config.LogLevel) {
		return fmt.Errorf("无效的日志级别: %s，支持的级别: %v", cm.config.LogLevel, validLogLevels)
	}

	// 验证日志格式
	validLogFormats := []string{"text", "json"}
	if !contains(validLogFormats, cm.config.LogFormat) {
		return fmt.Errorf("无效的日志格式: %s，支持的格式: %v", cm.config.LogFormat, validLogFormats)
	}

	return nil
}

// GetConfig 获取配置
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// GetViper 获取 viper 实例
func (cm *ConfigManager) GetViper() *viper.Viper {
	return cm.viper
}

// SaveConfigFile 保存配置到文件
func (cm *ConfigManager) SaveConfigFile(filename string) error {
	if filename == "" {
		filename = "yapi.yaml"
	}

	// 确保文件扩展名正确
	if !strings.HasSuffix(filename, ".yaml") && !strings.HasSuffix(filename, ".yml") {
		filename += ".yaml"
	}

	return cm.viper.WriteConfigAs(filename)
}

// PrintConfig 打印当前配置（隐藏敏感信息）
func (cm *ConfigManager) PrintConfig() {
	config := *cm.config

	// 隐藏敏感信息
	if config.Token != "" {
		if len(config.Token) > 8 {
			config.Token = config.Token[:4] + "****" + config.Token[len(config.Token)-4:]
		} else {
			config.Token = "****"
		}
	}

	fmt.Printf("当前配置:\n")
	fmt.Printf("  Base URL: %s\n", config.BaseURL)
	fmt.Printf("  Token: %s\n", config.Token)
	fmt.Printf("  Timeout: %d秒\n", config.Timeout)
	fmt.Printf("  Retry Count: %d\n", config.RetryCount)
	fmt.Printf("  Enable Cache: %t\n", config.EnableCache)
	fmt.Printf("  Cache TTL: %d秒\n", config.CacheTTL)
	fmt.Printf("  Cache Max Size: %d\n", config.CacheMaxSize)
	fmt.Printf("  Log Level: %s\n", config.LogLevel)
	fmt.Printf("  Log Format: %s\n", config.LogFormat)
	fmt.Printf("  Enable Metrics: %t\n", config.EnableMetrics)
	fmt.Printf("  Trusted Hosts: %v\n", config.TrustedHosts)

	if cm.viper.ConfigFileUsed() != "" {
		fmt.Printf("  Config File: %s\n", cm.viper.ConfigFileUsed())
	} else {
		fmt.Printf("  Config File: 未使用配置文件\n")
	}
}

// GetConfigSource 获取配置项的来源
func (cm *ConfigManager) GetConfigSource(key string) string {
	if cm.viper.IsSet(key) && cm.viper.GetString(key) != "" {
		return "环境变量/配置文件"
	}
	return "默认值"
}

// WatchConfig 监听配置文件变化
func (cm *ConfigManager) WatchConfig(callback func()) {
	cm.viper.WatchConfig()
	cm.viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件发生变化: %s\n", e.Name)

		// 重新加载配置
		if err := cm.viper.Unmarshal(cm.config); err != nil {
			fmt.Printf("重新加载配置失败: %v\n", err)
			return
		}

		// 重新验证配置
		if err := cm.validateConfig(); err != nil {
			fmt.Printf("配置验证失败: %v\n", err)
			return
		}

		fmt.Printf("配置已重新加载\n")

		// 执行回调函数
		if callback != nil {
			callback()
		}
	})
}

// GenerateExampleConfig 生成示例配置文件
func (cm *ConfigManager) GenerateExampleConfig(filename string) error {
	if filename == "" {
		filename = "yapi.example.yaml"
	}

	// 设置示例配置
	exampleConfig := map[string]interface{}{
		"base_url":       "http://your-yapi-server.com",
		"token":          "your_access_token",
		"timeout":        30,
		"retry_count":    3,
		"enable_cache":   true,
		"cache_ttl":      300,
		"cache_max_size": 1000,
		"log_level":      "info",
		"log_format":     "text",
		"enable_metrics": false,
		"trusted_hosts":  []string{"*.example.com", "localhost"},
	}

	// 创建临时 viper 实例用于生成配置文件
	v := viper.New()
	for key, value := range exampleConfig {
		v.Set(key, value)
	}

	return v.WriteConfigAs(filename)
}

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// 全局配置管理器实例
var configManager *ConfigManager

// GetConfigManager 获取全局配置管理器实例
func GetConfigManager() *ConfigManager {
	if configManager == nil {
		configManager = NewConfigManager()
	}
	return configManager
}

// LoadGlobalConfig 加载全局配置
func LoadGlobalConfig() (*Config, error) {
	cm := GetConfigManager()
	if err := cm.LoadConfig(); err != nil {
		return nil, err
	}
	return cm.GetConfig(), nil
}

// ValidateEnvironment 验证环境变量是否正确配置（兼容原有函数）
func ValidateEnvironment() error {
	config, err := LoadGlobalConfig()
	if err != nil {
		return err
	}

	if config.BaseURL == "" {
		return fmt.Errorf("YAPI_BASE_URL 环境变量未设置")
	}

	if config.Token == "" {
		return fmt.Errorf("YAPI_TOKEN 环境变量未设置")
	}

	return nil
}
