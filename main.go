package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// TwoSumRequest 定义两数之和请求的结构体
type TwoSumRequest struct {
	Numbers []int `json:"numbers"`
	Target  int   `json:"target"`
}

// TwoSumResponse 定义两数之和响应的结构体
type TwoSumResponse struct {
	Indices []int `json:"indices"`
	Found   bool  `json:"found"`
	Message string `json:"message"`
}

// twoSum 实现两数之和算法
// 给定一个整数数组 nums 和一个整数目标值 target，在该数组中找出和为目标值的两个整数，并返回它们的下标
func twoSum(nums []int, target int) []int {
	// 使用哈希表存储数值和对应的索引
	numMap := make(map[int]int)
	
	for i, num := range nums {
		complement := target - num
		// 如果补数存在于哈希表中，则找到了答案
		if index, exists := numMap[complement]; exists {
			return []int{index, i}
		}
		// 将当前数字和索引存入哈希表
		numMap[num] = i
	}
	
	// 没有找到答案
	return []int{}
}

// twoSumHandler 处理两数之和的HTTP请求
func twoSumHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if req.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	
	var request TwoSumRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		response := TwoSumResponse{
			Found:   false,
			Message: "请求格式错误: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 验证输入
	if len(request.Numbers) < 2 {
		response := TwoSumResponse{
			Found:   false,
			Message: "至少需要两个数字",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// 执行两数之和算法
	indices := twoSum(request.Numbers, request.Target)
	
	if len(indices) == 0 {
		response := TwoSumResponse{
			Found:   false,
			Message: "没有找到两个数的和等于目标值",
		}
		json.NewEncoder(w).Encode(response)
	} else {
		response := TwoSumResponse{
			Indices: indices,
			Found:   true,
			Message: fmt.Sprintf("找到答案：数组[%d] + 数组[%d] = %d + %d = %d", 
				indices[0], indices[1], 
				request.Numbers[indices[0]], request.Numbers[indices[1]], 
				request.Target),
		}
		json.NewEncoder(w).Encode(response)
	}
}

// twoSumSimple 简单的两数之和测试端点，通过URL参数传递
func twoSumSimple(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// 解析URL参数
	query := req.URL.Query()
	numsStr := query.Get("nums")  // 例如: "2,7,11,15"
	targetStr := query.Get("target")
	
	if numsStr == "" || targetStr == "" {
		http.Error(w, "缺少参数：nums（逗号分隔的数字）和target", http.StatusBadRequest)
		return
	}
	
	// 解析目标值
	target, err := strconv.Atoi(targetStr)
	if err != nil {
		http.Error(w, "target参数必须是整数", http.StatusBadRequest)
		return
	}
	
	// 解析数字数组（简单实现，假设格式正确）
	numsStrArr := []string{}
	if numsStr != "" {
		for _, s := range []rune(numsStr) {
			if s != ',' {
				if len(numsStrArr) == 0 || numsStrArr[len(numsStrArr)-1] == "" {
					numsStrArr = append(numsStrArr, string(s))
				} else {
					numsStrArr[len(numsStrArr)-1] += string(s)
				}
			} else {
				numsStrArr = append(numsStrArr, "")
			}
		}
	}
	
	nums := make([]int, 0, len(numsStrArr))
	for _, numStr := range numsStrArr {
		if numStr != "" {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				http.Error(w, "数字格式错误: "+numStr, http.StatusBadRequest)
				return
			}
			nums = append(nums, num)
		}
	}
	
	if len(nums) < 2 {
		http.Error(w, "至少需要两个数字", http.StatusBadRequest)
		return
	}
	
	// 执行两数之和算法
	indices := twoSum(nums, target)
	
	if len(indices) == 0 {
		response := TwoSumResponse{
			Found:   false,
			Message: "没有找到两个数的和等于目标值",
		}
		json.NewEncoder(w).Encode(response)
	} else {
		response := TwoSumResponse{
			Indices: indices,
			Found:   true,
			Message: fmt.Sprintf("找到答案：nums[%d] + nums[%d] = %d + %d = %d", 
				indices[0], indices[1], 
				nums[indices[0]], nums[indices[1]], 
				target),
		}
		json.NewEncoder(w).Encode(response)
	}
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
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/v1/ping", ping)
	
	// 两数之和相关端点
	http.HandleFunc("/twosum", twoSumHandler)        // POST请求，JSON格式
	http.HandleFunc("/twosum/simple", twoSumSimple)  // GET请求，URL参数

	fmt.Println("服务器启动在端口 8000")
	fmt.Println("两数之和端点:")
	fmt.Println("  POST /twosum - JSON格式请求")
	fmt.Println("  GET /twosum/simple?nums=2,7,11,15&target=9 - URL参数格式")
	
	http.ListenAndServe(":8000", nil)
}
