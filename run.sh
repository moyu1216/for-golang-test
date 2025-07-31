#!/bin/bash

# 旅行规划器启动脚本

echo "🌍 启动智能旅行规划器..."

# 检查环境变量
if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ 错误: 请设置 OPENAI_API_KEY 环境变量"
    echo "   例如: export OPENAI_API_KEY='your-api-key'"
    exit 1
fi

# 设置默认端口
export PORT=${PORT:-8000}

echo "✅ 环境检查通过"
echo "🔑 OpenAI API Key: ${OPENAI_API_KEY:0:10}..."
echo "🌐 服务端口: $PORT"

# 下载依赖
echo "📦 下载依赖..."
go mod tidy

# 运行应用
echo "🚀 启动服务..."
go run main.go
