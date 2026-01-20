package repository

import (
	"database/sql"
	"todo-api/internal/database"
	"todo-api/internal/models"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (id, username, email, password, is_admin, is_active, role_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)"
	_, err := r.db.Exec(query, user.UserID, user.Username, user.Email, user.Password, user.IsAdmin, user.IsActive, user.RoleID)
	return err
}

// GetUserByID retrieves a user by their ID with role and permission information
func (r *UserRepository) GetUserByID(id interface{}) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			r.role_id::TEXT as role_id_str, r.name, r.description, r.permission_id::TEXT as permission_id_str,
			p.id::TEXT as perm_id, p.name as perm_name, p.description as perm_description,
			p.view, p.create, p.update, p.delete
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		LEFT JOIN permissions p ON r.permission_id = p.id
		WHERE u.id = $1
	`

	var roleID sql.NullString
	var roleIDStr sql.NullString
	var roleName sql.NullString
	var roleDescription sql.NullString
	var permissionIDStr sql.NullString
	var permIDStr sql.NullString
	var permName sql.NullString
	var permDescription sql.NullString
	var permView sql.NullBool
	var permCreate sql.NullBool
	var permUpdate sql.NullBool
	var permDelete sql.NullBool

	err := r.db.QueryRow(query, id).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
		&roleIDStr, &roleName, &roleDescription, &permissionIDStr,
		&permIDStr, &permName, &permDescription, &permView, &permCreate, &permUpdate, &permDelete,
	)
	if err != nil {
		return nil, err
	}

	// Parse role_id from the user table
	if roleID.Valid {
		parsedRoleID, err := uuid.Parse(roleID.String)
		if err == nil {
			user.RoleID = &parsedRoleID
		}
	}

	// Populate role if it exists
	if roleIDStr.Valid && roleName.Valid {
		parsedRoleUUID, err := uuid.Parse(roleIDStr.String)
		if err == nil {
			role := &models.Role{
				RoleId:      parsedRoleUUID,
				Name:        roleName.String,
				Description: roleDescription.String,
			}

			// Parse permission_id if exists
			if permissionIDStr.Valid {
				parsedPermissionID, err := uuid.Parse(permissionIDStr.String)
				if err == nil {
					role.PermissionId = &parsedPermissionID
				}
			}

			// Populate permission object
			if permIDStr.Valid && permName.Valid {
				parsedPermID, err := uuid.Parse(permIDStr.String)
				if err == nil {
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
			}
			user.Role = role
		}
	}

	return user, nil
}

// GetUsersByRoleId retrieves a user by their role ID with role and permission information
func (r *UserRepository) GetUsersByRoleId(id interface{}) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			r.role_id::TEXT as role_id_str, r.name, r.description, r.permission_id::TEXT as permission_id_str,
			p.id::TEXT as perm_id, p.name as perm_name, p.description as perm_description,
			p.view, p.create, p.update, p.delete
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		LEFT JOIN permissions p ON r.permission_id = p.id
		WHERE u.role_id = $1
	`

	var roleID sql.NullString
	var roleIDStr sql.NullString
	var roleName sql.NullString
	var roleDescription sql.NullString
	var permissionIDStr sql.NullString
	var permIDStr sql.NullString
	var permName sql.NullString
	var permDescription sql.NullString
	var permView sql.NullBool
	var permCreate sql.NullBool
	var permUpdate sql.NullBool
	var permDelete sql.NullBool

	err := r.db.QueryRow(query, id).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
		&roleIDStr, &roleName, &roleDescription, &permissionIDStr,
		&permIDStr, &permName, &permDescription, &permView, &permCreate, &permUpdate, &permDelete,
	)
	if err != nil {
		return nil, err
	}

	// Parse role_id from the user table
	if roleID.Valid {
		parsedRoleID, err := uuid.Parse(roleID.String)
		if err == nil {
			user.RoleID = &parsedRoleID
		}
	}

	// Populate role if it exists
	if roleIDStr.Valid && roleName.Valid {
		parsedRoleUUID, err := uuid.Parse(roleIDStr.String)
		if err == nil {
			role := &models.Role{
				RoleId:      parsedRoleUUID,
				Name:        roleName.String,
				Description: roleDescription.String,
			}

			// Parse permission_id if exists
			if permissionIDStr.Valid {
				parsedPermissionID, err := uuid.Parse(permissionIDStr.String)
				if err == nil {
					role.PermissionId = &parsedPermissionID
				}
			}

			// Populate permission object if exists
			if permIDStr.Valid && permName.Valid {
				parsedPermID, err := uuid.Parse(permIDStr.String)
				if err == nil {
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
			}

			user.Role = role
		}
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email with role and permission information
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			r.role_id::TEXT as role_id_str, r.name, r.description, r.permission_id::TEXT as permission_id_str,
			p.id::TEXT as perm_id, p.name as perm_name, p.description as perm_description,
			p.view, p.create, p.update, p.delete
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		LEFT JOIN permissions p ON r.permission_id = p.id
		WHERE u.email = $1
	`

	var roleID sql.NullString
	var roleIDStr sql.NullString
	var roleName sql.NullString
	var roleDescription sql.NullString
	var permissionIDStr sql.NullString
	var permIDStr sql.NullString
	var permName sql.NullString
	var permDescription sql.NullString
	var permView sql.NullBool
	var permCreate sql.NullBool
	var permUpdate sql.NullBool
	var permDelete sql.NullBool

	err := r.db.QueryRow(query, email).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
		&roleIDStr, &roleName, &roleDescription, &permissionIDStr,
		&permIDStr, &permName, &permDescription, &permView, &permCreate, &permUpdate, &permDelete,
	)
	if err != nil {
		return nil, err
	}

	// Parse role_id from the user table
	if roleID.Valid {
		parsedRoleID, err := uuid.Parse(roleID.String)
		if err == nil {
			user.RoleID = &parsedRoleID
		}
	}

	// Populate role if it exists
	if roleIDStr.Valid && roleName.Valid {
		parsedRoleUUID, err := uuid.Parse(roleIDStr.String)
		if err == nil {
			role := &models.Role{
				RoleId:      parsedRoleUUID,
				Name:        roleName.String,
				Description: roleDescription.String,
			}

			// Parse permission_id if exists
			if permissionIDStr.Valid {
				parsedPermissionID, err := uuid.Parse(permissionIDStr.String)
				if err == nil {
					role.PermissionId = &parsedPermissionID
				}
			}

			// Populate permission object if exists
			if permIDStr.Valid && permName.Valid {
				parsedPermID, err := uuid.Parse(permIDStr.String)
				if err == nil {
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
			}

			user.Role = role
		}
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username with role and permission information
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			r.role_id::TEXT as role_id_str, r.name, r.description, r.permission_id::TEXT as permission_id_str,
			p.id::TEXT as perm_id, p.name as perm_name, p.description as perm_description,
			p.view, p.create, p.update, p.delete
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		LEFT JOIN permissions p ON r.permission_id = p.id
		WHERE u.username = $1
	`

	var roleID sql.NullString
	var roleIDStr sql.NullString
	var roleName sql.NullString
	var roleDescription sql.NullString
	var permissionIDStr sql.NullString
	var permIDStr sql.NullString
	var permName sql.NullString
	var permDescription sql.NullString
	var permView sql.NullBool
	var permCreate sql.NullBool
	var permUpdate sql.NullBool
	var permDelete sql.NullBool

	err := r.db.QueryRow(query, username).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
		&roleIDStr, &roleName, &roleDescription, &permissionIDStr,
		&permIDStr, &permName, &permDescription, &permView, &permCreate, &permUpdate, &permDelete,
	)
	if err != nil {
		return nil, err
	}

	// Parse role_id from the user table
	if roleID.Valid {
		parsedRoleID, err := uuid.Parse(roleID.String)
		if err == nil {
			user.RoleID = &parsedRoleID
		}
	}

	// Populate role if it exists
	if roleIDStr.Valid && roleName.Valid {
		parsedRoleUUID, err := uuid.Parse(roleIDStr.String)
		if err == nil {
			role := &models.Role{
				RoleId:      parsedRoleUUID,
				Name:        roleName.String,
				Description: roleDescription.String,
			}

			// Parse permission_id if exists
			if permissionIDStr.Valid {
				parsedPermissionID, err := uuid.Parse(permissionIDStr.String)
				if err == nil {
					role.PermissionId = &parsedPermissionID
				}
			}

			// Populate permission object if exists
			if permIDStr.Valid && permName.Valid {
				parsedPermID, err := uuid.Parse(permIDStr.String)
				if err == nil {
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
			}

			user.Role = role
		}
	}

	return user, nil
}

// GetAllUsers retrieves all users from the database with role and permission information
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			r.role_id::TEXT as role_id_str, r.name, r.description, r.permission_id::TEXT as permission_id_str,
			p.id::TEXT as perm_id, p.name as perm_name, p.description as perm_description,
			p.view, p.create, p.update, p.delete
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		LEFT JOIN permissions p ON r.permission_id = p.id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var roleID sql.NullString
		var roleIDStr sql.NullString
		var roleName sql.NullString
		var roleDescription sql.NullString
		var permissionIDStr sql.NullString
		var permIDStr sql.NullString
		var permName sql.NullString
		var permDescription sql.NullString
		var permView sql.NullBool
		var permCreate sql.NullBool
		var permUpdate sql.NullBool
		var permDelete sql.NullBool

		err := rows.Scan(
			&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
			&roleIDStr, &roleName, &roleDescription, &permissionIDStr,
			&permIDStr, &permName, &permDescription, &permView, &permCreate, &permUpdate, &permDelete,
		)
		if err != nil {
			return nil, err
		}

		// Parse role_id from the user table
		if roleID.Valid {
			parsedRoleID, err := uuid.Parse(roleID.String)
			if err == nil {
				user.RoleID = &parsedRoleID
			}
		}

		// Populate role if it exists
		if roleIDStr.Valid && roleName.Valid {
			parsedRoleUUID, err := uuid.Parse(roleIDStr.String)
			if err == nil {
				role := &models.Role{
					RoleId:      parsedRoleUUID,
					Name:        roleName.String,
					Description: roleDescription.String,
				}

				// Parse permission_id if exists
				if permissionIDStr.Valid {
					parsedPermissionID, err := uuid.Parse(permissionIDStr.String)
					if err == nil {
						role.PermissionId = &parsedPermissionID
					}
				}

				// Populate permission object if exists
				if permIDStr.Valid && permName.Valid {
					parsedPermID, err := uuid.Parse(permIDStr.String)
					if err == nil {
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
				}

				user.Role = role
			}
		}

		users = append(users, user)
	}
	return users, nil
}

