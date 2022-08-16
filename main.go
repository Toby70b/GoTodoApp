package main

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
	json.NewEncoder(writer).Encode(response)
}

func (controller *TodoController) returnSingleTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: returnSingleTodo")
	vars := mux.Vars(request)
	todoId := vars["id"]
	response, err := controller.todoService.ReturnSingleTodo(todoId)
	if err != nil {
		json.NewEncoder(writer).Encode(err.Error())
	} else {
		json.NewEncoder(writer).Encode(response)
	}
}

func (controller *TodoController) createNewTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: createNewTodo")
	reqBody, _ := io.ReadAll(request.Body)
	var todo models.Todo
	json.Unmarshal(reqBody, &todo)
	response, err := controller.todoService.CreateNewTodo(todo)
	if err != nil {
		json.NewEncoder(writer).Encode(err.Error())
	} else {
		json.NewEncoder(writer).Encode(response)
	}
}

func (controller *TodoController) deleteTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: deleteTodo")
	vars := mux.Vars(request)
	todoId := vars["todoId"]
	controller.todoService.DeleteTodo(todoId)
}

func (controller *TodoController) updateTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: updateTodo")
	reqBody, _ := io.ReadAll(request.Body)
	var todo models.Todo
	json.Unmarshal(reqBody, &todo)
	response, err := controller.todoService.UpdateTodo(todo)
	if err != nil {
		json.NewEncoder(writer).Encode(err.Error())
	} else {
		json.NewEncoder(writer).Encode(response)
	}

}

func handleRequests() {
	var todos = []models.Todo{
		{Id: "1", Title: "Hello", Desc: "Article Description", Completed: false},
		{Id: "2", Title: "Hello 2", Desc: "Article Description", Completed: false},
	}
	var controller = TodoController{todoService: services.TodoService{Todos: todos}}
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/todo", controller.createNewTodo).Methods("POST")
	myRouter.HandleFunc("/todo", controller.updateTodo).Methods("PUT")
	myRouter.HandleFunc("/todo", controller.returnAllTodos)
	myRouter.HandleFunc("/todo/{id}", controller.deleteTodo).Methods("DELETE")
	myRouter.HandleFunc("/todo/{id}", controller.returnSingleTodo)
	log.Fatalln(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v1.0 - Mux Routers")
	handleRequests()
}

/*
func returnJsonResponse(writer http.ResponseWriter, response models.Response) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(response.ResponseCode())

	if response.Body() != nil {
		err := json.NewEncoder(writer).Encode(response.Body())
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
}
*/
