package models

import "github.com/google/uuid"

// Permission action constants
const (
	PermView   = "view"
	PermCreate = "create"
	PermUpdate = "update"
	PermDelete = "delete"
)

// Predefined role names
const (
	RoleSuperAdmin = "Super Admin"
	RoleAdmin      = "Admin"
	RoleModerator  = "Moderator"
	RoleUser       = "User"
)

type Role struct {
	RoleId       uuid.UUID    `json:"role_id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	PermissionId *uuid.UUID   `json:"permission_id,omitempty"` // FK to permissions table
	Permission   *Permissions `json:"permission,omitempty"`    // Populated via join
}

type Permissions struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	View        bool      `json:"view"`
	Create      bool      `json:"create"`
	Update      bool      `json:"update"`
	Delete      bool      `json:"delete"`
}

// HasPermission checks if the role's permission allows a specific action
func (r *Role) HasPermission(action string) bool {
	if r.Permission == nil {
		return false
	}

	switch action {
	case PermView:
		return r.Permission.View
	case PermCreate:
		return r.Permission.Create
	case PermUpdate:
		return r.Permission.Update
	case PermDelete:
		return r.Permission.Delete
	default:
		return false
	}
}

// GetPredefinedRoles returns all predefined roles with their permissions
func GetPredefinedRoles() []Role {
	superAdminPermId := uuid.New()
	adminPermId := uuid.New()
	moderatorPermId := uuid.New()
	userPermId := uuid.New()

	return []Role{
		{
			RoleId:       uuid.New(),
			Name:         RoleSuperAdmin,
			Description:  "Full system access with all permissions",
			PermissionId: &superAdminPermId,
			Permission: &Permissions{
				Id:          superAdminPermId,
				Name:        "super_admin_permissions",
				Description: "Full access to all operations",
				View:        true,
				Create:      true,
				Update:      true,
				Delete:      true,
			},
		},
		{
			RoleId:       uuid.New(),
			Name:         RoleAdmin,
			Description:  "Administrative access with most permissions",
			PermissionId: &adminPermId,
			Permission: &Permissions{
				Id:          adminPermId,
				Name:        "admin_permissions",
				Description: "Admin level access",
				View:        true,
				Create:      true,
				Update:      true,
				Delete:      true,
			},
		},
		{
			RoleId:       uuid.New(),
			Name:         RoleModerator,
			Description:  "Content management and user viewing access",
			PermissionId: &moderatorPermId,
			Permission: &Permissions{
				Id:          moderatorPermId,
				Name:        "moderator_permissions",
				Description: "Moderator level access",
				View:        true,
				Create:      true,
				Update:      true,
				Delete:      false,
			},
		},
		{
			RoleId:       uuid.New(),
			Name:         RoleUser,
			Description:  "Basic user access with read permissions",
			PermissionId: &userPermId,
			Permission: &Permissions{
				Id:          userPermId,
				Name:        "user_permissions",
				Description: "Basic user access",
				View:        true,
				Create:      false,
				Update:      false,
				Delete:      false,
			},
		},
	}
}
