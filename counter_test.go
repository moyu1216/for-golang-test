package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCounterOperations(t *testing.T) {
	// 测试计数器基本操作
	counter.Reset()
	
	// 测试获取
	if counter.Get() != 0 {
		t.Errorf("期望计数器初始值为0，实际为%d", counter.Get())
	}
	
	// 测试增加
	counter.Increment()
	if counter.Get() != 1 {
		t.Errorf("期望计数器增加后为1，实际为%d", counter.Get())
	}
	
	// 测试减少
	counter.Decrement()
	if counter.Get() != 0 {
		t.Errorf("期望计数器减少后为0，实际为%d", counter.Get())
	}
	
	// 测试设置
	counter.Set(42)
	if counter.Get() != 42 {
		t.Errorf("期望计数器设置为42，实际为%d", counter.Get())
	}
	
	// 测试重置
	counter.Reset()
	if counter.Get() != 0 {
		t.Errorf("期望计数器重置后为0，实际为%d", counter.Get())
	}
}

func TestCounterAPI(t *testing.T) {
	// 重置计数器
	counter.Reset()
	
	// 测试GET /counter
	req, _ := http.NewRequest("GET", "/counter", nil)
	rr := httptest.NewRecorder()
	getCounter(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("期望状态码200，实际为%d", rr.Code)
	}
	
	var response APIResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	if !response.Success {
		t.Errorf("期望成功响应，实际为失败")
	}
	
	// 测试POST /counter/increment
	req, _ = http.NewRequest("POST", "/counter/increment", nil)
	rr = httptest.NewRecorder()
	incrementCounter(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("期望状态码200，实际为%d", rr.Code)
	}
	
	// 测试PUT /counter?value=10
	req, _ = http.NewRequest("PUT", "/counter?value=10", nil)
	rr = httptest.NewRecorder()
	setCounter(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("期望状态码200，实际为%d", rr.Code)
	}
	
	// 测试DELETE /counter
	req, _ = http.NewRequest("DELETE", "/counter", nil)
	rr = httptest.NewRecorder()
	resetCounter(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("期望状态码200，实际为%d", rr.Code)
	}
}

func TestCounterConcurrency(t *testing.T) {
	counter.Reset()
	
	// 并发测试
	done := make(chan bool)
	
	// 启动多个goroutine同时操作计数器
	for i := 0; i < 100; i++ {
		go func() {
			counter.Increment()
			done <- true
		}()
	}
	
	// 等待所有goroutine完成
	for i := 0; i < 100; i++ {
		<-done
	}
	
	// 验证最终值
	if counter.Get() != 100 {
		t.Errorf("期望并发操作后计数器为100，实际为%d", counter.Get())
	}
}