// GetUsersPaginated retrieves users with pagination support
func (r *UserRepository) GetUsersPaginated(limit, offset int) ([]models.User, error) {
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			r.role_id::TEXT as role_id_str, r.name, r.description, r.permission_id::TEXT as permission_id_str,
			p.id::TEXT as perm_id, p.name as perm_name, p.description as perm_description,
			p.view, p.create, p.update, p.delete
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		LEFT JOIN permissions p ON r.permission_id = p.id
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var roleID sql.NullString
		var roleIDStr sql.NullString
		var roleName sql.NullString
		var roleDescription sql.NullString
		var permissionIDStr sql.NullString
		var permIDStr sql.NullString
		var permName sql.NullString
		var permDescription sql.NullString
		var permView sql.NullBool
		var permCreate sql.NullBool
		var permUpdate sql.NullBool
		var permDelete sql.NullBool

		err := rows.Scan(
			&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
			&roleIDStr, &roleName, &roleDescription, &permissionIDStr,
			&permIDStr, &permName, &permDescription, &permView, &permCreate, &permUpdate, &permDelete,
		)
		if err != nil {
			return nil, err
		}

		// Parse role_id from the user table
		if roleID.Valid {
			parsedRoleID, err := uuid.Parse(roleID.String)
			if err == nil {
				user.RoleID = &parsedRoleID
			}
		}

		// Populate role if it exists
		if roleIDStr.Valid && roleName.Valid {
			parsedRoleUUID, err := uuid.Parse(roleIDStr.String)
			if err == nil {
				role := &models.Role{
					RoleId:      parsedRoleUUID,
					Name:        roleName.String,
					Description: roleDescription.String,
				}

				// Parse permission_id if exists
				if permissionIDStr.Valid {
					parsedPermissionID, err := uuid.Parse(permissionIDStr.String)
					if err == nil {
						role.PermissionId = &parsedPermissionID
					}
				}

				// Populate permission object if exists
				if permIDStr.Valid && permName.Valid {
					parsedPermID, err := uuid.Parse(permIDStr.String)
					if err == nil {
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
				}

				user.Role = role
			}
		}

		users = append(users, user)
	}
	return users, nil
}

