package models

import (
	"time"
)

// TodoStatus represents the state of a todo in the workflow
type TodoStatus string

const (
	StatusDraft    TodoStatus = "Draft"
	StatusReview   TodoStatus = "Review"
	StatusApproved TodoStatus = "Approved"
)

// Todo represents a task in the system
type TodoTask struct {
	ID          string
	Title       string
	Description string
	AssignedTo  string
	Status      TodoStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ReviewedBy  string
	ApprovedBy  string
}
