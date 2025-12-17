package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"todo-api/internal/database"
	"todo-api/internal/models"

	"github.com/google/uuid"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		db: database.DB,
	}
}

// CreateRole creates a new role in the database
func (r *RoleRepository) CreateRole(role *models.Role) error {
	if role.RoleId == uuid.Nil {
		role.RoleId = uuid.New()
	}

	query := `
		INSERT INTO roles (role_id, name, description, permissions)
		VALUES (@p1, @p2, @p3, @p4)
	`

	permissionsJSON, err := json.Marshal(role.Permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %w", err)
	}

	_, err = r.db.Exec(query, role.RoleId, role.Name, role.Description, string(permissionsJSON))
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// GetRoleByID retrieves a role by its ID
func (r *RoleRepository) GetRoleByID(roleID uuid.UUID) (*models.Role, error) {
	query := `
		SELECT CONVERT(VARCHAR(36), role_id) as role_id, name, description, permissions
		FROM roles
		WHERE role_id = @p1
	`

	role := &models.Role{}
	var roleIDStr string
	var permissionsJSON string
	err := r.db.QueryRow(query, roleID).Scan(
		&roleIDStr,
		&role.Name,
		&role.Description,
		&permissionsJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Parse the UUID string
	parsedUUID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse role_id: %w", err)
	}
	role.RoleId = parsedUUID

	err = json.Unmarshal([]byte(permissionsJSON), &role.Permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal permissions: %w", err)
	}

	return role, nil
}

// GetRoleByName retrieves a role by its name
func (r *RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	query := `
		SELECT CONVERT(VARCHAR(36), role_id) as role_id, name, description, permissions
		FROM roles
		WHERE name = @p1
	`

	role := &models.Role{}
	var roleIDStr string
	var permissionsJSON string
	err := r.db.QueryRow(query, name).Scan(
		&roleIDStr,
		&role.Name,
		&role.Description,
		&permissionsJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Parse the UUID string
	parsedUUID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse role_id: %w", err)
	}
	role.RoleId = parsedUUID

	err = json.Unmarshal([]byte(permissionsJSON), &role.Permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal permissions: %w", err)
	}

	return role, nil
}

// GetAllRoles retrieves all roles from the database
func (r *RoleRepository) GetAllRoles() ([]models.Role, error) {
	query := `
		SELECT CONVERT(VARCHAR(36), role_id) as role_id, name, description, permissions
		FROM roles
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		var roleIDStr string
		var permissionsJSON string
		err := rows.Scan(
			&roleIDStr,
			&role.Name,
			&role.Description,
			&permissionsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}

		// Parse the UUID string
		parsedUUID, err := uuid.Parse(roleIDStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse role_id: %w", err)
		}
		role.RoleId = parsedUUID

		err = json.Unmarshal([]byte(permissionsJSON), &role.Permissions)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal permissions: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// UpdateRole updates an existing role
func (r *RoleRepository) UpdateRole(role *models.Role) error {
	query := `
		UPDATE roles
		SET name = @p2, description = @p3, permissions = @p4
		WHERE role_id = @p1
	`

	permissionsJSON, err := json.Marshal(role.Permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %w", err)
	}

	result, err := r.db.Exec(query, role.RoleId, role.Name, role.Description, string(permissionsJSON))
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("role not found")
	}

	return nil
}

// DeleteRole deletes a role by its ID
func (r *RoleRepository) DeleteRole(roleID uuid.UUID) error {
	query := `DELETE FROM roles WHERE role_id = @p1`

	result, err := r.db.Exec(query, roleID)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("role not found")
	}

	return nil
}

// AssignRoleToUser assigns a role to a user
func (r *RoleRepository) AssignRoleToUser(userID, roleID uuid.UUID) error {
	query := `
		UPDATE users
		SET role_id = @p2
		WHERE id = @p1
	`

	result, err := r.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetUserPermissions retrieves permissions for a user by joining with roles
func (r *RoleRepository) GetUserPermissions(userID uuid.UUID) ([]string, error) {
	query := `
		SELECT r.permissions
		FROM users u
		JOIN roles r ON u.role_id = r.role_id
		WHERE u.id = @p1
	`

	var permissionsJSON string
	err := r.db.QueryRow(query, userID).Scan(&permissionsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	var permissions []string
	err = json.Unmarshal([]byte(permissionsJSON), &permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal permissions: %w", err)
	}

	return permissions, nil
}
