package main

import (
	"fmt"
)

func main() {
	fmt.Println("Rest API v1.0 - Mux Routers")
	todoController := InitializeTodoController()
	todoController.HandleRequests()
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
