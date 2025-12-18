package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Ccc-me/for-golang-test/db"
	"github.com/Ccc-me/for-golang-test/models"
	"github.com/gorilla/mux"
)

// CreateTaskHandler 创建任务
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "无效的JSON数据: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if task.Title == "" {
		http.Error(w, "title字段为必填", http.StatusBadRequest)
		return
	}

	// 设置默认值
	if task.Priority == "" {
		task.Priority = "medium"
	}
	if task.Status == "" {
		task.Status = "todo"
	}

	// 验证优先级和状态
	if !isValidPriority(task.Priority) {
		http.Error(w, "无效的优先级，必须是: low, medium, high", http.StatusBadRequest)
		return
	}
	if !isValidStatus(task.Status) {
		http.Error(w, "无效的状态，必须是: todo, in_progress, completed", http.StatusBadRequest)
		return
	}

	createdTask, err := models.CreateTask(db.DB, &task)
	if err != nil {
		http.Error(w, "创建任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

// GetTaskHandler 获取单个任务
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "无效的任务ID", http.StatusBadRequest)
		return
	}

	task, err := models.GetTaskByID(db.DB, id)
	if err != nil {
		http.Error(w, "任务不存在", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// ListTasksHandler 获取所有任务
func ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	priorityFilter := r.URL.Query().Get("priority")

	// 验证过滤参数
	if statusFilter != "" && !isValidStatus(statusFilter) {
		http.Error(w, "无效的状态值", http.StatusBadRequest)
		return
	}
	if priorityFilter != "" && !isValidPriority(priorityFilter) {
		http.Error(w, "无效的优先级值", http.StatusBadRequest)
		return
	}

	tasks, err := models.GetAllTasks(db.DB, statusFilter, priorityFilter)
	if err != nil {
		http.Error(w, "获取任务列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// UpdateTaskHandler 更新任务
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "只支持PUT方法", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "无效的任务ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "无效的JSON数据: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if task.Title == "" {
		http.Error(w, "title字段为必填", http.StatusBadRequest)
		return
	}

	// 验证优先级和状态
	if task.Priority != "" && !isValidPriority(task.Priority) {
		http.Error(w, "无效的优先级，必须是: low, medium, high", http.StatusBadRequest)
		return
	}
	if task.Status != "" && !isValidStatus(task.Status) {
		http.Error(w, "无效的状态，必须是: todo, in_progress, completed", http.StatusBadRequest)
		return
	}

	// 获取现有任务
	existingTask, err := models.GetTaskByID(db.DB, id)
	if err != nil {
		http.Error(w, "任务不存在", http.StatusNotFound)
		return
	}

	// 填充未提供的字段
	if task.Priority == "" {
		task.Priority = existingTask.Priority
	}
	if task.Status == "" {
		task.Status = existingTask.Status
	}
	if task.Description == "" {
		task.Description = existingTask.Description
	}
	if task.DueDate == nil {
		task.DueDate = existingTask.DueDate
	}

	updatedTask, err := models.UpdateTask(db.DB, id, &task)
	if err != nil {
		http.Error(w, "更新任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

// DeleteTaskHandler 删除任务
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "只支持DELETE方法", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "无效的任务ID", http.StatusBadRequest)
		return
	}

	// 检查任务是否存在
	_, err = models.GetTaskByID(db.DB, id)
	if err != nil {
		http.Error(w, "任务不存在", http.StatusNotFound)
		return
	}

	if err := models.DeleteTask(db.DB, id); err != nil {
		http.Error(w, "删除任务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// TaskStatsHandler 获取任务统计信息
func TaskStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := models.GetTaskStats(db.DB)
	if err != nil {
		http.Error(w, "获取统计信息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// isValidPriority 验证优先级是否有效
func isValidPriority(priority string) bool {
	return priority == "low" || priority == "medium" || priority == "high"
}

// isValidStatus 验证状态是否有效
func isValidStatus(status string) bool {
	return status == "todo" || status == "in_progress" || status == "completed"
}










