# Go 计数器服务

一个基于Go语言的高性能计数器HTTP服务，支持并发安全的计数操作。

## 🚀 功能特性

- ✅ **基础计数操作**: 增加、减少、设置、重置、获取
- 🔒 **并发安全**: 使用互斥锁保护共享状态
- 🌐 **RESTful API**: 标准化的HTTP接口
- 📊 **JSON响应**: 统一的响应格式
- ⚡ **高性能**: 内存存储，快速响应
- 🛡️ **错误处理**: 完善的输入验证和错误信息
- 🧪 **测试覆盖**: 单元测试和并发测试

## 📋 API 接口

### 获取当前计数值
```bash
GET /counter
```

**响应示例:**
```json
{
    "success": true,
    "data": {
        "value": 42
    },
    "message": "获取计数器成功"
}
```

### 增加计数
```bash
POST /counter/increment
```

### 减少计数
```bash
POST /counter/decrement
```

### 设置特定值
```bash
PUT /counter?value=100
```

### 重置计数器
```bash
DELETE /counter
```

## 🛠️ 运行方式

### 直接运行
```bash
go run main.go
```

### 编译运行
```bash
go build -o counter main.go
./counter
```

### 使用Docker
```bash
docker build -t counter-service .
docker run -p 8000:8000 counter-service
```

## 🧪 测试

运行单元测试：
```bash
go test -v
```

运行并发测试：
```bash
go test -race -v
```

## 📝 使用示例

### 使用curl测试API

```bash
# 获取当前值
curl http://localhost:8000/counter

# 增加计数
curl -X POST http://localhost:8000/counter/increment

# 减少计数
curl -X POST http://localhost:8000/counter/decrement

# 设置值为50
curl -X PUT "http://localhost:8000/counter?value=50"

# 重置计数器
curl -X DELETE http://localhost:8000/counter
```

### 使用JavaScript测试

```javascript
// 获取计数器值
fetch('http://localhost:8000/counter')
  .then(response => response.json())
  .then(data => console.log('当前值:', data.data.value));

// 增加计数
fetch('http://localhost:8000/counter/increment', {method: 'POST'})
  .then(response => response.json())
  .then(data => console.log('增加后:', data.data.value));
```

## 🔧 技术实现

- **语言**: Go 1.16+
- **并发控制**: sync.RWMutex
- **HTTP服务器**: net/http
- **数据格式**: JSON
- **存储方式**: 内存存储

## 📊 性能特点

- 支持高并发访问
- 内存存储，响应速度快
- 线程安全的计数器操作
- 标准化的错误处理

## 🎯 适用场景

- 网站访问计数器
- API调用次数统计
- 实时数据统计
- 微服务中的计数器组件
- 学习和测试Go并发编程