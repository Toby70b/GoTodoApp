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
	response := models.NewResponse(http.StatusOK, Articles)
	returnJsonResponse(writer, response)
}

func returnSingleArticle(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["id"]
	for _, article := range Articles {
		if article.Id == key {
			response := models.NewResponse(http.StatusOK, article)
			returnJsonResponse(writer, response)
		}
	}
}

func createNewArticle(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: createNewArticle")
	reqBody, _ := io.ReadAll(request.Body)
	var article models.Article
	json.Unmarshal(reqBody, &article)
	Articles = append(Articles, article)
	response := models.NewResponse(http.StatusCreated, article)
	returnJsonResponse(writer, response)
}

func deleteArticle(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: deleteArticle")
	vars := mux.Vars(request)
	id := vars["id"]

	for i, article := range Articles {
		if article.Id == id {
			//Articles equals all values before index (remember slices don't include value at the max index specified)
			//Plus all the values one index after the found index (remember slices do include the value at the min index)
			//the ... will pass the slice to the variadic function
			Articles = append(Articles[:i], Articles[i+1:]...)
		}
	}

	response := models.NewResponse(http.StatusOK, nil)
	returnJsonResponse(writer, response)
}

func updateArticle(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: updateArticle")

	vars := mux.Vars(request)
	key := vars["id"]

	reqBody, _ := io.ReadAll(request.Body)
	var requestArticle models.Article
	json.Unmarshal(reqBody, &requestArticle)

	for i, article := range Articles {
		if article.Id == key {
			Articles[i] = requestArticle
			response := models.NewResponse(http.StatusOK, Articles[i])
			returnJsonResponse(writer, response)
		}
	}

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")
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
