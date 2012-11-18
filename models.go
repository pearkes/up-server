package main

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beedb" // ORM
	"strings"
)

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
	_, err := getUrl(0)
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

// Insert a URL into the database
func addUrl(u string) (Url, error) {
	var newurl Url
	newurl.Url = u

	newurl.Checks = 0
	newurl.LastCheck = "01/01/12"
	check, err := checkUrl(newurl)
	if check != true {
		return newurl, err
	}
	err = orm.Save(&newurl)
	if err != nil {
		fmt.Println(err)
		return newurl, err
	}
	// Do an error check
	return newurl, err
}

// Delete a url
func deleteUrl(id int64) error {
	existurl, err := getUrl(id)
	_, err = orm.Where("id=$1", id).Delete(&existurl)
	return err
}

// Get a url
func getUrl(id int64) (Url, error) {
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
