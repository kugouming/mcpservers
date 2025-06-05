package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kugouming/mcpservers/helper"
	"gopkg.in/yaml.v3"
)

// 全局缓存模板
var defaultPromptTemplates []*PromptTemplate

func LoadAllPromptTemplates() ([]*PromptTemplate, error) {
	var templates []*PromptTemplate
	files, err := os.ReadDir(helper.GetConfigDir("prompts"))
	if err != nil {
		return nil, fmt.Errorf("读取模板目录失败: %w", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !(strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") || strings.HasSuffix(file.Name(), ".json")) {
			continue
		}
		path := filepath.Join(helper.GetConfigDir("prompts"), file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("读取模板文件失败: %s, %w", path, err)
		}
		var tpl PromptTemplate
		if err := yaml.Unmarshal(content, &tpl); err != nil {
			return nil, fmt.Errorf("解析模板失败: %s, %w", path, err)
		}
		if tpl.Name == "" {
			tpl.Name = strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		}
		templates = append(templates, &tpl)
	}
	defaultPromptTemplates = templates
	return templates, nil
}

// GetPromptTemplateByName 根据名称查找模板
func GetPromptTemplateByName(name string) *PromptTemplate {
	for _, tpl := range defaultPromptTemplates {
		if tpl.Name == name {
			return tpl
		}
	}
	return nil
}

// ReloadPromptTemplates 重新加载所有模板
func ReloadPromptTemplates() error {
	_, err := LoadAllPromptTemplates()
	return err
}
