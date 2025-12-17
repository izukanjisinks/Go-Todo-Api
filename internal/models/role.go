package models

import "github.com/google/uuid"

// Permission constants
const (
	// User Management
	PermUsersRead   = "users:read"
	PermUsersCreate = "users:create"
	PermUsersUpdate = "users:update"
	PermUsersDelete = "users:delete"

	// Content Management
	PermContentRead   = "content:read"
	PermContentCreate = "content:create"
	PermContentUpdate = "content:update"
	PermContentDelete = "content:delete"

	// System Settings
	PermSettingsRead   = "settings:read"
	PermSettingsUpdate = "settings:update"

	// Reports
	PermReportsView   = "reports:view"
	PermReportsExport = "reports:export"

	// Roles
	PermRolesManage = "roles:manage"
)

// Predefined role names
const (
	RoleSuperAdmin = "Super Admin"
	RoleAdmin      = "Admin"
	RoleModerator  = "Moderator"
	RoleUser       = "User"
)

type Role struct {
	RoleId      uuid.UUID `json:"role_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
}

// HasPermission checks if the role has a specific permission
func (r *Role) HasPermission(permission string) bool {
	for _, p := range r.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// GetPredefinedRoles returns all predefined roles with their permissions
func GetPredefinedRoles() []Role {
	return []Role{
		{
			RoleId:      uuid.New(),
			Name:        RoleSuperAdmin,
			Description: "Full system access with all permissions",
			Permissions: []string{
				PermUsersRead, PermUsersCreate, PermUsersUpdate, PermUsersDelete,
				PermContentRead, PermContentCreate, PermContentUpdate, PermContentDelete,
				PermSettingsRead, PermSettingsUpdate,
				PermReportsView, PermReportsExport,
				PermRolesManage,
			},
		},
		{
			RoleId:      uuid.New(),
			Name:        RoleAdmin,
			Description: "Administrative access with most permissions",
			Permissions: []string{
				PermUsersRead, PermUsersCreate, PermUsersUpdate, PermUsersDelete,
				PermContentRead, PermContentCreate, PermContentUpdate, PermContentDelete,
				PermSettingsRead,
				PermReportsView, PermReportsExport,
			},
		},
		{
			RoleId:      uuid.New(),
			Name:        RoleModerator,
			Description: "Content management and user viewing access",
			Permissions: []string{
				PermUsersRead,
				PermContentRead, PermContentCreate, PermContentUpdate, PermContentDelete,
			},
		},
		{
			RoleId:      uuid.New(),
			Name:        RoleUser,
			Description: "Basic user access with read permissions",
			Permissions: []string{
				PermContentRead,
				PermReportsView,
			},
		},
	}
}
