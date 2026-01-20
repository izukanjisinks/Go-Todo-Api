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
	GetUserPermissionsFunc func(userID uuid.UUID) (*models.Permissions, error)
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

func (m *MockRoleRepository) GetUserPermissions(userID uuid.UUID) (*models.Permissions, error) {
	if m.GetUserPermissionsFunc != nil {
		return m.GetUserPermissionsFunc(userID)
	}
	return &models.Permissions{}, nil
}

// Example test using the mock
func TestRoleService_CheckPermission(t *testing.T) {
	// Arrange
	userID := uuid.New()
	mockRepo := &MockRoleRepository{
		GetUserPermissionsFunc: func(uid uuid.UUID) (*models.Permissions, error) {
			if uid == userID {
				return &models.Permissions{
					Id:          uuid.New(),
					Name:        "test_permissions",
					Description: "Test permissions",
					View:        true,
					Create:      true,
					Update:      false,
					Delete:      false,
				}, nil
			}
			return &models.Permissions{}, nil
		},
	}
	service := NewRoleService(mockRepo)

	// Act - Test permission user has (view)
	hasPermission, err := service.CheckPermission(userID, "view")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !hasPermission {
		t.Error("Expected user to have view permission")
	}

	// Test permission user has (create)
	hasPermission, err = service.CheckPermission(userID, "create")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !hasPermission {
		t.Error("Expected user to have create permission")
	}

	// Test permission user doesn't have (delete)
	hasPermission, err = service.CheckPermission(userID, "delete")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if hasPermission {
		t.Error("Expected user to not have delete permission")
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
	permissionID := uuid.New()
	mockRepo := &MockRoleRepository{
		CreateRoleFunc: func(role *models.Role) error {
			capturedRole = role
			return nil
		},
	}
	service := NewRoleService(mockRepo)

	testRole := &models.Role{
		Name:         "Test Role",
		Description:  "A test role",
		PermissionId: &permissionID,
		Permission: &models.Permissions{
			Id:          permissionID,
			Name:        "test_permissions",
			Description: "Test permissions",
			View:        true,
			Create:      true,
			Update:      false,
			Delete:      false,
		},
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
	if capturedRole.PermissionId == nil {
		t.Error("Expected role to have a permission ID")
	}
	if capturedRole.Permission == nil {
		t.Error("Expected role to have permissions")
	}
	if capturedRole.Permission != nil && !capturedRole.Permission.View {
		t.Error("Expected role to have view permission")
	}
}

func TestRoleService_GetRoleByID(t *testing.T) {
	// Arrange
	roleID := uuid.New()
	permissionID := uuid.New()
	expectedRole := &models.Role{
		RoleId:       roleID,
		Name:         "Test Role",
		Description:  "A test role",
		PermissionId: &permissionID,
		Permission: &models.Permissions{
			Id:          permissionID,
			Name:        "test_permissions",
			Description: "Test permissions",
			View:        true,
			Create:      true,
			Update:      true,
			Delete:      false,
		},
	}

	mockRepo := &MockRoleRepository{
		GetRoleByIDFunc: func(id uuid.UUID) (*models.Role, error) {
			if id == roleID {
				return expectedRole, nil
			}
			return nil, errors.New("role not found")
		},
	}
	service := NewRoleService(mockRepo)

	// Act
	role, err := service.GetRoleByID(roleID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if role == nil {
		t.Fatal("Expected role to be returned")
	}
	if role.Name != "Test Role" {
		t.Errorf("Expected role name 'Test Role', got '%s'", role.Name)
	}
	if role.Permission == nil {
		t.Fatal("Expected role to have permissions")
	}
	if !role.Permission.View {
		t.Error("Expected role to have view permission")
	}
	if role.Permission.Delete {
		t.Error("Expected role to not have delete permission")
	}
}

func TestRoleService_UpdateRole(t *testing.T) {
	// Arrange
	roleID := uuid.New()
	permissionID := uuid.New()
	var updatedRole *models.Role

	mockRepo := &MockRoleRepository{
		UpdateRoleFunc: func(role *models.Role) error {
			updatedRole = role
			return nil
		},
	}
	service := NewRoleService(mockRepo)

	testRole := &models.Role{
		RoleId:       roleID,
		Name:         "Updated Role",
		Description:  "An updated role",
		PermissionId: &permissionID,
		Permission: &models.Permissions{
			Id:          permissionID,
			Name:        "updated_permissions",
			Description: "Updated permissions",
			View:        true,
			Create:      false,
			Update:      true,
			Delete:      false,
		},
	}

	// Act
	err := service.UpdateRole(testRole)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if updatedRole == nil {
		t.Fatal("Expected role to be captured")
	}
	if updatedRole.Name != "Updated Role" {
		t.Errorf("Expected role name 'Updated Role', got '%s'", updatedRole.Name)
	}
	if updatedRole.Permission == nil {
		t.Fatal("Expected role to have permissions")
	}
	if updatedRole.Permission.Create {
		t.Error("Expected role to not have create permission")
	}
}
