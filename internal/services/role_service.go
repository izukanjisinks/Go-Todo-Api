package services

import (
	"fmt"

	"todo-api/internal/interfaces"
	"todo-api/internal/models"

	"github.com/google/uuid"
)

type RoleService struct {
	repo interfaces.RoleRepositoryInterface
}

func NewRoleService(repo interfaces.RoleRepositoryInterface) *RoleService {
	return &RoleService{repo: repo}
}

// CreateRole creates a new role in the database
func (s *RoleService) CreateRole(role *models.Role) error {
	return s.repo.CreateRole(role)
}

// GetRoleByID retrieves a role by its ID
func (s *RoleService) GetRoleByID(roleID uuid.UUID) (*models.Role, error) {
	return s.repo.GetRoleByID(roleID)
}

// GetRoleByName retrieves a role by its name
func (s *RoleService) GetRoleByName(name string) (*models.Role, error) {
	return s.repo.GetRoleByName(name)
}

// GetAllRoles retrieves all roles from the database
func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	return s.repo.GetAllRoles()
}

// UpdateRole updates an existing role
func (s *RoleService) UpdateRole(role *models.Role) error {
	return s.repo.UpdateRole(role)
}

// DeleteRole deletes a role by its ID
func (s *RoleService) DeleteRole(roleID uuid.UUID) error {
	return s.repo.DeleteRole(roleID)
}

// AssignRoleToUser assigns a role to a user
func (s *RoleService) AssignRoleToUser(userID, roleID uuid.UUID) error {
	return s.repo.AssignRoleToUser(userID, roleID)
}

// InitializePredefinedRoles creates the predefined roles if they don't exist
func (s *RoleService) InitializePredefinedRoles() error {
	predefinedRoles := models.GetPredefinedRoles()

	for _, role := range predefinedRoles {
		// Check if role already exists
		existingRole, err := s.GetRoleByName(role.Name)
		if err == nil && existingRole != nil {
			// Role exists, skip
			continue
		}

		// Create the role
		err = s.CreateRole(&role)
		if err != nil {
			return fmt.Errorf("failed to initialize role %s: %w", role.Name, err)
		}
	}

	return nil
}

// CheckPermission checks if a user has a specific permission
func (s *RoleService) CheckPermission(userID uuid.UUID, permission string) (bool, error) {
	permissions, err := s.repo.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if p == permission {
			return true, nil
		}
	}

	return false, nil
}
