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
	tests := map[string]struct {
		todoId           string
		expectedCode     int
		expectedResponse interface{}
		mockSetup        func(mockedComponent *MockTodoServiceImpl)
	}{
		"Delete successful": {
			todoId:           "1",
			expectedCode:     http.StatusOK,
			expectedResponse: "Todo Deleted Successfully",
			mockSetup: func(mockedComponent *MockTodoServiceImpl) {
				mockedComponent.On("DeleteTodo", "1").Return()
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockTodoService := new(MockTodoServiceImpl)
			tt.mockSetup(mockTodoService)
			setupTodoController(mockTodoService)
			reqPathParams := map[string]string{
				"id": tt.todoId,
			}
			req := httptest.NewRequest(http.MethodDelete, "/"+tt.todoId, nil)
			req = mux.SetURLVars(req, reqPathParams)
			httpWriter := httptest.NewRecorder()
			todoController.DeleteTodo(httpWriter, req)
			res := httpWriter.Result()
			defer res.Body.Close()
			data, _ := io.ReadAll(res.Body)
			if httpWriter.Code != tt.expectedCode {
				t.Errorf("unexpected HTTP response code, expected [%v] but recieved [%v]", tt.expectedCode, httpWriter.Code)
			}
			expectedResponse, _ := json.Marshal(tt.expectedResponse)
			require.JSONEq(t, string(expectedResponse), string(data))
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	tests := map[string]struct {
		todoId           string
		requestBody      interface{}
		expectedCode     int
		expectedResponse interface{}
		mockSetup        func(mockedComponent *MockTodoServiceImpl)
	}{
		"Invalid Request Body": {
			requestBody:      "{invalid:json}}",
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: "Internal Server Error",
		},
		"Todo Not Found Is Not Created Successfully Due To Duplicate Id": {
			todoId: "1",
			requestBody: models.Todo{
				Id:        "1",
				Title:     "Bake cake",
				Desc:      "Bake a carrot cake for tomorrow's fate",
				Completed: false,
			},
			expectedCode:     http.StatusConflict,
			expectedResponse: "Todo with id [1] already exists",
			mockSetup: func(mockedComponent *MockTodoServiceImpl) {
				mockedComponent.On("UpdateTodo", mock.Anything).
					Return(models.Todo{}, errors.New("could not find todo with id [1]"))
				mockedComponent.On("CreateNewTodo", mock.Anything).
					Return(models.Todo{}, errors.New("todo with id [1] already exists"))
			},
		},
		"Todo Not Found Is Then Created Successfully": {
			todoId: "1",
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
				mockedComponent.On("UpdateTodo", mock.Anything).
					Return(models.Todo{}, errors.New("could not find todo with id [1]"))
				mockedComponent.On("CreateNewTodo", mock.Anything).
					Return(models.Todo{
						Id:        "1",
						Title:     "Bake cake",
						Desc:      "Bake a carrot cake for tomorrow's fate",
						Completed: false,
					}, nil)
			},
		},
		"Todo Updated Successfully": {
			requestBody: models.Todo{
				Id:        "1",
				Title:     "Bake cake",
				Desc:      "Bake a carrot cake for tomorrow's fate",
				Completed: false,
			},
			expectedCode: http.StatusOK,
			expectedResponse: models.Todo{
				Id:        "1",
				Title:     "Bake cake",
				Desc:      "Bake a carrot cake for tomorrow's fate",
				Completed: false,
			},
			mockSetup: func(mockedComponent *MockTodoServiceImpl) {
				mockedComponent.On("UpdateTodo", mock.Anything).Return(models.Todo{
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
			req := httptest.NewRequest(http.MethodPut, "/"+tt.todoId, bodyReader)
			reqPathParams := map[string]string{
				"id": tt.todoId,
			}
			req = mux.SetURLVars(req, reqPathParams)
			httpWriter := httptest.NewRecorder()
			todoController.UpdateTodo(httpWriter, req)
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
