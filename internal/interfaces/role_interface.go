package interfaces

import (
	"todo-api/internal/models"

	"github.com/google/uuid"
)

// RoleRepositoryInterface defines the contract for role data access
type RoleRepositoryInterface interface {
	CreateRole(role *models.Role) error
	GetRoleByID(roleID uuid.UUID) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	GetAllRoles() ([]models.Role, error)
	UpdateRole(role *models.Role) error
	DeleteRole(roleID uuid.UUID) error
	AssignRoleToUser(userID, roleID uuid.UUID) error
	GetUserPermissions(userID uuid.UUID) (*models.Permissions, error)
}
