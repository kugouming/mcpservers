package yapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// YapiClient YAPI客户端接口
type YapiClient interface {
	// GetInterfaces 获取项目接口列表
	GetInterfaces(projectID int) (*YapiInterfaceListResponse, error)
	// GetInterfaceDetail 获取接口详细信息
	GetInterfaceDetail(interfaceID int) (*YapiInterfaceDetailResponse, error)
	// GetProjectInfo 获取项目信息
	GetProjectInfo(projectID int) (*YapiProjectResponse, error)
}

// YapiInterface YAPI接口信息结构体
type YapiInterface struct {
	ID     int      `json:"_id"`    // 接口ID
	Title  string   `json:"title"`  // 接口标题
	Path   string   `json:"path"`   // 接口路径
	Method string   `json:"method"` // HTTP方法
	Status string   `json:"status"` // 接口状态
	Tag    []string `json:"tag"`    // 标签
}

// YapiInterfaceDetail YAPI接口详细信息结构体
type YapiInterfaceDetail struct {
	ID           int          `json:"_id"`            // 接口ID
	Title        string       `json:"title"`          // 接口标题
	Path         string       `json:"path"`           // 接口路径
	Method       string       `json:"method"`         // HTTP方法
	Status       string       `json:"status"`         // 接口状态
	Description  string       `json:"desc"`           // 接口描述
	Markdown     string       `json:"markdown"`       // Markdown文档
	ReqHeaders   []YapiHeader `json:"req_headers"`    // 请求头
	ReqQuery     []YapiParam  `json:"req_query"`      // 查询参数
	ReqBodyForm  []YapiParam  `json:"req_body_form"`  // 表单参数
	ReqBodyOther string       `json:"req_body_other"` // 其他请求体
	ResBody      string       `json:"res_body"`       // 响应体
	ResBodyType  string       `json:"res_body_type"`  // 响应体类型
	Tag          []string     `json:"tag"`            // 标签
	ProjectID    int          `json:"project_id"`     // 项目ID
	AddTime      int64        `json:"add_time"`       // 添加时间
	UpTime       int64        `json:"up_time"`        // 更新时间
}

// YapiHeader 请求头结构体
type YapiHeader struct {
	Name        string `json:"name"`     // 头部名称
	Value       string `json:"value"`    // 头部值
	Description string `json:"desc"`     // 描述
	Required    string `json:"required"` // 是否必需
}

// YapiParam 参数结构体
type YapiParam struct {
	Name        string `json:"name"`     // 参数名称
	Value       string `json:"value"`    // 参数值
	Type        string `json:"type"`     // 参数类型
	Description string `json:"desc"`     // 描述
	Required    string `json:"required"` // 是否必需
	Example     string `json:"example"`  // 示例值
}

// YapiProject 项目信息结构体
type YapiProject struct {
	ID          int    `json:"_id"`        // 项目ID
	Name        string `json:"name"`       // 项目名称
	Description string `json:"desc"`       // 项目描述
	BasePath    string `json:"basepath"`   // 基础路径
	GroupID     int    `json:"group_id"`   // 分组ID
	GroupName   string `json:"group_name"` // 分组名称
	Color       string `json:"color"`      // 颜色
	Icon        string `json:"icon"`       // 图标
	AddTime     int64  `json:"add_time"`   // 添加时间
	UpTime      int64  `json:"up_time"`    // 更新时间
}

// YapiResponse YAPI通用响应结构体
type YapiResponse struct {
	ErrCode int    `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 错误消息
	Data    any    `json:"data"`    // 数据
}

// YapiInterfaceListResponse 接口列表响应结构体
type YapiInterfaceListResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Data    struct {
		Count int             `json:"count"` // 总数
		Total int             `json:"total"` // 总计
		List  []YapiInterface `json:"list"`  // 接口列表
	} `json:"data"`
}

// YapiInterfaceDetailResponse 接口详情响应结构体
type YapiInterfaceDetailResponse struct {
	ErrCode int                 `json:"errcode"`
	ErrMsg  string              `json:"errmsg"`
	Data    YapiInterfaceDetail `json:"data"`
}

// YapiProjectResponse 项目信息响应结构体
type YapiProjectResponse struct {
	ErrCode int         `json:"errcode"`
	ErrMsg  string      `json:"errmsg"`
	Data    YapiProject `json:"data"`
}

// yapiClientImpl YAPI客户端实现
type yapiClientImpl struct {
	config     *Config      // 配置
	httpClient *http.Client // HTTP客户端
}

// NewYapiClient 创建新的YAPI客户端
func NewYapiClient(baseURL, token string) YapiClient {
	config := &Config{
		BaseURL: baseURL,
		Token:   token,
		Timeout: 30,
	}
	return &yapiClientImpl{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second, // 使用配置的超时时间
		},
	}
}

// NewYapiClientFromConfig 从配置创建YAPI客户端
func NewYapiClientFromConfig(config *Config) YapiClient {
	return &yapiClientImpl{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// GetInterfaces 获取项目接口列表
func (c *yapiClientImpl) GetInterfaces(projectID int) (*YapiInterfaceListResponse, error) {
	url := fmt.Sprintf("%s/api/interface/list?project_id=%d&token=%s", c.config.BaseURL, projectID, c.config.Token)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
	}

	var response YapiInterfaceListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("解析响应失败: %w\n\nDebug:\n%s\n\n%v", err, url, string(body))
	}

	if response.ErrCode != 0 {
		return nil, fmt.Errorf("YAPI错误: %s (错误码: %d)", response.ErrMsg, response.ErrCode)
	}

	return &response, nil
}

// GetInterfaceDetail 获取接口详细信息
func (c *yapiClientImpl) GetInterfaceDetail(interfaceID int) (*YapiInterfaceDetailResponse, error) {
	url := fmt.Sprintf("%s/api/interface/get?id=%d&token=%s", c.config.BaseURL, interfaceID, c.config.Token)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
	}

	var response YapiInterfaceDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("解析响应失败: %w\n\nDebug:\n%s\n\n%v", err, url, string(body))
	}

	if response.ErrCode != 0 {
		return nil, fmt.Errorf("YAPI错误: %s (错误码: %d)", response.ErrMsg, response.ErrCode)
	}

	return &response, nil
}

// GetProjectInfo 获取项目信息
func (c *yapiClientImpl) GetProjectInfo(projectID int) (*YapiProjectResponse, error) {
	url := fmt.Sprintf("%s/api/project/get?id=%d&token=%s", c.config.BaseURL, projectID, c.config.Token)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
	}

	var response YapiProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("解析响应失败: %w\n\nDebug:\n%s\n\n%v", err, url, string(body))
	}

	if response.ErrCode != 0 {
		return nil, fmt.Errorf("YAPI错误: %s (错误码: %d)", response.ErrMsg, response.ErrCode)
	}

	return &response, nil
}
