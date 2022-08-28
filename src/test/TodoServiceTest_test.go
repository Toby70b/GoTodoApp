package test

import (
	"TodoApp/src/main/models"
	"TodoApp/src/main/services"
	"log"
	"testing"
)

var todoService services.TodoService

func setupSuite(tb *testing.T) func(tb *testing.T) {
	log.Println("setup suite")
	var todos []models.Todo
	todoService = services.NewTodoService(todos)
	// Return a function to teardown the test
	return func(tb *testing.T) {
		log.Println("teardown suite")
	}
}

func TestReturnAllTodosNoTodoFound(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	actual := todoService.ReturnAllTodos()
	if len(actual) != 0 {
		t.Error("expected array of length", 0, "but received array of length", actual)
	}
}
