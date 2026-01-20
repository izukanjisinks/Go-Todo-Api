package repository

import (
	"database/sql"
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

// GetPermissionByName retrieves a permission by its name
func (r *RoleRepository) GetPermissionByName(name string) (*models.Permissions, error) {
	query := `
		SELECT id::TEXT as id, name, description, view, create, update, delete
		FROM permissions
		WHERE name = $1
	`

	var permission models.Permissions
	var idStr string

	err := r.db.QueryRow(query, name).Scan(&idStr, &permission.Name, &permission.Description, &permission.View, &permission.Create, &permission.Update, &permission.Delete)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse permission id: %w", err)
	}
	permission.Id = parsedID

	return &permission, nil
}

// CreatePermission creates a new permission in the database
func (r *RoleRepository) CreatePermission(perm *models.Permissions) error {
	// Check if permission already exists
	existing, err := r.GetPermissionByName(perm.Name)
	if err != nil {
		return err
	}
	if existing != nil {
		// Permission already exists, use its ID
		perm.Id = existing.Id
		return nil
	}

	if perm.Id == uuid.Nil {
		perm.Id = uuid.New()
	}

	query := `
		INSERT INTO permissions (id, name, description, view, create, update, delete)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.Exec(query, perm.Id, perm.Name, perm.Description, perm.View, perm.Create, perm.Update, perm.Delete)
	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	return nil
}

// CreateRole creates a new role in the database
func (r *RoleRepository) CreateRole(role *models.Role) error {
	// Check if role already exists by name
	existing, err := r.GetRoleByName(role.Name)
	if existing != nil {
		// Role already exists, just return without error
		role.RoleId = existing.RoleId
		return nil
	}
	if err != nil {
		return fmt.Errorf("error checking existing role: %w", err)
	}

	if role.RoleId == uuid.Nil {
		role.RoleId = uuid.New()
	}

	// If role has a Permission object, create the permission first
	if role.Permission != nil {
		err := r.CreatePermission(role.Permission)
		if err != nil {
			return fmt.Errorf("failed to create permission: %w", err)
		}
		// Ensure PermissionId is set to the Permission's ID
		role.PermissionId = &role.Permission.Id
	}

	query := `
		INSERT INTO roles (role_id, name, description, permission_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err = r.db.Exec(query, role.RoleId, role.Name, role.Description, role.PermissionId)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// GetRoleByID retrieves a role by its ID with permission details
func (r *RoleRepository) GetRoleByID(roleID uuid.UUID) (*models.Role, error) {
	query := `
		SELECT
			r.role_id::TEXT as role_id,
			r.name,
			r.description,
			r.permission_id::TEXT as permission_id,
			p.id::TEXT as perm_id,
			p.name as perm_name,
			p.description as perm_description,
			p.view,
			p.create,
			p.update,
			p.delete
		FROM roles r
		LEFT JOIN permissions p ON r.permission_id = p.id
		WHERE r.role_id = $1
	`

	role := &models.Role{}
	var roleIDStr string
	var permissionIDStr sql.NullString
	var permIDStr sql.NullString
	var permName sql.NullString
	var permDescription sql.NullString
	var permView sql.NullBool
	var permCreate sql.NullBool
	var permUpdate sql.NullBool
	var permDelete sql.NullBool

	err := r.db.QueryRow(query, roleID).Scan(
		&roleIDStr,
		&role.Name,
		&role.Description,
		&permissionIDStr,
		&permIDStr,
		&permName,
		&permDescription,
		&permView,
		&permCreate,
		&permUpdate,
		&permDelete,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Parse the role UUID string
	parsedRoleUUID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse role_id: %w", err)
	}
	role.RoleId = parsedRoleUUID

	// Parse permission_id if exists
	if permissionIDStr.Valid {
		parsedPermissionID, err := uuid.Parse(permissionIDStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse permission_id: %w", err)
		}
		role.PermissionId = &parsedPermissionID
	}

	// Populate permission object if exists
	if permIDStr.Valid && permName.Valid {
		parsedPermID, err := uuid.Parse(permIDStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse permission id: %w", err)
		}

		role.Permission = &models.Permissions{
			Id:          parsedPermID,
			Name:        permName.String,
			Description: permDescription.String,
			View:        permView.Bool,
			Create:      permCreate.Bool,
			Update:      permUpdate.Bool,
			Delete:      permDelete.Bool,
		}
	}

	return role, nil
}

// GetRoleByName retrieves a role by its name with permission details
func (r *RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	query := `
		SELECT
			r.role_id::TEXT as role_id,
			r.name,
			r.description,
			r.permission_id::TEXT as permission_id,
			p.id::TEXT as perm_id,
			p.name as perm_name,
			p.description as perm_description,
			p.view,
			p.create,
			p.update,
			p.delete
		FROM roles r
		LEFT JOIN permissions p ON r.permission_id = p.id
		WHERE r.name = $1
	`

	role := &models.Role{}
	var roleIDStr string
	var permissionIDStr sql.NullString
	var permIDStr sql.NullString
	var permName sql.NullString
	var permDescription sql.NullString
	var permView sql.NullBool
	var permCreate sql.NullBool
	var permUpdate sql.NullBool
	var permDelete sql.NullBool

	err := r.db.QueryRow(query, name).Scan(
		&roleIDStr,
		&role.Name,
		&role.Description,
		&permissionIDStr,
		&permIDStr,
		&permName,
		&permDescription,
		&permView,
		&permCreate,
		&permUpdate,
		&permDelete,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil when role not found (not an error)
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Parse the role UUID string
	parsedRoleUUID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse role_id: %w", err)
	}
	role.RoleId = parsedRoleUUID

	// Parse permission_id if exists
	if permissionIDStr.Valid {
		parsedPermissionID, err := uuid.Parse(permissionIDStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse permission_id: %w", err)
		}
		role.PermissionId = &parsedPermissionID
	}

	// Populate permission object if exists
	if permIDStr.Valid && permName.Valid {
		parsedPermID, err := uuid.Parse(permIDStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse permission id: %w", err)
		}

		role.Permission = &models.Permissions{
			Id:          parsedPermID,
			Name:        permName.String,
			Description: permDescription.String,
			View:        permView.Bool,
			Create:      permCreate.Bool,
			Update:      permUpdate.Bool,
			Delete:      permDelete.Bool,
		}
	}

	return role, nil
}

// GetAllRoles retrieves all roles from the database with their permissions
func (r *RoleRepository) GetAllRoles() ([]models.Role, error) {
	query := `
		SELECT
			r.role_id::TEXT as role_id,
			r.name,
			r.description,
			r.permission_id::TEXT as permission_id,
			p.id::TEXT as perm_id,
			p.name as perm_name,
			p.description as perm_description,
			p.view,
			p.create,
			p.update,
			p.delete
		FROM roles r
		LEFT JOIN permissions p ON r.permission_id = p.id
		ORDER BY r.name
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
		var permissionIDStr sql.NullString
		var permIDStr sql.NullString
		var permName sql.NullString
		var permDescription sql.NullString
		var permView sql.NullBool
		var permCreate sql.NullBool
		var permUpdate sql.NullBool
		var permDelete sql.NullBool

		err := rows.Scan(
			&roleIDStr,
			&role.Name,
			&role.Description,
			&permissionIDStr,
			&permIDStr,
			&permName,
			&permDescription,
			&permView,
			&permCreate,
			&permUpdate,
			&permDelete,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}

		// Parse the role UUID string
		parsedRoleUUID, err := uuid.Parse(roleIDStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse role_id: %w", err)
		}
		role.RoleId = parsedRoleUUID

		// Parse permission_id if exists
		if permissionIDStr.Valid {
			parsedPermissionID, err := uuid.Parse(permissionIDStr.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse permission_id: %w", err)
			}
			role.PermissionId = &parsedPermissionID
		}

		// Populate permission object if exists
		if permIDStr.Valid && permName.Valid {
			parsedPermID, err := uuid.Parse(permIDStr.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse permission id: %w", err)
			}

			role.Permission = &models.Permissions{
				Id:          parsedPermID,
				Name:        permName.String,
				Description: permDescription.String,
				View:        permView.Bool,
				Create:      permCreate.Bool,
				Update:      permUpdate.Bool,
				Delete:      permDelete.Bool,
			}
		}

		roles = append(roles, role)
	}

	return roles, nil
}

// UpdateRole updates an existing role
func (r *RoleRepository) UpdateRole(role *models.Role) error {
	query := `
		UPDATE roles
		SET name = $2, description = $3, permission_id = $4
		WHERE role_id = $1
	`

	result, err := r.db.Exec(query, role.RoleId, role.Name, role.Description, role.PermissionId)
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
	query := `DELETE FROM roles WHERE role_id = $1`

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
		SET role_id = $2
		WHERE id = $1
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

// GetUserPermissions retrieves permissions for a user by joining with roles and permissions
func (r *RoleRepository) GetUserPermissions(userID uuid.UUID) (*models.Permissions, error) {
	query := `
		SELECT
			p.id::TEXT as perm_id,
			p.name,
			p.description,
			p.view,
			p.create,
			p.update,
			p.delete
		FROM users u
		JOIN roles r ON u.role_id = r.role_id
		JOIN permissions p ON r.permission_id = p.id
		WHERE u.id = $1
	`

	var permIDStr string
	var permission models.Permissions

	err := r.db.QueryRow(query, userID).Scan(
		&permIDStr,
		&permission.Name,
		&permission.Description,
		&permission.View,
		&permission.Create,
		&permission.Update,
		&permission.Delete,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user or permissions not found")
		}
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Parse permission ID
	parsedPermID, err := uuid.Parse(permIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse permission id: %w", err)
	}
	permission.Id = parsedPermID

	return &permission, nil
}
