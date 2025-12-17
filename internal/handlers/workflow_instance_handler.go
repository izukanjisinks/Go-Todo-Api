package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todo-api/internal/repository"
	"todo-api/internal/services"
	"todo-api/pkg/utils"
)

// WorkflowInstanceHandler handles workflow instance operations (starting tasks, executing actions)
type WorkflowInstanceHandler struct {
	engine       *services.WorkflowEngine
	instanceRepo *repository.WorkflowInstanceRepository
}

func NewWorkflowInstanceHandler() *WorkflowInstanceHandler {
	return &WorkflowInstanceHandler{
		engine:       services.NewWorkflowEngine(),
		instanceRepo: repository.NewWorkflowInstanceRepository(),
	}
}

// StartTask creates a new workflow instance
func (h *WorkflowInstanceHandler) StartTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WorkflowID    string    `json:"workflow_id"`
		AssignedTo    string    `json:"assigned_to"`
		CurrentStepId string    `json:"current_step_id"`
		TodoId        string    `json:"todo_id"`
		UpdatedAt     time.Time `json:"updated_at"`
		CreatedAt     time.Time `json:"created_at"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error reading body")
		return
	}

	if req.WorkflowID == "" || req.AssignedTo == "" || req.TodoId == "" {
		utils.RespondError(w, http.StatusBadRequest, "workflow_id, assigned_to, and todo_id are required")
		return
	}

	instance, err := h.engine.StartWorkflow(req.WorkflowID, req.AssignedTo)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, instance)
}

// ExecuteAction executes a workflow action (transition)
func (h *WorkflowInstanceHandler) ExecuteAction(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("instance_id")
	if instanceID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Instance ID is required")
		return
	}

	var req struct {
		ActionName string `json:"action_name"`
		UserID     string `json:"user_id"`
		Comments   string `json:"comments"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error reading body")
		return
	}

	if req.ActionName == "" || req.UserID == "" {
		utils.RespondError(w, http.StatusBadRequest, "action_name and user_id are required")
		return
	}

	err = h.engine.ExecuteTransition(instanceID, req.ActionName, req.UserID, req.Comments)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Action executed successfully",
	})
}

// GetTask retrieves a task with details
func (h *WorkflowInstanceHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("instance_id")
	userID := r.URL.Query().Get("user_id")

	if instanceID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Instance ID is required")
		return
	}

	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	instanceDetails, err := h.engine.GetInstanceWithDetails(instanceID, userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, instanceDetails)
}

// GetAvailableActions retrieves available actions for a task
func (h *WorkflowInstanceHandler) GetAvailableActions(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("instance_id")
	userID := r.URL.Query().Get("user_id")

	if instanceID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Instance ID is required")
		return
	}

	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	actions, err := h.engine.GetAvailableActions(instanceID, userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, actions)
}

// GetTasksByUser retrieves all tasks assigned to a user
func (h *WorkflowInstanceHandler) GetTasksByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	instances, err := h.instanceRepo.GetInstancesByUser(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, instances)
}

// GetTasksByWorkflow retrieves all tasks for a workflow
func (h *WorkflowInstanceHandler) GetTasksByWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("workflow_id")
	if workflowID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Workflow ID is required")
		return
	}

	instances, err := h.instanceRepo.GetInstancesByWorkflow(workflowID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, instances)
}
