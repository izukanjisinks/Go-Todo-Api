package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todo-api/internal/interfaces"
	"todo-api/internal/models"
	"todo-api/internal/validations"
	"todo-api/pkg/utils"

	"github.com/google/uuid"
)

type UsersHandler struct {
	service interfaces.UserInterface
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUsersHandler(service interfaces.UserInterface) *UsersHandler {
	return &UsersHandler{
		service: service,
	}
}

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

	// Validate registration input
	err = validations.ValidateRegister(username, email, password)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser := &models.User{
		Username: username,
		Email:    email,
		IsAdmin:  isAdmin,
		Login: models.Login{
			Password: password,
		},
	}

	err = h.service.Register(newUser)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
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

	email := req.Email
	password := req.Password

	err = validations.ValidateLogin(email, password)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.Login(email, password)

	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

func (h *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	page := 1
	pageSize := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	paginatedResponse, err := h.service.GetUsersPaginated(page, pageSize)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, paginatedResponse)
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	updates := &models.User{}

	if err := utils.DecodeJson(r, updates); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	fmt.Printf("the updates are %v", updates)

	result, err := h.service.UpdateUser(updates)
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
