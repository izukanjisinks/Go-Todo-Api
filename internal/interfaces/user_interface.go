package interfaces

import (
	"todo-api/internal/models"
)

// UserInterface defines the business logic contract for user operations
type UserInterface interface {
	Register(user *models.User) error
	Login(email, password string) (map[string]interface{}, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id interface{}) (*models.User, error)
	UpdateUser(updates *models.User) (*models.User, error)
	DeleteUser(id int) error
}
