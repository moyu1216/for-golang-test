package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

// Tool 定义工具接口
type Tool interface {
	Name() string
	Description() string
	Execute(args map[string]interface{}) (string, error)
}

// TravelAgent 旅游智能代理结构
type TravelAgent struct {
	client    *openai.Client
	tools     map[string]Tool
	sessionID string
	context   []openai.ChatCompletionMessage
	maxSteps  int
}

// NewTravelAgent 创建新的旅游代理
func NewTravelAgent(apiKey string) *TravelAgent {
	return &TravelAgent{
		client:    openai.NewClient(apiKey),
		tools:     make(map[string]Tool),
		sessionID: uuid.New().String(),
		context:   make([]openai.ChatCompletionMessage, 0),
		maxSteps:  10,
	}
}

// RegisterTool 注册工具
func (agent *TravelAgent) RegisterTool(tool Tool) {
	agent.tools[tool.Name()] = tool
}

// GetToolDefinitions 获取工具定义供LLM使用
func (agent *TravelAgent) GetToolDefinitions() []openai.Tool {
	tools := make([]openai.Tool, 0)
	for _, tool := range agent.tools {
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: openai.FunctionDefinition{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "用户查询或操作参数",
						},
					},
					"required": []string{"query"},
				},
			},
		})
	}
	return tools
}

// Process 处理用户输入，实现ReAct模式
func (agent *TravelAgent) Process(userInput string) (string, error) {
	// 添加用户消息到上下文
	agent.context = append(agent.context, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userInput,
	})

	// 添加系统提示词
	if len(agent.context) == 1 {
		systemMessage := openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleSystem,
			Content: `你是一个专业的旅游规划助手。你可以帮助用户规划旅行、查询天气、推荐景点、搜索酒店等。

你有以下工具可以使用：
- weather_query: 查询天气信息
- attraction_recommend: 推荐景点
- hotel_search: 搜索酒店
- route_planning: 规划路线
- food_recommend: 推荐美食

请根据用户需求选择合适的工具，并提供专业的建议。回复应该友好、详细且实用。`,
		}
		// 将系统消息插入到开头
		agent.context = append([]openai.ChatCompletionMessage{systemMessage}, agent.context...)
	}

	steps := 0
	for steps < agent.maxSteps {
		steps++

		// 调用LLM获取响应
		resp, err := agent.client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: agent.context,
				Tools:    agent.GetToolDefinitions(),
			},
		)

		if err != nil {
			return "", fmt.Errorf("调用LLM失败: %w", err)
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("LLM没有返回响应")
		}

		choice := resp.Choices[0]
		message := choice.Message

		// 添加助手消息到上下文
		agent.context = append(agent.context, message)

		// 检查是否有工具调用
		if len(message.ToolCalls) > 0 {
			// 执行工具调用
			for _, toolCall := range message.ToolCalls {
				toolName := toolCall.Function.Name
				
				if tool, exists := agent.tools[toolName]; exists {
					// 解析参数
					var args map[string]interface{}
					if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
						log.Printf("解析工具参数失败: %v", err)
						continue
					}

					// 执行工具
					result, err := tool.Execute(args)
					if err != nil {
						result = fmt.Sprintf("工具执行失败: %v", err)
					}

					// 添加工具结果到上下文
					agent.context = append(agent.context, openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						Content:    result,
						ToolCallID: toolCall.ID,
					})
				}
			}
		} else {
			// 没有工具调用，返回最终回复
			return message.Content, nil
		}
	}

	return "抱歉，处理您的请求时达到了最大步数限制。请重新开始对话。", nil
}

// ClearContext 清空对话上下文
func (agent *TravelAgent) ClearContext() {
	agent.context = make([]openai.ChatCompletionMessage, 0)
}

// GetContext 获取当前对话上下文
func (agent *TravelAgent) GetContext() []openai.ChatCompletionMessage {
	return agent.context
}