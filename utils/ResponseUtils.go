package utils

import (
	"encoding/json"
	"net/http"
)

// ReturnJsonResponse sends a HTTP response back to the client, if a body is provided it will be serialized into JSON format.
//
// ReturnJsonResponse receives a writer, used to return the response to the client, a HTTP code to be returned as part of
// the response header, and an optional response body to be included in the response
func ReturnJsonResponse(writer http.ResponseWriter, httpCode int, responseBody any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpCode)
	err := json.NewEncoder(writer).Encode(responseBody)
	if err != nil {
		panic(err)
	}
}
