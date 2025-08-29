package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// Hello World 处理函数
func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World! 你好世界!\n")
}

// 显示请求头信息
func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	fmt.Fprintf(w, "当前时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

// ping 健康检查
func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "时区: %s\n", os.Getenv("TZ"))
	fmt.Fprintf(w, "当前时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "状态: 正常运行\n")
}

func main() {
	fmt.Println("启动 Hello World 服务器...")
	fmt.Println("访问 http://localhost:8000/hello 查看 Hello World")
	fmt.Println("访问 http://localhost:8000/headers 查看请求头")
	fmt.Println("访问 http://localhost:8000/v1/ping 进行健康检查")

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/v1/ping", ping)

	fmt.Println("服务器正在端口 8000 上运行...")
	http.ListenAndServe(":8000", nil)
}
