// internal/routes/user_routes.go
package routes

import (
	"net/http"
	"strings"
	"todo-api/internal/handlers"
	"todo-api/internal/models"
)

func RegisterUserRoutes(userHandler *handlers.UsersHandler) {
	http.HandleFunc("/register", withAuthAndPermission(userHandler.Register, models.PermCreate))
	http.HandleFunc("/users", withAuthAndPermission(userHandler.GetUsers, models.PermView))
	http.HandleFunc("/users/update", withAuthAndPermission(userHandler.UpdateUser, models.PermUpdate))
	http.HandleFunc("/logout", withAuth(handlers.Logout))
	http.HandleFunc("/protected", withAuth(handlers.Protected))
}

func RegisterRoleRoutes(roleHandler *handlers.RoleHandler, userHandler *handlers.UsersHandler) {
	http.HandleFunc("/roles", withAuthAndPermission(roleHandler.RolesHandler, models.PermView))
	http.HandleFunc("/roles/", withAuthAndPermission(roleHandler.RolesHandler, models.PermCreate))

	// User role assignment routes
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/role") {
			withAuthAndPermission(roleHandler.AssignRoleHandler, models.PermUpdate)(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/permissions") {
			withAuthAndPermission(roleHandler.GetUserPermissionsHandler, models.PermView)(w, r)
		} else {
			withAuthAndPermission(userHandler.GetUsers, models.PermView)(w, r)
		}
	})
}
