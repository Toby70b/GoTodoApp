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

func TestReturnAllTodos(t *testing.T) {

	tests := map[string]struct {
		expectedCode     int
		expectedResponse []models.Todo
		mockSetup        func(mockedComponent *MockTodoServiceImpl)
	}{
		"Return Single Todo": {
			expectedCode: http.StatusOK,
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
			expectedCode: http.StatusOK,
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
			expectedCode:     http.StatusOK,
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

func TestReturnSingleTodo(t *testing.T) {

	tests := map[string]struct {
		todoIdPathParam  string
		expectedCode     int
		expectedResponse interface{}
		mockSetup        func(mockedComponent *MockTodoServiceImpl)
	}{
		"No Todo With Matching Id Found": {
			todoIdPathParam:  "999",
			expectedCode:     http.StatusNotFound,
			expectedResponse: "Could not find todo with id [999]",
			mockSetup: func(mockedComponent *MockTodoServiceImpl) {
				mockedComponent.On("ReturnSingleTodo", "999").
					Return(models.Todo{}, errors.New("could not find todo with id [999]"))
			},
		},
		"Todo With Matching Id Found": {
			todoIdPathParam: "1",
			expectedCode:    http.StatusOK,
			expectedResponse: models.Todo{
				Id:        "1",
				Title:     "Bake cake",
				Desc:      "Bake a carrot cake for tomorrow's fate",
				Completed: false,
			},
			mockSetup: func(mockedComponent *MockTodoServiceImpl) {
				mockedComponent.On("ReturnSingleTodo", "1").Return(
					models.Todo{
						Id:        "1",
						Title:     "Bake cake",
						Desc:      "Bake a carrot cake for tomorrow's fate",
						Completed: false,
					}, nil)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockTodoService := new(MockTodoServiceImpl)
			tt.mockSetup(mockTodoService)
			setupTodoController(mockTodoService)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.todoIdPathParam, nil)
			reqPathParams := map[string]string{
				"id": tt.todoIdPathParam,
			}
			req = mux.SetURLVars(req, reqPathParams)
			httpWriter := httptest.NewRecorder()

			todoController.ReturnSingleTodo(httpWriter, req)
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

func TestCreateNewTodo(t *testing.T) {

	tests := map[string]struct {
		requestBody      interface{}
		expectedCode     int
		expectedResponse interface{}
		mockSetup        func(mockedComponent *MockTodoServiceImpl)
	}{
		"Request Body Invalid": {
			requestBody:      "{invalid:json}}",
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: "Internal Server Error",
		},
		"Todo With Matching Id Already Exists": {
			requestBody: models.Todo{
				Id:        "1",
				Title:     "Bake cake",
				Desc:      "Bake a carrot cake for tomorrow's fate",
				Completed: false,
			},
			expectedCode:     http.StatusConflict,
			expectedResponse: "Todo with id [1] already exists",
			mockSetup: func(mockedComponent *MockTodoServiceImpl) {
				mockedComponent.On("CreateNewTodo", mock.Anything).
					Return(models.Todo{}, errors.New("todo with id [1] already exists"))
			},
		},
		"Todo Created Successfully": {
			requestBody: models.Todo{
				Id:        "1",
				Title:     "Bake cake",
				Desc:      "Bake a carrot cake for tomorrow's fate",
				Completed: false,
			},
			expectedCode: http.StatusCreated,
			expectedResponse: models.Todo{
				Id:        "1",
				Title:     "Bake cake",
				Desc:      "Bake a carrot cake for tomorrow's fate",
				Completed: false,
			},
			mockSetup: func(mockedComponent *MockTodoServiceImpl) {
				mockedComponent.On("CreateNewTodo", mock.Anything).Return(models.Todo{
					Id:        "1",
					Title:     "Bake cake",
					Desc:      "Bake a carrot cake for tomorrow's fate",
					Completed: false,
				}, nil)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockTodoService := new(MockTodoServiceImpl)

			if tt.mockSetup != nil {
				tt.mockSetup(mockTodoService)
			}

			setupTodoController(mockTodoService)

			mockTodoJson, _ := json.Marshal(tt.requestBody)
			bodyReader := strings.NewReader(string(mockTodoJson))
			req := httptest.NewRequest(http.MethodPost, "/", bodyReader)
			httpWriter := httptest.NewRecorder()
			todoController.CreateNewTodo(httpWriter, req)

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
