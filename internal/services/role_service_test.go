package services

import (
	"errors"
	"testing"

	"todo-api/internal/models"

	"github.com/google/uuid"
)

// MockRoleRepository is a mock implementation of RoleRepositoryInterface for testing
type MockRoleRepository struct {
	CreateRoleFunc         func(role *models.Role) error
	GetRoleByIDFunc        func(roleID uuid.UUID) (*models.Role, error)
	GetRoleByNameFunc      func(name string) (*models.Role, error)
	GetAllRolesFunc        func() ([]models.Role, error)
	UpdateRoleFunc         func(role *models.Role) error
	DeleteRoleFunc         func(roleID uuid.UUID) error
	AssignRoleToUserFunc   func(userID, roleID uuid.UUID) error
	GetUserPermissionsFunc func(userID uuid.UUID) ([]string, error)
}

func (m *MockRoleRepository) CreateRole(role *models.Role) error {
	if m.CreateRoleFunc != nil {
		return m.CreateRoleFunc(role)
	}
	return nil
}

func (m *MockRoleRepository) GetRoleByID(roleID uuid.UUID) (*models.Role, error) {
	if m.GetRoleByIDFunc != nil {
		return m.GetRoleByIDFunc(roleID)
	}
	return nil, nil
}

func (m *MockRoleRepository) GetRoleByName(name string) (*models.Role, error) {
	if m.GetRoleByNameFunc != nil {
		return m.GetRoleByNameFunc(name)
	}
	return nil, nil
}

func (m *MockRoleRepository) GetAllRoles() ([]models.Role, error) {
	if m.GetAllRolesFunc != nil {
		return m.GetAllRolesFunc()
	}
	return []models.Role{}, nil
}

func (m *MockRoleRepository) UpdateRole(role *models.Role) error {
	if m.UpdateRoleFunc != nil {
		return m.UpdateRoleFunc(role)
	}
	return nil
}

func (m *MockRoleRepository) DeleteRole(roleID uuid.UUID) error {
	if m.DeleteRoleFunc != nil {
		return m.DeleteRoleFunc(roleID)
	}
	return nil
}

func (m *MockRoleRepository) AssignRoleToUser(userID, roleID uuid.UUID) error {
	if m.AssignRoleToUserFunc != nil {
		return m.AssignRoleToUserFunc(userID, roleID)
	}
	return nil
}

func (m *MockRoleRepository) GetUserPermissions(userID uuid.UUID) ([]string, error) {
	if m.GetUserPermissionsFunc != nil {
		return m.GetUserPermissionsFunc(userID)
	}
	return []string{}, nil
}

// Example test using the mock
func TestRoleService_CheckPermission(t *testing.T) {
	// Arrange
	userID := uuid.New()
	mockRepo := &MockRoleRepository{
		GetUserPermissionsFunc: func(uid uuid.UUID) ([]string, error) {
			if uid == userID {
				return []string{"users:read", "users:create"}, nil
			}
			return []string{}, nil
		},
	}
	service := NewRoleService(mockRepo)

	// Act
	hasPermission, err := service.CheckPermission(userID, "users:read")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !hasPermission {
		t.Error("Expected user to have permission")
	}

	// Test permission user doesn't have
	hasPermission, err = service.CheckPermission(userID, "users:delete")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if hasPermission {
		t.Error("Expected user to not have permission")
	}
}

func TestRoleService_InitializePredefinedRoles(t *testing.T) {
	// Arrange
	createdRoles := []string{}
	mockRepo := &MockRoleRepository{
		GetRoleByNameFunc: func(name string) (*models.Role, error) {
			// Simulate role doesn't exist
			return nil, errors.New("role not found")
		},
		CreateRoleFunc: func(role *models.Role) error {
			createdRoles = append(createdRoles, role.Name)
			return nil
		},
	}
	service := NewRoleService(mockRepo)

	// Act
	err := service.InitializePredefinedRoles()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(createdRoles) != 4 {
		t.Errorf("Expected 4 roles to be created, got %d", len(createdRoles))
	}

	// Verify all predefined roles were created
	expectedRoles := []string{"Super Admin", "Admin", "Moderator", "User"}
	for _, expected := range expectedRoles {
		found := false
		for _, created := range createdRoles {
			if created == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected role '%s' to be created", expected)
		}
	}
}

func TestRoleService_CreateRole(t *testing.T) {
	// Arrange
	var capturedRole *models.Role
	mockRepo := &MockRoleRepository{
		CreateRoleFunc: func(role *models.Role) error {
			capturedRole = role
			return nil
		},
	}
	service := NewRoleService(mockRepo)

	testRole := &models.Role{
		Name:        "Test Role",
		Description: "A test role",
		Permissions: []string{"test:read", "test:write"},
	}

	// Act
	err := service.CreateRole(testRole)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if capturedRole == nil {
		t.Fatal("Expected role to be captured")
	}
	if capturedRole.Name != "Test Role" {
		t.Errorf("Expected role name 'Test Role', got '%s'", capturedRole.Name)
	}
}
