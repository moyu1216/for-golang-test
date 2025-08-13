package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Counter 计数器结构体
type Counter struct {
	Value     int       `json:"value"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	mutex     sync.RWMutex
}

// CounterManager 计数器管理器
type CounterManager struct {
	counters map[string]*Counter
	mutex    sync.RWMutex
}

// 全局计数器管理器
var counterManager = &CounterManager{
	counters: make(map[string]*Counter),
}

// NewCounter 创建新计数器
func NewCounter(name string) *Counter {
	now := time.Now()
	return &Counter{
		Value:     0,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Increment 递增计数器
func (c *Counter) Increment(step int) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Value += step
	c.UpdatedAt = time.Now()
	return c.Value
}

// Decrement 递减计数器
func (c *Counter) Decrement(step int) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Value -= step
	c.UpdatedAt = time.Now()
	return c.Value
}

// Reset 重置计数器
func (c *Counter) Reset() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Value = 0
	c.UpdatedAt = time.Now()
	return c.Value
}

// GetValue 获取当前值
func (c *Counter) GetValue() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Value
}

// SetValue 设置特定值
func (c *Counter) SetValue(value int) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Value = value
	c.UpdatedAt = time.Now()
	return c.Value
}

// GetOrCreateCounter 获取或创建计数器
func (cm *CounterManager) GetOrCreateCounter(name string) *Counter {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if counter, exists := cm.counters[name]; exists {
		return counter
	}
	
	counter := NewCounter(name)
	cm.counters[name] = counter
	return counter
}

// GetAllCounters 获取所有计数器
func (cm *CounterManager) GetAllCounters() map[string]*Counter {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	result := make(map[string]*Counter)
	for name, counter := range cm.counters {
		result[name] = counter
	}
	return result
}

// DeleteCounter 删除计数器
func (cm *CounterManager) DeleteCounter(name string) bool {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if _, exists := cm.counters[name]; exists {
		delete(cm.counters, name)
		return true
	}
	return false
}

// API 路由处理器

// getCounter 获取计数器当前值
func getCounter(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		name = "default"
	}
	
	counter := counterManager.GetOrCreateCounter(name)
	
	c.JSON(http.StatusOK, gin.H{
		"name":       counter.Name,
		"value":      counter.GetValue(),
		"created_at": counter.CreatedAt,
		"updated_at": counter.UpdatedAt,
	})
}

// incrementCounter 递增计数器
func incrementCounter(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		name = "default"
	}
	
	step := 1
	if stepStr := c.Query("step"); stepStr != "" {
		if s, err := strconv.Atoi(stepStr); err == nil {
			step = s
		}
	}
	
	counter := counterManager.GetOrCreateCounter(name)
	newValue := counter.Increment(step)
	
	c.JSON(http.StatusOK, gin.H{
		"name":       counter.Name,
		"value":      newValue,
		"step":       step,
		"operation":  "increment",
		"updated_at": counter.UpdatedAt,
	})
}

// decrementCounter 递减计数器
func decrementCounter(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		name = "default"
	}
	
	step := 1
	if stepStr := c.Query("step"); stepStr != "" {
		if s, err := strconv.Atoi(stepStr); err == nil {
			step = s
		}
	}
	
	counter := counterManager.GetOrCreateCounter(name)
	newValue := counter.Decrement(step)
	
	c.JSON(http.StatusOK, gin.H{
		"name":       counter.Name,
		"value":      newValue,
		"step":       step,
		"operation":  "decrement",
		"updated_at": counter.UpdatedAt,
	})
}

// resetCounter 重置计数器
func resetCounter(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		name = "default"
	}
	
	counter := counterManager.GetOrCreateCounter(name)
	newValue := counter.Reset()
	
	c.JSON(http.StatusOK, gin.H{
		"name":       counter.Name,
		"value":      newValue,
		"operation":  "reset",
		"updated_at": counter.UpdatedAt,
	})
}

// setCounter 设置计数器值
func setCounter(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		name = "default"
	}
	
	valueStr := c.Query("value")
	if valueStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "value parameter is required",
		})
		return
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid value parameter",
		})
		return
	}
	
	counter := counterManager.GetOrCreateCounter(name)
	newValue := counter.SetValue(value)
	
	c.JSON(http.StatusOK, gin.H{
		"name":       counter.Name,
		"value":      newValue,
		"operation":  "set",
		"updated_at": counter.UpdatedAt,
	})
}

// getAllCounters 获取所有计数器
func getAllCounters(c *gin.Context) {
	counters := counterManager.GetAllCounters()
	
	result := make(map[string]interface{})
	for name, counter := range counters {
		result[name] = gin.H{
			"name":       counter.Name,
			"value":      counter.GetValue(),
			"created_at": counter.CreatedAt,
			"updated_at": counter.UpdatedAt,
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"counters": result,
		"total":    len(counters),
	})
}

// deleteCounter 删除计数器
func deleteCounter(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "counter name is required",
		})
		return
	}
	
	success := counterManager.DeleteCounter(name)
	if success {
		c.JSON(http.StatusOK, gin.H{
			"message": "counter deleted successfully",
			"name":    name,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "counter not found",
			"name":  name,
		})
	}
}

// 健康检查
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Counter service is running",
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}

func main() {
	// 设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)
	
	r := gin.Default()
	
	// 添加CORS支持
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	// 健康检查路由
	r.GET("/ping", ping)
	r.GET("/v1/ping", ping)
	
	// 计数器API路由
	api := r.Group("/api/v1")
	{
		// 获取所有计数器
		api.GET("/counters", getAllCounters)
		
		// 计数器操作（使用默认计数器）
		api.GET("/counter", getCounter)
		api.POST("/counter/increment", incrementCounter)
		api.POST("/counter/decrement", decrementCounter)
		api.POST("/counter/reset", resetCounter)
		api.PUT("/counter/set", setCounter)
		
		// 具名计数器操作
		api.GET("/counter/:name", getCounter)
		api.POST("/counter/:name/increment", incrementCounter)
		api.POST("/counter/:name/decrement", decrementCounter)
		api.POST("/counter/:name/reset", resetCounter)
		api.PUT("/counter/:name/set", setCounter)
		api.DELETE("/counter/:name", deleteCounter)
	}
	
	// 启动服务器
	r.Run(":8000")
}
