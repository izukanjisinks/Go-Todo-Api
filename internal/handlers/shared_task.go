package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/pkg/utils"

	"github.com/google/uuid"
)

type SharedTaskHandler struct {
	repo *repository.SharedTaskRepository
}

func NewSharedTaskHandler() *SharedTaskHandler {
	return &SharedTaskHandler{
		repo: repository.NewSharedTaskRepository(),
	}
}

func (h *SharedTaskHandler) GetAllSharedTasks(w http.ResponseWriter, r *http.Request) {
	sharedTasks, err := h.repo.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, sharedTasks)
}

func (h *SharedTaskHandler) CreateSharedTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var newSharedTask models.SharedTask

	err := json.NewDecoder(r.Body).Decode(&newSharedTask)
	fmt.Println("new shared task", newSharedTask)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error reading body")
		return
	}

	if newSharedTask.OwnerID == 0 || newSharedTask.SharedWithID == 0 || newSharedTask.TodoID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input submitted")
		return
	}

	newSharedTask.ID = uuid.New().String()

	err = h.repo.Create(&newSharedTask)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, newSharedTask)
}

func (h *SharedTaskHandler) SharedTasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.GetAllSharedTasks(w, r)
	case "POST":
		h.CreateSharedTask(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *SharedTaskHandler) GetSharedTaskById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/shared-tasks/"):]

	sharedTask, err := h.repo.GetById(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if sharedTask == nil {
		utils.RespondError(w, http.StatusNotFound, "Shared task not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, sharedTask)
}

func (h *SharedTaskHandler) DeleteSharedTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/shared-tasks/"):]

	rowsAffected, err := h.repo.Delete(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if rowsAffected == 0 {
		utils.RespondError(w, http.StatusNotFound, "Shared task not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Deleted shared task with id " + id + " successfully",
	})
}

func (h *SharedTaskHandler) SharedTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.GetSharedTaskById(w, r)
	case "DELETE":
		h.DeleteSharedTask(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *SharedTaskHandler) GetSharedTasksByOwnerId(w http.ResponseWriter, r *http.Request) {
	// Extract owner_id from query parameter
	ownerIDStr := r.URL.Query().Get("owner_id")
	if ownerIDStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "owner_id query parameter is required")
		return
	}

	// Convert owner_id to int
	var ownerID int
	_, err := fmt.Sscanf(ownerIDStr, "%d", &ownerID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid owner_id format")
		return
	}

	// Get shared tasks from repository
	sharedTasks, err := h.repo.GetByOwnerId(ownerID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, sharedTasks)
}

func (h *SharedTaskHandler) GetSharedTasksById(w http.ResponseWriter, r *http.Request) {
	// Extract id from query parameter
	sharedWithIDStr := r.URL.Query().Get("id")
	if sharedWithIDStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "id query parameter is required")
		return
	}

	// Convert id to int
	var sharedWithID int
	_, err := fmt.Sscanf(sharedWithIDStr, "%d", &sharedWithID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid id format")
		return
	}

	// Get todos shared with the user from repository
	sharedTodos, err := h.repo.GetTodosBySharedId(sharedWithID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, sharedTodos)
}

func (h *SharedTaskHandler) GetSharedTasksByTodoId(w http.ResponseWriter, r *http.Request) {
	// Extract todo_id from query parameter
	todoID := r.URL.Query().Get("todo_id")
	if todoID == "" {
		utils.RespondError(w, http.StatusBadRequest, "todo_id query parameter is required")
		return
	}

	// Get shared tasks from repository
	sharedTasks, err := h.repo.GetByTodoId(todoID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, sharedTasks)
}
