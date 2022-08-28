package main

import (
	"TodoApp/src/main/controllers"
	"TodoApp/src/main/services"
	"github.com/google/wire"
)

func InitializeTodoController() controllers.TodoController {
	wire.Build(controllers.NewTodoController, services.NewTodoService)
	return controllers.TodoController{}
}
