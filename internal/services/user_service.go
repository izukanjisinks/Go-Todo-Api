package services

import (
	"errors"
	"fmt"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/pkg/utils"

	"github.com/google/uuid"
)

type UserService struct {
	repo     *repository.UserRepository
	roleRepo *repository.RoleRepository
}

// NewUserService creates a new user service with dependency injection
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo:     repo,
		roleRepo: repository.NewRoleRepository(),
	}
}

// Register handles user registration business logic
func (s *UserService) Register(user *models.User) error {
	// Check if username already exists
	exists, err := s.repo.UserExists(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already exists")
	}

	// Check if email already exists
	emailExists, err := s.repo.EmailExists(user.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Generate UUID for the user
	user.UserID = uuid.New()

	// Set user as active by default
	user.IsActive = true

	// Auto-assign default "User" role if no role is set
	if user.RoleID == nil {
		defaultRole, err := s.roleRepo.GetRoleByName(models.RoleUser)
		if err == nil && defaultRole != nil {
			user.RoleID = &defaultRole.RoleId
		}
		// If default role doesn't exist, user will be created without a role
		// They can be assigned a role later by an admin
	}

	// Create the user
	return s.repo.CreateUser(user)
}

// Login handles user authentication business logic
func (s *UserService) Login(email, password string) (map[string]interface{}, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user account is active
	if !user.IsActive {
		return nil, errors.New("account is inactive, please contact administrator")
	}

	// Verify password
	if utils.ComparePasswords(user.Password, password) != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.Email, user.UserID)
	if err != nil {
		return nil, errors.New("error generating authentication token")
	}

	// Return response
	response := map[string]interface{}{
		"token":    token,
		"user_id":  user.UserID,
		"username": user.Username,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
	}

	return response, nil
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

// GetUsersPaginated retrieves users with pagination
func (s *UserService) GetUsersPaginated(page, pageSize int) (*models.PaginatedResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	totalCount, err := s.repo.GetUsersCount()
	if err != nil {
		return nil, err
	}

	// Get paginated users
	users, err := s.repo.GetUsersPaginated(pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Create paginated response
	return models.NewPaginatedResponse(users, page, pageSize, totalCount), nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id interface{}) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

// GetUserByRoleId retrieves a user by role ID
func (s *UserService) GetUserByRoleId(roleId uuid.UUID) (*models.User, error) {
	return s.repo.GetUsersByRoleId(roleId)
}

// UpdateUser handles user update business logic with validation
func (s *UserService) UpdateUser(updates *models.User) (*models.User, error) {

	// Validate email uniqueness if being changed
	if updates.Email != "" {
		existingUser, err := s.repo.GetUserByEmail(updates.Email)
		fmt.Println("existing user", existingUser)
		if err == nil && existingUser.UserID != updates.UserID {
			return nil, errors.New("email already in use by another user")
		}
	}

	// Validate username uniqueness if being changed
	if updates.Username != "" {
		existingUser, err := s.repo.GetUserByUsername(updates.Username)
		if err == nil && existingUser.UserID != updates.UserID {
			return nil, errors.New("username already in use by another user")
		}
	}

	// Update the user
	return s.repo.UpdateUser(updates)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
