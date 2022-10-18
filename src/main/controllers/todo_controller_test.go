package controllers

import (
	"TodoApp/src/main/models"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
	todoController = NewTodoController(service)
}

func setupMock(mockTodoService *mock.Mock, methodName string, mockResponse interface{}) {
	mockTodoService.On(methodName).Return(mockResponse)
}

func getHttpResponse(t *testing.T, res *http.Response) []byte {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}

	if err != nil {
		t.Errorf("error when reading HTTP response: [%v]", err)
	}
	return data
}

type testMock struct {
	mockedService      mock.Mock
	serviceMockMethods []mockMethodDetails
}

type mockMethodDetails struct {
	name         string
	mockResponse interface{}
}

func TestReturnAllTodos(t *testing.T) {

	tests := map[string]struct {
		expectedCode     int
		expectedResponse []models.Todo
		mockSetup        func(mockedComponent *MockTodoServiceImpl)
	}{
		"Return Single Todo": {
			expectedCode: 200,
			expectedResponse: []models.Todo{
				{
					Id:        "1",
					Title:     "Bake cake",
					Desc:      "Bake a carrot cake for tomorrow's fate",
					Completed: false,
				},
			},
			mockSetup: func(mockedComponent *MockTodoServiceImpl) {
				mockedComponent.On("ReturnAllTodos").Return([]models.Todo{
					{
						Id:        "1",
						Title:     "Bake cake",
						Desc:      "Bake a carrot cake for tomorrow's fate",
						Completed: false,
					},
				})
			},
		},
		"Return Multiple Todos": {
			expectedCode: 200,
			expectedResponse: []models.Todo{
				{
					Id:        "1",
					Title:     "Bake cake",
					Desc:      "Bake a carrot cake for tomorrow's fate",
					Completed: false,
				},
				{
					Id:        "2",
					Title:     "Iron shirts",
					Desc:      "Iron shirts that are in the dryer",
					Completed: false,
				},
				{
					Id:        "3",
					Title:     "Walk dog",
					Desc:      "Walk the dog around the town",
					Completed: false,
				},
			},
			mockSetup: func(mockedComponent *MockTodoServiceImpl) {
				mockedComponent.On("ReturnAllTodos").Return([]models.Todo{
					{
						Id:        "1",
						Title:     "Bake cake",
						Desc:      "Bake a carrot cake for tomorrow's fate",
						Completed: false,
					},
					{
						Id:        "2",
						Title:     "Iron shirts",
						Desc:      "Iron shirts that are in the dryer",
						Completed: false,
					},
					{
						Id:        "3",
						Title:     "Walk dog",
						Desc:      "Walk the dog around the town",
						Completed: false,
					},
				})
			},
		},

		"No Todos Found": {
			expectedCode:     200,
			expectedResponse: []models.Todo{},
			mockSetup: func(mockedComponent *MockTodoServiceImpl) {
				mockedComponent.On("ReturnAllTodos").Return([]models.Todo{})
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockTodoService := new(MockTodoServiceImpl)
			tt.mockSetup(mockTodoService)
			setupTodoController(mockTodoService)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			httpWriter := httptest.NewRecorder()

			todoController.ReturnAllTodos(httpWriter, req)
			res := httpWriter.Result()
			defer res.Body.Close()
			data := getHttpResponse(t, res)
			if httpWriter.Code != tt.expectedCode {
				t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", tt.expectedCode, httpWriter.Code)
			}
			expectedResponse, _ := json.Marshal(tt.expectedResponse)
			require.JSONEq(t, string(expectedResponse), string(data))
		})
	}

}

func TestReturnSingleTodoTodoFound(t *testing.T) {
	mockTodoService := new(MockTodoServiceImpl)
	var mockTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a carrot cake for tomorrow's fate",
		Completed: false,
	}
	mockTodoService.On("ReturnSingleTodo", "1").Return(mockTodo, nil)
	setupTodoController(mockTodoService)
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
	mockTodoService := new(MockTodoServiceImpl)
	mockTodoService.On("ReturnSingleTodo", "999").Return(models.Todo{}, errors.New("could not find todo with id [999]"))
	setupTodoController(mockTodoService)
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
	mockTodoService := new(MockTodoServiceImpl)
	var mockTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a carrot cake for tomorrow's fate",
		Completed: false,
	}
	mockTodoService.On("CreateNewTodo", mockTodo).Return(mockTodo, nil)
	setupTodoController(mockTodoService)

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
	mockTodoService := new(MockTodoServiceImpl)

	setupTodoController(mockTodoService)
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
	mockTodoService := new(MockTodoServiceImpl)
	var mockTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a carrot cake for tomorrow's fate",
		Completed: false,
	}
	mockTodoService.On("CreateNewTodo", mockTodo).Return(models.Todo{}, errors.New("todo with id [1] already exists"))
	setupTodoController(mockTodoService)

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
	mockTodoService := new(MockTodoServiceImpl)
	mockTodoService.On("DeleteTodo", "1").Return()
	setupTodoController(mockTodoService)
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
	mockTodoService := new(MockTodoServiceImpl)
	var mockUpdatedTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a lemon cake for tomorrow's fate",
		Completed: true,
	}
	mockTodoService.On("UpdateTodo", mockUpdatedTodo).Return(mockUpdatedTodo, nil)
	setupTodoController(mockTodoService)
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
	mockTodoService := new(MockTodoServiceImpl)

	setupTodoController(mockTodoService)
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
	mockTodoService := new(MockTodoServiceImpl)
	var mockUpdatedTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a lemon cake for tomorrow's fate",
		Completed: true,
	}
	mockTodoService.On("UpdateTodo", mockUpdatedTodo).Return(models.Todo{}, errors.New("could not find todo with id [1]"))
	mockTodoService.On("CreateNewTodo", mockUpdatedTodo).Return(mockUpdatedTodo, nil)
	setupTodoController(mockTodoService)
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
	mockTodoService := new(MockTodoServiceImpl)
	var mockUpdatedTodo = models.Todo{
		Id:        "1",
		Title:     "Bake cake",
		Desc:      "Bake a lemon cake for tomorrow's fate",
		Completed: true,
	}
	mockTodoService.On("UpdateTodo", mockUpdatedTodo).Return(models.Todo{}, errors.New("could not find todo with id [1]"))
	mockTodoService.On("CreateNewTodo", mockUpdatedTodo).Return(models.Todo{}, errors.New("todo with id [1] already exists"))
	setupTodoController(mockTodoService)
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
