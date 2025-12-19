// internal/routes/user_routes.go
package routes

import (
	"net/http"
	"strings"
	"todo-api/internal/handlers"
	"todo-api/internal/models"
)

func RegisterUserRoutes(userHandler *handlers.UsersHandler) {
	http.HandleFunc("/users", withAuthAndPermission(userHandler.GetUsers, models.PermUsersRead))
	http.HandleFunc("/users/update", withAuthAndPermission(userHandler.UpdateUser, models.PermUsersUpdate))
	http.HandleFunc("/logout", withAuth(handlers.Logout))
	http.HandleFunc("/protected", withAuth(handlers.Protected))
}

func RegisterRoleRoutes(roleHandler *handlers.RoleHandler, userHandler *handlers.UsersHandler) {
	http.HandleFunc("/roles", withAuthAndPermission(roleHandler.RolesHandler, models.PermRolesManage))
	http.HandleFunc("/roles/", withAuthAndPermission(roleHandler.RolesHandler, models.PermRolesManage))

	// User role assignment routes
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/role") {
			withAuthAndPermission(roleHandler.AssignRoleHandler, models.PermRolesManage)(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/permissions") {
			withAuthAndPermission(roleHandler.GetUserPermissionsHandler, models.PermUsersRead)(w, r)
		} else {
			withAuthAndPermission(userHandler.GetUsers, models.PermUsersRead)(w, r)
		}
	})
}
