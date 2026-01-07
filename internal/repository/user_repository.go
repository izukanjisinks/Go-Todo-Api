package repository

import (
	"database/sql"
	"encoding/json"
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
	query := "INSERT INTO users (id, username, email, password, is_admin, is_active, role_id, created_at, updated_at) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, GETDATE(), GETDATE())"
	_, err := r.db.Exec(query, user.UserID, user.Username, user.Email, user.Password, user.IsAdmin, user.IsActive, user.RoleID)
	return err
}

// GetUserByID retrieves a user by their ID with role information
func (r *UserRepository) GetUserByID(id interface{}) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			CONVERT(VARCHAR(36), r.role_id) as role_id_str, r.name, r.description, r.permissions
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		WHERE u.id = @p1
	`

	var roleID sql.NullString
	var roleIDStr sql.NullString
	var roleName sql.NullString
	var roleDescription sql.NullString
	var permissionsJSON sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
		&roleIDStr, &roleName, &roleDescription, &permissionsJSON,
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
		parsedUUID, err := uuid.Parse(roleIDStr.String)
		if err == nil {
			role := &models.Role{
				RoleId:      parsedUUID,
				Name:        roleName.String,
				Description: roleDescription.String,
			}

			if permissionsJSON.Valid {
				var permissions []string
				if err := json.Unmarshal([]byte(permissionsJSON.String), &permissions); err == nil {
					role.Permissions = permissions
				}
			}

			user.Role = role
		}
	}

	return user, nil
}

// GetUsersByRoleId retrieves a user by their role ID with role information
func (r *UserRepository) GetUsersByRoleId(id interface{}) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			CONVERT(VARCHAR(36), r.role_id) as role_id_str, r.name, r.description, r.permissions
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		WHERE u.role_id = @p1
	`

	var roleID sql.NullString
	var roleIDStr sql.NullString
	var roleName sql.NullString
	var roleDescription sql.NullString
	var permissionsJSON sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
		&roleIDStr, &roleName, &roleDescription, &permissionsJSON,
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
		parsedUUID, err := uuid.Parse(roleIDStr.String)
		if err == nil {
			role := &models.Role{
				RoleId:      parsedUUID,
				Name:        roleName.String,
				Description: roleDescription.String,
			}

			if permissionsJSON.Valid {
				var permissions []string
				if err := json.Unmarshal([]byte(permissionsJSON.String), &permissions); err == nil {
					role.Permissions = permissions
				}
			}

			user.Role = role
		}
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email with role information
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			CONVERT(VARCHAR(36), r.role_id) as role_id_str, r.name, r.description, r.permissions
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		WHERE u.email = @p1
	`

	var roleID sql.NullString
	var roleIDStr sql.NullString
	var roleName sql.NullString
	var roleDescription sql.NullString
	var permissionsJSON sql.NullString

	err := r.db.QueryRow(query, email).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
		&roleIDStr, &roleName, &roleDescription, &permissionsJSON,
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
		parsedUUID, err := uuid.Parse(roleIDStr.String)
		if err == nil {
			role := &models.Role{
				RoleId:      parsedUUID,
				Name:        roleName.String,
				Description: roleDescription.String,
			}

			if permissionsJSON.Valid {
				var permissions []string
				if err := json.Unmarshal([]byte(permissionsJSON.String), &permissions); err == nil {
					role.Permissions = permissions
				}
			}

			user.Role = role
		}
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username with role information
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			CONVERT(VARCHAR(36), r.role_id) as role_id_str, r.name, r.description, r.permissions
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		WHERE u.username = @p1
	`

	var roleID sql.NullString
	var roleIDStr sql.NullString
	var roleName sql.NullString
	var roleDescription sql.NullString
	var permissionsJSON sql.NullString

	err := r.db.QueryRow(query, username).Scan(
		&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
		&roleIDStr, &roleName, &roleDescription, &permissionsJSON,
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
		parsedUUID, err := uuid.Parse(roleIDStr.String)
		if err == nil {
			role := &models.Role{
				RoleId:      parsedUUID,
				Name:        roleName.String,
				Description: roleDescription.String,
			}

			if permissionsJSON.Valid {
				var permissions []string
				if err := json.Unmarshal([]byte(permissionsJSON.String), &permissions); err == nil {
					role.Permissions = permissions
				}
			}

			user.Role = role
		}
	}

	return user, nil
}

// GetAllUsers retrieves all users from the database with role information
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `
		SELECT
			u.id, u.username, u.email, u.password, u.is_admin, u.is_active, u.created_at, u.updated_at, u.role_id,
			CONVERT(VARCHAR(36), r.role_id) as role_id_str, r.name, r.description, r.permissions
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
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
		var permissionsJSON sql.NullString

		err := rows.Scan(
			&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &roleID,
			&roleIDStr, &roleName, &roleDescription, &permissionsJSON,
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
			parsedUUID, err := uuid.Parse(roleIDStr.String)
			if err == nil {
				role := &models.Role{
					RoleId:      parsedUUID,
					Name:        roleName.String,
					Description: roleDescription.String,
				}

				if permissionsJSON.Valid {
					var permissions []string
					if err := json.Unmarshal([]byte(permissionsJSON.String), &permissions); err == nil {
						role.Permissions = permissions
					}
				}

				user.Role = role
			}
		}

		users = append(users, user)
	}
	return users, nil
}

// UpdateUser updates an existing user in the database and returns the updated user
func (r *UserRepository) UpdateUser(user *models.User) (*models.User, error) {
	query := `UPDATE users SET username = @p1, email = @p2, is_admin = @p3, is_active = @p4, updated_at = GETDATE() WHERE id = @p5`
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
	query := `DELETE FROM users WHERE id = @p1`
	_, err := r.db.Exec(query, id)
	return err
}

// UserExists checks if a user with the given username exists
func (r *UserRepository) UserExists(username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = @p1`
	err := r.db.QueryRow(query, username).Scan(&count)
	return count > 0, err
}

// EmailExists checks if a user with the given email exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = @p1`
	err := r.db.QueryRow(query, email).Scan(&count)
	return count > 0, err
}
