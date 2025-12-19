// internal/routes/todo_routes.go
package routes

import (
	"net/http"
	"todo-api/internal/handlers"
	"todo-api/internal/models"
)

func RegisterTodoRoutes(todoHandler *handlers.TodoHandler) {
	http.HandleFunc("/todos", withAuthAndPermission(todoHandler.TodosHandler, models.PermContentRead))
	http.HandleFunc("/todos/", withAuthAndPermission(todoHandler.TodoByIdHandler, models.PermContentRead))
	http.HandleFunc("/todos/user", withAuthAndPermission(todoHandler.GetTodosByUserId, models.PermContentRead))
}
