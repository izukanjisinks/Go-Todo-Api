package handlers

import (
	"encoding/json"
	"net/http"

	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/pkg/utils"

	"github.com/google/uuid"
)

type TodoWorkflowHandler struct {
	repo *repository.TodoWorkflow
}

func NewTodoWorkflowHandler() *TodoWorkflowHandler {
	return &TodoWorkflowHandler{
		repo: repository.NewTodoWorkflow(),
	}
}

// CreateTodoTask creates a new todo task in draft status
func (h *TodoWorkflowHandler) CreateTodoTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		AssignedTo  string `json:"assigned_to"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error reading body")
		return
	}

	if req.Title == "" || req.AssignedTo == "" {
		utils.RespondError(w, http.StatusBadRequest, "Title and assigned_to are required")
		return
	}

	id := uuid.New().String()

	todo, err := h.repo.CreateTodo(id, req.Title, req.Description, req.AssignedTo)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, todo)
}

// SubmitForReview submits a todo for review
func (h *TodoWorkflowHandler) SubmitForReview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	submittedBy := r.PathValue("submitted_by")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "Todo ID is required")
		return
	}

	if submittedBy == "" {
		utils.RespondError(w, http.StatusBadRequest, "Submitted by is required")
		return
	}

	err := h.repo.SubmitForReview(id, submittedBy)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Todo submitted for review successfully",
	})
}

// ApproveTodo approves a todo
func (h *TodoWorkflowHandler) ApproveTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	approvedBy := r.PathValue("approved_by")

	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "Todo ID is required")
		return
	}

	if approvedBy == "" {
		utils.RespondError(w, http.StatusBadRequest, "Approved by is required")
		return
	}

	err := h.repo.ApprovedTodo(id, approvedBy)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Todo approved successfully",
	})
}

// RejectTodo rejects a todo
func (h *TodoWorkflowHandler) RejectTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rejectedBy := r.PathValue("rejected_by")

	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "Todo ID is required")
		return
	}

	if rejectedBy == "" {
		utils.RespondError(w, http.StatusBadRequest, "Rejected by is required")
		return
	}

	err := h.repo.RejectTodo(id, rejectedBy)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Todo rejected successfully",
	})
}

// GetTodosByUser gets all todos assigned to a user
func (h *TodoWorkflowHandler) GetTodosByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	todos, err := h.repo.GetTodosByUser(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}

// GetTodosByStatus gets all todos with a specific status
func (h *TodoWorkflowHandler) GetTodosByStatus(w http.ResponseWriter, r *http.Request) {
	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "status query parameter is required")
		return
	}

	status := models.TodoStatus(statusStr)
	// Validate status
	if status != models.StatusDraft && status != models.StatusReview && status != models.StatusApproved {
		utils.RespondError(w, http.StatusBadRequest, "Invalid status. Must be Draft, Review, or Approved")
		return
	}

	todos, err := h.repo.GetTodosByStatus(status)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}
