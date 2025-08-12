#!/bin/bash

# 旅游智能助手启动脚本

echo "🎒 旅游智能助手 (Travel Agent)"
echo "============================="

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误：未找到Go环境，请先安装Go 1.21+"
    exit 1
fi

# 下载依赖
echo "📦 检查依赖..."
go mod tidy > /dev/null 2>&1

# 检查参数
MODE=${1:-cli}

case $MODE in
    "cli"|"")
        echo "💬 启动命令行模式..."
        go run .
        ;;
    "web")
        echo "🌐 启动Web服务模式..."
        echo "📖 请访问 http://localhost:8000"
        go run . web
        ;;
    "demo")
        echo "🎭 启动演示模式..."
        go run . demo
        ;;
    "build")
        echo "🔨 构建可执行文件..."
        go build -o travel-agent .
        echo "✅ 构建完成: travel-agent"
        echo "使用方式:"
        echo "  ./travel-agent        # CLI模式"
        echo "  ./travel-agent web    # Web模式"
        echo "  ./travel-agent demo   # 演示模式"
        ;;
    "help"|"-h"|"--help")
        echo "使用方法: $0 [模式]"
        echo ""
        echo "可用模式:"
        echo "  cli    - 命令行交互模式 (默认)"
        echo "  web    - Web浏览器界面模式"
        echo "  demo   - 功能演示模式"
        echo "  build  - 构建可执行文件"
        echo "  help   - 显示此帮助信息"
        echo ""
        echo "环境变量:"
        echo "  OPENAI_API_KEY - OpenAI API密钥 (可选，用于真实LLM调用)"
        echo ""
        echo "示例:"
        echo "  $0 demo              # 快速体验功能"
        echo "  $0 web               # 启动Web界面"
        echo "  OPENAI_API_KEY=sk-xxx $0 cli  # 使用真实API"
        ;;
    *)
        echo "❌ 错误：未知模式 '$MODE'"
        echo "运行 '$0 help' 查看使用说明"
        exit 1
        ;;
esac
