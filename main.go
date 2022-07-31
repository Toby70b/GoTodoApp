package main

import (
	"TodoApp/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

var Articles []models.Article

func homePage(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(writer, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homepage")
}

func returnAllArticles(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(writer).Encode(Articles)
}

func returnSingleArticle(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["id"]
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(writer).Encode(article)
		}
	}
}

func createNewArticle(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: createNewArticle")
	reqBody, _ := io.ReadAll(request.Body)
	var article models.Article
	json.Unmarshal(reqBody, &article)
	Articles = append(Articles, article)
	json.NewEncoder(writer).Encode(article)
}

func deleteArticle(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: deleteArticle")
	vars := mux.Vars(request)
	id := vars["id"]

	for index, article := range Articles {
		if article.Id == id {
			//Articles equals all values before index (remember slices don't include value at the max index specified)
			//Plus all the values one index after the found index (remember slices do include the value at the min index)
			//the ... will pass the slice to the variadic function
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	log.Fatalln(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Articles = []models.Article{
		{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}
