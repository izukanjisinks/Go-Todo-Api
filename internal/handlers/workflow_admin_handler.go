package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/pkg/utils"

	"github.com/google/uuid"
)

// WorkflowAdminHandler handles workflow administration (creating workflows, steps, transitions)
type WorkflowAdminHandler struct {
	repo *repository.WorkflowRepository
}

func NewWorkflowAdminHandler() *WorkflowAdminHandler {
	return &WorkflowAdminHandler{
		repo: repository.NewWorkflowRepository(),
	}
}

// CreateWorkflow creates a new workflow template
func (h *WorkflowAdminHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedBy   string `json:"created_by"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error reading body")
		return
	}

	if req.Name == "" || req.CreatedBy == "" {
		utils.RespondError(w, http.StatusBadRequest, "Name and created_by are required")
		return
	}

	workflow := &models.Workflow{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = h.repo.CreateWorkflow(workflow)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, workflow)
}

// GetWorkflow retrieves a workflow by ID
func (h *WorkflowAdminHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "Workflow ID is required")
		return
	}

	workflow, err := h.repo.GetWorkflow(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, workflow)
}

// GetAllWorkflows retrieves all active workflows
func (h *WorkflowAdminHandler) GetAllWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := h.repo.GetAllWorkflows()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, workflows)
}

// CreateStep creates a new workflow step
func (h *WorkflowAdminHandler) CreateStep(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("workflow_id")
	if workflowID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Workflow ID is required")
		return
	}

	var req struct {
		StepName     string   `json:"step_name"`
		StepOrder    int      `json:"step_order"`
		Initial      bool     `json:"initial"`
		Final        bool     `json:"final"`
		AllowedRoles []string `json:"allowed_roles"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error reading body")
		return
	}

	if req.StepName == "" {
		utils.RespondError(w, http.StatusBadRequest, "Step name is required")
		return
	}

	step := &models.WorkflowStep{
		ID:           uuid.New().String(),
		WorkflowID:   workflowID,
		StepName:     req.StepName,
		StepOrder:    req.StepOrder,
		Initial:      req.Initial,
		Final:        req.Final,
		AllowedRoles: req.AllowedRoles,
		CreatedAt:    time.Now(),
	}

	err = h.repo.CreateStep(step)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, step)
}

// GetWorkflowSteps retrieves all steps for a workflow
func (h *WorkflowAdminHandler) GetWorkflowSteps(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("workflow_id")
	if workflowID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Workflow ID is required")
		return
	}

	steps, err := h.repo.GetWorkflowSteps(workflowID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, steps)
}

// CreateTransition creates a new workflow transition
func (h *WorkflowAdminHandler) CreateTransition(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("workflow_id")
	if workflowID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Workflow ID is required")
		return
	}

	var req struct {
		FromStepID     string `json:"from_step_id"`
		ToStepID       string `json:"to_step_id"`
		ActionName     string `json:"action_name"`
		ConditionType  string `json:"condition_type"`
		ConditionValue string `json:"condition_value"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error reading body")
		return
	}

	if req.FromStepID == "" || req.ToStepID == "" || req.ActionName == "" {
		utils.RespondError(w, http.StatusBadRequest, "from_step_id, to_step_id, and action_name are required")
		return
	}

	transition := &models.WorkflowTransition{
		ID:             uuid.New().String(),
		WorkflowID:     workflowID,
		FromStepID:     req.FromStepID,
		ToStepID:       req.ToStepID,
		ActionName:     req.ActionName,
		ConditionType:  req.ConditionType,
		ConditionValue: req.ConditionValue,
		CreatedAt:      time.Now(),
	}

	err = h.repo.CreateTransition(transition)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, transition)
}

// GetWorkflowTransitions retrieves all transitions for a workflow
func (h *WorkflowAdminHandler) GetWorkflowTransitions(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("workflow_id")
	if workflowID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Workflow ID is required")
		return
	}

	transitions, err := h.repo.GetTransitions(workflowID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, transitions)
}
