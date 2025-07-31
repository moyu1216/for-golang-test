package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sashabaranov/go-openai"
)

// TravelPlanner 旅行规划器结构
type TravelPlanner struct {
	client *openai.Client
}

// Plan 规划步骤
type Plan struct {
	Steps []string `json:"steps"`
}

// TravelRequest 旅行请求
type TravelRequest struct {
	Query     string `json:"query"`
	Streaming bool   `json:"streaming"`
}

// TravelResponse 旅行响应
type TravelResponse struct {
	Plan      *Plan  `json:"plan,omitempty"`
	Response  string `json:"response,omitempty"`
	Error     string `json:"error,omitempty"`
	Completed bool   `json:"completed"`
}

// 升级器配置
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewTravelPlanner 创建新的旅行规划器
func NewTravelPlanner(apiKey string) *TravelPlanner {
	client := openai.NewClient(apiKey)
	return &TravelPlanner{
		client: client,
	}
}

// isTravelRelated 检查问题是否与旅行相关
func (tp *TravelPlanner) isTravelRelated(query string) bool {
	travelKeywords := []string{
		"旅行", "旅游", "度假", "景点", "酒店", "机票", "行程", "攻略",
		"travel", "trip", "vacation", "tour", "hotel", "flight", "itinerary",
		"destination", "景点", "美食", "文化", "历史", "自然", "城市",
		"country", "city", "beach", "mountain", "museum", "restaurant",
	}
	
	queryLower := strings.ToLower(query)
	for _, keyword := range travelKeywords {
		if strings.Contains(queryLower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// createPlan 创建旅行计划
func (tp *TravelPlanner) createPlan(ctx context.Context, query string) (*Plan, error) {
	systemPrompt := `你是一个专业的旅行规划师。请根据用户的需求制定详细的旅行计划。
请按照以下格式返回JSON：
{
  "steps": [
    "第一步：确定目的地和旅行时间",
    "第二步：制定预算计划",
    "第三步：预订交通和住宿",
    "第四步：规划具体行程",
    "第五步：准备必要物品和证件"
  ]
}

请确保计划详细、实用且符合用户的具体需求。`

	resp, err := tp.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return nil, err
	}

	var plan Plan
	content := resp.Choices[0].Message.Content
	
	// 尝试解析JSON
	if err := json.Unmarshal([]byte(content), &plan); err != nil {
		// 如果不是JSON格式，创建默认计划
		plan = Plan{
			Steps: []string{
				"分析旅行需求",
				"制定初步计划",
				"完善细节",
				"提供建议",
			},
		}
	}

	return &plan, nil
}

// executePlan 执行旅行计划
func (tp *TravelPlanner) executePlan(ctx context.Context, query string, plan *Plan) (string, error) {
	stepsText := ""
	for i, step := range plan.Steps {
		stepsText += fmt.Sprintf("%d. %s\n", i+1, step)
	}

	systemPrompt := `你是一个专业的旅行规划师。基于已制定的计划，为用户提供详细的旅行建议和实施方案。
请提供具体、实用的建议，包括：
- 具体的推荐和选择
- 实用的提示和注意事项
- 预算和时间安排建议
- 可能遇到的问题和解决方案

请用中文回答，确保内容详细且易于理解。`

	resp, err := tp.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(`用户需求：%s

已制定的计划：
%s

请基于以上计划提供详细的实施建议。`, query, stepsText),
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// streamResponse 流式响应
func (tp *TravelPlanner) streamResponse(ctx context.Context, query string, plan *Plan, conn *websocket.Conn) error {
	stepsText := ""
	for i, step := range plan.Steps {
		stepsText += fmt.Sprintf("%d. %s\n", i+1, step)
	}

	systemPrompt := `你是一个专业的旅行规划师。基于已制定的计划，为用户提供详细的旅行建议和实施方案。
请提供具体、实用的建议，包括：
- 具体的推荐和选择
- 实用的提示和注意事项
- 预算和时间安排建议
- 可能遇到的问题和解决方案

请用中文回答，确保内容详细且易于理解。`

	stream, err := tp.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(`用户需求：%s

已制定的计划：
%s

请基于以上计划提供详细的实施建议。`, query, stepsText),
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err != nil {
			break
		}

		if len(response.Choices) > 0 {
			content := response.Choices[0].Delta.Content
			if content != "" {
				response := TravelResponse{
					Response:  content,
					Completed: false,
				}
				
				if err := conn.WriteJSON(response); err != nil {
					return err
				}
			}
		}
	}

	// 发送完成信号
	response := TravelResponse{
		Completed: true,
	}
	return conn.WriteJSON(response)
}

// handleTravelRequest 处理旅行请求
func (tp *TravelPlanner) handleTravelRequest(c *gin.Context) {
	// 检查是否是WebSocket升级请求
	if c.GetHeader("Upgrade") == "websocket" {
		tp.handleWebSocketRequest(c)
		return
	}

	var req TravelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, TravelResponse{
			Error: "无效的请求格式",
		})
		return
	}

	// 检查是否与旅行相关
	if !tp.isTravelRelated(req.Query) {
		c.JSON(http.StatusBadRequest, TravelResponse{
			Error: "抱歉，我只能回答与旅行相关的问题。请询问关于旅行、旅游、度假、景点、酒店、机票等问题。",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 第一步：制定计划
	plan, err := tp.createPlan(ctx, req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, TravelResponse{
			Error: "创建计划时出错: " + err.Error(),
		})
		return
	}

	// 非流式响应
	response, err := tp.executePlan(ctx, req.Query, plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, TravelResponse{
			Error: "执行计划时出错: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TravelResponse{
		Plan:      plan,
		Response:  response,
		Completed: true,
	})
}

// handleWebSocketRequest 处理WebSocket请求
func (tp *TravelPlanner) handleWebSocketRequest(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// 读取客户端消息
	var req TravelRequest
	if err := conn.ReadJSON(&req); err != nil {
		conn.WriteJSON(TravelResponse{
			Error: "读取请求失败: " + err.Error(),
		})
		return
	}

	// 检查是否与旅行相关
	if !tp.isTravelRelated(req.Query) {
		conn.WriteJSON(TravelResponse{
			Error: "抱歉，我只能回答与旅行相关的问题。请询问关于旅行、旅游、度假、景点、酒店、机票等问题。",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 第一步：制定计划
	plan, err := tp.createPlan(ctx, req.Query)
	if err != nil {
		conn.WriteJSON(TravelResponse{
			Error: "创建计划时出错: " + err.Error(),
		})
		return
	}

	// 发送计划
	planResponse := TravelResponse{
		Plan:      plan,
		Completed: false,
	}
	if err := conn.WriteJSON(planResponse); err != nil {
		return
	}

	// 流式执行计划
	if err := tp.streamResponse(ctx, req.Query, plan, conn); err != nil {
		conn.WriteJSON(TravelResponse{
			Error: "流式响应出错: " + err.Error(),
		})
	}
}

// healthCheck 健康检查
func (tp *TravelPlanner) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "travel-planner",
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置 OPENAI_API_KEY 环境变量")
	}

	planner := NewTravelPlanner(apiKey)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.Default()

	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// 静态文件服务
	router.StaticFile("/", "./index.html")
	router.StaticFile("/index.html", "./index.html")

	// 路由
	router.GET("/health", planner.healthCheck)
	router.POST("/travel/plan", planner.handleTravelRequest)
	router.GET("/travel/plan", planner.handleTravelRequest) // 支持WebSocket

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("旅行规划器服务启动在端口 %s", port)
	log.Printf("访问 http://localhost:%s 开始使用旅行规划器", port)
	log.Fatal(router.Run(":" + port))
}
