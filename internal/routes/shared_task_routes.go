// internal/routes/shared_task_routes.go
package routes

import (
	"net/http"
	"todo-api/internal/handlers"
	"todo-api/internal/models"
)

func RegisterSharedTaskRoutes(sharedTaskHandler *handlers.SharedTaskHandler) {
	http.HandleFunc("/shared-tasks", withAuthAndPermission(sharedTaskHandler.SharedTasksHandler, models.PermContentRead))
	http.HandleFunc("/shared-tasks/", withAuthAndPermission(sharedTaskHandler.SharedTaskByIdHandler, models.PermContentRead))
	http.HandleFunc("/shared-tasks/owner", withAuthAndPermission(sharedTaskHandler.GetSharedTasksByOwnerId, models.PermContentRead))
	http.HandleFunc("/shared-tasks/id", withAuthAndPermission(sharedTaskHandler.GetSharedTasksById, models.PermContentRead))
	http.HandleFunc("/shared-tasks/todo", withAuthAndPermission(sharedTaskHandler.GetSharedTasksByTodoId, models.PermContentRead))
}
