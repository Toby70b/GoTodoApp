package test

import (
	"TodoApp/src/main/models"
	"TodoApp/src/main/services"
	"log"
	"testing"
)

var todoService services.TodoService

func setupTest(tb *testing.T) {
	log.Println("setup test")
	var todos []models.Todo
	todoService = services.NewTodoService(todos)
}

func TestReturnAllTodosNoTodoFound(t *testing.T) {
	setupTest(t)
	actualTodos := todoService.ReturnAllTodos()
	if len(actualTodos) != 0 {
		t.Error("expected array of length", 0, "but received array of length", len(actualTodos))
	}
}

func TestReturnAllTodosOneTodoFound(t *testing.T) {
	setupTest(t)
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
	setupTest(t)
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

func TestReturnSingleTodoNoTodoFoundEmptyArray(t *testing.T) {
	setupTest(t)
	actualTodo, _ := todoService.ReturnSingleTodo("1")
	if actualTodo != (models.Todo{Completed: false}) {
		t.Error("expected Todo to be empty, instead was:", actualTodo)
	}
}

func TestReturnSingleTodoNoTodoFoundReturnError(t *testing.T) {
	setupTest(t)
	_, err := todoService.ReturnSingleTodo("1")
	if err == nil {
		t.Error("expected error but no error returned")
	}
	if err.Error() != "could not find todo with id [1]" {
		t.Error("expected error message was not found, instead was:", err.Error())
	}
}

func TestReturnSingleTodoNoTodoFoundWrongId(t *testing.T) {
	setupTest(t)
	expectedTodo := models.Todo{
		Id:        "1",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}
	todoService.Todos = append(todoService.Todos, expectedTodo)
	actualTodo, _ := todoService.ReturnSingleTodo("2")
	if actualTodo != (models.Todo{Completed: false}) {
		t.Error("expected Todo to be empty, instead was:", actualTodo)
	}
}

func TestReturnSingleTodoTodoFound(t *testing.T) {
	setupTest(t)
	expectedTodo := models.Todo{
		Id:        "1",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}
	todoService.Todos = append(todoService.Todos, expectedTodo)
	actualTodo, _ := todoService.ReturnSingleTodo("1")
	if actualTodo != expectedTodo {
		t.Error("expected todo:", expectedTodo, "but received todo:", actualTodo)
	}
}
