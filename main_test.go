package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTravelPlanner(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("OPENAI_API_KEY", "test-key")
	
	planner := NewTravelPlanner("test-key")

	// 测试旅行相关问题检测
	t.Run("测试旅行相关问题检测", func(t *testing.T) {
		travelQueries := []string{
			"我想去日本旅行",
			"推荐一些旅游景点",
			"如何预订酒店",
			"travel to Paris",
			"vacation planning",
		}

		for _, query := range travelQueries {
			if !planner.isTravelRelated(query) {
				t.Errorf("查询 '%s' 应该被识别为旅行相关问题", query)
			}
		}

		nonTravelQueries := []string{
			"如何做菜",
			"数学公式",
			"编程问题",
			"weather today",
		}

		for _, query := range nonTravelQueries {
			if planner.isTravelRelated(query) {
				t.Errorf("查询 '%s' 不应该被识别为旅行相关问题", query)
			}
		}
	})

	// 测试HTTP请求处理
	t.Run("测试HTTP请求处理", func(t *testing.T) {
		// 跳过需要真实API的测试
		t.Skip("跳过需要真实OpenAI API的测试")

		// 创建测试服务器
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/travel/plan", planner.handleTravelRequest)

		// 测试旅行相关请求
		travelRequest := TravelRequest{
			Query:     "我想去东京旅行5天",
			Streaming: false,
		}
		jsonData, _ := json.Marshal(travelRequest)

		req, _ := http.NewRequest("POST", "/travel/plan", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，得到 %d", http.StatusOK, w.Code)
		}

		// 测试非旅行相关请求
		nonTravelRequest := TravelRequest{
			Query:     "如何做菜",
			Streaming: false,
		}
		jsonData, _ = json.Marshal(nonTravelRequest)

		req, _ = http.NewRequest("POST", "/travel/plan", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("期望状态码 %d，得到 %d", http.StatusBadRequest, w.Code)
		}
	})

	// 测试计划创建
	t.Run("测试计划创建", func(t *testing.T) {
		plan := &Plan{
			Steps: []string{
				"第一步：确定目的地和旅行时间",
				"第二步：制定预算计划",
				"第三步：预订交通和住宿",
			},
		}

		if len(plan.Steps) == 0 {
			t.Error("计划应该包含步骤")
		}

		// 验证步骤格式
		for i, step := range plan.Steps {
			if !strings.Contains(step, "第") && !strings.Contains(step, "step") {
				t.Errorf("步骤 %d 应该包含步骤标识", i+1)
			}
		}
	})

	// 测试响应结构
	t.Run("测试响应结构", func(t *testing.T) {
		response := TravelResponse{
			Plan: &Plan{
				Steps: []string{"测试步骤"},
			},
			Response:  "测试响应",
			Completed: true,
		}

		if response.Plan == nil {
			t.Error("响应应该包含计划")
		}

		if response.Response == "" {
			t.Error("响应应该包含内容")
		}

		if !response.Completed {
			t.Error("响应应该标记为完成")
		}
	})
}

// 基准测试
func BenchmarkTravelRelatedCheck(b *testing.B) {
	planner := NewTravelPlanner("test-key")
	query := "我想去日本东京旅行5天，预算5000元"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		planner.isTravelRelated(query)
	}
}

// 示例函数
func ExampleTravelPlanner_isTravelRelated() {
	planner := NewTravelPlanner("test-key")
	
	// 旅行相关查询
	fmt.Println(planner.isTravelRelated("我想去日本旅行"))
	fmt.Println(planner.isTravelRelated("推荐旅游景点"))
	
	// 非旅行相关查询
	fmt.Println(planner.isTravelRelated("如何做菜"))
	fmt.Println(planner.isTravelRelated("数学问题"))
	
	// Output:
	// true
	// true
	// false
	// false
}