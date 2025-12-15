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
	http.HandleFunc("/health", middleware.CORS(handlers.HealthHandler))
	http.HandleFunc("/todos", middleware.CORS(todoHandler.TodosHandler))
	http.HandleFunc("/todos/", middleware.CORS(todoHandler.TodoByIdHandler))
	http.HandleFunc("/users", middleware.CORS(userHandler.GetUsers))
	http.HandleFunc("/register", middleware.CORS(userHandler.Register))
	http.HandleFunc("/login", middleware.CORS(userHandler.Login))
	http.HandleFunc("/logout", middleware.CORS(handlers.Logout))
	http.HandleFunc("/protected", middleware.CORS(handlers.Protected))
	http.HandleFunc("/todos/user", middleware.CORS(todoHandler.GetTodosByUserId))
	http.HandleFunc("/shared-tasks", middleware.CORS(sharedTaskHandler.SharedTasksHandler))
	http.HandleFunc("/shared-tasks/", middleware.CORS(sharedTaskHandler.SharedTaskByIdHandler))
	http.HandleFunc("/shared-tasks/owner", middleware.CORS(sharedTaskHandler.GetSharedTasksByOwnerId))
	http.HandleFunc("/shared-tasks/id", middleware.CORS(sharedTaskHandler.GetSharedTasksById))
	http.HandleFunc("/shared-tasks/todo", middleware.CORS(sharedTaskHandler.GetSharedTasksByTodoId))

	// Workflow routes (old hardcoded workflow)
	http.HandleFunc("POST /workflow/todos", middleware.CORS(todoWorkflowHandler.CreateTodoTask))
	http.HandleFunc("GET /workflow/todos/user", middleware.CORS(todoWorkflowHandler.GetTodosByUser))
	http.HandleFunc("GET /workflow/todos/status", middleware.CORS(todoWorkflowHandler.GetTodosByStatus))
	http.HandleFunc("POST /workflow/todos/{id}/submit/{submitted_by}", middleware.CORS(todoWorkflowHandler.SubmitForReview))
	http.HandleFunc("POST /workflow/todos/{id}/approve/{approved_by}", middleware.CORS(todoWorkflowHandler.ApproveTodo))
	http.HandleFunc("POST /workflow/todos/{id}/reject/{rejected_by}", middleware.CORS(todoWorkflowHandler.RejectTodo))

	// Dynamic Workflow Admin routes (for creating workflows, steps, transitions)
	http.HandleFunc("OPTIONS /api/workflows", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/workflows", middleware.CORS(workflowAdminHandler.CreateWorkflow))
	http.HandleFunc("GET /api/workflows", middleware.CORS(workflowAdminHandler.GetAllWorkflows))
	http.HandleFunc("OPTIONS /api/workflows/{id}", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/workflows/{id}", middleware.CORS(workflowAdminHandler.GetWorkflow))
	http.HandleFunc("OPTIONS /api/workflows/{workflow_id}/steps", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/workflows/{workflow_id}/steps", middleware.CORS(workflowAdminHandler.CreateStep))
	http.HandleFunc("GET /api/workflows/{workflow_id}/steps", middleware.CORS(workflowAdminHandler.GetWorkflowSteps))
	http.HandleFunc("OPTIONS /api/workflows/{workflow_id}/transitions", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/workflows/{workflow_id}/transitions", middleware.CORS(workflowAdminHandler.CreateTransition))
	http.HandleFunc("GET /api/workflows/{workflow_id}/transitions", middleware.CORS(workflowAdminHandler.GetWorkflowTransitions))

	// Dynamic Workflow Instance routes (for running tasks through workflows)
	http.HandleFunc("OPTIONS /api/tasks", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/tasks", middleware.CORS(workflowInstanceHandler.StartTask))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/tasks/{instance_id}", middleware.CORS(workflowInstanceHandler.GetTask))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}/execute", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("POST /api/tasks/{instance_id}/execute", middleware.CORS(workflowInstanceHandler.ExecuteAction))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}/actions", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/tasks/{instance_id}/actions", middleware.CORS(workflowInstanceHandler.GetAvailableActions))
	http.HandleFunc("OPTIONS /api/tasks/{instance_id}/history", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("OPTIONS /api/tasks/user", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/tasks/user", middleware.CORS(workflowInstanceHandler.GetTasksByUser))
	http.HandleFunc("OPTIONS /api/workflows/{workflow_id}/tasks", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {}))
	http.HandleFunc("GET /api/workflows/{workflow_id}/tasks", middleware.CORS(workflowInstanceHandler.GetTasksByWorkflow))

	// Start server
	port := ":" + cfg.ServerPort
	fmt.Printf("🚀 Server listening on port %s\n", cfg.ServerPort)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
