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
	query := "INSERT INTO users (username, email, password, is_admin, session_token, csrf_token) VALUES (@p1, @p2, @p3, @p4, @p5, @p6)"
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.IsAdmin, user.SessionToken, user.CSRFToken)
	return err
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, email, password, is_admin, session_token, csrf_token FROM users WHERE id = @p1"
	err := r.db.QueryRow(query, id).Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.SessionToken, &user.CSRFToken)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, is_admin, session_token, csrf_token FROM users WHERE username = @p1`
	err := r.db.QueryRow(query, username).Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.SessionToken, &user.CSRFToken)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllUsers retrieves all users from the database
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, username, email, password, is_admin, session_token, csrf_token FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.SessionToken, &user.CSRFToken)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUser updates an existing user in the database
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `UPDATE users SET username = @p1, email = @p2, password = @p3, is_admin = @p4, session_token = @p5, csrf_token = @p6 WHERE id = @p7`
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.IsAdmin, user.SessionToken, user.CSRFToken, user.UserID)
	return err
}

// UpdateUserTokens updates only the session and CSRF tokens for a user
func (r *UserRepository) UpdateUserTokens(username, sessionToken, csrfToken string) error {
	query := `UPDATE users SET session_token = @p1, csrf_token = @p2 WHERE username = @p3`
	_, err := r.db.Exec(query, sessionToken, csrfToken, username)
	return err
}

// ClearUserTokens clears the session and CSRF tokens for a user (useful for logout)
func (r *UserRepository) ClearUserTokens(username string) error {
	query := `UPDATE users SET session_token = '', csrf_token = '' WHERE username = @p1`
	_, err := r.db.Exec(query, username)
	return err
}

// DeleteUser deletes a user from the database
func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = @p1`
	_, err := r.db.Exec(query, id)
	return err
}

// UserExists checks if a user with the given username exists
func (r *UserRepository) UserExists(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = @p1)`
	err := r.db.QueryRow(query, username).Scan(&exists)
	return exists, err
}
