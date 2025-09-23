// 这是一个简单的Go HTTP服务器示例，用于演示基本的Web路由处理
package main

// 导入必要的标准库包
import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// hello函数：处理/hello路由，返回简单的hello消息
func hello(w http.ResponseWriter, req *http.Request) {
	// 向响应写入hello消息
	fmt.Fprintf(w, "hello\n")
}

// headers函数：处理/headers路由，打印请求头信息并显示当前时间
func headers(w http.ResponseWriter, req *http.Request) {
	// 遍历请求头并输出
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	// 输出当前服务器时间
	fmt.Fprintf(w, "NowTime: %s", time.Now().Format("2006-01-02 15:04:05"))
}

// ping函数：处理/v1/ping路由，用于健康检查，返回时区和当前时间
func ping(w http.ResponseWriter, req *http.Request) {
	// 输出环境变量TZ（时区）
	fmt.Fprintf(w, "TZ:%s", os.Getenv("TZ"))
	// 输出当前时间（带换行符）
	fmt.Fprintf(w, "NowTime: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

// main函数：程序入口，设置路由并启动HTTP服务器
func main() {
	// 输出Hello World消息
	fmt.Println("Hello, World!")

	// 注册路由处理函数
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/v1/ping", ping)

	// 启动服务器，监听8000端口
	http.ListenAndServe(":8000", nil)
}
