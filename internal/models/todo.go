package models

import "time"

type Todo struct {
	Id              string    `json:"id"`
	TaskName        string    `json:"task_name"`
	TaskDescription string    `json:"task_description"`
	Completed       bool      `json:"completed"`
	UserID          int       `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
