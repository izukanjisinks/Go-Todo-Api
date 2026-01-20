package examples

import (
	"net/http"

	"todo-api/internal/middleware"
	"todo-api/internal/models"
)

// Example of how to use RBAC middleware in your routes

func ExampleRoutes() {
	// Example 1: Require a specific permission
	http.Handle("/api/users", middleware.RequirePermission(models.PermView)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handler code - only users with "view" permission can access
			w.Write([]byte("User list"))
		}),
	))

	// Example 2: Require any of multiple permissions
	http.Handle("/api/content", middleware.RequireAnyPermission(
		models.PermView,
		models.PermCreate,
	)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handler code - users with either permission can access
			w.Write([]byte("Content"))
		}),
	))

	// Example 3: Require all permissions
	http.Handle("/api/admin/settings", middleware.RequireAllPermissions(
		models.PermView,
		models.PermUpdate,
	)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handler code - users must have both permissions
			w.Write([]byte("Settings"))
		}),
	))

	// Example 4: Require a specific role
	http.Handle("/api/admin/roles", middleware.RequireRole(models.RoleSuperAdmin)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handler code - only Super Admins can access
			w.Write([]byte("Role management"))
		}),
	))

	// Example 5: Check permission programmatically in handler
	http.HandleFunc("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*models.User)
		if !ok || user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodGet:
			if !user.HasPermission(models.PermView) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			// Handle GET
		case http.MethodPost:
			if !user.HasPermission(models.PermCreate) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			// Handle POST
		case http.MethodPut:
			if !user.HasPermission(models.PermUpdate) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			// Handle PUT
		case http.MethodDelete:
			if !user.HasPermission(models.PermDelete) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			// Handle DELETE
		}
	})
}

// Example of using RoleService
func ExampleRoleService() {
	// Assuming you have a RoleService instance
	// roleService := services.NewRoleService(db)

	// Initialize predefined roles (run once during app startup)
	// err := roleService.InitializePredefinedRoles()

	// Assign a role to a user
	// err := roleService.AssignRoleToUser(userID, roleID)

	// Check if a user has a permission
	// hasPermission, err := roleService.CheckPermission(userID, models.PermUsersRead)

	// Get all roles
	// roles, err := roleService.GetAllRoles()

	// Create a custom role
	// customRole := &models.Role{
	// 	Name:        "Content Editor",
	// 	Description: "Can edit content but not delete",
	// 	Permissions: []string{
	// 		models.PermContentRead,
	// 		models.PermContentCreate,
	// 		models.PermContentUpdate,
	// 	},
	// }
	// err := roleService.CreateRole(customRole)
}
