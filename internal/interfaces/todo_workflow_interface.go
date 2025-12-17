package interfaces

import "todo-api/internal/models"

// TodoWorkflowInterface defines the business logic contract for todo workflow operations
type TodoWorkflowInterface interface {
	CreateTodoTask(todo *models.TodoTask) error
	GetTodosByUser(userID string) ([]models.TodoTask, error)
	GetTodosByStatus(status models.TodoStatus) ([]models.TodoTask, error)
	SubmitForReview(todoID string, submittedBy string) error
	ApproveTodo(todoID string, approvedBy string) error
	RejectTodo(todoID string, rejectedBy string) error
}
