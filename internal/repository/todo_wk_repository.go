package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"todo-api/internal/database"
	"todo-api/internal/models"
)

type TodoWorkflow struct {
	db *sql.DB
}

func NewTodoWorkflow() *TodoWorkflow {
	return &TodoWorkflow{
		db: database.DB,
	}
}

// CreateTodo creates a new todo in draft status
func (tw *TodoWorkflow) CreateTodo(id, title, description, assignedTo string) (*models.TodoTask, error) {
	if id == "" || title == "" || assignedTo == "" {
		return nil, errors.New("id, title, and assignedTo are required")
	}

	// Check if todo already exists
	existing, err := tw.getTodo(id)
	if err == nil && existing != nil {
		return nil, errors.New("todo with this ID already exists")
	}

	todo := &models.TodoTask{
		ID:          id,
		Title:       title,
		Description: description,
		AssignedTo:  assignedTo,
		Status:      models.StatusDraft,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = tw.db.Exec(`INSERT INTO todo_tasks (id, title, description, assigned_to, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		todo.ID, todo.Title, todo.Description, todo.AssignedTo, todo.Status, todo.CreatedAt, todo.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return todo, nil
}

func (tw *TodoWorkflow) SubmitForReview(id, submittedBy string) error {

	todo, err := tw.getTodo(id)

	if err != nil {
		return err
	}

	if todo.Status != models.StatusDraft {
		return fmt.Errorf("todo must be in draft status to submit for review, current status: %s", todo.Status)
	}

	if todo.AssignedTo != submittedBy {
		return errors.New("only the assigned user can submit the todo for review")
	}

	_, err = tw.db.Exec(`UPDATE todo_tasks SET status = $1, reviewed_by = $2, updated_at = $3 WHERE id = $4`,
		models.StatusReview, submittedBy, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to submit todo for review: %w", err)
	}

	return nil
}

func (tw *TodoWorkflow) ApprovedTodo(id, approvedBy string) error {
	todo, err := tw.getTodo(id)

	if err != nil {
		return err
	}

	if todo.Status != models.StatusReview {
		return fmt.Errorf("todo must be in review status to approve, current status: %s", todo.Status)
	}

	if todo.ReviewedBy != approvedBy {
		return errors.New("only the reviewer can approve the todo")
	}

	_, err = tw.db.Exec(`UPDATE todo_tasks SET status = $1, approved_by = $2, updated_at = $3 WHERE id = $4`,
		models.StatusApproved, approvedBy, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to approve todo: %w", err)
	}

	return nil
}

// the one rejecting the todo should be the approver or reviewer of course
func (tw *TodoWorkflow) RejectTodo(id, rejectedBy string) error {
	todo, err := tw.getTodo(id)

	if err != nil {
		return err
	}

	if todo.Status != models.StatusReview {
		return fmt.Errorf("todo must be in review status to reject, current status: %s", todo.Status)
	}

	_, err = tw.db.Exec(`UPDATE todo_tasks SET status = $1, reviewed_by = NULL, updated_at = $2 WHERE id = $3`,
		models.StatusDraft, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to reject todo: %w", err)
	}

	return nil
}

// GetTodosByUser returns all todos assigned to a user
func (tw *TodoWorkflow) GetTodosByUser(userID string) ([]*models.TodoTask, error) {
	rows, err := tw.db.Query(`SELECT id, title, description, assigned_to, status, reviewed_by, approved_by, created_at, updated_at 
		FROM todo_tasks WHERE assigned_to = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos by user: %w", err)
	}
	defer rows.Close()

	var todos []*models.TodoTask
	for rows.Next() {
		todo := &models.TodoTask{}
		var reviewedBy, approvedBy sql.NullString
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.AssignedTo, &todo.Status,
			&reviewedBy, &approvedBy, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		if reviewedBy.Valid {
			todo.ReviewedBy = reviewedBy.String
		}
		if approvedBy.Valid {
			todo.ApprovedBy = approvedBy.String
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// GetTodosByStatus returns all todos with a specific status
func (tw *TodoWorkflow) GetTodosByStatus(status models.TodoStatus) ([]*models.TodoTask, error) {
	rows, err := tw.db.Query(`SELECT id, title, description, assigned_to, status, reviewed_by, approved_by, created_at, updated_at 
		FROM todo_tasks WHERE status = $1`, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos by status: %w", err)
	}
	defer rows.Close()

	var todos []*models.TodoTask
	for rows.Next() {
		todo := &models.TodoTask{}
		var reviewedBy, approvedBy sql.NullString
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.AssignedTo, &todo.Status,
			&reviewedBy, &approvedBy, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		if reviewedBy.Valid {
			todo.ReviewedBy = reviewedBy.String
		}
		if approvedBy.Valid {
			todo.ApprovedBy = approvedBy.String
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// getTodo is a helper to retrieve and validate todo existence
func (tw *TodoWorkflow) getTodo(id string) (*models.TodoTask, error) {
	todo := &models.TodoTask{}
	var reviewedBy, approvedBy sql.NullString
	err := tw.db.QueryRow(`SELECT id, title, description, assigned_to, status, reviewed_by, approved_by, created_at, updated_at 
		FROM todo_tasks WHERE id = $1`, id).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.AssignedTo, &todo.Status,
		&reviewedBy, &approvedBy, &todo.CreatedAt, &todo.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("todo not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	if reviewedBy.Valid {
		todo.ReviewedBy = reviewedBy.String
	}
	if approvedBy.Valid {
		todo.ApprovedBy = approvedBy.String
	}

	return todo, nil
}
