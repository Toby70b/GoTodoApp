package test

import (
	"TodoApp/src/main/models"
	"TodoApp/src/main/services"
	"log"
	"testing"
)

var todoService services.TodoService

func setupTest(tb *testing.T) func(tb *testing.T) {
	log.Println("setup test")
	var todos []models.Todo
	todoService = services.NewTodoService(todos)
	// Return a function to teardown the test
	return func(tb *testing.T) {
		log.Println("teardown test")
	}
}

func TestReturnAllTodosNoTodoFound(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	actualTodos := todoService.ReturnAllTodos()
	if len(actualTodos) != 0 {
		t.Error("expected array of length", 0, "but received array of length", len(actualTodos))
	}
}

func TestReturnAllTodosOneTodoFound(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expectedTodo := models.Todo{
		Id:        "1",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}

	todoService.Todos = append(todoService.Todos, expectedTodo)

	actualTodos := todoService.ReturnAllTodos()
	if len(actualTodos) != 1 {
		t.Error("expected array of length", 1, "but received array of length", len(actualTodos))
	}
	if actualTodos[0] != expectedTodo {
		t.Error("expected todo:", expectedTodo, "but received todo:", actualTodos[0])
	}
}

func TestReturnAllTodosMultipleTodosFound(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	expectedTodo1 := models.Todo{
		Id:        "1",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}

	expectedTodo2 := models.Todo{
		Id:        "2",
		Title:     "Example Title 2",
		Desc:      "Example Description ",
		Completed: false,
	}

	todoService.Todos = append(todoService.Todos, []models.Todo{expectedTodo1, expectedTodo2}...)

	actualTodos := todoService.ReturnAllTodos()
	if len(actualTodos) != 2 {
		t.Error("expected array of length", 2, "but received array of length", len(actualTodos))
	}
	if actualTodos[0] != expectedTodo1 {
		t.Error("expected todo:", expectedTodo1, "but received todo:", actualTodos[0])
	}
	if actualTodos[1] != expectedTodo2 {
		t.Error("expected todo:", expectedTodo2, "but received todo:", actualTodos[1])
	}
}
