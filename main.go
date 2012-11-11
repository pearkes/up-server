package main

import (
	"code.google.com/p/gorest"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beedb"
	"github.com/bmizerany/pq"
	"net/http"
	"os"
	"strings"
	"time"
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
	Id        int       `PK`
	Url       string    `json:"url,omitempty"`
	Checks    int       `json:"checks,omitempty"`
	LastCheck time.Time `json:"last_check,omitempty"`
}

// The base response object
type BaseResponse struct {
	Message string `json:"message,omitempty"`
	Error   bool   `json:"error,omitempty"`
	Urls    []Url  `json:"urls,omitempty"`
}

// ORM //
// Database init

var orm beedb.Model

func openDb() *sql.DB {
	var database_url = getDatabaseUrl()
	fmt.Println("DATABASEURL: " + database_url)
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
	// Probe the table, if not, create it.
	_, err := getUrl(0)
	// Hacky check to know when to make a non-existant table
	if strings.Contains(err.Error(), "does not exist") {
		// Make the table
		// In the future, this should probably be done in a safer way.
		// I don't really know how we might do this better, so leaving
		// for now.
		db.Exec("CREATE TABLE url ( id SERIAL NOT NULL, url varchar NOT NULL, checks int, last_check date, CONSTRAINT url_pkey PRIMARY KEY (id) ) WITH (OIDS=FALSE);")
	}
}

// Start the service
func main() {
	// Initialize the orm
	initOrm()
	// Register the RESTful service
	gorest.RegisterService(new(UpService))
	// Handle HTTP requests via the gorest service
	http.Handle("/", gorest.Handle())
	// Serve the people
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

// Insert a URL into the database
func addUrl(url string) (Url, error) {
	var newurl Url
	newurl.Url = url
	newurl.Checks = 0
	newurl.LastCheck = time.Now()
	err := orm.Save(&newurl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(newurl)
	// Do an error check
	return newurl, err
}

// Get a url
func getUrl(id int) (Url, error) {
	var existurl Url
	err := orm.Where("id=$1", 1).Find(&existurl)
	if err != nil {
		fmt.Println(err)
	}
	return existurl, err
}

// Get all of the urls in the system
func getUrls() ([]Url, error) {
	var urls []Url
	err := orm.FindAll(&urls)
	return urls, err
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
