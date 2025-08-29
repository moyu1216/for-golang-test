package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

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

// 计数器结构体
type Counter struct {
	mu    sync.Mutex
	value int
}

// 全局计数器实例
var counter = &Counter{}

// 增加计数器
func (c *Counter) Increment() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

// 获取当前计数
func (c *Counter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// 重置计数器
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
}

// 计数器端点 - 增加计数
func counterIncrement(w http.ResponseWriter, req *http.Request) {
	count := counter.Increment()
	fmt.Fprintf(w, "计数器已增加，当前值: %d\n", count)
}

// 计数器端点 - 获取当前计数
func counterGet(w http.ResponseWriter, req *http.Request) {
	count := counter.Get()
	fmt.Fprintf(w, "当前计数器值: %d\n", count)
}

// 计数器端点 - 重置计数器
func counterReset(w http.ResponseWriter, req *http.Request) {
	counter.Reset()
	fmt.Fprintf(w, "计数器已重置为 0\n")
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/v1/ping", ping)
	
	// 计数器相关的路由
	http.HandleFunc("/counter/increment", counterIncrement)
	http.HandleFunc("/counter/get", counterGet)
	http.HandleFunc("/counter/reset", counterReset)

	fmt.Println("服务器启动在端口 8000")
	fmt.Println("计数器API端点:")
	fmt.Println("  GET /counter/increment - 增加计数器")
	fmt.Println("  GET /counter/get - 获取当前计数")
	fmt.Println("  GET /counter/reset - 重置计数器")
	
	http.ListenAndServe(":8000", nil)
}
