package main

import (
	"code.google.com/p/gorest"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Configuration //

//Get the Port from the environment so we can run on Heroku
func getPort() string {
	return ":" + os.Getenv("PORT")
}

//Get the secret key for interacting with the application
func getSecret() string {
	return os.Getenv("SECRET")
}

// Responses //

// The base response
type response struct {
	Message string `json:"message,omitempty"`
	Error   bool   `json:"error,omitempty"`
}

// Start the service
func main() {
	gorest.RegisterService(new(UpService)) //Register our service
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(getPort(), nil)
}

// Helpers //

//Helper to encode JSON responses and catch encoding errors
func encodeJson(r response) string {
	j, err := json.MarshalIndent(r, "", "  ")

	// Catch JSON encoding errors
	if err != nil {
		fmt.Println("error encoding json:", err)
	}

	return string(j)
}

// Service  //

// The main service
type UpService struct {
	gorest.RestService `root:"/" consumes:"application/json" produces:"application/json"`
	home               gorest.EndPoint `method:"GET" path:"/" output:"string"`
}

// Routes //

// The home page service
func (serv UpService) Home() string {
	r := response{
		Message: "hello world",
	}
	return encodeJson(r)
}
