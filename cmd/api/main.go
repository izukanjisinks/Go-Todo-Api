package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/internal/config"
	"todo-api/internal/database"
	"todo-api/internal/handlers"
	"todo-api/internal/middleware"
)

//var users = map[string]Login{}

// withAuth wraps a handler with both CORS and JWT authentication middleware
func withAuth(handler http.HandlerFunc) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(http.HandlerFunc(handler)).ServeHTTP(w, r)
	})
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	err = database.Connect(cfg.GetConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler()
	userHandler := handlers.NewUserRepository()
	sharedTaskHandler := handlers.NewSharedTaskHandler()
	todoWorkflowHandler := handlers.NewTodoWorkflowHandler()
	workflowAdminHandler := handlers.NewWorkflowAdminHandler()
	workflowInstanceHandler := handlers.NewWorkflowInstanceHandler()

	// Setup routes
	// Public routes
	http.HandleFunc("/health", middleware.CORS(handlers.HealthHandler))
	http.HandleFunc("/register", middleware.CORS(userHandler.Register))
	http.HandleFunc("/login", middleware.CORS(userHandler.Login))

	// Protected routes (require JWT authentication)
	http.HandleFunc("/todos", withAuth(todoHandler.TodosHandler))
	http.HandleFunc("/todos/", withAuth(todoHandler.TodoByIdHandler))
	http.HandleFunc("/todos/user", withAuth(todoHandler.GetTodosByUserId))
	http.HandleFunc("/users", withAuth(userHandler.GetUsers))
	http.HandleFunc("/users/update", withAuth(userHandler.UpdateUser))
	http.HandleFunc("/logout", withAuth(handlers.Logout))
	http.HandleFunc("/protected", withAuth(handlers.Protected))
	http.HandleFunc("/shared-tasks", withAuth(sharedTaskHandler.SharedTasksHandler))
	http.HandleFunc("/shared-tasks/", withAuth(sharedTaskHandler.SharedTaskByIdHandler))
	http.HandleFunc("/shared-tasks/owner", withAuth(sharedTaskHandler.GetSharedTasksByOwnerId))
	http.HandleFunc("/shared-tasks/id", withAuth(sharedTaskHandler.GetSharedTasksById))
	http.HandleFunc("/shared-tasks/todo", withAuth(sharedTaskHandler.GetSharedTasksByTodoId))

	// Workflow routes (old hardcoded workflow) - Protected
	http.HandleFunc("POST /workflow/todos", withAuth(todoWorkflowHandler.CreateTodoTask))
	http.HandleFunc("GET /workflow/todos/user", withAuth(todoWorkflowHandler.GetTodosByUser))
	http.HandleFunc("GET /workflow/todos/status", withAuth(todoWorkflowHandler.GetTodosByStatus))
	http.HandleFunc("POST /workflow/todos/{id}/submit/{submitted_by}", withAuth(todoWorkflowHandler.SubmitForReview))
	http.HandleFunc("POST /workflow/todos/{id}/approve/{approved_by}", withAuth(todoWorkflowHandler.ApproveTodo))
	http.HandleFunc("POST /workflow/todos/{id}/reject/{rejected_by}", withAuth(todoWorkflowHandler.RejectTodo))

	// Dynamic Workflow Admin routes (for creating workflows, steps, transitions) - Protected
	http.HandleFunc("OPTIONS /api/workflows", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/workflows", withAuth(workflowAdminHandler.CreateWorkflow))
	http.HandleFunc("GET /api/workflows", withAuth(workflowAdminHandler.GetAllWorkflows))
	http.HandleFunc("OPTIONS /api/workflows/{id}", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/workflows/{id}", withAuth(workflowAdminHandler.GetWorkflow))
	http.HandleFunc("OPTIONS /api/workflows/{workflow_id}/steps", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/workflows/{workflow_id}/steps", withAuth(workflowAdminHandler.CreateStep))
	http.HandleFunc("GET /api/workflows/{workflow_id}/steps", withAuth(workflowAdminHandler.GetWorkflowSteps))
	http.HandleFunc("OPTIONS /api/workflows/{workflow_id}/transitions", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/workflows/{workflow_id}/transitions", withAuth(workflowAdminHandler.CreateTransition))
	http.HandleFunc("GET /api/workflows/{workflow_id}/transitions", withAuth(workflowAdminHandler.GetWorkflowTransitions))

	// Dynamic Workflow Instance routes (for running tasks through workflows) - Protected
	http.HandleFunc("OPTIONS /api/tasks", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/tasks", withAuth(workflowInstanceHandler.StartTask))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/tasks/{instance_id}", withAuth(workflowInstanceHandler.GetTask))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}/execute", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/tasks/{instance_id}/execute", withAuth(workflowInstanceHandler.ExecuteAction))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}/actions", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/tasks/{instance_id}/actions", withAuth(workflowInstanceHandler.GetAvailableActions))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}/history", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("OPTIONS /api/tasks/user", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/tasks/user", withAuth(workflowInstanceHandler.GetTasksByUser))
	http.HandleFunc("OPTIONS /api/workflows/{workflow_id}/tasks", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/workflows/{workflow_id}/tasks", withAuth(workflowInstanceHandler.GetTasksByWorkflow))

	// Start server
	port := ":" + cfg.ServerPort
	fmt.Printf("🚀 Server listening on port %s\n", cfg.ServerPort)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
