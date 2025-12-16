package interfaces

import "todo-api/internal/models"

// SharedTaskInterface defines the business logic contract for shared task operations
type SharedTaskInterface interface {
	CreateSharedTask(sharedTask *models.SharedTask) error
	GetSharedTaskByID(id string) (*models.SharedTask, error)
	GetSharedTasksByOwnerID(ownerID int) ([]models.SharedTask, error)
	GetSharedTasksByTodoID(todoID string) ([]models.SharedTask, error)
	GetAllSharedTasks() ([]models.SharedTask, error)
	UpdateSharedTask(id string, sharedTask *models.SharedTask) (int64, error)
	DeleteSharedTask(id string) (int64, error)
}
