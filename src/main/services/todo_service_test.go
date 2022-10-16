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

func TestCreateNewTodo(t *testing.T) {
	tests := map[string]struct {
		prerequisite         []models.Todo
		input                models.Todo
		expected             models.Todo
		errorExpected        bool
		expectedErrorMessage string
	}{
		"Validation Error": {
			prerequisite: []models.Todo{},
			input: models.Todo{
				Title:     "Example Title",
				Desc:      "Example Description",
				Completed: false,
			},
			expected:             models.Todo{},
			errorExpected:        true,
			expectedErrorMessage: "todo Id cannot be null",
		},
		"Duplicate Id Error": {
			prerequisite: []models.Todo{
				{
					Id:        "1",
					Title:     "Example Title",
					Desc:      "Example Description",
					Completed: false,
				},
			},
			input: models.Todo{
				Id:        "1",
				Title:     "Example Title",
				Desc:      "Example Description",
				Completed: false,
			},
			expected:             models.Todo{},
			errorExpected:        true,
			expectedErrorMessage: "todo with id [1] already exists",
		},
		"Create Todo Successfully": {
			prerequisite: []models.Todo{},
			input: models.Todo{
				Id:        "1",
				Title:     "Example Title",
				Desc:      "Example Description",
				Completed: false,
			},
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
			var numOfTodosAfterPreReq = len(todoService.Todos)
			actual, err := todoService.CreateNewTodo(tt.input)
			diff := cmp.Diff(tt.expected, actual)
			if diff != "" {
				t.Fatalf(diff)
			}
			if tt.errorExpected {
				if len(todoService.Todos) != numOfTodosAfterPreReq {
					t.Fatalf("Number of persisted Todos has changed unexpectedly")
				}
				if err == nil {
					t.Fatalf("Error expected but none occured")
				} else if err.Error() != tt.expectedErrorMessage {
					t.Fatalf("Error message not as expected, expected [%v] but was [%v]", tt.expectedErrorMessage, err.Error())
				}

			} else if !tt.errorExpected && err != nil {
				t.Fatalf("Error occured when none expected: [%v]", err)
			}
			if !tt.errorExpected && len(todoService.Todos) <= numOfTodosAfterPreReq {
				t.Fatalf("Number of has not increased as expected")
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	tests := map[string]struct {
		prerequisite []models.Todo
		input        string
		expected     []models.Todo
	}{

		"Successful deletion": {
			prerequisite: []models.Todo{

				{
					Id:        "1",
					Title:     "Example Title",
					Desc:      "Example Description",
					Completed: false,
				},
				{
					Id:        "2",
					Title:     "Example Title",
					Desc:      "Example Description",
					Completed: false,
				},
			},
			input: "1",
			expected: []models.Todo{
				{
					Id:        "2",
					Title:     "Example Title",
					Desc:      "Example Description",
					Completed: false,
				},
			},
		},
		"Non-Matching Id Does Not Delete Anything": {
			prerequisite: []models.Todo{
				{
					Id:        "1",
					Title:     "Example Title",
					Desc:      "Example Description",
					Completed: false,
				},
				{
					Id:        "2",
					Title:     "Example Title",
					Desc:      "Example Description",
					Completed: false,
				},
			},
			input: "3",
			expected: []models.Todo{
				{
					Id:        "1",
					Title:     "Example Title",
					Desc:      "Example Description",
					Completed: false,
				},
				{
					Id:        "2",
					Title:     "Example Title",
					Desc:      "Example Description",
					Completed: false,
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			setupTest()
			todoService.Todos = append(todoService.Todos, tt.prerequisite...)
			todoService.DeleteTodo(tt.input)
			diff := cmp.Diff(tt.expected, todoService.Todos)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
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
