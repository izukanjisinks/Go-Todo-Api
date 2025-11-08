package repository

import (
	"database/sql"
	"todo-api/internal/database"
	"todo-api/internal/models"
)

type SharedTaskRepository struct {
	db *sql.DB
}

func NewSharedTaskRepository() *SharedTaskRepository {
	return &SharedTaskRepository{
		db: database.DB,
	}
}

func (r *SharedTaskRepository) GetAll() ([]models.SharedTask, error) {
	rows, err := r.db.Query(`SELECT CAST(id AS VARCHAR(36)), owner_id, shared_with_id, CAST(todo_id AS VARCHAR(36)) FROM shared_tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sharedTasks []models.SharedTask
	for rows.Next() {
		var sharedTask models.SharedTask
		err := rows.Scan(&sharedTask.ID, &sharedTask.OwnerID, &sharedTask.SharedWithID, &sharedTask.TodoID)
		if err != nil {
			return nil, err
		}
		sharedTasks = append(sharedTasks, sharedTask)
	}

	return sharedTasks, nil
}

func (r *SharedTaskRepository) GetById(id string) (*models.SharedTask, error) {
	var sharedTask models.SharedTask
	err := r.db.QueryRow(`SELECT CAST(id AS VARCHAR(36)), owner_id, shared_with_id, CAST(todo_id AS VARCHAR(36)) FROM shared_tasks WHERE id = @p1`, id).
		Scan(&sharedTask.ID, &sharedTask.OwnerID, &sharedTask.SharedWithID, &sharedTask.TodoID)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &sharedTask, nil
}

func (r *SharedTaskRepository) Create(sharedTask *models.SharedTask) error {
	_, err := r.db.Exec(`INSERT INTO shared_tasks (id, owner_id, shared_with_id, todo_id) VALUES (@p1, @p2, @p3, @p4)`,
		sharedTask.ID, sharedTask.OwnerID, sharedTask.SharedWithID, sharedTask.TodoID)
	return err
}

func (r *SharedTaskRepository) Delete(id string) (int64, error) {
	result, err := r.db.Exec("DELETE FROM shared_tasks WHERE id = @p1", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

func (r *SharedTaskRepository) GetByOwnerId(ownerID int) ([]models.SharedTask, error) {
	rows, err := r.db.Query(`SELECT CAST(id AS VARCHAR(36)), owner_id, shared_with_id, CAST(todo_id AS VARCHAR(36)) FROM shared_tasks WHERE owner_id = @p1`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sharedTasks []models.SharedTask
	for rows.Next() {
		var sharedTask models.SharedTask
		err := rows.Scan(&sharedTask.ID, &sharedTask.OwnerID, &sharedTask.SharedWithID, &sharedTask.TodoID)
		if err != nil {
			return nil, err
		}
		sharedTasks = append(sharedTasks, sharedTask)
	}

	return sharedTasks, nil
}

func (r *SharedTaskRepository) GetTodosBySharedId(sharedWithID int) ([]models.SharedTodoWithOwner, error) {
	rows, err := r.db.Query(`
		SELECT 
			CAST(t.id AS VARCHAR(36)) AS todo_id,
			t.task_name,
			t.task_description,
			u.username AS owner_username,
			s.shared_with_id,
			t.created_at,
			t.updated_at
		FROM todos t
		JOIN shared_tasks s ON t.id = s.todo_id
		JOIN users u ON s.owner_id = u.id
		WHERE s.shared_with_id = @p1
	`, sharedWithID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sharedTodos []models.SharedTodoWithOwner
	for rows.Next() {
		var sharedTodo models.SharedTodoWithOwner
		err := rows.Scan(&sharedTodo.TodoID, &sharedTodo.TaskName, &sharedTodo.TaskDescription, &sharedTodo.OwnerUsername, &sharedTodo.SharedWithID, &sharedTodo.CreatedAt, &sharedTodo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		sharedTodos = append(sharedTodos, sharedTodo)
	}

	return sharedTodos, nil
}

func (r *SharedTaskRepository) GetByTodoId(todoID string) ([]models.SharedTask, error) {
	rows, err := r.db.Query(`SELECT CAST(id AS VARCHAR(36)), owner_id, shared_with_id, CAST(todo_id AS VARCHAR(36)) FROM shared_tasks WHERE todo_id = @p1`, todoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sharedTasks []models.SharedTask
	for rows.Next() {
		var sharedTask models.SharedTask
		err := rows.Scan(&sharedTask.ID, &sharedTask.OwnerID, &sharedTask.SharedWithID, &sharedTask.TodoID)
		if err != nil {
			return nil, err
		}
		sharedTasks = append(sharedTasks, sharedTask)
	}

	return sharedTasks, nil
}
