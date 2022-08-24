package services

import (
	"TodoApp/models"
	"errors"
	"fmt"
)

// A TodoService represents a Service class responsible for functionality relating to Todo items
//
// Contains an array Todos which acts as a in-memory DB for persisting Todo items
type TodoService struct {
	Todos []models.Todo
}

// NewTodoService creates a new TodoController object. This is used by Wire when starting the API to perform the
// necessary dependency injection
func NewTodoService(todos []models.Todo) TodoService {
	return TodoService{todos}
}

// ReturnAllTodos returns all Todo items currently persisted within the DB
func (service *TodoService) ReturnAllTodos() []models.Todo {
	return service.Todos
}

// ReturnSingleTodo returns a single Todo item, identified via the id param. If no Todo item is found with a matching
// Id then an error is returned
func (service *TodoService) ReturnSingleTodo(id string) (models.Todo, error) {
	for _, todo := range service.Todos {
		if todo.Id == id {
			return todo, nil
		}
	}
	return models.Todo{}, errors.New(fmt.Sprintf("could not find todo with id [%s]", id))
}

// CreateNewTodo persists a new Todo item in the DB. If a existing Todo item with an id matching that of the new Todo item
// is found within the DB then an error will be returned
// The Todo item passed as a parameter must include an id
func (service *TodoService) CreateNewTodo(newTodo models.Todo) (models.Todo, error) {
	err := validateTodo(newTodo)
	if err != nil {
		return models.Todo{}, err
	}

	fmt.Println("Endpoint Hit: createNewTodo")
	for _, todo := range service.Todos {
		if todo.Id == newTodo.Id {
			return models.Todo{}, errors.New(fmt.Sprintf("Todo with id [%s] already exists", newTodo.Id))
		}
	}
	service.Todos = append(service.Todos, newTodo)
	return newTodo, nil
}

// DeleteTodo removes a Todo item from the DB with an id matching that of the id provided as a parameter
func (service *TodoService) DeleteTodo(id string) {
	fmt.Println("Endpoint Hit: deleteTodo")
	for i, todo := range service.Todos {
		if todo.Id == id {
			//Todos equals all values before index (remember slices don't include value at the max index specified)
			//Plus all the values one index after the found index (remember slices do include the value at the min index)
			//the ... will pass the slice to the variadic function
			service.Todos = append(service.Todos[:i], service.Todos[i+1:]...)
		}
	}
}

// UpdateTodo updates a Todo item with a id matching that of the Todo item pass as a parameter. If a Todo item with
// an id matching that of the Todo item passed as a parameter cannot be found then an error will be returned.
//
// The Todo item passed as a parameter must include an id
func (service *TodoService) UpdateTodo(newTodo models.Todo) (models.Todo, error) {
	err := validateTodo(newTodo)
	if err != nil {
		return models.Todo{}, err
	}
	fmt.Println("Endpoint Hit: updateTodo")
	for i, todo := range service.Todos {
		if todo.Id == newTodo.Id {
			service.Todos[i] = newTodo
			return newTodo, nil
		}
	}
	return models.Todo{}, errors.New(fmt.Sprintf("could not find todo with id [%s]", newTodo.Id))
}

// validateTodo applies validation rules against a Todo object to confirm it is valid
func validateTodo(todo models.Todo) error {
	if todo.Id == "" {
		return errors.New("todo Id cannot be null")
	}
	return nil
}
