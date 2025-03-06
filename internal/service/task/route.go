package task

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
)

type Handler struct {
	log *slog.Logger
}

var (
	tasks     = []Task{}
	idCounter = 1
	mu        sync.Mutex
)

func NewTaskHandler(logger *slog.Logger) *Handler {
	return &Handler{
		log: logger,
	}
}

func (h *Handler) SetupRoute(router *http.ServeMux) {
	router.HandleFunc("/tasks", h.HandleTask)
	router.HandleFunc("/tasks/", h.HandleTaskByID)
}
func (h *Handler) HandleTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.log.Info("current tasks", slog.Any("tasks", tasks))
		json.NewEncoder(w).Encode(tasks)
	case http.MethodPost:
		var newTask Task
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		mu.Lock()
		newTask.ID = idCounter
		idCounter++
		tasks = append(tasks, newTask)
		mu.Unlock()
		w.WriteHeader(http.StatusCreated)
		h.log.Info("created newTask", slog.Any("newTask", newTask))
		json.NewEncoder(w).Encode(newTask)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPut:
		var updatedTask Task
		if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Title = updatedTask.Title
				tasks[i].Done = updatedTask.Done
				h.log.Info("updated task", slog.Any("updatedTasks", tasks[i]))
				json.NewEncoder(w).Encode(&tasks[i])
				return
			}
		}
		http.Error(w, "Task not found", http.StatusNotFound)
	case http.MethodDelete:
		mu.Lock()
		defer mu.Unlock()
		for i, task := range tasks {
			if task.ID == id {
				h.log.Info("deleted task", slog.Any("updatedTasks", tasks[i]))
				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "Task deleted"}`))
				return
			}
		}
		http.Error(w, "Task not found", http.StatusNotFound)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
