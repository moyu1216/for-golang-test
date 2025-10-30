package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ccc-me/for-golang-test/models"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	Repo *models.TaskRepository
}

func NewTaskHandler(repo *models.TaskRepository) *TaskHandler {
	return &TaskHandler{Repo: repo}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")
	sortBy := r.URL.Query().Get("sort_by") // created_at, due_date, priority

	tasks, err := h.Repo.GetAll(status, priority)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 排序
	if sortBy != "" {
		tasks = sortTasks(tasks, sortBy)
	}

	respondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "无效的任务ID")
		return
	}

	task, err := h.Repo.GetByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "任务未找到")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if task.Title == "" {
		respondError(w, http.StatusBadRequest, "标题不能为空")
		return
	}

	if task.Status == "" {
		task.Status = "pending"
	}
	if task.Priority == "" {
		task.Priority = "medium"
	}

	if err := h.Repo.Create(&task); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "无效的任务ID")
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	task.ID = id
	if err := h.Repo.Update(&task); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedTask, err := h.Repo.GetByID(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, updatedTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "无效的任务ID")
		return
	}

	if err := h.Repo.Delete(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) SearchTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, "搜索关键词不能为空")
		return
	}

	tasks, err := h.Repo.Search(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Repo.GetStats()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

func (h *TaskHandler) BatchDelete(w http.ResponseWriter, r *http.Request) {
	var ids struct {
		IDs []int `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		respondError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	for _, id := range ids.IDs {
		if err := h.Repo.Delete(id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) BatchUpdateStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs    []int  `json:"ids"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if req.Status == "" {
		respondError(w, http.StatusBadRequest, "状态不能为空")
		return
	}

	for _, id := range req.IDs {
		task, err := h.Repo.GetByID(id)
		if err != nil {
			continue
		}
		task.Status = req.Status
		if err := h.Repo.Update(task); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}

func sortTasks(tasks []models.Task, sortBy string) []models.Task {
	switch strings.ToLower(sortBy) {
	case "created_at":
		for i := 0; i < len(tasks)-1; i++ {
			for j := i + 1; j < len(tasks); j++ {
				if tasks[i].CreatedAt.After(tasks[j].CreatedAt) {
					tasks[i], tasks[j] = tasks[j], tasks[i]
				}
			}
		}
	case "due_date":
		for i := 0; i < len(tasks)-1; i++ {
			for j := i + 1; j < len(tasks); j++ {
				if tasks[i].DueDate == nil && tasks[j].DueDate != nil {
					tasks[i], tasks[j] = tasks[j], tasks[i]
				} else if tasks[i].DueDate != nil && tasks[j].DueDate != nil {
					if tasks[i].DueDate.After(*tasks[j].DueDate) {
						tasks[i], tasks[j] = tasks[j], tasks[i]
					}
				}
			}
		}
	case "priority":
		priorityOrder := map[string]int{
			"high":   3,
			"medium": 2,
			"low":    1,
		}
		for i := 0; i < len(tasks)-1; i++ {
			for j := i + 1; j < len(tasks); j++ {
				if priorityOrder[tasks[i].Priority] < priorityOrder[tasks[j].Priority] {
					tasks[i], tasks[j] = tasks[j], tasks[i]
				}
			}
		}
	}
	return tasks
}