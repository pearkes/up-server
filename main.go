package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beedb"
	"github.com/bmizerany/pq"
	"net/http"
	"os"
	"strings"
)

// Configuration //

// Get the Postgres Database Url.
func getDatabaseUrl() string {
	var database_url = os.Getenv("DATABASE_URL")
	// Set a default database url if there is nothing in the environemnt
	if database_url == "" {
		// Postgres.app uses this variable to set-up postgres.
		user := os.Getenv("USER")
		// Inject the username into the connection string
		database_url = "postgres://" + user + "@localhost/up?sslmode=disable"
		// Let the user know they are using a default
		fmt.Println("--- INFO: No DATABASE_URL env var detected, defaulting to " + database_url)
	}
	conn_str, err := pq.ParseURL(database_url)
	if err != nil {
		panic("Unable to Parse DB URL connection string: " + database_url)
	}
	return conn_str
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
	Id        int    `PK`
	Url       string `json:"url,omitempty"`
	Checks    int    `json:"checks,omitempty"`
	LastCheck string `json:"last_check,omitempty"`
}

// The base response object
type BaseResponse struct {
	Message string `json:"message,omitempty"`
	Error   bool   `json:"error,omitempty"`
	Url     Url    `json:"url,omitempty"`
	Urls    []Url  `json:"urls,omitempty"`
}

// ORM //
// Database init

var orm beedb.Model

func openDb() *sql.DB {
	var database_url = getDatabaseUrl()
	db, err := sql.Open("postgres", database_url)
	if err != nil {
		panic("Unable to open database: " + database_url)
	}
	return db
}

func initOrm() {
	// Connect to the DB
	db := openDb()
	orm = beedb.New(db, "pg")
	orm.SetPK("id")
	// Probe the table, if not, create it.
	_, err := getUrl("0")
	// Hacky check to know when to make a non-existant table
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			// Make the table
			//
			// In the future, this should probably be done in a safer way.
			// I don't really know how we might do this better, so leaving
			// for now. But don't want any extra steps when starting up.
			//
			db.Exec("CREATE TABLE url ( id SERIAL NOT NULL, url varchar NOT NULL, checks int, last_check date, CONSTRAINT url_pkey PRIMARY KEY (id) ) WITH (OIDS=FALSE);")
			fmt.Println("--- INFO: No url table found, creating one...")
		}
	}
}

// Start the service
func main() {
	// Initialize the orm
	initOrm()
	// Initalize the server
	initServer()
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

// Insert a URL into the database
func addUrl(url string) (Url, error) {
	var newurl Url
	newurl.Url = url
	newurl.Checks = 0
	newurl.LastCheck = "01/01/12"
	err := orm.Save(&newurl)
	if err != nil {
		fmt.Println(err)
	}
	// Do an error check
	return newurl, err
}

// Get a url
func getUrl(id string) (Url, error) {
	var existurl Url
	err := orm.Where("id=$1", id).Find(&existurl)
	return existurl, err
}

// Get all of the urls in the system
func getUrls() ([]Url, error) {
	var urls []Url
	err := orm.FindAll(&urls)
	return urls, err
}

// Responses //

// Builds up the Urls in a response object
func UrlsResponse() []Url {
	urls, err := getUrls()
	if err != nil {
		fmt.Println(err)
	}
	return urls
}

// Builds up a Url response object
func UrlResponse(id string) Url {
	// temp for testing, will be real urls from db
	url, err := getUrl(id)
	if err != nil {
		fmt.Println(err)
	}
	return url
}

// Builds up a Add Url response object
func AddUrlResponse(u string) Url {
	// temp for testing, will be real urls from db
	url, err := addUrl(u)
	if err != nil {
		fmt.Println(err)
	}
	return url
}

// Routes //

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	resp := BaseResponse{
		Message: "I will keep you up.",
		Url:     Url{Url: "github.com/pearkes/up"},
	}
	fmt.Println("200 " + r.URL.Path)
	fmt.Fprintf(w, encodeJson(resp))
}

func UrlsHandler(w http.ResponseWriter, r *http.Request) {
	resp := BaseResponse{
		Urls: UrlsResponse(),
	}
	fmt.Println("200 " + r.URL.Path)
	fmt.Fprintf(w, encodeJson(resp))
}

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]
	urlresponse := UrlResponse(id)
	if urlresponse.Id == 0 {
		http.NotFound(w, r)
		fmt.Println("404 " + r.URL.Path)
		return
	}
	resp := BaseResponse{
		Url: urlresponse,
	}
	fmt.Println("200 " + r.URL.Path)
	fmt.Fprintf(w, encodeJson(resp))
}

// Initalize the web server
func initServer() {
	// Register the handlers
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/urls", UrlsHandler)
	http.HandleFunc("/url/", UrlHandler)
	// Serve the people
	fmt.Println("Starting web service...")
	http.ListenAndServe(getPort(), nil)
}
