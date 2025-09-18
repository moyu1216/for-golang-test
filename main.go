package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// 全局计数器变量和互斥锁
var (
	counter int
	mutex   sync.RWMutex
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

// 计数器增加接口
func counterIncrement(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	
	mutex.Lock()
	counter++
	currentValue := counter
	mutex.Unlock()
	
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"message": "计数器增加成功", "counter": %d}`, currentValue)
}

// 计数器查看接口
func counterGet(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}
	
	mutex.RLock()
	currentValue := counter
	mutex.RUnlock()
	
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"counter": %d}`, currentValue)
}

// 计数器重置接口
func counterReset(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	
	mutex.Lock()
	counter = 0
	mutex.Unlock()
	
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"message": "计数器重置成功", "counter": 0}`)
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/v1/ping", ping)
	
	// 计数器相关路由
	http.HandleFunc("/counter", counterGet)           // GET 查看计数器
	http.HandleFunc("/counter/increment", counterIncrement) // POST 增加计数器
	http.HandleFunc("/counter/reset", counterReset)   // POST 重置计数器

	fmt.Println("服务器启动在端口 :8000")
	fmt.Println("计数器API:")
	fmt.Println("  GET  /counter            - 查看当前计数")
	fmt.Println("  POST /counter/increment  - 增加计数")
	fmt.Println("  POST /counter/reset      - 重置计数")
	http.ListenAndServe(":8000", nil)
}
