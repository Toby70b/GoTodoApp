package services

import (
	"TodoApp/src/main/models"
	"github.com/google/go-cmp/cmp"
	"testing"
)

var todoService *TodoServiceImpl

func setupTest() {
	todoService = NewTodoServiceImpl([]models.Todo{})
}

func TestReturnAllTodos(t *testing.T) {
	tests := map[string]struct {
		prerequisite []models.Todo
		expected     []models.Todo
	}{
		"Return No Todos": {
			prerequisite: []models.Todo{},
			expected:     []models.Todo{},
		},
		"Return Single Todo": {
			prerequisite: []models.Todo{{Id: "1", Title: "Example Title", Desc: "Example Description", Completed: false}},
			expected:     []models.Todo{{Id: "1", Title: "Example Title", Desc: "Example Description", Completed: false}},
		},
		"Return Multiple Todos": {
			prerequisite: []models.Todo{
				{
					Id: "1", Title: "Example Title", Desc: "Example Description", Completed: false,
				},
				{
					Id: "2", Title: "Example Title 2", Desc: "Example Description", Completed: false,
				},
			},
			expected: []models.Todo{
				{
					Id: "1", Title: "Example Title", Desc: "Example Description", Completed: false,
				},
				{
					Id: "2", Title: "Example Title 2", Desc: "Example Description", Completed: false,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			setupTest()
			todoService.Todos = append(todoService.Todos, tt.prerequisite...)
			actual := todoService.ReturnAllTodos()
			diff := cmp.Diff(tt.expected, actual)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestReturnSingleTodo(t *testing.T) {
	tests := map[string]struct {
		prerequisite         []models.Todo
		input                string
		expected             models.Todo
		errorExpected        bool
		expectedErrorMessage string
	}{
		"No Todo Found Empty Array": {
			prerequisite:         []models.Todo{},
			input:                "1",
			expected:             models.Todo{},
			errorExpected:        true,
			expectedErrorMessage: "could not find todo with id [1]",
		},
		"No Todo Found Wrong Id": {
			prerequisite: []models.Todo{
				{
					Id:        "1",
					Title:     "Example Title",
					Desc:      "Example Description",
					Completed: false,
				},
			},
			input:                "2",
			expected:             models.Todo{},
			errorExpected:        true,
			expectedErrorMessage: "could not find todo with id [2]",
		},
		"Todo With Matching Id Found": {
			prerequisite: []models.Todo{
				{
					Id:        "1",
					Title:     "Example Title",
					Desc:      "Example Description",
					Completed: false,
				},
			},
			input: "1",
			expected: models.Todo{
				Id:        "1",
				Title:     "Example Title",
				Desc:      "Example Description",
				Completed: false,
			},
			errorExpected: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			setupTest()
			todoService.Todos = append(todoService.Todos, tt.prerequisite...)
			actual, err := todoService.ReturnSingleTodo(tt.input)
			diff := cmp.Diff(tt.expected, actual)
			if diff != "" {
				t.Fatalf(diff)
			}
			if tt.errorExpected {
				if err == nil {
					t.Fatalf("Error expected but none occured")
				} else if err.Error() != tt.expectedErrorMessage {
					t.Fatalf("Error message not as expected, expected [%v] but was [%v]", tt.expectedErrorMessage, err.Error())
				}

			} else if !tt.errorExpected && err != nil {
				t.Fatalf("Error occured when none expected: [%v]", err)
			}
		})
	}
}

func TestCreateNewTodoValidationError(t *testing.T) {
	setupTest()
	newTodo := models.Todo{
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}
	_, err := todoService.CreateNewTodo(newTodo)
	if err == nil {
		t.Error("expected error but no error returned")
	}
	if err.Error() != "todo Id cannot be null" {
		t.Error("expected error message was not found, instead was:", err.Error())
	}
	if len(todoService.Todos) != 0 {
		t.Error("expected array of length", 0, "but received array of length", len(todoService.Todos))
	}
}

func TestCreateNewTodoDuplicateIdError(t *testing.T) {
	setupTest()
	newTodo := models.Todo{
		Id:        "1",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}
	todoService.Todos = append(todoService.Todos, newTodo)
	_, err := todoService.CreateNewTodo(newTodo)
	if err == nil {
		t.Error("expected error but no error returned")
	}
	if err.Error() != "todo with id [1] already exists" {
		t.Error("expected error message was not found, instead was:", err.Error())
	}
	if len(todoService.Todos) != 1 {
		t.Error("expected array of length", 1, "but received array of length", len(todoService.Todos))
	}
}

func TestCreateNewTodoNewTodoSuccessfullyCreated(t *testing.T) {
	setupTest()
	newTodo := models.Todo{
		Id:        "1",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}
	actualNewTodo, err := todoService.CreateNewTodo(newTodo)
	if err != nil {
		t.Error("expected no error but error returned")
	}
	if len(todoService.Todos) != 1 {
		t.Error("expected array of length", 1, "but received array of length", len(todoService.Todos))
	}
	if newTodo != actualNewTodo {
		t.Error("expected todo:", newTodo, "but received todo:", actualNewTodo)
	}
}

func TestDeleteTodoSuccessfulDeletion(t *testing.T) {
	setupTest()
	expectedTodo := models.Todo{
		Id:        "1",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}
	todoService.Todos = append(todoService.Todos, expectedTodo)
	todoService.DeleteTodo("1")
	if len(todoService.Todos) != 0 {
		t.Error("expected array of length", 0, "but received array of length", len(todoService.Todos))
	}
}

func TestDeleteTodoSuccessfulDeletionDoesntDeleteNonMatchingTodos(t *testing.T) {
	setupTest()
	expectedTodo1 := models.Todo{
		Id:        "1",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}
	expectedTodo2 := models.Todo{
		Id:        "2",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}
	todoService.Todos = append(todoService.Todos, []models.Todo{expectedTodo1, expectedTodo2}...)

	todoService.DeleteTodo("1")
	if len(todoService.Todos) != 1 {
		t.Error("expected array of length", 1, "but received array of length", len(todoService.Todos))
	}
	if todoService.Todos[0] != expectedTodo2 {
		t.Error("expected todo:", expectedTodo2, "but received todo:", todoService.Todos[0])
	}
}

func TestUpdateTodoValidationError(t *testing.T) {
	setupTest()
	newTodo := models.Todo{
		Id:        "1",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}
	UpdatedTodo := models.Todo{
		Title:     "Updated Example Title",
		Desc:      "Updated Example Description",
		Completed: true,
	}
	todoService.Todos = append(todoService.Todos, newTodo)
	_, err := todoService.UpdateTodo(UpdatedTodo)
	if err == nil {
		t.Error("expected error but no error returned")
	}
	if err.Error() != "todo Id cannot be null" {
		t.Error("expected error message was not found, instead was:", err.Error())
	}
}

func TestUpdateTodoNoTodoFound(t *testing.T) {
	setupTest()
	UpdatedTodo := models.Todo{
		Id:        "1",
		Title:     "Updated Example Title",
		Desc:      "Updated Example Description",
		Completed: true,
	}
	_, err := todoService.UpdateTodo(UpdatedTodo)
	if err == nil {
		t.Error("expected error but no error returned")
	}
	if err.Error() != "could not find todo with id [1]" {
		t.Error("expected error message was not found, instead was:", err.Error())
	}
}

func TestUpdateTodoSuccessfully(t *testing.T) {
	setupTest()
	newTodo := models.Todo{
		Id:        "1",
		Title:     "Example Title",
		Desc:      "Example Description",
		Completed: false,
	}
	UpdatedTodo := models.Todo{
		Id:        "1",
		Title:     "Updated Example Title",
		Desc:      "Updated Example Description",
		Completed: true,
	}
	todoService.Todos = append(todoService.Todos, newTodo)
	_, err := todoService.UpdateTodo(UpdatedTodo)
	if err != nil {
		t.Error("expected no error but error returned")
	}
	if todoService.Todos[0] != UpdatedTodo {
		t.Error("expected todo:", UpdatedTodo, "but received todo:", todoService.Todos[0])
	}
}
