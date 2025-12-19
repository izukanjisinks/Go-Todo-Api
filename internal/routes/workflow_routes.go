// internal/routes/workflow_routes.go
package routes

import (
	"net/http"
	"todo-api/internal/handlers"
	"todo-api/internal/middleware"
	"todo-api/internal/models"
)

func RegisterWorkflowRoutes(
	todoWorkflowHandler *handlers.TodoWorkflowHandler,
	workflowAdminHandler *handlers.WorkflowAdminHandler,
	workflowInstanceHandler *handlers.WorkflowInstanceHandler,
) {
	// Old hardcoded workflow routes
	http.HandleFunc("POST /workflow/todos", withAuthAndPermission(todoWorkflowHandler.CreateTodoTask, models.PermContentCreate))
	http.HandleFunc("GET /workflow/todos/user", withAuthAndPermission(todoWorkflowHandler.GetTodosByUser, models.PermContentRead))
	http.HandleFunc("GET /workflow/todos/status", withAuthAndPermission(todoWorkflowHandler.GetTodosByStatus, models.PermContentRead))
	http.HandleFunc("POST /workflow/todos/{id}/submit/{submitted_by}", withAuthAndPermission(todoWorkflowHandler.SubmitForReview, models.PermContentUpdate))
	http.HandleFunc("POST /workflow/todos/{id}/approve/{approved_by}", withAuthAndPermission(todoWorkflowHandler.ApproveTodo, models.PermContentUpdate))
	http.HandleFunc("POST /workflow/todos/{id}/reject/{rejected_by}", withAuthAndPermission(todoWorkflowHandler.RejectTodo, models.PermContentUpdate))

	// Dynamic Workflow Admin routes (for creating workflows, steps, transitions)
	http.HandleFunc("OPTIONS /api/workflows", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/workflows", withAuthAndPermission(workflowAdminHandler.CreateWorkflow, models.PermContentCreate))
	http.HandleFunc("GET /api/workflows", withAuthAndPermission(workflowAdminHandler.GetAllWorkflows, models.PermContentRead))
	http.HandleFunc("OPTIONS /api/workflows/{id}", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/workflows/{id}", withAuthAndPermission(workflowAdminHandler.GetWorkflow, models.PermContentRead))
	http.HandleFunc("OPTIONS /api/workflows/{workflow_id}/steps", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/workflows/{workflow_id}/steps", withAuthAndPermission(workflowAdminHandler.CreateStep, models.PermContentCreate))
	http.HandleFunc("GET /api/workflows/{workflow_id}/steps", withAuthAndPermission(workflowAdminHandler.GetWorkflowSteps, models.PermContentRead))
	http.HandleFunc("OPTIONS /api/workflows/{workflow_id}/transitions", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/workflows/{workflow_id}/transitions", withAuthAndPermission(workflowAdminHandler.CreateTransition, models.PermContentCreate))
	http.HandleFunc("GET /api/workflows/{workflow_id}/transitions", withAuthAndPermission(workflowAdminHandler.GetWorkflowTransitions, models.PermContentRead))

	// Dynamic Workflow Instance routes
	http.HandleFunc("OPTIONS /api/tasks", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/tasks", withAuthAndPermission(workflowInstanceHandler.StartTask, models.PermContentCreate))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/tasks/{instance_id}", withAuthAndPermission(workflowInstanceHandler.GetTask, models.PermContentRead))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}/execute", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/tasks/{instance_id}/execute", withAuthAndPermission(workflowInstanceHandler.ExecuteAction, models.PermContentUpdate))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}/actions", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/tasks/{instance_id}/actions", withAuthAndPermission(workflowInstanceHandler.GetAvailableActions, models.PermContentRead))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}/history", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("OPTIONS /api/tasks/user", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/tasks/user", withAuthAndPermission(workflowInstanceHandler.GetTasksByUser, models.PermContentRead))
	http.HandleFunc("OPTIONS /api/workflows/{workflow_id}/tasks", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/workflows/{workflow_id}/tasks", withAuthAndPermission(workflowInstanceHandler.GetTasksByWorkflow, models.PermContentRead))
}
