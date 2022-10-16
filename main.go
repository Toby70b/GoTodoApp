package main

import (
	"fmt"
)

func main() {
	fmt.Println("Rest API v1.0 - Mux Routers")
	todoController := InitializeTodoController()
	todoController.HandleRequests()
}
