package services

import (
	"todo-api/internal/models"
	"todo-api/internal/repository"
)

type TodoService struct {
	repo *repository.TodoRepository
}

// NewTodoService creates a new todo service with dependency injection
func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// CreateTodo handles todo creation business logic
func (s *TodoService) CreateTodo(todo *models.Todo) error {
	return s.repo.Create(todo)
}

// GetTodoByID retrieves a todo by ID
func (s *TodoService) GetTodoByID(id string) (*models.Todo, error) {
	return s.repo.GetById(id)
}

// GetTodosByUserID retrieves all todos for a specific user
func (s *TodoService) GetTodosByUserID(userID string) ([]models.Todo, error) {
	return s.repo.GetByUserId(userID)
}

// GetAllTodos retrieves all todos
func (s *TodoService) GetAllTodos() ([]models.Todo, error) {
	return s.repo.GetAll()
}

// UpdateTodo handles todo update business logic
func (s *TodoService) UpdateTodo(id string, todo *models.Todo) (int64, error) {
	return s.repo.Update(id, todo)
}

// DeleteTodo handles todo deletion business logic
func (s *TodoService) DeleteTodo(id string) (int64, error) {
	return s.repo.Delete(id)
}
