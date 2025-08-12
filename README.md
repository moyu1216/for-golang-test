# 🎒 旅游智能助手 (Travel Agent)

一个基于Go语言开发的智能旅游助手，采用ReAct模式实现，支持多种旅游服务功能。

## ✨ 主要功能

### 🔧 核心能力
- **🌤️ 天气查询** - 获取目的地实时天气信息
- **🏛️ 景点推荐** - 推荐热门景点和特色路线  
- **🍜 美食推荐** - 推荐当地特色美食和餐厅
- **🏨 酒店搜索** - 查找不同价位的住宿选择
- **🗺️ 路线规划** - 制定个性化旅游行程

### 🎯 技术特色
- **ReAct模式**: 实现推理(Reasoning)和行动(Acting)的智能循环
- **工具调用**: 支持多种旅游工具的动态调用
- **多模式运行**: 支持CLI、Web服务和演示模式
- **模拟模式**: 无需API密钥即可体验完整功能

## 🏗️ 系统架构

```
旅游Agent架构
├── Agent框架 (agent.go)
│   ├── TravelAgent - 核心代理结构
│   ├── Tool接口 - 工具标准接口
│   └── ReAct处理循环
├── 工具集合 (tools.go)
│   ├── WeatherTool - 天气查询
│   ├── AttractionTool - 景点推荐
│   ├── HotelTool - 酒店搜索
│   ├── RouteTool - 路线规划
│   └── FoodTool - 美食推荐
└── 用户界面 (main.go)
    ├── CLI模式 - 命令行交互
    ├── Web模式 - 浏览器界面
    └── Demo模式 - 功能演示
```

## 🚀 快速开始

### 环境要求
- Go 1.21+
- （可选）OpenAI API Key用于真实LLM调用

### 安装依赖
```bash
go mod download
```

### 运行方式

#### 1. 命令行模式（推荐新手）
```bash
go run . 
```
支持交互式对话，输入`help`查看使用指南

#### 2. 演示模式（快速体验）
```bash
go run . demo
```
自动展示各项功能的示例对话

#### 3. Web服务模式（图形界面）
```bash
go run . web
```
访问 http://localhost:8000 使用Web界面

### 配置OpenAI API（可选）
```bash
export OPENAI_API_KEY="your-api-key-here"
go run . 
```

## 📖 使用示例

### CLI模式对话示例
```
👤 您: 我想去北京旅游，帮我规划一下

🤖 助手: 🗺️ 为您规划 '北京旅游' 的旅游路线：

📋 推荐行程：北京一日游

• morning：09:00 天安门广场 → 10:30 故宫博物院
• afternoon：14:00 北海公园 → 16:00 南锣鼓巷  
• evening：18:00 王府井大街 (晚餐购物)
🚌 交通建议：地铁+步行，建议购买一日交通卡

💡 出行小贴士：
   🚇 优先使用公共交通，环保又经济
   ⏰ 合理安排时间，避免行程过于紧张
   📍 相近景点可安排在同一天游览
```

### Web界面功能
- 🎨 现代化UI设计，支持移动端
- 💬 实时对话交互
- 🔄 一键清空对话历史
- 📱 示例问题快速输入

### API接口
```bash
# 发送消息
curl -X POST http://localhost:8000/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "推荐上海的景点"}'

# 清空对话
curl -X POST http://localhost:8000/api/clear
```

## 🔨 开发指南

### 添加新工具
1. 实现`Tool`接口：
```go
type MyTool struct{}

func (t *MyTool) Name() string {
    return "my_tool"
}

func (t *MyTool) Description() string {
    return "工具描述"
}

func (t *MyTool) Execute(args map[string]interface{}) (string, error) {
    // 工具逻辑实现
    return "结果", nil
}
```

2. 注册工具：
```go
agent.RegisterTool(&MyTool{})
```

### 自定义数据源
当前使用模拟数据，实际部署时可以：
- 集成真实天气API (如OpenWeatherMap)
- 连接旅游数据库
- 对接酒店预订系统
- 集成地图服务API

## 📚 技术细节

### ReAct模式实现
```go
func (agent *TravelAgent) Process(userInput string) (string, error) {
    // 1. 添加用户输入到上下文
    // 2. 调用LLM获取响应
    // 3. 检查是否需要工具调用
    // 4. 执行工具并添加结果到上下文
    // 5. 重复直到获得最终答案
}
```

### 工具调用机制
- 动态工具注册和发现
- 标准化参数传递
- 错误处理和重试机制
- 结果格式化和展示

### 多模态支持
- CLI: 适合开发调试
- Web: 适合最终用户
- API: 适合集成到其他系统

## 🔮 扩展可能

### 高级功能
- 🎯 个性化推荐算法
- 🗃️ 用户偏好记忆
- 🌐 多语言支持
- 📊 数据分析和洞察

### 企业级特性  
- 🔐 用户认证和授权
- 📈 使用统计和监控
- 🎚️ 负载均衡和缓存
- 📱 移动端APP

### 集成选项
- 💳 支付系统集成
- 📧 邮件通知服务
- 📱 短信提醒功能
- 🗂️ CRM系统对接

## 🤝 贡献指南

欢迎提交Issue和Pull Request来改进这个项目！

### 开发环境设置
```bash
git clone <repository-url>
cd for-golang-test
go mod download
go run . demo  # 验证安装
```

### 代码规范
- 遵循Go标准代码风格
- 添加适当的注释和文档
- 编写单元测试
- 提交前运行`go fmt`

## 📄 许可证

本项目采用MIT许可证 - 查看LICENSE文件了解详情

## 🙏 致谢

- OpenAI GPT模型提供智能对话能力
- Gin框架提供Web服务支持
- Go语言社区的工具和库支持

---

**开始您的智能旅游之旅吧！** 🌟