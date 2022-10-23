package controllers

import (
	"TodoApp/src/main/models"
	"TodoApp/src/main/services"
	"TodoApp/src/main/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

// A TodoController represents a REST controller for handling HTTP requests to the API under the "todo/" URI
type TodoController struct {
	todoService services.TodoService
}

// NewTodoController creates a new TodoController object. This is used by Wire when starting the API to perform the
// necessary dependency injection
func NewTodoController(todoService services.TodoService) TodoController {
	return TodoController{todoService}
}

// ReturnAllTodos returns all todos items persisted within the DB
func (controller *TodoController) ReturnAllTodos(writer http.ResponseWriter, _ *http.Request) {
	fmt.Println("Endpoint Hit: returnAllTodos")
	todos := controller.todoService.ReturnAllTodos()
	utils.ReturnJsonResponse(writer, http.StatusOK, todos)
}

// ReturnSingleTodo returns a single todo item persisted within the DB with an id matching the id passed as a path parameter.
// The path param is accessed via the map within request parameter. If an existing todo item with an id matching that of
// the new todo item is not found, and error will be returned instead
func (controller *TodoController) ReturnSingleTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: returnSingleTodo")
	vars := mux.Vars(request)
	todoId := vars["id"]
	todo, err := controller.todoService.ReturnSingleTodo(todoId)

	if err != nil {
		log.Println(err.Error())
		utils.ReturnJsonResponse(writer, http.StatusNotFound, fmt.Sprintf("could not find todo with id [%s]", todoId))
	} else {
		utils.ReturnJsonResponse(writer, http.StatusOK, todo)

	}
}

// CreateNewTodo creates a new todo item and persist it within the DB. If an existing todo item with an id matching
// that of the new todo item is found, and error will be returned instead
func (controller *TodoController) CreateNewTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: createNewTodo")
	reqBody, _ := io.ReadAll(request.Body)
	var todo models.Todo
	err := json.Unmarshal(reqBody, &todo)
	if err != nil {
		log.Println("Error deserializing the request", err)
		utils.ReturnJsonResponse(writer, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	response, err := controller.todoService.CreateNewTodo(todo)
	if err != nil {
		log.Println(err.Error())
		utils.ReturnJsonResponse(writer, http.StatusConflict, fmt.Sprintf("Todo with id [%s] already exists", todo.Id))
	} else {
		utils.ReturnJsonResponse(writer, http.StatusCreated, response)
	}
}

// DeleteTodo removes a todo item persisted within the DB with an id matching the id passed as a path parameter.
// The path param is accessed via the map within request parameter
func (controller *TodoController) DeleteTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: deleteTodo")
	vars := mux.Vars(request)
	todoId := vars["id"]
	controller.todoService.DeleteTodo(todoId)
	utils.ReturnJsonResponse(writer, http.StatusOK, "Todo Deleted Successfully")
}

// UpdateTodo modifies an existing todo item with the details from the todo item passed in the request. If an existing
// todo item with an id matching that of the new todo item is not found, and error will be returned instead
func (controller *TodoController) UpdateTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: updateTodo")
	reqBody, _ := io.ReadAll(request.Body)
	var todo models.Todo
	err := json.Unmarshal(reqBody, &todo)
	if err != nil {
		log.Println("Error deserializing the request", err)
		utils.ReturnJsonResponse(writer, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	response, err := controller.todoService.UpdateTodo(todo)
	if err != nil {
		log.Printf("Failed to find existing todo item with id [%v] attempting to create new todo item\n", todo.Id)
		response, err := controller.todoService.CreateNewTodo(todo)
		if err != nil {
			log.Println(err.Error())
			utils.ReturnJsonResponse(writer, http.StatusConflict, fmt.Sprintf("Todo with id [%s] already exists", todo.Id))
		} else {
			utils.ReturnJsonResponse(writer, http.StatusCreated, response)
		}
	} else {
		utils.ReturnJsonResponse(writer, http.StatusOK, response)
	}
}

// HandleRequests initializes a new MUX router to receive requests under the "todo/" URI and handles them by calling
// methods within TodoController
func (controller TodoController) HandleRequests() {
	fmt.Println("Starting TodoController...")
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/todo", controller.CreateNewTodo).Methods("POST")
	myRouter.HandleFunc("/todo", controller.UpdateTodo).Methods("PUT")
	myRouter.HandleFunc("/todo", controller.ReturnAllTodos).Methods("GET")
	myRouter.HandleFunc("/todo/{id}", controller.DeleteTodo).Methods("DELETE")
	myRouter.HandleFunc("/todo/{id}", controller.ReturnSingleTodo).Methods("GET")
	fmt.Println("TodoController Listening...")
	log.Fatalln(http.ListenAndServe(":10000", myRouter))

}
