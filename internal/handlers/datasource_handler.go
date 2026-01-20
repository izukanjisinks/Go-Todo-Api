package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"todo-api/internal/models"
	"todo-api/internal/services"
)

type DataSourceHandler struct {
	service *services.DataSourceService
}

func NewDataSourceHandler(service *services.DataSourceService) *DataSourceHandler {
	return &DataSourceHandler{service: service}
}

// GetAllDataSources handles GET /api/data-sources
// Returns metadata about all available data sources
func (h *DataSourceHandler) GetAllDataSources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	dataSources := h.service.GetAllDataSources()

	h.sendSuccess(w, dataSources)
}

// GetDataSourceData handles GET /api/data-sources/:dataSourceId
// Returns data for a specific data source based on widget_type query parameter
func (h *DataSourceHandler) GetDataSourceData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	// Extract data source ID from path
	// Path: /api/data-sources/{dataSourceId}
	path := strings.TrimPrefix(r.URL.Path, "/api/data-sources/")
	dataSourceID := strings.TrimSuffix(path, "/")

	if dataSourceID == "" {
		h.sendError(w, "Data source ID is required", "MISSING_PARAMETER", http.StatusBadRequest)
		return
	}

	// Get widget_type query parameter
	widgetType := r.URL.Query().Get("widget_type")
	if widgetType == "" {
		h.sendError(w, "Missing required parameter: widget_type", "MISSING_PARAMETER", http.StatusBadRequest)
		return
	}

	// Validate widget type
	if widgetType != "pie_chart" && widgetType != "table" {
		h.sendError(w, "Invalid widget_type. Must be 'pie_chart' or 'table'", "INVALID_WIDGET_TYPE", http.StatusBadRequest)
		return
	}

	// Check if data source exists
	ds := h.service.GetDataSourceByID(dataSourceID)
	if ds == nil {
		h.sendError(w, "Data source not found: "+dataSourceID, "DATA_SOURCE_NOT_FOUND", http.StatusNotFound)
		return
	}

	// Fetch data
	data, err := h.service.GetDataSourceData(dataSourceID, widgetType)
	if err != nil {
		h.sendError(w, err.Error(), "INTERNAL_ERROR", http.StatusInternalServerError)
		return
	}

	h.sendSuccess(w, data)
}

// DataSourcesHandler is the main router for /api/data-sources endpoints
func (h *DataSourceHandler) DataSourcesHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/data-sources")
	path = strings.TrimPrefix(path, "/")

	// If path is empty, list all data sources
	if path == "" {
		h.GetAllDataSources(w, r)
		return
	}

	// Otherwise, get specific data source data
	h.GetDataSourceData(w, r)
}

// Helper methods

func (h *DataSourceHandler) sendSuccess(w http.ResponseWriter, data interface{}) {
	response := models.DataSourceResponse{
		Success: true,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *DataSourceHandler) sendError(w http.ResponseWriter, message string, code string, statusCode int) {
	response := models.DataSourceResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
