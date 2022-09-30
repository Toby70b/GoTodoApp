package test

import (
	"TodoApp/src/main/controllers"
	"TodoApp/src/main/models"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO look at test tables
var todoController controllers.TodoController

type MockTodoServiceImpl struct {
	mock.Mock
	Todos []models.Todo
}

func (service MockTodoServiceImpl) ReturnAllTodos() []models.Todo {
	args := service.Called()
	return args.Get(0).([]models.Todo)
}

func (service MockTodoServiceImpl) ReturnSingleTodo(id string) (models.Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (service MockTodoServiceImpl) CreateNewTodo(newTodo models.Todo) (models.Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (service MockTodoServiceImpl) DeleteTodo(id string) {
	//TODO implement me
	panic("implement me")
}

func (service MockTodoServiceImpl) UpdateTodo(newTodo models.Todo) (models.Todo, error) {
	//TODO implement me
	panic("implement me")
}

func setupTodoController(service *MockTodoServiceImpl) {
	log.Println("setup test")
	todoController = controllers.NewTodoController(service)
}

func TestReturnAllTodosSingleTodo(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	var mockTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a carrot cake for tomorrow's fate",
		Completed: false,
	}
	testObj.On("ReturnAllTodos").Return([]models.Todo{mockTodo})

	setupTodoController(testObj)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	todoController.ReturnAllTodos(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	expectedResponse, _ := json.Marshal([]models.Todo{mockTodo})
	require.JSONEq(t, string(expectedResponse), string(data))
}

func TestReturnAllTodosMultipleTodos(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	var mockTodos = []models.Todo{
		{
			Id:        "1",
			Title:     "Bake cake",
			Desc:      "Bake a carrot cake for tomorrow's fate",
			Completed: false,
		},
		{
			Id:        "2",
			Title:     "Iron shirts",
			Desc:      "Iron shirts within dryer",
			Completed: false,
		},
		{
			Id:        "3",
			Title:     "Walk dog",
			Desc:      "Walk the dog around the town",
			Completed: false,
		}}
	testObj.On("ReturnAllTodos").Return(mockTodos)

	setupTodoController(testObj)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	todoController.ReturnAllTodos(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	expectedResponse, _ := json.Marshal(mockTodos)
	require.JSONEq(t, string(expectedResponse), string(data))
}

func TestReturnAllTodosMultipleNilTodo(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	testObj.On("ReturnAllTodos").Return([]models.Todo{})

	setupTodoController(testObj)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	todoController.ReturnAllTodos(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	expectedResponse, _ := json.Marshal([]models.Todo{})
	require.JSONEq(t, string(expectedResponse), string(data))
}
