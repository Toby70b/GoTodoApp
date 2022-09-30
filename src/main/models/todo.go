package models

// Todo a Todo item. Composed of the following fields:
//
// Id: A unique identifier of the todo item
//
// Title: A short description of the todo item
//
// Desc: A longer, more detailed description of the todo item
//
// Completed: boolean value indicating whether the todo item has been completed or not
type Todo struct {
	Id        string `json:"Id"`
	Title     string `json:"Title"`
	Desc      string `json:"Desc"`
	Completed bool   `json:"Completed"`
}
