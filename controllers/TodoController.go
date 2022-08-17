package controllers

import (
	"TodoApp/models"
	"TodoApp/services"
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

func (controller *TodoController) returnAllTodos(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: returnAllTodos")
	response := controller.todoService.ReturnAllTodos()
	returnResponse(writer, response)
}

func (controller *TodoController) returnSingleTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: returnSingleTodo")
	vars := mux.Vars(request)
	todoId := vars["id"]
	todo, err := controller.todoService.ReturnSingleTodo(todoId)

	if err != nil {
		returnResponse(writer, err.Error())
	} else {
		returnResponse(writer, todo)

	}
}

func (controller *TodoController) createNewTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: createNewTodo")
	reqBody, _ := io.ReadAll(request.Body)
	var todo models.Todo
	err := json.Unmarshal(reqBody, &todo)
	if err != nil {
		log.Println("Error deserializing the request", err)
		returnResponse(writer, "Internal Server Error")
	}
	response, err := controller.todoService.CreateNewTodo(todo)
	if err != nil {
		returnResponse(writer, err.Error())
	} else {
		returnResponse(writer, response)
	}
}

func (controller *TodoController) deleteTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: deleteTodo")
	vars := mux.Vars(request)
	todoId := vars["id"]
	controller.todoService.DeleteTodo(todoId)
}

func (controller *TodoController) updateTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: updateTodo")
	reqBody, _ := io.ReadAll(request.Body)
	var todo models.Todo
	err := json.Unmarshal(reqBody, &todo)
	if err != nil {
		log.Println("Error deserializing the request", err)
		returnResponse(writer, "Internal Server Error")
	}
	response, err := controller.todoService.UpdateTodo(todo)
	if err != nil {
		returnResponse(writer, err.Error())
	} else {
		returnResponse(writer, response)
	}

}

func HandleRequests() {
	fmt.Println("Starting TodoController...")
	var todos []models.Todo
	var controller = TodoController{todoService: services.TodoService{Todos: todos}}
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/todo", controller.createNewTodo).Methods("POST")
	myRouter.HandleFunc("/todo", controller.updateTodo).Methods("PUT")
	myRouter.HandleFunc("/todo", controller.returnAllTodos).Methods("GET")
	myRouter.HandleFunc("/todo/{id}", controller.deleteTodo).Methods("DELETE")
	myRouter.HandleFunc("/todo/{id}", controller.returnSingleTodo).Methods("GET")
	fmt.Println("TodoController Listening...")
	log.Fatalln(http.ListenAndServe(":10000", myRouter))

}

func returnResponse(writer http.ResponseWriter, body any) {
	err := json.NewEncoder(writer).Encode(body)
	if err != nil {
		panic(err)
	}
}
