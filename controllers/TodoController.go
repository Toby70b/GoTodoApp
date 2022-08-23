package controllers

import (
	"TodoApp/models"
	"TodoApp/services"
	"TodoApp/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type TodoController struct {
	todoService services.TodoService
}

func NewTodoController(todoService services.TodoService) TodoController {
	return TodoController{todoService}
}

func (controller *TodoController) returnAllTodos(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: returnAllTodos")
	todos := controller.todoService.ReturnAllTodos()
	utils.ReturnJsonResponse(writer, http.StatusOK, todos)
}

func (controller *TodoController) returnSingleTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: returnSingleTodo")
	vars := mux.Vars(request)
	todoId := vars["id"]
	todo, err := controller.todoService.ReturnSingleTodo(todoId)

	if err != nil {
		utils.ReturnJsonResponse(writer, http.StatusNotFound, err.Error())
	} else {
		utils.ReturnJsonResponse(writer, http.StatusOK, todo)

	}
}

func (controller *TodoController) createNewTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: createNewTodo")
	reqBody, _ := io.ReadAll(request.Body)
	var todo models.Todo
	err := json.Unmarshal(reqBody, &todo)
	if err != nil {
		log.Println("Error deserializing the request", err)
		http.Error(writer, "Internal Server Error", 500)
	}
	response, err := controller.todoService.CreateNewTodo(todo)
	if err != nil {
		utils.ReturnJsonResponse(writer, http.StatusConflict, err.Error())
	} else {
		utils.ReturnJsonResponse(writer, http.StatusCreated, response)

	}
}

func (controller *TodoController) deleteTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: deleteTodo")
	vars := mux.Vars(request)
	todoId := vars["id"]
	controller.todoService.DeleteTodo(todoId)
	utils.ReturnJsonResponse(writer, http.StatusOK, "Todo Deleted Successfully")
}

func (controller *TodoController) updateTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: updateTodo")
	reqBody, _ := io.ReadAll(request.Body)
	var todo models.Todo
	err := json.Unmarshal(reqBody, &todo)
	if err != nil {
		log.Println("Error deserializing the request", err)
		http.Error(writer, "Internal Server Error", 500)
	}
	response, err := controller.todoService.UpdateTodo(todo)
	if err != nil {
		utils.ReturnJsonResponse(writer, http.StatusNotFound, err.Error())
	} else {
		utils.ReturnJsonResponse(writer, http.StatusOK, response)
	}

}

func (controller TodoController) HandleRequests() {
	fmt.Println("Starting TodoController...")
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/todo", controller.createNewTodo).Methods("POST")
	myRouter.HandleFunc("/todo", controller.updateTodo).Methods("PUT")
	myRouter.HandleFunc("/todo", controller.returnAllTodos).Methods("GET")
	myRouter.HandleFunc("/todo/{id}", controller.deleteTodo).Methods("DELETE")
	myRouter.HandleFunc("/todo/{id}", controller.returnSingleTodo).Methods("GET")
	fmt.Println("TodoController Listening...")
	log.Fatalln(http.ListenAndServe(":10000", myRouter))

}
