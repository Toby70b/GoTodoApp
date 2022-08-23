package main

import (
	"TodoApp/controllers"
	"TodoApp/services"
	"github.com/google/wire"
)

func InitializeTodoController() controllers.TodoController {
	wire.Build(controllers.NewTodoController, services.NewTodoService)
	return controllers.TodoController{}
}
