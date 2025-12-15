package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/pkg/utils"

	"github.com/google/uuid"
)

type UsersHandler struct {
	repo *repository.UserRepository
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserRepository() *UsersHandler {
	return &UsersHandler{
		repo: repository.NewUserRepository(),
	}
}

var users = map[string]models.User{}

func (h *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		er := http.StatusMethodNotAllowed
		http.Error(w, "Invalid method", er)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	username := req.Username
	email := req.Email
	password := req.Password
	isAdmin := req.IsAdmin

	if username == "" || email == "" || password == "" {
		er := http.StatusBadRequest
		http.Error(w, "Invalid input, username, email and password are required", er)
		return
	}

	// Check if username already exists
	usernameExists, err := h.repo.UserExists(username)
	if err != nil {
		http.Error(w, "Error checking username", http.StatusInternalServerError)
		return
	}
	if usernameExists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Check if email already exists
	emailExists, err := h.repo.EmailExists(email)
	if err != nil {
		http.Error(w, "Error checking email", http.StatusInternalServerError)
		return
	}
	if emailExists {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}
	hashedPassword, _ := utils.HashPassword(password)

	newUser := &models.User{
		UserID:   uuid.New(),
		Username: username,
		Email:    email,
		IsAdmin:  isAdmin,
		Login: models.Login{
			Password: hashedPassword,
		},
	}

	err = h.repo.CreateUser(newUser)

	fmt.Println("error encountered: ", err)

	if err != nil {
		http.Error(w, "Failed to create user in database", http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, newUser)
}

func (h *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		er := http.StatusMethodNotAllowed
		http.Error(w, "Invalid method", er)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	username := req.Username
	password := req.Password

	user, err := h.repo.GetUserByUsername(username)

	if err != nil {
		er := http.StatusUnauthorized
		http.Error(w, "User not found", er)
		return
	}

	if utils.ComparePasswords(user.Password, password) != nil || (user.Password == "" && password != "") {
		er := http.StatusUnauthorized
		http.Error(w, "Invalid user credentails entered", er)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.Email, user.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error generating authentication token")
		return
	}

	// Return JWT token in response
	response := map[string]interface{}{
		"token":    token,
		"user_id":  user.UserID,
		"username": user.Username,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

func (h *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	users, err := h.repo.GetAllUsers()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, users)
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user := &models.User{}

	if err := utils.DecodeJson(r, user); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	result, err := h.repo.UpdateUser(user)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, result)
}

// Logout - With JWT, logout is handled client-side by removing the token
// This endpoint is kept for API consistency but doesn't need to do much
func Logout(w http.ResponseWriter, r *http.Request) {
	// With JWT, the client simply discards the token
	// Optionally, you could implement a token blacklist here for added security
	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully.",
	})
}

// Protected - Example protected endpoint using JWT authentication
// Note: This should be wrapped with the JWTAuth middleware in your router
func Protected(w http.ResponseWriter, r *http.Request) {
	// Extract user info from context (set by JWTAuth middleware)
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	email, _ := r.Context().Value("userEmail").(string)

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Welcome to the protected resource",
		"user_id": userID,
		"email":   email,
	})
}
