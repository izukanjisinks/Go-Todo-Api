package routes

import (
	"net/http"
	"todo-api/internal/handlers"
	"todo-api/internal/middleware"
)

func RegisterPublicRoutes(userHandler *handlers.UsersHandler) {
	http.HandleFunc("/health", middleware.CORS(handlers.HealthHandler))
	http.HandleFunc("/register", middleware.CORS(userHandler.Register))
	http.HandleFunc("/login", middleware.CORS(userHandler.Login))
}
