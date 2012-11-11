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

// URL Response
type UrlResponse struct {
	Id        int    `json:"id,omitempty"`
	Url       string `json:"url,omitempty"`
	Checks    int    `json:"checks,omitempty"`
	LastCheck string `json:"last_check,omitempty"`
}

// The base response
type BaseResponse struct {
	Message      string        `json:"message,omitempty"`
	Error        bool          `json:"error,omitempty"`
	UrlsResponse []UrlResponse `json:"urls,omitempty"`
}

// Start the service
func main() {
	gorest.RegisterService(new(UpService))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(getPort(), nil)
}

// Helpers //

//Helper to encode JSON responses and catch encoding errors
func encodeJson(r BaseResponse) string {
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
	urls               gorest.EndPoint `method:"GET" path:"/urls" output:"string"`
}

// Routes //

// The home page route
func (serv UpService) Home() string {
	r := BaseResponse{Message: "I will keep you up."}
	return encodeJson(r)
}

// Builds up the Urls in a response object
func buildUrlsResponse() []UrlResponse {
	urls := make([]UrlResponse, 3)
	urls[0] = UrlResponse{
		Id:     1,
		Url:    "http://google.com",
		Checks: 42,
	}
	urls[1] = UrlResponse{
		Id:     2,
		Url:    "http://facebook.com",
		Checks: 50,
	}
	urls[2] = UrlResponse{
		Id:     3,
		Url:    "http://yahoo.com",
		Checks: 32,
	}
	return urls
}

// The urls page route
func (serv UpService) Urls() string {
	r := BaseResponse{
		UrlsResponse: buildUrlsResponse(),
	}
	return encodeJson(r)
}
