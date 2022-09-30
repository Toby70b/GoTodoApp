package controllers

import (
	"TodoApp/src/main/models"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO look at test tables
var todoController TodoController

type MockTodoServiceImpl struct {
	mock.Mock
	Todos []models.Todo
}

func (service MockTodoServiceImpl) ReturnAllTodos() []models.Todo {
	args := service.Called()
	return args.Get(0).([]models.Todo)
}

func (service MockTodoServiceImpl) ReturnSingleTodo(id string) (models.Todo, error) {
	args := service.Called()
	return args.Get(0).(models.Todo), args.Get(1).(error)
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
	todoController = NewTodoController(service)
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
	httpWriter := httptest.NewRecorder()

	todoController.ReturnAllTodos(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}

	if httpWriter.Code != http.StatusOK {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusOK, httpWriter.Code)
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
	httpWriter := httptest.NewRecorder()

	todoController.ReturnAllTodos(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusOK {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusOK, httpWriter.Code)
	}
	expectedResponse, _ := json.Marshal(mockTodos)
	require.JSONEq(t, string(expectedResponse), string(data))
}

func TestReturnAllTodosMultipleNilTodo(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	testObj.On("ReturnAllTodos").Return([]models.Todo{})

	setupTodoController(testObj)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	httpWriter := httptest.NewRecorder()

	todoController.ReturnAllTodos(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusOK {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusOK, httpWriter.Code)
	}
	expectedResponse, _ := json.Marshal([]models.Todo{})
	require.JSONEq(t, string(expectedResponse), string(data))
}

func TestReturnSingleTodoTodoFound(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	var mockTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a carrot cake for tomorrow's fate",
		Completed: false,
	}
	testObj.On("ReturnSingleTodo").Return(mockTodo)
	setupTodoController(testObj)
	req := httptest.NewRequest(http.MethodGet, "/1", nil)
	httpWriter := httptest.NewRecorder()

	todoController.ReturnSingleTodo(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusOK {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusOK, httpWriter.Code)
	}
	expectedResponse, _ := json.Marshal(mockTodo)
	require.JSONEq(t, string(expectedResponse), string(data))
}

func TestReturnSingleTodoTodoNotFound(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	testObj.On("ReturnSingleTodo").Return(models.Todo{}, errors.New("could not find todo with id [999]"))
	setupTodoController(testObj)
	req := httptest.NewRequest(http.MethodGet, "/999", nil)
	reqPathParams := map[string]string{
		"id": "999",
	}
	req = mux.SetURLVars(req, reqPathParams)
	httpWriter := httptest.NewRecorder()

	todoController.ReturnSingleTodo(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusNotFound {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusNotFound, httpWriter.Code)
	}
	var expectedResponse = "\"could not find todo with id [999]\""
	require.JSONEq(t, expectedResponse, string(data))
}
