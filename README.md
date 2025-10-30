# 任务管理器 (Task Manager)

一个基于 Go 和 SQLite 的简单任务管理系统，提供 RESTful API 和 Web 界面。

## 功能特性

- ✅ 任务的 CRUD 操作（创建、读取、更新、删除）
- ✅ 任务搜索功能（按标题和描述搜索）
- ✅ 任务过滤（按状态、优先级）
- ✅ 任务排序（按创建时间、截止日期、优先级）
- ✅ 批量操作（批量删除、批量更新状态）
- ✅ 任务统计信息
- ✅ 现代化 Web 界面
- ✅ CORS 支持，支持跨域请求

## 技术栈

- **后端**: Go 1.16+
- **路由**: gorilla/mux
- **数据库**: SQLite3
- **前端**: HTML5 + CSS3 + JavaScript (原生)

## 快速开始

### 前置要求

- Go 1.16 或更高版本
- SQLite3

### 安装依赖

```bash
go mod tidy
```

### 运行服务器

```bash
go run main.go
```

或者使用提供的脚本：

```bash
./run.sh
```

服务器将在 `http://localhost:8000` 启动。

### 访问 Web 界面

打开浏览器访问：`http://localhost:8000`

## API 文档

### 基础 URL

```
http://localhost:8000/api/tasks
```

### 端点列表

#### 1. 获取任务列表

```
GET /api/tasks
```

**查询参数**（可选）:
- `status`: 过滤状态 (pending/in_progress/completed)
- `priority`: 过滤优先级 (high/medium/low)
- `sort_by`: 排序方式 (created_at/due_date/priority)

**示例**:
```bash
curl http://localhost:8000/api/tasks?status=pending&sort_by=created_at
```

**响应**:
```json
[
  {
    "id": 1,
    "title": "完成项目文档",
    "description": "编写 API 文档",
    "status": "pending",
    "priority": "high",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z",
    "due_date": "2024-01-15T18:00:00Z"
  }
]
```

#### 2. 创建任务

```
POST /api/tasks
```

**请求体**:
```json
{
  "title": "任务标题",
  "description": "任务描述",
  "status": "pending",
  "priority": "medium",
  "due_date": "2024-01-15T18:00:00Z"
}
```

**必填字段**: `title`

**可选字段**: `description`, `status`, `priority`, `due_date`

**默认值**:
- `status`: "pending"
- `priority`: "medium"

**响应**: 返回创建的任务对象（状态码 201）

#### 3. 获取单个任务

```
GET /api/tasks/{id}
```

**示例**:
```bash
curl http://localhost:8000/api/tasks/1
```

**响应**: 返回任务对象（状态码 200）

#### 4. 更新任务

```
PUT /api/tasks/{id}
```

**请求体**: 同创建任务

**响应**: 返回更新后的任务对象（状态码 200）

#### 5. 删除任务

```
DELETE /api/tasks/{id}
```

**响应**: 无内容（状态码 204）

#### 6. 搜索任务

```
GET /api/tasks/search?q={查询关键词}
```

**示例**:
```bash
curl http://localhost:8000/api/tasks/search?q=文档
```

**响应**: 返回匹配的任务列表（状态码 200）

#### 7. 获取统计信息

```
GET /api/tasks/stats
```

**响应**:
```json
{
  "pending": 5,
  "in_progress": 3,
  "completed": 10
}
```

#### 8. 批量删除

```
POST /api/tasks/batch/delete
```

**请求体**:
```json
{
  "ids": [1, 2, 3]
}
```

**响应**: 无内容（状态码 204）

#### 9. 批量更新状态

```
POST /api/tasks/batch/status
```

**请求体**:
```json
{
  "ids": [1, 2, 3],
  "status": "completed"
}
```

**响应**: 无内容（状态码 204）

## 数据结构

### Task

```json
{
  "id": 1,
  "title": "任务标题",
  "description": "任务描述",
  "status": "pending",
  "priority": "medium",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z",
  "due_date": "2024-01-15T18:00:00Z"
}
```

### 字段说明

- `id`: 任务 ID（自动生成）
- `title`: 任务标题（必填）
- `description`: 任务描述（可选）
- `status`: 任务状态，可选值：`pending`, `in_progress`, `completed`
- `priority`: 优先级，可选值：`low`, `medium`, `high`
- `created_at`: 创建时间（自动生成）
- `updated_at`: 更新时间（自动更新）
- `due_date`: 截止日期（可选，ISO 8601 格式）

## 错误处理

所有错误响应都遵循统一格式：

```json
{
  "error": "错误消息"
}
```

常见 HTTP 状态码：
- `200`: 成功
- `201`: 创建成功
- `204`: 成功（无内容）
- `400`: 请求错误
- `404`: 资源未找到
- `500`: 服务器错误

## 项目结构

```
├── main.go          # 主程序，集成所有路由
├── models/
│   └── task.go      # 任务模型和数据访问层
├── handlers/
│   └── task.go      # 任务 HTTP 处理器
├── db/
│   └── db.go        # 数据库连接和初始化
├── static/          # 前端静态文件
│   ├── index.html
│   ├── css/
│   │   └── style.css
│   └── js/
│       └── app.js
├── go.mod           # Go 模块定义
└── README.md        # 项目文档
```

## 开发

### 构建

```bash
go build -o task-manager main.go
```

### 运行测试

```bash
go test ./...
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！