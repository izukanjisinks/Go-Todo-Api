// internal/routes/datasource_routes.go
package routes

import (
	"net/http"
	"todo-api/internal/handlers"
	"todo-api/internal/models"
)

func RegisterDataSourceRoutes(dataSourceHandler *handlers.DataSourceHandler) {
	// GET /api/data-sources - List all available data sources
	// GET /api/data-sources/:id?widget_type=X - Get data for specific data source
	http.HandleFunc("/api/data-sources", withAuthAndPermission(dataSourceHandler.DataSourcesHandler, models.PermView))
	http.HandleFunc("/api/data-sources/", withAuthAndPermission(dataSourceHandler.DataSourcesHandler, models.PermView))
}
