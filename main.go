package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ccc-me/for-golang-test/db"
	"github.com/Ccc-me/for-golang-test/handlers"
	"github.com/Ccc-me/for-golang-test/models"
	"github.com/gorilla/mux"
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

// CORS中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	defer db.CloseDB()

	// 创建任务仓库和处理器
	taskRepo := models.NewTaskRepository(db.DB)
	taskHandler := handlers.NewTaskHandler(taskRepo)

	// 创建路由器
	r := mux.NewRouter()

	// 应用CORS中间件
	r.Use(corsMiddleware)

	// 旧的路由
	r.HandleFunc("/hello", hello)
	r.HandleFunc("/headers", headers)
	r.HandleFunc("/v1/ping", ping)

	// API路由
	api := r.PathPrefix("/api/tasks").Subrouter()
	api.HandleFunc("", taskHandler.GetTasks).Methods("GET")
	api.HandleFunc("", taskHandler.CreateTask).Methods("POST")
	api.HandleFunc("/stats", taskHandler.GetStats).Methods("GET")
	api.HandleFunc("/search", taskHandler.SearchTasks).Methods("GET")
	api.HandleFunc("/batch/delete", taskHandler.BatchDelete).Methods("POST")
	api.HandleFunc("/batch/status", taskHandler.BatchUpdateStatus).Methods("POST")
	api.HandleFunc("/{id}", taskHandler.GetTask).Methods("GET")
	api.HandleFunc("/{id}", taskHandler.UpdateTask).Methods("PUT")
	api.HandleFunc("/{id}", taskHandler.DeleteTask).Methods("DELETE")

	// 静态文件服务
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Println("服务器启动在端口 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
