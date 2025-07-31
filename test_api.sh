#!/bin/bash

# API测试脚本

BASE_URL="http://localhost:8000"

echo "🧪 开始API测试..."

# 测试健康检查
echo "1. 测试健康检查..."
curl -s "$BASE_URL/health" | jq '.' || echo "健康检查失败"

echo -e "\n2. 测试旅行规划API（非流式）..."
curl -s -X POST "$BASE_URL/travel/plan" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "我想去日本东京旅行5天，预算5000元",
    "streaming": false
  }' | jq '.' || echo "旅行规划API测试失败"

echo -e "\n3. 测试非旅行相关问题..."
curl -s -X POST "$BASE_URL/travel/plan" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "如何做菜",
    "streaming": false
  }' | jq '.' || echo "非旅行问题测试失败"

echo -e "\n✅ API测试完成"