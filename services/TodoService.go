package services

import (
	"TodoApp/models"
	"errors"
	"fmt"
)

type TodoService struct {
	Todos []models.Todo
}

func (service *TodoService) ReturnAllTodos() []models.Todo {
	return service.Todos
}

func (service *TodoService) ReturnSingleTodo(id string) (models.Todo, error) {
	for _, todo := range service.Todos {
		if todo.Id == id {
			return todo, nil
		}
	}
	return models.Todo{}, errors.New(fmt.Sprintf("Could not find todo with id [%s]", id))
}

/*
func (service *TodoService) createNewTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: createNewTodo")
	reqBody, _ := io.ReadAll(request.Body)
	var todo models.Todo
	json.Unmarshal(reqBody, &todo)
	Todos = append(Todos, todo)
	response := models.NewResponse(http.StatusCreated, todo)
	returnJsonResponse(writer, response)
}

func (service *TodoService) deleteTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: deleteTodo")
	vars := mux.Vars(request)
	id := vars["id"]

	for i, todo := range Todos {
		if todo.Id == id {
			//Todos equals all values before index (remember slices don't include value at the max index specified)
			//Plus all the values one index after the found index (remember slices do include the value at the min index)
			//the ... will pass the slice to the variadic function
			Todos = append(Todos[:i], Todos[i+1:]...)
		}
	}

	response := models.NewResponse(http.StatusOK, nil)
	returnJsonResponse(writer, response)
}

func (service *TodoService) updateTodo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: updateTodo")

	vars := mux.Vars(request)
	key := vars["id"]

	reqBody, _ := io.ReadAll(request.Body)
	var requestTodo models.Todo
	json.Unmarshal(reqBody, &requestTodo)

	for i, todo := range Todos {
		if todo.Id == key {
			Todos[i] = requestTodo
			response := models.NewResponse(http.StatusOK, Todos[i])
			returnJsonResponse(writer, response)
		}
	}
*/
