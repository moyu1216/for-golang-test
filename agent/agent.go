package agent

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Agent 智能体接口定义
type Agent interface {
	// 处理用户输入并返回响应
	Process(ctx context.Context, input string, sessionID string) (*Response, error)
	// 获取Agent能力描述
	GetCapabilities() []string
	// 重置会话状态
	ResetSession(sessionID string) error
}

// Response Agent响应结构
type Response struct {
	Type        ResponseType `json:"type"`        // 响应类型
	Content     string       `json:"content"`     // 响应内容
	Data        interface{}  `json:"data"`        // 附加数据
	Suggestions []string     `json:"suggestions"` // 建议操作
	SessionID   string       `json:"session_id"`  // 会话ID
	Timestamp   time.Time    `json:"timestamp"`   // 时间戳
}

// ResponseType 响应类型枚举
type ResponseType string

const (
	ResponseTypeText             ResponseType = "text"
	ResponseTypeRecommendation   ResponseType = "recommendation"
	ResponseTypeItinerary        ResponseType = "itinerary"
	ResponseTypeBooking          ResponseType = "booking"
	ResponseTypeWeather          ResponseType = "weather"
	ResponseTypeNavigation       ResponseType = "navigation"
	ResponseTypeError            ResponseType = "error"
)

// Session 会话管理
type Session struct {
	ID            string                 `json:"id"`
	UserPreferences map[string]interface{} `json:"user_preferences"`
	ConversationHistory []Message         `json:"conversation_history"`
	CurrentState   string                 `json:"current_state"`
	CreatedAt     time.Time              `json:"created_at"`
	LastUpdated   time.Time              `json:"last_updated"`
}

// Message 消息结构
type Message struct {
	Role      string    `json:"role"`      // user, assistant, system
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Intent 意图识别结果
type Intent struct {
	Name       string                 `json:"name"`        // 意图名称
	Confidence float64                `json:"confidence"`  // 置信度
	Entities   map[string]interface{} `json:"entities"`    // 实体提取结果
}

// BaseAgent 基础Agent实现
type BaseAgent struct {
	name          string
	capabilities  []string
	sessionStore  map[string]*Session
	intentHandler map[string]IntentHandler
}

// IntentHandler 意图处理器接口
type IntentHandler func(ctx context.Context, intent *Intent, session *Session) (*Response, error)

// NewBaseAgent 创建基础Agent
func NewBaseAgent(name string) *BaseAgent {
	return &BaseAgent{
		name:          name,
		capabilities:  []string{},
		sessionStore:  make(map[string]*Session),
		intentHandler: make(map[string]IntentHandler),
	}
}

// RegisterIntentHandler 注册意图处理器
func (a *BaseAgent) RegisterIntentHandler(intentName string, handler IntentHandler) {
	a.intentHandler[intentName] = handler
}

// GetCapabilities 获取Agent能力
func (a *BaseAgent) GetCapabilities() []string {
	return a.capabilities
}

// GetOrCreateSession 获取或创建会话
func (a *BaseAgent) GetOrCreateSession(sessionID string) *Session {
	session, exists := a.sessionStore[sessionID]
	if !exists {
		session = &Session{
			ID:                  sessionID,
			UserPreferences:     make(map[string]interface{}),
			ConversationHistory: []Message{},
			CurrentState:        "welcome",
			CreatedAt:          time.Now(),
			LastUpdated:        time.Now(),
		}
		a.sessionStore[sessionID] = session
	}
	return session
}

// ResetSession 重置会话
func (a *BaseAgent) ResetSession(sessionID string) error {
	delete(a.sessionStore, sessionID)
	log.Printf("Session %s reset", sessionID)
	return nil
}

// AddMessage 添加消息到会话历史
func (a *BaseAgent) AddMessage(session *Session, role, content string) {
	message := Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}
	session.ConversationHistory = append(session.ConversationHistory, message)
	session.LastUpdated = time.Now()
}

// AnalyzeIntent 意图识别（简化版本）
func (a *BaseAgent) AnalyzeIntent(input string) *Intent {
	// 这里实现简化的意图识别逻辑
	// 在实际项目中，这里会调用NLP模型或规则引擎
	
	intent := &Intent{
		Entities: make(map[string]interface{}),
	}
	
	// 简单的关键词匹配
	switch {
	case containsKeywords(input, []string{"推荐", "景点", "旅游", "去哪"}):
		intent.Name = "destination_recommendation"
		intent.Confidence = 0.8
	case containsKeywords(input, []string{"行程", "计划", "安排", "规划"}):
		intent.Name = "itinerary_planning"
		intent.Confidence = 0.8
	case containsKeywords(input, []string{"天气", "气温", "下雨"}):
		intent.Name = "weather_query"
		intent.Confidence = 0.9
	case containsKeywords(input, []string{"酒店", "住宿", "预订", "订房"}):
		intent.Name = "hotel_booking"
		intent.Confidence = 0.8
	case containsKeywords(input, []string{"路线", "导航", "怎么走", "交通"}):
		intent.Name = "navigation"
		intent.Confidence = 0.8
	case containsKeywords(input, []string{"帮助", "功能", "能做什么"}):
		intent.Name = "help"
		intent.Confidence = 0.9
	default:
		intent.Name = "general_chat"
		intent.Confidence = 0.5
	}
	
	return intent
}

// containsKeywords 检查输入是否包含关键词
func containsKeywords(input string, keywords []string) bool {
	for _, keyword := range keywords {
		if contains(input, keyword) {
			return true
		}
	}
	return false
}

// contains 简单的字符串包含检查
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr || 
		      indexOf(s, substr) >= 0)))
}

// indexOf 查找子字符串位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Process 处理用户输入
func (a *BaseAgent) Process(ctx context.Context, input string, sessionID string) (*Response, error) {
	session := a.GetOrCreateSession(sessionID)
	a.AddMessage(session, "user", input)
	
	// 意图识别
	intent := a.AnalyzeIntent(input)
	
	// 查找并执行意图处理器
	if handler, exists := a.intentHandler[intent.Name]; exists {
		response, err := handler(ctx, intent, session)
		if err != nil {
			return nil, fmt.Errorf("处理意图失败: %w", err)
		}
		
		a.AddMessage(session, "assistant", response.Content)
		response.SessionID = sessionID
		response.Timestamp = time.Now()
		return response, nil
	}
	
	// 默认响应
	response := &Response{
		Type:        ResponseTypeText,
		Content:     "抱歉，我不太理解您的意思。您可以问我关于旅游推荐、行程规划、天气查询等问题。",
		Suggestions: []string{"推荐一些景点", "帮我规划行程", "查询天气", "我需要帮助"},
		SessionID:   sessionID,
		Timestamp:   time.Now(),
	}
	
	a.AddMessage(session, "assistant", response.Content)
	return response, nil
}