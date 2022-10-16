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
	"strings"
	"testing"
)

// TODO attempt to refactor the below into table tests
// TODO can setup be done before each test without explicitly putting it in each function?
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
	args := service.Called(id)
	if args.Error(1) == nil {
		return args.Get(0).(models.Todo), nil
	}
	return args.Get(0).(models.Todo), args.Get(1).(error)
}

func (service MockTodoServiceImpl) CreateNewTodo(newTodo models.Todo) (models.Todo, error) {
	args := service.Called(newTodo)
	if args.Error(1) == nil {
		return args.Get(0).(models.Todo), nil
	}
	return args.Get(0).(models.Todo), args.Get(1).(error)
}

func (service MockTodoServiceImpl) DeleteTodo(id string) {
	return
}

func (service MockTodoServiceImpl) UpdateTodo(newTodo models.Todo) (models.Todo, error) {
	args := service.Called(newTodo)
	if args.Error(1) == nil {
		return args.Get(0).(models.Todo), nil
	}
	return args.Get(0).(models.Todo), args.Get(1).(error)
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
	testObj.On("ReturnSingleTodo", "1").Return(mockTodo, nil)
	setupTodoController(testObj)
	req := httptest.NewRequest(http.MethodGet, "/1", nil)
	reqPathParams := map[string]string{
		"id": "1",
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
	if httpWriter.Code != http.StatusOK {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusOK, httpWriter.Code)
	}
	expectedResponse, _ := json.Marshal(mockTodo)
	require.JSONEq(t, string(expectedResponse), string(data))
}

func TestReturnSingleTodoTodoNotFound(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	testObj.On("ReturnSingleTodo", "999").Return(models.Todo{}, errors.New("could not find todo with id [999]"))
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

func TestCreateNewTodoSuccessfully(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	var mockTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a carrot cake for tomorrow's fate",
		Completed: false,
	}
	testObj.On("CreateNewTodo", mockTodo).Return(mockTodo, nil)
	setupTodoController(testObj)

	mockTodoJson, _ := json.Marshal(mockTodo)
	bodyReader := strings.NewReader(string(mockTodoJson))
	req := httptest.NewRequest(http.MethodPost, "/", bodyReader)
	httpWriter := httptest.NewRecorder()
	todoController.CreateNewTodo(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusCreated {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusCreated, httpWriter.Code)
	}
	expectedResponse, _ := json.Marshal(mockTodo)
	require.JSONEq(t, string(expectedResponse), string(data))
}

func TestCreateNewTodoInvalidRequestBody(t *testing.T) {
	testObj := new(MockTodoServiceImpl)

	setupTodoController(testObj)
	bodyReader := strings.NewReader("{invalid:jsonsjs}}")
	req := httptest.NewRequest(http.MethodPost, "/", bodyReader)
	httpWriter := httptest.NewRecorder()
	todoController.CreateNewTodo(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusInternalServerError {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusInternalServerError, httpWriter.Code)
	}
	var expectedResponse = "\"Internal Server Error\""
	require.JSONEq(t, expectedResponse, string(data))
}

func TestCreateNewTodoTodoWithDuplicateIdAlreadyExists(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	var mockTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a carrot cake for tomorrow's fate",
		Completed: false,
	}
	testObj.On("CreateNewTodo", mockTodo).Return(models.Todo{}, errors.New("todo with id [1] already exists"))
	setupTodoController(testObj)

	mockTodoJson, _ := json.Marshal(mockTodo)
	bodyReader := strings.NewReader(string(mockTodoJson))
	req := httptest.NewRequest(http.MethodPost, "/", bodyReader)
	httpWriter := httptest.NewRecorder()
	todoController.CreateNewTodo(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusConflict {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusConflict, httpWriter.Code)
	}
	var expectedResponse = "\"Todo with id [1] already exists\""
	require.JSONEq(t, expectedResponse, string(data))
}

