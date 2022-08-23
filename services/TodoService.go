package services

import (
	"TodoApp/models"
	"errors"
	"fmt"
)

type TodoService struct {
	Todos []models.Todo
}

func NewTodoService(todos []models.Todo) TodoService {
	return TodoService{todos}
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
	return models.Todo{}, errors.New(fmt.Sprintf("could not find todo with id [%s]", id))
}

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

func validateTodo(todo models.Todo) error {
	if todo.Id == "" {
		return errors.New("todo Id cannot be null")
	}
	return nil
}
