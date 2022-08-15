package main

import (
	"TodoApp/models"
	"TodoApp/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

func returnSingleTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: returnSingleTodo")
}

func createNewTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: createNewTodo")

	//returnJsonResponse(writer, response)
}

func deleteTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: deleteTodo")

	//returnJsonResponse(writer, response)
}

func updateTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: updateTodo")

}

func handleRequests() {
	var todos = []models.Todo{
		{Id: "1", Title: "Hello", Desc: "Article Description", Completed: false},
		{Id: "2", Title: "Hello 2", Desc: "Article Description", Completed: false},
	}
	var controller = TodoController{todoService: services.TodoService{Todos: todos}}
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/todo", createNewTodo).Methods("POST")
	myRouter.HandleFunc("/todo", controller.returnAllTodos)
	myRouter.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")
	myRouter.HandleFunc("/todo/{id}", updateTodo).Methods("PUT")
	myRouter.HandleFunc("/todo/{id}", returnSingleTodo)
	log.Fatalln(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
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
