package middleware

import (
	"context"
	"net/http"

	"todo-api/internal/models"
)

// RequirePermission is a middleware that checks if the authenticated user has a specific permission
func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context (set by auth middleware)
			user, ok := r.Context().Value("user").(*models.User)
			if !ok || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has the required permission
			if !user.HasPermission(permission) {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission checks if the user has at least one of the specified permissions
func RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*models.User)
			if !ok || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has any of the required permissions
			hasPermission := false
			for _, perm := range permissions {
				if user.HasPermission(perm) {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAllPermissions checks if the user has all of the specified permissions
func RequireAllPermissions(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*models.User)
			if !ok || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has all required permissions
			for _, perm := range permissions {
				if !user.HasPermission(perm) {
					http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole checks if the user has a specific role
func RequireRole(roleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*models.User)
			if !ok || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if user.Role == nil || user.Role.Name != roleName {
				http.Error(w, "Forbidden: role required", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AttachPermissionsToContext is a helper to attach permissions to the request context
func AttachPermissionsToContext(ctx context.Context, permissions []string) context.Context {
	return context.WithValue(ctx, "permissions", permissions)
}

// GetPermissionsFromContext retrieves permissions from the request context
func GetPermissionsFromContext(ctx context.Context) []string {
	permissions, ok := ctx.Value("permissions").([]string)
	if !ok {
		return []string{}
	}
	return permissions
}
