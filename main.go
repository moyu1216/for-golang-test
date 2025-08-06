package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Counter 结构体，包含计数器值和互斥锁以确保并发安全
type Counter struct {
	value int
	mutex sync.RWMutex
}

// 全局计数器实例
var counter = &Counter{value: 0}

// 增加计数器
func (c *Counter) Increment() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value++
	return c.value
}

// 减少计数器
func (c *Counter) Decrement() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value--
	return c.value
}

// 获取当前值
func (c *Counter) GetValue() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.value
}

// 重置计数器
func (c *Counter) Reset() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value = 0
	return c.value
}

// 设置计数器为指定值
func (c *Counter) SetValue(value int) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value = value
	return c.value
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
	fmt.Fprintf(w, "TZ:%s\n", os.Getenv("TZ"))
	fmt.Fprintf(w, "NowTime: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

// 计数器相关的HTTP处理函数

// 获取当前计数器值
func getCounter(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	value := counter.GetValue()
	fmt.Fprintf(w, `{"value": %d, "message": "当前计数器值"}`, value)
}

// 增加计数器
func incrementCounter(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	value := counter.Increment()
	fmt.Fprintf(w, `{"value": %d, "message": "计数器已增加"}`, value)
}

// 减少计数器
func decrementCounter(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	value := counter.Decrement()
	fmt.Fprintf(w, `{"value": %d, "message": "计数器已减少"}`, value)
}

// 重置计数器
func resetCounter(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	value := counter.Reset()
	fmt.Fprintf(w, `{"value": %d, "message": "计数器已重置"}`, value)
}

// 设置计数器为指定值
func setCounter(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// 从URL参数获取值
	valueStr := req.URL.Query().Get("value")
	if valueStr == "" {
		http.Error(w, `{"error": "缺少value参数"}`, http.StatusBadRequest)
		return
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		http.Error(w, `{"error": "value参数必须是整数"}`, http.StatusBadRequest)
		return
	}
	
	newValue := counter.SetValue(value)
	fmt.Fprintf(w, `{"value": %d, "message": "计数器已设置为指定值"}`, newValue)
}

// 计数器主页，显示使用说明
func counterHome(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Go 计数器</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; }
        .counter-display { text-align: center; font-size: 48px; font-weight: bold; color: #007bff; margin: 30px 0; }
        .buttons { text-align: center; margin: 20px 0; }
        .btn { background: #007bff; color: white; border: none; padding: 10px 20px; margin: 5px; border-radius: 5px; cursor: pointer; font-size: 16px; }
        .btn:hover { background: #0056b3; }
        .btn.danger { background: #dc3545; }
        .btn.danger:hover { background: #c82333; }
        .api-info { margin-top: 30px; }
        .endpoint { background: #f8f9fa; padding: 15px; margin: 10px 0; border-left: 4px solid #007bff; }
        .method { font-weight: bold; color: #28a745; }
    </style>
    <script>
        function updateCounter() {
            fetch('/counter')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('counter-value').textContent = data.value;
                });
        }
        
        function increment() {
            fetch('/counter/increment', {method: 'POST'})
                .then(response => response.json())
                .then(data => {
                    document.getElementById('counter-value').textContent = data.value;
                });
        }
        
        function decrement() {
            fetch('/counter/decrement', {method: 'POST'})
                .then(response => response.json())
                .then(data => {
                    document.getElementById('counter-value').textContent = data.value;
                });
        }
        
        function reset() {
            fetch('/counter/reset', {method: 'POST'})
                .then(response => response.json())
                .then(data => {
                    document.getElementById('counter-value').textContent = data.value;
                });
        }
        
        function setValue() {
            const value = document.getElementById('set-value').value;
            fetch('/counter/set?value=' + value, {method: 'POST'})
                .then(response => response.json())
                .then(data => {
                    document.getElementById('counter-value').textContent = data.value;
                    document.getElementById('set-value').value = '';
                });
        }
        
        // 页面加载时获取当前值
        window.onload = function() {
            updateCounter();
        }
    </script>
</head>
<body>
    <div class="container">
        <h1>🔢 Go 计数器应用</h1>
        
        <div class="counter-display">
            当前值: <span id="counter-value">0</span>
        </div>
        
        <div class="buttons">
            <button class="btn" onclick="increment()">➕ 增加</button>
            <button class="btn" onclick="decrement()">➖ 减少</button>
            <button class="btn danger" onclick="reset()">🔄 重置</button>
        </div>
        
        <div class="buttons">
            <input type="number" id="set-value" placeholder="输入数值" style="padding: 8px; margin-right: 10px;">
            <button class="btn" onclick="setValue()">📝 设置值</button>
        </div>
        
        <div class="api-info">
            <h2>📡 API 端点</h2>
            
            <div class="endpoint">
                <span class="method">GET</span> <code>/counter</code><br>
                获取当前计数器值
            </div>
            
            <div class="endpoint">
                <span class="method">POST</span> <code>/counter/increment</code><br>
                增加计数器（+1）
            </div>
            
            <div class="endpoint">
                <span class="method">POST</span> <code>/counter/decrement</code><br>
                减少计数器（-1）
            </div>
            
            <div class="endpoint">
                <span class="method">POST</span> <code>/counter/reset</code><br>
                重置计数器为0
            </div>
            
            <div class="endpoint">
                <span class="method">POST</span> <code>/counter/set?value=数值</code><br>
                设置计数器为指定值
            </div>
        </div>
    </div>
</body>
</html>`
	fmt.Fprint(w, html)
}

func main() {
	// 原有的端点
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/v1/ping", ping)
	
	// 计数器端点
	http.HandleFunc("/", counterHome)                    // 计数器主页
	http.HandleFunc("/counter", getCounter)              // 获取当前值
	http.HandleFunc("/counter/increment", incrementCounter) // 增加
	http.HandleFunc("/counter/decrement", decrementCounter) // 减少
	http.HandleFunc("/counter/reset", resetCounter)      // 重置
	http.HandleFunc("/counter/set", setCounter)          // 设置值
	
	fmt.Println("🚀 服务器启动在端口 8000")
	fmt.Println("📱 访问 http://localhost:8000 查看计数器界面")
	fmt.Println("🔗 API端点:")
	fmt.Println("   GET  /counter - 获取当前值")
	fmt.Println("   POST /counter/increment - 增加计数器")
	fmt.Println("   POST /counter/decrement - 减少计数器")
	fmt.Println("   POST /counter/reset - 重置计数器")
	fmt.Println("   POST /counter/set?value=N - 设置计数器为N")
	
	http.ListenAndServe(":8000", nil)
}
