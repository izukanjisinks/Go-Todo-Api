package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"todo-api/internal/config"
	"todo-api/internal/database"
	"todo-api/internal/handlers"
	"todo-api/internal/middleware"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/internal/services"
)

func seedDefaultAdmin(userService *services.UserService, roleService *services.RoleService) {
	users, _ := userService.GetAllUsers()
	hasSuperAdmin := false

	for _, user := range users {
		if user.Role != nil && user.Role.Name == models.RoleSuperAdmin {
			hasSuperAdmin = true
			break
		}
	}

	if !hasSuperAdmin {
		log.Println("No super admin found, creating default super admin...")
	}

	superAdmin, err := roleService.GetRoleByName(models.RoleSuperAdmin)

	if err != nil {
		log.Println("Warning: Could not create super admin role, super admin role not found...")
	}

	adminUser := &models.User{
		Username: "admin",
		Email:    "admin@backend.com",
		RoleID:   &superAdmin.RoleId,
		IsAdmin:  true,
		Login: models.Login{
			Password: "Admin@123",
		},
	}

	err = userService.Register(adminUser)

	if err != nil {
		log.Printf("Warning, could not create default admin user %v\n", err)
	} else {
		log.Println("Default admin user created successfully!")
		log.Println("Email: admin@backend.com")
		log.Println("Password: Admin@123")
	}

}

// withAuth wraps a handler with both CORS and JWT authentication middleware
func withAuth(handler http.HandlerFunc) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(http.HandlerFunc(handler)).ServeHTTP(w, r)
	})
}

// withAuthAndPermission wraps a handler with CORS, JWT auth, and permission check
func withAuthAndPermission(handler http.HandlerFunc, permission string) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(
			middleware.RequirePermission(permission)(
				http.HandlerFunc(handler),
			),
		).ServeHTTP(w, r)
	})
}

// withAuthAndAnyPermission wraps a handler with CORS, JWT auth, and any permission check
func withAuthAndAnyPermission(handler http.HandlerFunc, permissions ...string) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(
			middleware.RequireAnyPermission(permissions...)(
				http.HandlerFunc(handler),
			),
		).ServeHTTP(w, r)
	})
}

