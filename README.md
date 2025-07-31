# 🌍 智能旅行规划器

基于 **Plan and Execute** 框架的Go语言旅行规划器，提供专业的旅行建议和规划服务。

## ✨ 特性

- 🎯 **Plan and Execute 框架**：先制定详细计划，再执行具体建议
- 🚀 **双模式输出**：支持流式和非流式两种响应模式
- 🎨 **现代化UI**：美观的Web界面，支持实时交互
- 🔒 **智能过滤**：仅回答与旅行相关的问题
- 🌐 **WebSocket支持**：实时流式响应
- 🧪 **完整测试**：包含单元测试和基准测试

## 🏗️ 架构设计

### Plan and Execute 框架

1. **Plan阶段**：分析用户需求，制定详细的旅行计划步骤
2. **Execute阶段**：基于计划提供具体的实施建议和推荐

### 核心组件

- `TravelPlanner`：主要规划器结构
- `Plan`：旅行计划数据结构
- `TravelRequest/Response`：API请求响应结构
- WebSocket处理器：流式响应支持

## 🚀 快速开始

### 环境要求

- Go 1.21+
- OpenAI API Key

### 安装依赖

```bash
go mod tidy
```

### 设置环境变量

```bash
export OPENAI_API_KEY="your-openai-api-key"
export PORT="8000"  # 可选，默认8000
```

### 运行服务

```bash
go run main.go
```

服务启动后，访问 `http://localhost:8000` 即可使用Web界面。

## 📖 API 文档

### 健康检查

```bash
GET /health
```

响应：
```json
{
  "status": "healthy",
  "service": "travel-planner",
  "time": "2024-01-01 12:00:00"
}
```

### 旅行规划请求

#### 非流式模式

```bash
POST /travel/plan
Content-Type: application/json

{
  "query": "我想去日本东京旅行5天，预算5000元",
  "streaming": false
}
```

响应：
```json
{
  "plan": {
    "steps": [
      "第一步：确定目的地和旅行时间",
      "第二步：制定预算计划",
      "第三步：预订交通和住宿",
      "第四步：规划具体行程",
      "第五步：准备必要物品和证件"
    ]
  },
  "response": "详细的旅行建议...",
  "completed": true
}
```

#### 流式模式（WebSocket）

```javascript
const ws = new WebSocket('ws://localhost:8000/travel/plan');
ws.send(JSON.stringify({
  "query": "我想去日本东京旅行5天，预算5000元",
  "streaming": true
}));

ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  // 处理流式响应
};
```

## 🧪 测试

### 运行测试

```bash
go test -v
```

### 运行基准测试

```bash
go test -bench=.
```

### 运行示例

```bash
go test -run=Example
```

## 🎨 前端界面

项目包含一个现代化的Web界面，支持：

- 📝 旅行需求输入
- 🔄 流式/非流式模式切换
- 📋 计划步骤展示
- 💬 实时响应显示
- 🎨 响应式设计

## 🔧 配置选项

### 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `OPENAI_API_KEY` | OpenAI API密钥 | 必需 |
| `PORT` | 服务端口 | 8000 |

### 旅行关键词

系统会自动识别以下关键词来判断问题是否与旅行相关：

- 中文：旅行、旅游、度假、景点、酒店、机票、行程、攻略等
- 英文：travel、trip、vacation、tour、hotel、flight、itinerary等

## 📁 项目结构

```
.
├── main.go                 # 主程序入口
├── test_travel_planner.go  # 测试文件
├── index.html             # Web界面
├── go.mod                 # Go模块文件
├── DockerFile             # Docker配置
├── run.sh                 # 运行脚本
└── README.md              # 项目文档
```

## 🐳 Docker 部署

### 构建镜像

```bash
docker build -t travel-planner .
```

### 运行容器

```bash
docker run -p 8000:8000 -e OPENAI_API_KEY=your-key travel-planner
```

## 🔍 使用示例

### 示例1：国内旅行规划

**输入**：我想去云南大理旅行3天，预算2000元

**输出**：
- 计划：包含5个详细步骤
- 建议：具体的行程安排、住宿推荐、美食指南等

### 示例2：国际旅行规划

**输入**：我想去欧洲旅行2周，预算3万元

**输出**：
- 计划：包含签证、交通、住宿、行程等步骤
- 建议：详细的欧洲旅行攻略

## 🤝 贡献

欢迎提交Issue和Pull Request来改进项目！

## 📄 许可证

MIT License

## 🔗 相关链接

- [OpenAI API](https://platform.openai.com/)
- [Gin Framework](https://gin-gonic.com/)
- [WebSocket](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)