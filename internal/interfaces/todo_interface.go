package interfaces

import "todo-api/internal/models"

// TodoInterface defines the business logic contract for todo operations
type TodoInterface interface {
	CreateTodo(todo *models.Todo) error
	GetTodoByID(id string) (*models.Todo, error)
	GetTodosByUserID(userID string) ([]models.Todo, error)
	GetAllTodos() ([]models.Todo, error)
	UpdateTodo(id string, todo *models.Todo) (int64, error)
	DeleteTodo(id string) (int64, error)
}
