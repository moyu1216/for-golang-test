package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Counter 计数器结构体
type Counter struct {
	Value int64     `json:"value"`
	mutex sync.RWMutex
}

// APIResponse 标准API响应格式
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// 全局计数器实例
var counter = &Counter{Value: 0}

// 计数器方法实现
func (c *Counter) Get() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Value
}

func (c *Counter) Increment() int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Value++
	return c.Value
}

func (c *Counter) Decrement() int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Value--
	return c.Value
}

func (c *Counter) Set(value int64) int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Value = value
	return c.Value
}

func (c *Counter) Reset() int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Value = 0
	return c.Value
}

// 辅助函数：发送JSON响应
func sendJSONResponse(w http.ResponseWriter, success bool, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := APIResponse{
		Success: success,
		Data:    data,
		Message: message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// 计数器API处理函数
func getCounter(w http.ResponseWriter, req *http.Request) {
	value := counter.Get()
	sendJSONResponse(w, true, map[string]int64{"value": value}, "获取计数器成功", http.StatusOK)
}

func incrementCounter(w http.ResponseWriter, req *http.Request) {
	value := counter.Increment()
	sendJSONResponse(w, true, map[string]int64{"value": value}, "计数器增加成功", http.StatusOK)
}

func decrementCounter(w http.ResponseWriter, req *http.Request) {
	value := counter.Decrement()
	sendJSONResponse(w, true, map[string]int64{"value": value}, "计数器减少成功", http.StatusOK)
}

func setCounter(w http.ResponseWriter, req *http.Request) {
	if req.Method != "PUT" {
		sendJSONResponse(w, false, nil, "方法不允许", http.StatusMethodNotAllowed)
		return
	}
	
	valueStr := req.URL.Query().Get("value")
	if valueStr == "" {
		sendJSONResponse(w, false, nil, "缺少value参数", http.StatusBadRequest)
		return
	}
	
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		sendJSONResponse(w, false, nil, "无效的数值", http.StatusBadRequest)
		return
	}
	
	newValue := counter.Set(value)
	sendJSONResponse(w, true, map[string]int64{"value": newValue}, "计数器设置成功", http.StatusOK)
}

func resetCounter(w http.ResponseWriter, req *http.Request) {
	value := counter.Reset()
	sendJSONResponse(w, true, map[string]int64{"value": value}, "计数器重置成功", http.StatusOK)
}

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	fmt.Fprintf(w, "NowTime: %s", time.Now().Format("2006-01-02 15:04:05"))
}

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "TZ:%s", os.Getenv("TZ"))
	fmt.Fprintf(w, "NowTime: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	// 原有端点
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/v1/ping", ping)
	
	// 计数器API端点
	http.HandleFunc("/counter", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			getCounter(w, req)
		case "PUT":
			setCounter(w, req)
		case "DELETE":
			resetCounter(w, req)
		default:
			sendJSONResponse(w, false, nil, "方法不允许", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/counter/increment", incrementCounter)
	http.HandleFunc("/counter/decrement", decrementCounter)

	fmt.Println("🚀 服务器启动在端口 8000")
	fmt.Println("📊 计数器API:")
	fmt.Println("  GET    /counter          - 获取当前计数值")
	fmt.Println("  POST   /counter/increment - 增加计数")
	fmt.Println("  POST   /counter/decrement - 减少计数")
	fmt.Println("  PUT    /counter?value=N  - 设置特定值")
	fmt.Println("  DELETE /counter          - 重置计数器")
	
	http.ListenAndServe(":8000", nil)
}
