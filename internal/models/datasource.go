package models

// DataSourceMetadata describes a data source and its capabilities
type DataSourceMetadata struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Category          string   `json:"category"` // "audit", "risk", "compliance"
	CompatibleWidgets []string `json:"compatible_widgets"`
	RequiresEntity    bool     `json:"requires_entity"`
}

// PieChartSlice represents a single slice in a pie chart
type PieChartSlice struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

// TableColumn defines a column in a table widget
type TableColumn struct {
	Key    string `json:"key"`
	Header string `json:"header"`
	Width  string `json:"width,omitempty"`
}

// TableData represents data for a table widget
type TableData struct {
	Columns []TableColumn          `json:"columns"`
	Rows    []map[string]interface{} `json:"rows"`
}

// DataSourceResponse wraps the response for data source endpoints
type DataSourceResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// Standard color palettes
var SeverityColors = map[string]string{
	"critical": "#7c3aed",
	"high":     "#ef4444",
	"medium":   "#f59e0b",
	"low":      "#22c55e",
}

var StatusColors = map[string]string{
	"open":        "#ef4444",
	"in_progress": "#3b82f6",
	"resolved":    "#10b981",
	"closed":      "#22c55e",
}

var PriorityColors = map[string]string{
	"critical": "#7c3aed",
	"high":     "#ef4444",
	"medium":   "#f59e0b",
	"low":      "#22c55e",
}

// GetAllDataSources returns metadata for all available data sources
func GetAllDataSources() []DataSourceMetadata {
	return []DataSourceMetadata{
		{
			ID:                "todos_by_priority",
			Name:              "Todos by Priority",
			Description:       "Distribution of todos by priority level",
			Category:          "todos",
			CompatibleWidgets: []string{"pie_chart", "table"},
			RequiresEntity:    false,
		},
		{
			ID:                "todos_by_status",
			Name:              "Todos by Status",
			Description:       "Distribution of todos by completion status",
			Category:          "todos",
			CompatibleWidgets: []string{"pie_chart", "table"},
			RequiresEntity:    false,
		},
		{
			ID:                "todos_list",
			Name:              "Todos List",
			Description:       "List of all todos with details",
			Category:          "todos",
			CompatibleWidgets: []string{"table"},
			RequiresEntity:    false,
		},
		{
			ID:                "users_by_role",
			Name:              "Users by Role",
			Description:       "Distribution of users by their assigned role",
			Category:          "users",
			CompatibleWidgets: []string{"pie_chart", "table"},
			RequiresEntity:    false,
		},
		{
			ID:                "user_activity",
			Name:              "User Activity",
			Description:       "Recent user activity and statistics",
			Category:          "users",
			CompatibleWidgets: []string{"table"},
			RequiresEntity:    false,
		},
	}
}
