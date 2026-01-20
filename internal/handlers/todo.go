package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"todo-api/internal/interfaces"
	"todo-api/internal/models"
	"todo-api/pkg/utils"

	"github.com/google/uuid"
)

type TodoHandler struct {
	service interfaces.TodoInterface
}

func NewTodoHandler(service interfaces.TodoInterface) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetAllTodos()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo models.Todo

	fmt.Println("body received", r.Body)

	err := json.NewDecoder(r.Body).Decode(&newTodo)

	fmt.Println("user id", newTodo.UserID)
	fmt.Println("todo id", newTodo.Id)
	fmt.Println("todo task", newTodo.TaskName)
	fmt.Println("todo description", newTodo.TaskDescription)

	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error reading body")
		return
	}

	if newTodo.TaskName == "" {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input submitted")
		return
	}

	newTodo.Id = uuid.New().String()

	err = h.service.CreateTodo(&newTodo)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, newTodo)
}

func (h *TodoHandler) TodosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.GetAllTodos(w, r)
	case "POST":
		h.CreateTodo(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/todos/"):]

	todo, err := h.service.GetTodoByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if todo == nil {
		utils.RespondError(w, http.StatusNotFound, "Todo not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/todos/"):]

	var updatedTodo models.Todo
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error reading body")
		return
	}

	err = json.Unmarshal(body, &updatedTodo)
	if err != nil || updatedTodo.TaskName == "" {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input submitted")
		return
	}

	rowsAffected, err := h.service.UpdateTodo(id, &updatedTodo)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if rowsAffected == 0 {
		utils.RespondError(w, http.StatusNotFound, "Todo not found")
		return
	}

	updatedTodo.Id = id
	utils.RespondJSON(w, http.StatusOK, updatedTodo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/todos/"):]

	rowsAffected, err := h.service.DeleteTodo(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if rowsAffected == 0 {
		utils.RespondError(w, http.StatusNotFound, "Todo not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Deleted todo with id " + id + " successfully",
	})
}

func (h *TodoHandler) TodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.GetTodoById(w, r)
	case "PUT":
		h.UpdateTodo(w, r)
	case "DELETE":
		h.DeleteTodo(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *TodoHandler) GetTodosByUserId(w http.ResponseWriter, r *http.Request) {
	// Extract user_id from query parameter
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	// Convert user_id to int
	userID := userIDStr

	// Get todos from repository
	todos, err := h.service.GetTodosByUserID(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}
