package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 检查环境变量
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("警告：OPENAI_API_KEY 环境变量未设置，将使用模拟模式")
		apiKey = "demo-key"
	}

	// 创建旅游Agent
	agent := NewTravelAgent(apiKey)

	// 注册工具
	agent.RegisterTool(&WeatherTool{})
	agent.RegisterTool(&AttractionTool{})
	agent.RegisterTool(&HotelTool{})
	agent.RegisterTool(&RouteTool{})
	agent.RegisterTool(&FoodTool{})

	fmt.Println("🎒 旅游智能助手启动成功！")
	fmt.Println("📋 可用功能：")
	fmt.Println("   • 天气查询 - 获取目的地天气信息")
	fmt.Println("   • 景点推荐 - 推荐热门景点和路线")
	fmt.Println("   • 酒店搜索 - 查找合适的住宿")
	fmt.Println("   • 路线规划 - 制定旅游行程")
	fmt.Println("   • 美食推荐 - 推荐当地特色美食")
	fmt.Println()

	// 检查启动模式
	if len(os.Args) > 1 && os.Args[1] == "web" {
		startWebServer(agent)
	} else if len(os.Args) > 1 && os.Args[1] == "demo" {
		runDemo(agent)
	} else {
		startCLI(agent)
	}
}

// 启动Web服务
func startWebServer(agent *TravelAgent) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 静态文件和模板
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// 主页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "旅游智能助手",
		})
	})

	// API接口
	r.POST("/api/chat", func(c *gin.Context) {
		var request struct {
			Message string `json:"message"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if request.Message == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Message cannot be empty"})
			return
		}

		// 处理消息
		response, err := agent.Process(request.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": response,
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// 清空对话
	r.POST("/api/clear", func(c *gin.Context) {
		agent.ClearContext()
		c.JSON(http.StatusOK, gin.H{"message": "对话已清空"})
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// 兼容原有的接口
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello from travel agent\n")
	})

	r.GET("/v1/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "TZ:%s\nNowTime: %s\n", 
			os.Getenv("TZ"), 
			time.Now().Format("2006-01-02 15:04:05"))
	})

	fmt.Println("🌐 Web服务启动在 http://localhost:8000")
	fmt.Println("📖 访问 http://localhost:8000 使用Web界面")
	fmt.Println("🔗 API接口: POST /api/chat")
	log.Fatal(r.Run(":8000"))
}

// 启动CLI模式
func startCLI(agent *TravelAgent) {
	fmt.Println("💬 命令行模式启动（输入 'exit' 退出，'clear' 清空对话）")
	fmt.Println(strings.Repeat("=", 50))

	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("👤 您: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("👋 再见！祝您旅途愉快！")
			break
		}

		if input == "clear" {
			agent.ClearContext()
			fmt.Println("✅ 对话历史已清空")
			continue
		}

		if input == "help" {
			printHelp()
			continue
		}

		fmt.Print("🤖 助手: ")
		
		// 如果没有OpenAI API Key，使用模拟模式
		if os.Getenv("OPENAI_API_KEY") == "" {
			response := handleDemoQuery(input)
			fmt.Println(response)
		} else {
			response, err := agent.Process(input)
			if err != nil {
				fmt.Printf("❌ 错误: %v\n", err)
				continue
			}
			fmt.Println(response)
		}
		fmt.Println()
	}
}

// 运行演示模式
func runDemo(agent *TravelAgent) {
	fmt.Println("🎭 演示模式 - 展示旅游Agent功能")
	fmt.Println(strings.Repeat("=", 40))

	demoQueries := []string{
		"北京的天气怎么样？",
		"推荐上海的景点",
		"广州有什么好吃的？",
		"帮我规划北京一日游路线",
		"上海有什么好的酒店推荐？",
	}

	for i, query := range demoQueries {
		fmt.Printf("\n📝 示例 %d: %s\n", i+1, query)
		fmt.Println("🤖 助手回复:")
		
		response := handleDemoQuery(query)
		fmt.Println(response)
		
		if i < len(demoQueries)-1 {
			fmt.Println("\n" + strings.Repeat("-", 40))
			time.Sleep(2 * time.Second)
		}
	}
	
	fmt.Println("\n✨ 演示完成！可以使用以下命令启动：")
	fmt.Println("   ./main          - CLI模式")
	fmt.Println("   ./main web      - Web服务模式")
}

// 处理演示查询
func handleDemoQuery(query string) string {
	query = strings.ToLower(query)
	
	// 模拟各种工具的响应
	if strings.Contains(query, "天气") || strings.Contains(query, "weather") {
		tool := &WeatherTool{}
		response, _ := tool.Execute(map[string]interface{}{"query": query})
		return response
	}
	
	if strings.Contains(query, "景点") || strings.Contains(query, "推荐") && !strings.Contains(query, "美食") && !strings.Contains(query, "酒店") {
		tool := &AttractionTool{}
		response, _ := tool.Execute(map[string]interface{}{"query": query})
		return response
	}
	
	if strings.Contains(query, "美食") || strings.Contains(query, "吃") {
		tool := &FoodTool{}
		response, _ := tool.Execute(map[string]interface{}{"query": query})
		return response
	}
	
	if strings.Contains(query, "路线") || strings.Contains(query, "规划") || strings.Contains(query, "游") {
		tool := &RouteTool{}
		response, _ := tool.Execute(map[string]interface{}{"query": query})
		return response
	}
	
	if strings.Contains(query, "酒店") || strings.Contains(query, "住宿") {
		tool := &HotelTool{}
		response, _ := tool.Execute(map[string]interface{}{"query": query})
		return response
	}
	
	return "你好！我是您的旅游助手。我可以帮您:\n• 查询天气信息\n• 推荐景点和路线\n• 搜索酒店住宿\n• 规划旅游行程\n• 推荐当地美食\n\n请告诉我您想了解什么？"
}

// 打印帮助信息
func printHelp() {
	fmt.Println(`
🎯 旅游智能助手使用指南

📋 主要功能：
   • 天气查询 - "北京天气怎么样？"
   • 景点推荐 - "推荐上海的景点"
   • 美食推荐 - "广州有什么好吃的？"
   • 酒店搜索 - "上海有什么好酒店？"
   • 路线规划 - "帮我规划北京一日游"

💡 使用技巧：
   • 可以组合多个需求，如"去北京旅游，推荐景点和美食"
   • 支持多轮对话，会记住上下文
   • 输入城市名称获得更精准的建议

🔧 命令：
   • help  - 显示此帮助
   • clear - 清空对话历史
   • exit  - 退出程序

🌐 Web模式：
   运行 './main web' 启动Web界面服务
`)
}
