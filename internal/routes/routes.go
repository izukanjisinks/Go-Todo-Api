package routes

import "todo-api/internal/handlers"

func RegisterRoutes(
	todoHandler *handlers.TodoHandler,
	userHandler *handlers.UsersHandler,
	sharedTaskHandler *handlers.SharedTaskHandler,
	roleHandler *handlers.RoleHandler,
	todoWorkflowHandler *handlers.TodoWorkflowHandler,
	workflowAdminHandler *handlers.WorkflowAdminHandler,
	workflowInstanceHandler *handlers.WorkflowInstanceHandler,

) {
	// Register all routes
	RegisterPublicRoutes(userHandler)
	RegisterTodoRoutes(todoHandler)
	RegisterUserRoutes(userHandler)
	RegisterSharedTaskRoutes(sharedTaskHandler)
	RegisterRoleRoutes(roleHandler, userHandler)
	RegisterWorkflowRoutes(todoWorkflowHandler, workflowAdminHandler, workflowInstanceHandler)

}
