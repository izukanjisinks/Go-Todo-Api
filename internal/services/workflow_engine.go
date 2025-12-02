package services

import (
	"fmt"
	"time"
	"todo-api/internal/models"
	"todo-api/internal/repository"

	"github.com/google/uuid"
)

type WorkflowEngine struct {
	workflowRepo *repository.WorkflowRepository
	instanceRepo *repository.WorkflowInstanceRepository
}

func NewWorkflowEngine() *WorkflowEngine {
	return &WorkflowEngine{
		workflowRepo: repository.NewWorkflowRepository(),
		instanceRepo: repository.NewWorkflowInstanceRepository(),
	}
}

// StartWorkflow creates a new workflow instance at the start step
func (e *WorkflowEngine) StartWorkflow(workflowID, title, description, taskData, assignedTo, createdBy string) (*models.WorkflowInstance, error) {
	// Get workflow to ensure it exists and is active
	workflow, err := e.workflowRepo.GetWorkflow(workflowID)
	if err != nil {
		return nil, fmt.Errorf("workflow not found: %w", err)
	}

	if !workflow.IsActive {
		return nil, fmt.Errorf("workflow is not active")
	}

	// Get the start step
	startStep, err := e.workflowRepo.GetStartStep(workflowID)
	if err != nil {
		return nil, fmt.Errorf("start step not found: %w", err)
	}

	// Create the instance
	instance := &models.WorkflowInstance{
		ID:            uuid.New().String(),
		WorkflowID:    workflowID,
		CurrentStepID: startStep.ID,
		Title:         title,
		Description:   description,
		TaskData:      taskData,
		AssignedTo:    assignedTo,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = e.instanceRepo.CreateInstance(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	// Create initial history record
	history := &models.WorkflowHistory{
		ID:          uuid.New().String(),
		InstanceID:  instance.ID,
		FromStepID:  nil, // No previous step
		ToStepID:    startStep.ID,
		ActionTaken: "created",
		PerformedBy: createdBy,
		Comments:    "Workflow instance created",
		Timestamp:   time.Now(),
	}

	err = e.instanceRepo.CreateHistory(history)
	if err != nil {
		return nil, fmt.Errorf("failed to create history: %w", err)
	}

	return instance, nil
}

// ExecuteTransition moves an instance from one step to another
func (e *WorkflowEngine) ExecuteTransition(instanceID, actionName, userID, comments string) error {
	// Get the instance
	instance, err := e.instanceRepo.GetInstance(instanceID)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	// Find the transition
	transition, err := e.workflowRepo.FindTransition(instance.WorkflowID, instance.CurrentStepID, actionName)
	if err != nil {
		return fmt.Errorf("invalid action for current step: %w", err)
	}

	// Validate the transition
	canTransition, err := e.ValidateTransition(instance, transition, userID)
	if err != nil {
		return err
	}
	if !canTransition {
		return fmt.Errorf("user not authorized to perform this action")
	}

	// Execute the transition
	oldStepID := instance.CurrentStepID
	err = e.instanceRepo.UpdateInstanceStep(instanceID, transition.ToStepID)
	if err != nil {
		return fmt.Errorf("failed to update instance step: %w", err)
	}

	// Create history record
	history := &models.WorkflowHistory{
		ID:          uuid.New().String(),
		InstanceID:  instanceID,
		FromStepID:  &oldStepID,
		ToStepID:    transition.ToStepID,
		ActionTaken: actionName,
		PerformedBy: userID,
		Comments:    comments,
		Timestamp:   time.Now(),
	}

	err = e.instanceRepo.CreateHistory(history)
	if err != nil {
		return fmt.Errorf("failed to create history: %w", err)
	}

	// Check if we reached an end step
	toStep, err := e.workflowRepo.GetStep(transition.ToStepID)
	if err == nil && toStep.Final {
		// Could trigger completion hooks here
		// e.g., send notifications, update external systems, etc.
	}

	return nil
}

// ValidateTransition checks if a user can perform a transition
func (e *WorkflowEngine) ValidateTransition(instance *models.WorkflowInstance, transition *models.WorkflowTransition, userID string) (bool, error) {
	// Check condition type
	switch transition.ConditionType {
	case "assigned_user_only":
		// Only the assigned user can perform this action
		return instance.AssignedTo == userID, nil

	case "creator_only":
		// Only the creator can perform this action
		return instance.CreatedBy == userID, nil

	case "not_assigned_user":
		// Anyone except the assigned user (e.g., for approval by someone else)
		return instance.AssignedTo != userID, nil

	case "any_user":
		// Any authenticated user can perform this action
		return true, nil

	case "":
		// No condition specified, allow by default
		return true, nil

	default:
		// Unknown condition type
		return false, fmt.Errorf("unknown condition type: %s", transition.ConditionType)
	}
}

// GetAvailableActions returns the actions a user can take on an instance
func (e *WorkflowEngine) GetAvailableActions(instanceID, userID string) ([]models.AvailableAction, error) {
	// Get the instance
	instance, err := e.instanceRepo.GetInstance(instanceID)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}

	// Get available transitions from current step
	transitions, err := e.workflowRepo.GetAvailableTransitions(instance.WorkflowID, instance.CurrentStepID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transitions: %w", err)
	}

	var actions []models.AvailableAction
	for _, transition := range transitions {
		// Check if user can perform this transition
		canPerform, err := e.ValidateTransition(instance, transition, userID)
		if err != nil || !canPerform {
			continue
		}

		// Get the target step name
		toStep, err := e.workflowRepo.GetStep(transition.ToStepID)
		if err != nil {
			continue
		}

		actions = append(actions, models.AvailableAction{
			ActionName:   transition.ActionName,
			ToStepName:   toStep.StepName,
			TransitionID: transition.ID,
		})
	}

	return actions, nil
}

// GetInstanceWithDetails returns an instance with current step and available actions
func (e *WorkflowEngine) GetInstanceWithDetails(instanceID, userID string) (*models.WorkflowInstanceWithDetails, error) {
	// Get the instance
	instance, err := e.instanceRepo.GetInstance(instanceID)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}

	// Get current step
	currentStep, err := e.workflowRepo.GetStep(instance.CurrentStepID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current step: %w", err)
	}

	// Get workflow
	workflow, err := e.workflowRepo.GetWorkflow(instance.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	// Get available actions
	actions, err := e.GetAvailableActions(instanceID, userID)
	if err != nil {
		actions = []models.AvailableAction{} // Empty if error
	}

	return &models.WorkflowInstanceWithDetails{
		WorkflowInstance: *instance,
		CurrentStepName:  currentStep.StepName,
		WorkflowName:     workflow.Name,
		AvailableActions: actions,
	}, nil
}
