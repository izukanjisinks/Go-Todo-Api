package repository

import (
	"database/sql"
	"fmt"
	"todo-api/internal/database"
	"todo-api/internal/models"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		db: database.DB,
	}
}

func (r *TodoRepository) GetAll() ([]models.Todo, error) {
	rows, err := r.db.Query(`SELECT CAST(id AS VARCHAR(36)), task_name, task_description, completed, user_id, created_at, updated_at FROM todos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.Id, &todo.TaskName, &todo.TaskDescription, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *TodoRepository) GetById(userID string) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.QueryRow(`SELECT CAST(id AS VARCHAR(36)), task_name, task_description, completed, user_id, created_at, updated_at FROM todos WHERE id = $1`, userID).
		Scan(&todo.Id, &todo.TaskName, &todo.TaskDescription, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepository) Create(todo *models.Todo) error {
	fmt.Println("the todo", todo)
	_, err := r.db.Exec(`INSERT INTO todos (id, task_name, task_description, completed, user_id) VALUES ($1, $2, $3, $4, $5)`,
		todo.Id, todo.TaskName, todo.TaskDescription, todo.Completed, todo.UserID)
	return err
}

func (r *TodoRepository) Update(id string, todo *models.Todo) (int64, error) {
	result, err := r.db.Exec(`UPDATE todos SET task_name = $1, task_description = $2, completed = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`,
		todo.TaskName, todo.TaskDescription, todo.Completed, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

func (r *TodoRepository) Delete(id string) (int64, error) {
	result, err := r.db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

func (r *TodoRepository) GetByUserId(userID string) ([]models.Todo, error) {
	rows, err := r.db.Query(`SELECT CAST(id AS VARCHAR(36)), task_name, task_description, completed, user_id, created_at, updated_at FROM todos WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.Id, &todo.TaskName, &todo.TaskDescription, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}
