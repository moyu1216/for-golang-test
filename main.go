package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// 全局计数器变量
var (
	counter int
	mu      sync.Mutex
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, World!\n")
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

// 计数器增加端点
func counterIncrement(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	counter++
	currentCount := counter
	mu.Unlock()
	
	fmt.Fprintf(w, "计数器已增加，当前值: %d\n", currentCount)
}

// 获取当前计数端点
func counterGet(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	currentCount := counter
	mu.Unlock()
	
	fmt.Fprintf(w, "当前计数器值: %d\n", currentCount)
}

// 重置计数器端点
func counterReset(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	counter = 0
	mu.Unlock()
	
	fmt.Fprintf(w, "计数器已重置为: 0\n")
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/v1/ping", ping)
	
	// 计数器相关路由
	http.HandleFunc("/counter/increment", counterIncrement)
	http.HandleFunc("/counter/get", counterGet)
	http.HandleFunc("/counter/reset", counterReset)

	fmt.Println("服务器启动在端口 8000")
	fmt.Println("可用端点:")
	fmt.Println("  /hello - Hello World")
	fmt.Println("  /counter/increment - 增加计数器")
	fmt.Println("  /counter/get - 获取当前计数")
	fmt.Println("  /counter/reset - 重置计数器")
	
	http.ListenAndServe(":8000", nil)
}
