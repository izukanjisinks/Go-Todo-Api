package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/internal/config"
	"todo-api/internal/database"
	"todo-api/internal/handlers"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/internal/routes"
	"todo-api/internal/services"
)

func seedDefaultAdmin(userService *services.UserService, roleService *services.RoleService) {

	superAdmin, err := roleService.GetRoleByName(models.RoleSuperAdmin)
	if err != nil || superAdmin == nil {
		log.Println("Warning: Could not find super admin role, skipping default admin creation...")
		return
	}

	user, _ := userService.GetUserByRoleId(superAdmin.RoleId)
	hasSuperAdmin := false

	if user != nil {
		if user.Role != nil && user.Role.Name == models.RoleSuperAdmin {
			hasSuperAdmin = true
		}
	}

	if !hasSuperAdmin {
		log.Println("No super admin found, creating default super admin...")
	} else {
		log.Println("Super admin user already exists, skipping creation...")
		return
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

	// Register all routes using the routes package
	routes.RegisterRoutes(
		todoHandler,
		userHandler,
		sharedTaskHandler,
		roleHandler,
		todoWorkflowHandler,
		workflowAdminHandler,
		workflowInstanceHandler,
	)

	// Start server
	port := ":" + cfg.ServerPort
	fmt.Printf("🚀 Server listening on port %s\n", cfg.ServerPort)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