// withAuthAndAllPermissions wraps a handler with CORS, JWT auth, and all permissions check
func withAuthAndAllPermissions(handler http.HandlerFunc, permissions ...string) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(
			middleware.RequireAllPermissions(permissions...)(
				http.HandlerFunc(handler),
			),
		).ServeHTTP(w, r)
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

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	todoRepo := repository.NewTodoRepository()
	roleRepo := repository.NewRoleRepository()

	// Initialize services with repository dependencies
	userService := services.NewUserService(userRepo)
	todoService := services.NewTodoService(todoRepo)
	roleService := services.NewRoleService(roleRepo)

	// Initialize predefined roles (run once at startup)
	err = roleService.InitializePredefinedRoles()
	if err != nil {
		log.Println("Warning: Failed to initialize predefined roles:", err)
	}

	seedDefaultAdmin(userService, roleService)

	// Initialize handlers with service dependencies
	todoHandler := handlers.NewTodoHandler(todoService)
	userHandler := handlers.NewUsersHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	sharedTaskHandler := handlers.NewSharedTaskHandler()
	todoWorkflowHandler := handlers.NewTodoWorkflowHandler()
	workflowAdminHandler := handlers.NewWorkflowAdminHandler()
	workflowInstanceHandler := handlers.NewWorkflowInstanceHandler()

	// Setup routes
	// Public routes
	http.HandleFunc("/health", middleware.CORS(handlers.HealthHandler))
	http.HandleFunc("/register", middleware.CORS(userHandler.Register))
	http.HandleFunc("/login", middleware.CORS(userHandler.Login))

	// Protected routes (require JWT authentication + permissions)
	// Todos - mapped to content permissions
	http.HandleFunc("/todos", withAuthAndPermission(todoHandler.TodosHandler, models.PermContentRead))
	http.HandleFunc("/todos/", withAuthAndPermission(todoHandler.TodoByIdHandler, models.PermContentRead))
	http.HandleFunc("/todos/user", withAuthAndPermission(todoHandler.GetTodosByUserId, models.PermContentRead))
	// Users - require users permissions
	http.HandleFunc("/users", withAuthAndPermission(userHandler.GetUsers, models.PermUsersRead))
	http.HandleFunc("/users/update", withAuthAndPermission(userHandler.UpdateUser, models.PermUsersUpdate))
	http.HandleFunc("/logout", withAuth(handlers.Logout))
	http.HandleFunc("/protected", withAuth(handlers.Protected))
	// Shared tasks routes - mapped to content permissions
	http.HandleFunc("/shared-tasks", withAuthAndPermission(sharedTaskHandler.SharedTasksHandler, models.PermContentRead))
	http.HandleFunc("/shared-tasks/", withAuthAndPermission(sharedTaskHandler.SharedTaskByIdHandler, models.PermContentRead))
	http.HandleFunc("/shared-tasks/owner", withAuthAndPermission(sharedTaskHandler.GetSharedTasksByOwnerId, models.PermContentRead))
	http.HandleFunc("/shared-tasks/id", withAuthAndPermission(sharedTaskHandler.GetSharedTasksById, models.PermContentRead))
	http.HandleFunc("/shared-tasks/todo", withAuthAndPermission(sharedTaskHandler.GetSharedTasksByTodoId, models.PermContentRead))

	// RBAC Role Management routes - Protected (require roles:manage permission)
	http.HandleFunc("/roles", withAuthAndPermission(roleHandler.RolesHandler, models.PermRolesManage))
	http.HandleFunc("/roles/", withAuthAndPermission(roleHandler.RolesHandler, models.PermRolesManage))

	// User role assignment routes - must come before /users/ to avoid conflicts
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a role assignment or permissions request
		if strings.HasSuffix(r.URL.Path, "/role") {
			withAuthAndPermission(roleHandler.AssignRoleHandler, models.PermRolesManage)(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/permissions") {
			withAuthAndPermission(roleHandler.GetUserPermissionsHandler, models.PermUsersRead)(w, r)
		} else {
			// Default user handler
			withAuthAndPermission(userHandler.GetUsers, models.PermUsersRead)(w, r)
		}
	})

	// Workflow routes (old hardcoded workflow) - Protected, mapped to content permissions
	http.HandleFunc("POST /workflow/todos", withAuthAndPermission(todoWorkflowHandler.CreateTodoTask, models.PermContentCreate))
	http.HandleFunc("GET /workflow/todos/user", withAuthAndPermission(todoWorkflowHandler.GetTodosByUser, models.PermContentRead))
	http.HandleFunc("GET /workflow/todos/status", withAuthAndPermission(todoWorkflowHandler.GetTodosByStatus, models.PermContentRead))
	http.HandleFunc("POST /workflow/todos/{id}/submit/{submitted_by}", withAuthAndPermission(todoWorkflowHandler.SubmitForReview, models.PermContentUpdate))
	http.HandleFunc("POST /workflow/todos/{id}/approve/{approved_by}", withAuthAndPermission(todoWorkflowHandler.ApproveTodo, models.PermContentUpdate))
	http.HandleFunc("POST /workflow/todos/{id}/reject/{rejected_by}", withAuthAndPermission(todoWorkflowHandler.RejectTodo, models.PermContentUpdate))

	// Dynamic Workflow Admin routes (for creating workflows, steps, transitions) - Protected
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

	// Dynamic Workflow Instance routes (for running tasks through workflows) - Protected
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

	// Start server
	port := ":" + cfg.ServerPort
	fmt.Printf("🚀 Server listening on port %s\n", cfg.ServerPort)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