// GetUsersCount returns the total count of users
func (r *UserRepository) GetUsersCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// UpdateUser updates an existing user in the database and returns the updated user
func (r *UserRepository) UpdateUser(user *models.User) (*models.User, error) {
	query := `UPDATE users SET username = $1, email = $2, is_admin = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`
	_, err := r.db.Exec(query, user.Username, user.Email, user.IsAdmin, user.IsActive, user.UserID)
	if err != nil {
		return nil, err
	}

	// Fetch and return the updated user
	return r.GetUserByID(user.UserID)
}

// UpdateUserTokens is deprecated - JWT tokens are stateless and don't need database storage
// Kept for backward compatibility but does nothing
func (r *UserRepository) UpdateUserTokens(username, sessionToken, csrfToken string) error {
	// No-op: JWT tokens are not stored in the database
	return nil
}

// ClearUserTokens is deprecated - JWT tokens are stateless
// Kept for backward compatibility but does nothing
func (r *UserRepository) ClearUserTokens(username string) error {
	// No-op: JWT tokens are not stored in the database
	return nil
}

// DeleteUser deletes a user from the database
func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// UserExists checks if a user with the given username exists
func (r *UserRepository) UserExists(username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&count)
	return count > 0, err
}

// EmailExists checks if a user with the given email exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&count)
	return count > 0, err
}
