package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"todo-api/internal/database"
	"todo-api/internal/models"
)

type DataSourceRepository struct {
	db *sql.DB
}

func NewDataSourceRepository() *DataSourceRepository {
	return &DataSourceRepository{
		db: database.DB,
	}
}

// GetTodosByPriority returns todos grouped by priority
func (r *DataSourceRepository) GetTodosByPriority() ([]models.PieChartSlice, error) {
	query := `
		SELECT
			COALESCE(priority, 'medium') as priority,
			COUNT(*) as count
		FROM todos
		GROUP BY priority
		ORDER BY
			CASE priority
				WHEN 'critical' THEN 1
				WHEN 'high' THEN 2
				WHEN 'medium' THEN 3
				WHEN 'low' THEN 4
				ELSE 5
			END
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos by priority: %w", err)
	}
	defer rows.Close()

	var slices []models.PieChartSlice
	for rows.Next() {
		var priority string
		var count float64

		if err := rows.Scan(&priority, &count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Capitalize first letter for label
		label := strings.ToUpper(string(priority[0])) + priority[1:]

		slices = append(slices, models.PieChartSlice{
			Label: label,
			Value: count,
			Color: models.PriorityColors[strings.ToLower(priority)],
		})
	}

	return slices, nil
}

// GetTodosByPriorityTable returns todos by priority in table format
func (r *DataSourceRepository) GetTodosByPriorityTable() (*models.TableData, error) {
	query := `
		SELECT
			COALESCE(priority, 'medium') as priority,
			COUNT(*) as count,
			ROUND(COUNT(*) * 100.0 / NULLIF(SUM(COUNT(*)) OVER(), 0), 0) as percentage
		FROM todos
		GROUP BY priority
		ORDER BY
			CASE priority
				WHEN 'critical' THEN 1
				WHEN 'high' THEN 2
				WHEN 'medium' THEN 3
				WHEN 'low' THEN 4
				ELSE 5
			END
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos by priority: %w", err)
	}
	defer rows.Close()

	tableData := &models.TableData{
		Columns: []models.TableColumn{
			{Key: "priority", Header: "Priority", Width: "40%"},
			{Key: "count", Header: "Count", Width: "30%"},
			{Key: "percentage", Header: "Percentage", Width: "30%"},
		},
		Rows: []map[string]interface{}{},
	}

	for rows.Next() {
		var priority string
		var count int
		var percentage float64

		if err := rows.Scan(&priority, &count, &percentage); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		label := strings.ToUpper(string(priority[0])) + priority[1:]

		tableData.Rows = append(tableData.Rows, map[string]interface{}{
			"priority":   label,
			"count":      count,
			"percentage": fmt.Sprintf("%.0f%%", percentage),
		})
	}

	return tableData, nil
}

// GetTodosByStatus returns todos grouped by completion status
func (r *DataSourceRepository) GetTodosByStatus() ([]models.PieChartSlice, error) {
	query := `
		SELECT
			CASE WHEN completed THEN 'Completed' ELSE 'Pending' END as status,
			COUNT(*) as count
		FROM todos
		GROUP BY completed
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos by status: %w", err)
	}
	defer rows.Close()

	statusColors := map[string]string{
		"Completed": "#22c55e",
		"Pending":   "#f59e0b",
	}

	var slices []models.PieChartSlice
	for rows.Next() {
		var status string
		var count float64

		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		slices = append(slices, models.PieChartSlice{
			Label: status,
			Value: count,
			Color: statusColors[status],
		})
	}

	return slices, nil
}

// GetTodosByStatusTable returns todos by status in table format
func (r *DataSourceRepository) GetTodosByStatusTable() (*models.TableData, error) {
	query := `
		SELECT
			CASE WHEN completed THEN 'Completed' ELSE 'Pending' END as status,
			COUNT(*) as count,
			ROUND(COUNT(*) * 100.0 / NULLIF(SUM(COUNT(*)) OVER(), 0), 0) as percentage
		FROM todos
		GROUP BY completed
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos by status: %w", err)
	}
	defer rows.Close()

	tableData := &models.TableData{
		Columns: []models.TableColumn{
			{Key: "status", Header: "Status", Width: "40%"},
			{Key: "count", Header: "Count", Width: "30%"},
			{Key: "percentage", Header: "Percentage", Width: "30%"},
		},
		Rows: []map[string]interface{}{},
	}

	for rows.Next() {
		var status string
		var count int
		var percentage float64

		if err := rows.Scan(&status, &count, &percentage); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		tableData.Rows = append(tableData.Rows, map[string]interface{}{
			"status":     status,
			"count":      count,
			"percentage": fmt.Sprintf("%.0f%%", percentage),
		})
	}

	return tableData, nil
}

