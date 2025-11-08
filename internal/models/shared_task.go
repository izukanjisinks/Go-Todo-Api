package models

import "time"

type SharedTask struct {
	ID           string `json:"id"`
	OwnerID      int    `json:"owner_id"`
	SharedWithID int    `json:"shared_with_id"`
	TodoID       string `json:"todo_id"`
}

type SharedTodoWithOwner struct {
	TodoID          string    `json:"todo_id"`
	TaskName        string    `json:"task_name"`
	TaskDescription string    `json:"task_description"`
	OwnerUsername   string    `json:"owner_username"`
	SharedWithID    int       `json:"shared_with_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
