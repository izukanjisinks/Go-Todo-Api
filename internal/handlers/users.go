package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"todo-api/internal/models"
	"todo-api/internal/repository"
	"todo-api/pkg/utils"
)

type UsersHandler struct {
	repo *repository.UserRepository
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	password := req.Password

	if username == "" || password == "" {
		er := http.StatusBadRequest
		http.Error(w, "Invalid input", er)
		return
	}

	if _, ok := users[username]; ok {
		er := http.StatusConflict
		http.Error(w, "Username already exists", er)
		return
	}
	hashedPassword, _ := utils.HashPassword(password)

	newUser := &models.User{
		Username: username,
		Login: models.Login{
			HashedPassword: hashedPassword,
			SessionToken:   "",
			CSRFToken:      "",
		},
	}

	err = h.repo.CreateUser(newUser)

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

	fmt.Println(username, password)

	user, err := h.repo.GetUserByUsername(username)

	if err != nil {
		er := http.StatusUnauthorized
		http.Error(w, "User not found", er)
		return
	}

	if utils.ComparePasswords(user.HashedPassword, password) != nil || (user.HashedPassword == "" && password != "") {
		er := http.StatusUnauthorized
		http.Error(w, "Invalid user credentails entered", er)
		return
	}

	sessionToken := utils.GenerateSessionToken(32)
	csrfToken := utils.GenerateSessionToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	//store tokens
	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken

	err = h.repo.UpdateUserTokens(username, sessionToken, csrfToken)

	if err != nil {
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
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

func Logout(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	
	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: false,
	})

	// Clear tokens from database
	if username != "" {
		userHandler := NewUserRepository()
		userHandler.repo.UpdateUserTokens(username, "", "")
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

func Protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Invalid method")
		return
	}

	username := r.FormValue("username")

	// Get session token from cookie
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Verify with database
	userHandler := NewUserRepository()
	user, err := userHandler.repo.GetUserByUsername(username)
	if err != nil || user.SessionToken != st.Value {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Verify CSRF token if provided
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf != "" && csrf != user.CSRFToken {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": username + " welcome to the protected resource"})
}