func TestDeleteTodo(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	testObj.On("DeleteTodo", "1").Return()
	setupTodoController(testObj)
	reqPathParams := map[string]string{
		"id": "1",
	}
	req := httptest.NewRequest(http.MethodDelete, "/1", nil)
	req = mux.SetURLVars(req, reqPathParams)
	httpWriter := httptest.NewRecorder()
	todoController.DeleteTodo(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusOK {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusOK, httpWriter.Code)
	}
	var expectedResponse = "\"Todo Deleted Successfully\""
	require.JSONEq(t, expectedResponse, string(data))
}

func TestUpdateTodoTodoUpdatedSuccessfully(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	var mockUpdatedTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a lemon cake for tomorrow's fate",
		Completed: true,
	}
	testObj.On("UpdateTodo", mockUpdatedTodo).Return(mockUpdatedTodo, nil)
	setupTodoController(testObj)
	mockTodoJson, _ := json.Marshal(mockUpdatedTodo)
	bodyReader := strings.NewReader(string(mockTodoJson))
	req := httptest.NewRequest(http.MethodPut, "/1", bodyReader)
	reqPathParams := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, reqPathParams)
	httpWriter := httptest.NewRecorder()
	todoController.UpdateTodo(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusOK {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusOK, httpWriter.Code)
	}
	expectedResponse, _ := json.Marshal(mockUpdatedTodo)
	require.JSONEq(t, string(expectedResponse), string(data))
}

func TestUpdateTodoInvalidRequestBody(t *testing.T) {
	testObj := new(MockTodoServiceImpl)

	setupTodoController(testObj)
	bodyReader := strings.NewReader("{invalid:jsonsjs}}")
	req := httptest.NewRequest(http.MethodPatch, "/", bodyReader)
	httpWriter := httptest.NewRecorder()
	todoController.UpdateTodo(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusInternalServerError {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusInternalServerError, httpWriter.Code)
	}
	var expectedResponse = "\"Internal Server Error\""
	require.JSONEq(t, expectedResponse, string(data))
}

func TestUpdateTodoTodoNotFoundCreatedSuccessfully(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	var mockUpdatedTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a lemon cake for tomorrow's fate",
		Completed: true,
	}
	testObj.On("UpdateTodo", mockUpdatedTodo).Return(models.Todo{}, errors.New("could not find todo with id [1]"))
	testObj.On("CreateNewTodo", mockUpdatedTodo).Return(mockUpdatedTodo, nil)
	setupTodoController(testObj)
	mockTodoJson, _ := json.Marshal(mockUpdatedTodo)
	bodyReader := strings.NewReader(string(mockTodoJson))
	req := httptest.NewRequest(http.MethodPut, "/1", bodyReader)
	reqPathParams := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, reqPathParams)
	httpWriter := httptest.NewRecorder()
	todoController.UpdateTodo(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusCreated {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusCreated, httpWriter.Code)
	}
	expectedResponse, _ := json.Marshal(mockUpdatedTodo)
	require.JSONEq(t, string(expectedResponse), string(data))
}

func TestUpdateTodoTodoNotFoundCreatedUnSuccessfullyDueToDuplicateId(t *testing.T) {
	testObj := new(MockTodoServiceImpl)
	var mockUpdatedTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a lemon cake for tomorrow's fate",
		Completed: true,
	}
	testObj.On("UpdateTodo", mockUpdatedTodo).Return(models.Todo{}, errors.New("could not find todo with id [1]"))
	testObj.On("CreateNewTodo", mockUpdatedTodo).Return(models.Todo{}, errors.New("todo with id [1] already exists"))
	setupTodoController(testObj)
	mockTodoJson, _ := json.Marshal(mockUpdatedTodo)
	bodyReader := strings.NewReader(string(mockTodoJson))
	req := httptest.NewRequest(http.MethodPut, "/1", bodyReader)
	reqPathParams := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, reqPathParams)
	httpWriter := httptest.NewRecorder()
	todoController.UpdateTodo(httpWriter, req)
	res := httpWriter.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	if httpWriter.Code != http.StatusConflict {
		t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", http.StatusConflict, httpWriter.Code)
	}
	var expectedResponse = "\"Todo with id [1] already exists\""
	require.JSONEq(t, expectedResponse, string(data))
}
