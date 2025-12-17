package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"todo-api/internal/models"
	"todo-api/internal/services"

	"github.com/google/uuid"
)

type RoleHandler struct {
	roleService *services.RoleService
}

func NewRoleHandler(roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func (h *RoleHandler) RolesHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/roles")
	path = strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		if path == "" {
			h.GetAllRoles(w, r)
		} else {
			h.GetRoleByID(w, r, path)
		}
	case http.MethodPost:
		if path == "" {
			h.CreateRole(w, r)
		}
	case http.MethodPut:
		if path != "" {
			h.UpdateRole(w, r, path)
		}
	case http.MethodDelete:
		if path != "" {
			h.DeleteRole(w, r, path)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request, roleIDStr string) {
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	role, err := h.roleService.GetRoleByID(roleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if role.Name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}

	if len(role.Permissions) == 0 {
		http.Error(w, "At least one permission is required", http.StatusBadRequest)
		return
	}

	if err := h.roleService.CreateRole(&role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request, roleIDStr string) {
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role.RoleId = roleID

	if err := h.roleService.UpdateRole(&role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request, roleIDStr string) {
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	if err := h.roleService.DeleteRole(roleID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UserRoleHandler handles /users/{userId}/role
func (h *RoleHandler) UserRoleHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	userIDStr := parts[0]
	if r.Method == http.MethodPost && parts[1] == "role" {
		h.AssignRoleToUser(w, r, userIDStr)
	} else if r.Method == http.MethodGet && parts[1] == "permissions" {
		h.GetUserPermissions(w, r, userIDStr)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// AssignRoleToUser handles POST /users/{userId}/role
func (h *RoleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request, userIDStr string) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var request struct {
		RoleID uuid.UUID `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.roleService.AssignRoleToUser(userID, request.RoleID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Role assigned successfully",
	})
}

// GetUserPermissions handles GET /api/users/{userId}/permissions
func (h *RoleHandler) GetUserPermissions(w http.ResponseWriter, r *http.Request, userIDStr string) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*models.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Users can only view their own permissions unless they have users:read permission
	if user.UserID != userID && !user.HasPermission(models.PermUsersRead) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get the target user's role
	role, err := h.roleService.GetRoleByID(*user.RoleID)
	if err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"user_id":     userID,
		"role":        role.Name,
		"permissions": role.Permissions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AssignRoleHandler handles POST /users/{userId}/role with specific routing
func (h *RoleHandler) AssignRoleHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path like /users/{id}/role
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	path = strings.TrimSuffix(path, "/role")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.AssignRoleToUser(w, r, path)
}

// GetUserPermissionsHandler handles GET /users/{userId}/permissions with specific routing
func (h *RoleHandler) GetUserPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path like /users/{id}/permissions
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	path = strings.TrimSuffix(path, "/permissions")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.GetUserPermissions(w, r, path)
}
