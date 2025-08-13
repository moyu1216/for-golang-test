# 计数器服务 (Counter Service)

一个基于Go和Gin框架的高性能计数器服务，支持多个命名计数器的并发操作。

## 功能特性

- ✅ 支持多个命名计数器
- ✅ 线程安全的并发操作
- ✅ RESTful API接口
- ✅ 递增/递减操作（支持自定义步长）
- ✅ 重置计数器
- ✅ 设置特定值
- ✅ 获取所有计数器状态
- ✅ 删除计数器
- ✅ CORS支持
- ✅ 健康检查接口

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8000` 启动

### 3. 健康检查

```bash
curl http://localhost:8000/ping
```

## API 接口

### 基础路由

- `GET /ping` - 健康检查
- `GET /v1/ping` - 健康检查（版本化）

### 计数器API (`/api/v1`)

#### 1. 获取所有计数器
```bash
GET /api/v1/counters
```

响应示例：
```json
{
  "counters": {
    "default": {
      "name": "default",
      "value": 5,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:05:00Z"
    },
    "page_views": {
      "name": "page_views", 
      "value": 100,
      "created_at": "2024-01-01T09:00:00Z",
      "updated_at": "2024-01-01T10:03:00Z"
    }
  },
  "total": 2
}
```

#### 2. 获取计数器当前值
```bash
# 获取默认计数器
GET /api/v1/counter

# 获取指定名称的计数器
GET /api/v1/counter/{name}
```

响应示例：
```json
{
  "name": "default",
  "value": 10,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:05:00Z"
}
```

#### 3. 递增计数器
```bash
# 递增默认计数器（步长=1）
POST /api/v1/counter/increment

# 递增指定计数器（步长=1）
POST /api/v1/counter/{name}/increment

# 指定步长递增
POST /api/v1/counter/{name}/increment?step=5
```

响应示例：
```json
{
  "name": "page_views",
  "value": 15,
  "step": 5,
  "operation": "increment",
  "updated_at": "2024-01-01T10:05:00Z"
}
```

#### 4. 递减计数器
```bash
# 递减默认计数器（步长=1）
POST /api/v1/counter/decrement

# 递减指定计数器（步长=1）
POST /api/v1/counter/{name}/decrement

# 指定步长递减
POST /api/v1/counter/{name}/decrement?step=3
```

#### 5. 重置计数器
```bash
# 重置默认计数器为0
POST /api/v1/counter/reset

# 重置指定计数器为0
POST /api/v1/counter/{name}/reset
```

响应示例：
```json
{
  "name": "page_views",
  "value": 0,
  "operation": "reset",
  "updated_at": "2024-01-01T10:06:00Z"
}
```

#### 6. 设置计数器特定值
```bash
# 设置默认计数器值
PUT /api/v1/counter/set?value=100

# 设置指定计数器值
PUT /api/v1/counter/{name}/set?value=50
```

响应示例：
```json
{
  "name": "downloads",
  "value": 100,
  "operation": "set",
  "updated_at": "2024-01-01T10:07:00Z"
}
```

#### 7. 删除计数器
```bash
DELETE /api/v1/counter/{name}
```

响应示例：
```json
{
  "message": "counter deleted successfully",
  "name": "old_counter"
}
```

## 使用示例

### 创建和操作页面访问计数器

```bash
# 1. 递增页面访问量
curl -X POST "http://localhost:8000/api/v1/counter/page_views/increment"

# 2. 查看当前访问量
curl "http://localhost:8000/api/v1/counter/page_views"

# 3. 增加10次访问
curl -X POST "http://localhost:8000/api/v1/counter/page_views/increment?step=10"

# 4. 设置为特定值
curl -X PUT "http://localhost:8000/api/v1/counter/page_views/set?value=1000"
```

### 创建和操作下载计数器

```bash
# 1. 增加下载次数
curl -X POST "http://localhost:8000/api/v1/counter/downloads/increment?step=5"

# 2. 查看所有计数器
curl "http://localhost:8000/api/v1/counters"

# 3. 重置下载计数器
curl -X POST "http://localhost:8000/api/v1/counter/downloads/reset"
```

## 特性说明

### 线程安全
所有计数器操作都是线程安全的，支持高并发访问。

### 自动创建
当访问不存在的计数器时，系统会自动创建它。

### 内存存储
计数器数据存储在内存中，服务重启后数据会丢失。如需持久化，可以扩展添加数据库支持。

### CORS支持
服务已配置CORS支持，可以被前端应用直接调用。

## 部署

### Docker部署
项目包含Dockerfile，可以直接构建Docker镜像：

```bash
docker build -t counter-service .
docker run -p 8000:8000 counter-service
```

### 生产环境
建议配置：
- 使用反向代理（如Nginx）
- 添加速率限制
- 配置日志记录
- 考虑数据持久化

## 扩展建议

1. **数据持久化**: 添加Redis或数据库支持
2. **统计功能**: 添加计数器使用统计
3. **权限控制**: 添加API认证和授权
4. **监控告警**: 集成Prometheus等监控系统
5. **集群支持**: 支持多实例部署和数据同步

## 许可证

MIT License