// GetTodosList returns list of todos for table widget
func (r *DataSourceRepository) GetTodosList() (*models.TableData, error) {
	query := `
		SELECT
			t.id,
			t.title,
			COALESCE(t.priority, 'medium') as priority,
			CASE WHEN t.completed THEN 'Completed' ELSE 'Pending' END as status,
			u.username as owner,
			t.created_at
		FROM todos t
		LEFT JOIN users u ON t.user_id = u.id
		ORDER BY t.created_at DESC
		LIMIT 100
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos list: %w", err)
	}
	defer rows.Close()

	tableData := &models.TableData{
		Columns: []models.TableColumn{
			{Key: "id", Header: "ID", Width: "10%"},
			{Key: "title", Header: "Title", Width: "30%"},
			{Key: "priority", Header: "Priority", Width: "15%"},
			{Key: "status", Header: "Status", Width: "15%"},
			{Key: "owner", Header: "Owner", Width: "15%"},
			{Key: "created_at", Header: "Created", Width: "15%"},
		},
		Rows: []map[string]interface{}{},
	}

	for rows.Next() {
		var id int
		var title, priority, status string
		var owner sql.NullString
		var createdAt sql.NullTime

		if err := rows.Scan(&id, &title, &priority, &status, &owner, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		ownerStr := "Unassigned"
		if owner.Valid {
			ownerStr = owner.String
		}

		createdStr := ""
		if createdAt.Valid {
			createdStr = createdAt.Time.Format("2006-01-02")
		}

		tableData.Rows = append(tableData.Rows, map[string]interface{}{
			"id":         id,
			"title":      title,
			"priority":   strings.ToUpper(string(priority[0])) + priority[1:],
			"status":     status,
			"owner":      ownerStr,
			"created_at": createdStr,
		})
	}

	return tableData, nil
}

// GetUsersByRole returns users grouped by role
func (r *DataSourceRepository) GetUsersByRole() ([]models.PieChartSlice, error) {
	query := `
		SELECT
			COALESCE(r.name, 'No Role') as role_name,
			COUNT(*) as count
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		GROUP BY r.name
		ORDER BY count DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by role: %w", err)
	}
	defer rows.Close()

	// Color palette for roles
	colors := []string{"#3b82f6", "#22c55e", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899"}
	colorIndex := 0

	var slices []models.PieChartSlice
	for rows.Next() {
		var roleName string
		var count float64

		if err := rows.Scan(&roleName, &count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		slices = append(slices, models.PieChartSlice{
			Label: roleName,
			Value: count,
			Color: colors[colorIndex%len(colors)],
		})
		colorIndex++
	}

	return slices, nil
}

// GetUsersByRoleTable returns users by role in table format
func (r *DataSourceRepository) GetUsersByRoleTable() (*models.TableData, error) {
	query := `
		SELECT
			COALESCE(r.name, 'No Role') as role_name,
			COUNT(*) as count,
			ROUND(COUNT(*) * 100.0 / NULLIF(SUM(COUNT(*)) OVER(), 0), 0) as percentage
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		GROUP BY r.name
		ORDER BY count DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by role: %w", err)
	}
	defer rows.Close()

	tableData := &models.TableData{
		Columns: []models.TableColumn{
			{Key: "role", Header: "Role", Width: "40%"},
			{Key: "count", Header: "Count", Width: "30%"},
			{Key: "percentage", Header: "Percentage", Width: "30%"},
		},
		Rows: []map[string]interface{}{},
	}

	for rows.Next() {
		var roleName string
		var count int
		var percentage float64

		if err := rows.Scan(&roleName, &count, &percentage); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		tableData.Rows = append(tableData.Rows, map[string]interface{}{
			"role":       roleName,
			"count":      count,
			"percentage": fmt.Sprintf("%.0f%%", percentage),
		})
	}

	return tableData, nil
}

// GetUserActivity returns recent user activity
func (r *DataSourceRepository) GetUserActivity() (*models.TableData, error) {
	query := `
		SELECT
			u.id::TEXT as user_id,
			u.username,
			u.email,
			COALESCE(r.name, 'No Role') as role_name,
			u.is_active,
			u.created_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		ORDER BY u.created_at DESC
		LIMIT 50
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query user activity: %w", err)
	}
	defer rows.Close()

	tableData := &models.TableData{
		Columns: []models.TableColumn{
			{Key: "username", Header: "Username", Width: "20%"},
			{Key: "email", Header: "Email", Width: "25%"},
			{Key: "role", Header: "Role", Width: "20%"},
			{Key: "status", Header: "Status", Width: "15%"},
			{Key: "joined", Header: "Joined", Width: "20%"},
		},
		Rows: []map[string]interface{}{},
	}

	for rows.Next() {
		var userID, username, email, roleName string
		var isActive bool
		var createdAt sql.NullTime

		if err := rows.Scan(&userID, &username, &email, &roleName, &isActive, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		status := "Inactive"
		if isActive {
			status = "Active"
		}

		joinedStr := ""
		if createdAt.Valid {
			joinedStr = createdAt.Time.Format("2006-01-02")
		}

		tableData.Rows = append(tableData.Rows, map[string]interface{}{
			"username": username,
			"email":    email,
			"role":     roleName,
			"status":   status,
			"joined":   joinedStr,
		})
	}

	return tableData, nil
}
