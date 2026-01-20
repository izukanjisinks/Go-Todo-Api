package models

import (
	"time"

	"github.com/google/uuid"
)

type Login struct {
	Password string `json:"-"`
}

type User struct {
	UserID    uuid.UUID  `json:"user_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	IsAdmin   bool       `json:"is_admin"`       // Deprecated: Use RoleID instead
	RoleID    *uuid.UUID `json:"role_id"`        // Foreign key to roles table
	Role      *Role      `json:"role,omitempty"` // Role relationship (populated via join)
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Login
}

// HasPermission checks if the user has a specific permission through their role
func (u *User) HasPermission(permission string) bool {
	if u.Role != nil {
		return u.Role.HasPermission(permission)
	}
	// Fallback to IsAdmin for backward compatibility
	// if u.IsAdmin {
	// 	return true
	// }
	return false
}
