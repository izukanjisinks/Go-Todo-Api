package repository

import (
	"database/sql"
	"todo-api/internal/database"
	"todo-api/internal/models"
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
	query := "INSERT INTO users (id, username, email, password, is_admin) VALUES (@p1, @p2, @p3, @p4, @p5)"
	_, err := r.db.Exec(query, user.UserID, user.Username, user.Email, user.Password, user.IsAdmin)
	return err
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, email, password, is_admin FROM users WHERE id = @p1"
	err := r.db.QueryRow(query, id).Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, is_admin FROM users WHERE username = @p1`
	err := r.db.QueryRow(query, username).Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllUsers retrieves all users from the database
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, username, email, password, is_admin FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUser updates an existing user in the database
func (r *UserRepository) UpdateUser(user *models.User) (result sql.Result, err error) {
	query := `UPDATE users SET username = @p1, email = @p2, password = @p3, is_admin = @p4 WHERE id = @p5`
	result, err = r.db.Exec(query, user.Username, user.Email, user.Password, user.IsAdmin, user.UserID)
	return result, err
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
