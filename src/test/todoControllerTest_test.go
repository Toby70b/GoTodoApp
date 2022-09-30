package test

import (
	"TodoApp/src/main/controllers"
	"TodoApp/src/main/models"
	"github.com/stretchr/testify/mock"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestReturnAllTodos(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	var mockTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a carrot cake for tomorrow's fate",
		Completed: false,
	}
	testObj.On("ReturnAllTodos").Return()

	setupTodoController(testObj)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	todoController.ReturnAllTodos(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "ABC" {
		t.Errorf("expected ABC got %v", string(data))
	}
}
