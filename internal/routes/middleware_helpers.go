// internal/routes/middleware_helpers.go
package routes

import (
	"net/http"
	"todo-api/internal/middleware"
)

// withAuth wraps a handler with both CORS and JWT authentication middleware
func withAuth(handler http.HandlerFunc) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(http.HandlerFunc(handler)).ServeHTTP(w, r)
	})
}

// withAuthAndPermission wraps a handler with CORS, JWT auth, and permission check
func withAuthAndPermission(handler http.HandlerFunc, permission string) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(
			middleware.RequirePermission(permission)(
				http.HandlerFunc(handler),
			),
		).ServeHTTP(w, r)
	})
}

// withAuthAndAnyPermission wraps a handler with CORS, JWT auth, and any permission check
func withAuthAndAnyPermission(handler http.HandlerFunc, permissions ...string) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(
			middleware.RequireAnyPermission(permissions...)(
				http.HandlerFunc(handler),
			),
		).ServeHTTP(w, r)
	})
}

// withAuthAndAllPermissions wraps a handler with CORS, JWT auth, and all permissions check
func withAuthAndAllPermissions(handler http.HandlerFunc, permissions ...string) http.HandlerFunc {
	return middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(
			middleware.RequireAllPermissions(permissions...)(
				http.HandlerFunc(handler),
			),
		).ServeHTTP(w, r)
	})
}
