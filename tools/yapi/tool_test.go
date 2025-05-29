package yapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockYapiServer 模拟YAPI服务器
func MockYapiServer() *httptest.Server {
	mux := http.NewServeMux()

	// 模拟接口列表API
	mux.HandleFunc("/api/interface/list", func(w http.ResponseWriter, r *http.Request) {
		projectID := r.URL.Query().Get("project_id")
		token := r.URL.Query().Get("token")

		if token != "test_token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if projectID == "123" {
			response := YapiInterfaceListResponse{
				ErrCode: 0,
				ErrMsg:  "成功",
				Data: struct {
					Count int             `json:"count"`
					Total int             `json:"total"`
					List  []YapiInterface `json:"list"`
				}{
					Count: 2,
					Total: 2,
					List: []YapiInterface{
						{
							ID:     1001,
							Title:  "获取用户信息",
							Path:   "/api/user/info",
							Method: "GET",
							Status: "done",
							Tag:    []string{"用户", "基础"},
						},
						{
							ID:     1002,
							Title:  "创建用户",
							Path:   "/api/user/create",
							Method: "POST",
							Status: "done",
							Tag:    []string{"用户", "管理"},
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// 模拟接口详情API
	mux.HandleFunc("/api/interface/get", func(w http.ResponseWriter, r *http.Request) {
		interfaceID := r.URL.Query().Get("id")
		token := r.URL.Query().Get("token")

		if token != "test_token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if interfaceID == "1001" {
			response := YapiInterfaceDetailResponse{
				ErrCode: 0,
				ErrMsg:  "成功",
				Data: YapiInterfaceDetail{
					ID:          1001,
					Title:       "获取用户信息",
					Path:        "/api/user/info",
					Method:      "GET",
					Status:      "done",
					Description: "根据用户ID获取用户详细信息",
					Markdown:    "# 获取用户信息\n\n这是一个获取用户信息的接口",
					ReqHeaders: []YapiHeader{
						{
							Name:        "Authorization",
							Value:       "Bearer token",
							Description: "认证令牌",
							Required:    "1",
						},
					},
					ReqQuery: []YapiParam{
						{
							Name:        "user_id",
							Type:        "number",
							Description: "用户ID",
							Required:    "1",
							Example:     "123",
						},
					},
					ReqBodyForm:  []YapiParam{},
					ReqBodyOther: "",
					ResBody:      `{"code": 0, "message": "success", "data": {"id": 123, "name": "测试用户"}}`,
					ResBodyType:  "json",
					Tag:          []string{"用户", "基础"},
					ProjectID:    123,
					AddTime:      1640995200, // 2022-01-01 00:00:00
					UpTime:       1640995200,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// 模拟项目信息API
	mux.HandleFunc("/api/project/get", func(w http.ResponseWriter, r *http.Request) {
		projectID := r.URL.Query().Get("id")
		token := r.URL.Query().Get("token")

		if token != "test_token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if projectID == "123" {
			response := YapiProjectResponse{
				ErrCode: 0,
				ErrMsg:  "成功",
				Data: YapiProject{
					ID:          123,
					Name:        "测试项目",
					Description: "这是一个测试项目",
					BasePath:    "/api",
					GroupID:     1,
					GroupName:   "测试分组",
					Color:       "blue",
					Icon:        "user",
					AddTime:     1640995200,
					UpTime:      1640995200,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	return httptest.NewServer(mux)
}

func TestYapiClient_GetInterfaces(t *testing.T) {
	// 启动模拟服务器
	mockServer := MockYapiServer()
	defer mockServer.Close()

	// 创建客户端
	client := NewYapiClient(mockServer.URL, "test_token")

	// 测试成功情况
	t.Run("成功获取接口列表", func(t *testing.T) {
		response, err := client.GetInterfaces(123)
		require.NoError(t, err)
		assert.Equal(t, 0, response.ErrCode)
		assert.Equal(t, 2, response.Data.Count)
		assert.Equal(t, 2, response.Data.Total)
		assert.Len(t, response.Data.List, 2)

		// 验证第一个接口
		firstInterface := response.Data.List[0]
		assert.Equal(t, 1001, firstInterface.ID)
		assert.Equal(t, "获取用户信息", firstInterface.Title)
		assert.Equal(t, "/api/user/info", firstInterface.Path)
		assert.Equal(t, "GET", firstInterface.Method)
		assert.Equal(t, "done", firstInterface.Status)
	})

	// 测试失败情况
	t.Run("项目不存在", func(t *testing.T) {
		_, err := client.GetInterfaces(999)
		assert.Error(t, err)
	})
}

func TestYapiClient_GetInterfaceDetail(t *testing.T) {
	// 启动模拟服务器
	mockServer := MockYapiServer()
	defer mockServer.Close()

	// 创建客户端
	client := NewYapiClient(mockServer.URL, "test_token")

	// 测试成功情况
	t.Run("成功获取接口详情", func(t *testing.T) {
		response, err := client.GetInterfaceDetail(1001)
		require.NoError(t, err)
		assert.Equal(t, 0, response.ErrCode)

		detail := response.Data
		assert.Equal(t, 1001, detail.ID)
		assert.Equal(t, "获取用户信息", detail.Title)
		assert.Equal(t, "/api/user/info", detail.Path)
		assert.Equal(t, "GET", detail.Method)
		assert.Equal(t, "根据用户ID获取用户详细信息", detail.Description)
		assert.Len(t, detail.ReqHeaders, 1)
		assert.Len(t, detail.ReqQuery, 1)

		// 验证请求头
		assert.Equal(t, "Authorization", detail.ReqHeaders[0].Name)
		assert.Equal(t, "Bearer token", detail.ReqHeaders[0].Value)

		// 验证查询参数
		assert.Equal(t, "user_id", detail.ReqQuery[0].Name)
		assert.Equal(t, "number", detail.ReqQuery[0].Type)
	})

	// 测试失败情况
	t.Run("接口不存在", func(t *testing.T) {
		_, err := client.GetInterfaceDetail(9999)
		assert.Error(t, err)
	})
}

func TestYapiClient_GetProjectInfo(t *testing.T) {
	// 启动模拟服务器
	mockServer := MockYapiServer()
	defer mockServer.Close()

	// 创建客户端
	client := NewYapiClient(mockServer.URL, "test_token")

	// 测试成功情况
	t.Run("成功获取项目信息", func(t *testing.T) {
		response, err := client.GetProjectInfo(123)
		require.NoError(t, err)
		assert.Equal(t, 0, response.ErrCode)

		project := response.Data
		assert.Equal(t, 123, project.ID)
		assert.Equal(t, "测试项目", project.Name)
		assert.Equal(t, "这是一个测试项目", project.Description)
		assert.Equal(t, "/api", project.BasePath)
		assert.Equal(t, 1, project.GroupID)
		assert.Equal(t, "测试分组", project.GroupName)
	})

	// 测试失败情况
	t.Run("项目不存在", func(t *testing.T) {
		_, err := client.GetProjectInfo(999)
		assert.Error(t, err)
	})
}

func TestMCPIntegration(t *testing.T) {
	// 启动模拟服务器
	mockServer := MockYapiServer()
	defer mockServer.Close()

	// 设置环境变量
	os.Setenv("YAPI_BASE_URL", mockServer.URL)
	os.Setenv("YAPI_TOKEN", "test_token")
	defer func() {
		os.Unsetenv("YAPI_BASE_URL")
		os.Unsetenv("YAPI_TOKEN")
	}()

	// 创建MCP服务器
	mcpServer := server.NewMCPServer("Test YAPI Server", "1.0.0")

	// 注册YAPI工具
	RegisterTool(mcpServer)

	t.Run("验证工具注册成功", func(t *testing.T) {
		// 验证YAPI客户端已初始化
		assert.NotNil(t, yapiClient, "YAPI客户端应该已初始化")
	})
}

func TestFormatJSONResponse(t *testing.T) {
	t.Run("格式化普通对象", func(t *testing.T) {
		data := map[string]any{
			"name": "测试",
			"id":   123,
		}

		result := FormatJSONResponse(data)
		assert.Contains(t, result, "测试")
		assert.Contains(t, result, "123")
		// 验证JSON格式正确
		var parsed map[string]any
		err := json.Unmarshal([]byte(result), &parsed)
		assert.NoError(t, err)
	})

	t.Run("格式化空对象", func(t *testing.T) {
		data := map[string]any{}

		result := FormatJSONResponse(data)
		assert.Equal(t, "{}", result)
	})
}

// BenchmarkYapiClient 性能测试
func BenchmarkYapiClient_GetInterfaces(b *testing.B) {
	mockServer := MockYapiServer()
	defer mockServer.Close()

	client := NewYapiClient(mockServer.URL, "test_token")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetInterfaces(123)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkYapiClient_GetInterfaceDetail(b *testing.B) {
	mockServer := MockYapiServer()
	defer mockServer.Close()

	client := NewYapiClient(mockServer.URL, "test_token")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetInterfaceDetail(1001)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ExampleYapiClient_GetInterfaces 使用示例
func ExampleYapiClient_GetInterfaces() {
	// 创建YAPI客户端
	client := NewYapiClient("http://yapi.example.com", "your_token_here")

	// 获取项目123的接口列表
	response, err := client.GetInterfaces(123)
	if err != nil {
		panic(err)
	}

	// 打印接口信息
	for _, iface := range response.Data.List {
		println("接口:", iface.Title, "路径:", iface.Path, "方法:", iface.Method)
	}
}

// ExampleYapiClient_GetInterfaceDetail 使用示例
func ExampleYapiClient_GetInterfaceDetail() {
	// 创建YAPI客户端
	client := NewYapiClient("http://yapi.example.com", "your_token_here")

	// 获取接口详情
	response, err := client.GetInterfaceDetail(1001)
	if err != nil {
		panic(err)
	}

	// 打印接口详情
	detail := response.Data
	println("接口标题:", detail.Title)
	println("接口描述:", detail.Description)
	println("请求方法:", detail.Method)
	println("请求路径:", detail.Path)
}
