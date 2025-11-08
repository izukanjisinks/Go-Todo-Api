package handlers

import (
	"net/http"
	"todo-api/internal/models"
	"todo-api/pkg/utils"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	health := models.HealthResponse{
		Status:  "OK",
		Message: "API health is ok",
	}

	utils.RespondJSON(w, http.StatusOK, health)
}
