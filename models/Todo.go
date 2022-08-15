package models

type Todo struct {
	Id        string `json:"Id"`
	Title     string `json:"Title"`
	Desc      string `json:"Desc"`
	Completed bool   `json:"Completed"`
}
