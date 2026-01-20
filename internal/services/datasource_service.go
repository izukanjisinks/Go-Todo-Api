package services

import (
	"fmt"

	"todo-api/internal/models"
	"todo-api/internal/repository"
)

type DataSourceService struct {
	repo *repository.DataSourceRepository
}

func NewDataSourceService(repo *repository.DataSourceRepository) *DataSourceService {
	return &DataSourceService{repo: repo}
}

// GetAllDataSources returns metadata for all available data sources
func (s *DataSourceService) GetAllDataSources() []models.DataSourceMetadata {
	return models.GetAllDataSources()
}

// GetDataSourceByID returns metadata for a specific data source
func (s *DataSourceService) GetDataSourceByID(id string) *models.DataSourceMetadata {
	dataSources := models.GetAllDataSources()
	for _, ds := range dataSources {
		if ds.ID == id {
			return &ds
		}
	}
	return nil
}

// GetDataSourceData fetches data for a specific data source based on widget type
func (s *DataSourceService) GetDataSourceData(dataSourceID string, widgetType string) (interface{}, error) {
	// Validate widget type
	if widgetType != "pie_chart" && widgetType != "table" {
		return nil, fmt.Errorf("invalid widget_type: must be 'pie_chart' or 'table'")
	}

	// Check if data source exists
	ds := s.GetDataSourceByID(dataSourceID)
	if ds == nil {
		return nil, fmt.Errorf("data source not found: %s", dataSourceID)
	}

	// Check if widget type is compatible
	isCompatible := false
	for _, w := range ds.CompatibleWidgets {
		if w == widgetType {
			isCompatible = true
			break
		}
	}
	if !isCompatible {
		return nil, fmt.Errorf("widget type '%s' is not compatible with data source '%s'", widgetType, dataSourceID)
	}

	// Fetch data based on data source ID and widget type
	switch dataSourceID {
	case "todos_by_priority":
		if widgetType == "pie_chart" {
			return s.repo.GetTodosByPriority()
		}
		return s.repo.GetTodosByPriorityTable()

	case "todos_by_status":
		if widgetType == "pie_chart" {
			return s.repo.GetTodosByStatus()
		}
		return s.repo.GetTodosByStatusTable()

	case "todos_list":
		return s.repo.GetTodosList()

	case "users_by_role":
		if widgetType == "pie_chart" {
			return s.repo.GetUsersByRole()
		}
		return s.repo.GetUsersByRoleTable()

	case "user_activity":
		return s.repo.GetUserActivity()

	default:
		return nil, fmt.Errorf("data source not implemented: %s", dataSourceID)
	}
}
