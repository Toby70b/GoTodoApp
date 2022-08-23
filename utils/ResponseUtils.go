package utils

import (
	"encoding/json"
	"net/http"
)

func ReturnJsonResponse(writer http.ResponseWriter, httpCode int, responseBody any) {
	writer.Header().Set("Content-Type", "application/json")
	response := responseBody
	writer.WriteHeader(httpCode)

	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		panic(err)
	}
}
