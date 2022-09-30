package main

import (
	"TodoApp/src/main/controllers"
	"TodoApp/src/main/models"
	"TodoApp/src/main/services"
	"github.com/google/wire"
)

func InitializeTodoController() controllers.TodoController {
	todoServiceImpl := provideTodoServiceImpl()
	todoController := controllers.NewTodoController(todoServiceImpl)
	return todoController
}

// wire.go:

func provideTodoServiceImpl() *services.TodoServiceImpl {
	var todos []models.Todo
	return services.NewTodoServiceImpl(todos)
}

var Set = wire.NewSet(
	provideTodoServiceImpl, wire.Bind(new(services.TodoService), new(*services.TodoServiceImpl)))
