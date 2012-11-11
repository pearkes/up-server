package main

import (
	"code.google.com/p/gorest"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/bmizerany/pq"
	"net/http"
	"os"
)

// Configuration //

// Get the Postgres Database Url.
func getDatabase() string {
	var database_url = os.Getenv("DATABASE_URL")
	// Set a default database url if there is nothing in the environemnt
	if database_url == "" {
		database_url = ""
		fmt.Println("--- INFO: No DATABASE_URL env var detected, defaulting to " + database_url)
	}
	return database_url
}

// Get the Port from the environment so we can run on Heroku
func getPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("--- INFO: No PORT env var detected, defaulting to " + port)
	}
	return ":" + port
}

// Get the secret key for interacting with the application
func getSecret() string {
	var secret = os.Getenv("SECRET")
	if secret == "" {
		panic("You must set a SECRET in your environment.\n\n $ export SECRET=foo")
	}
	return secret
}

// Objects //

// URL Object
type Url struct {
	Id        int    `json:"id,omitempty"`
	Url       string `json:"url,omitempty"`
	Checks    int    `json:"checks,omitempty"`
	LastCheck string `json:"last_check,omitempty"`
}

// The base response object
type BaseResponse struct {
	Message string `json:"message,omitempty"`
	Error   bool   `json:"error,omitempty"`
	Urls    []Url  `json:"urls,omitempty"`
}

// Serialize to URL
func serializeUrl(r sql.Row) {
	// Mutate row to Url
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

// Core Logic //

// Retrive URL from the database
func getUrl(id string) {
	// Open db connection

	// Query for the id, get back the object and tranfsform
	// it to a url object

	// close connection
}

// Place a URL in the database

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
func buildUrlsResponse() []Url {
	// temp for testing, will be real urls from db
	urls := make([]Url, 3)
	//
	urls[0] = Url{
		Id:     1,
		Url:    "http://google.com",
		Checks: 42,
	}
	urls[1] = Url{
		Id:     2,
		Url:    "http://facebook.com",
		Checks: 50,
	}
	urls[2] = Url{
		Id:     3,
		Url:    "http://yahoo.com",
		Checks: 32,
	}
	return urls
}

// The urls page route
func (serv UpService) Urls() string {
	r := BaseResponse{
		Urls: buildUrlsResponse(),
	}
	return encodeJson(r)
}
