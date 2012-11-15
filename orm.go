package up

// ORM //

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beedb" // ORM
	"strings"
)

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
