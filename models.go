package main

import (
	"database/sql"
	"fmt"
	"github.com/pearkes/hood" // ORM
	"strings"
	"time"
)

// Objects //

// URL Object, which can be used in a response
type Url struct {
	Id              int       `sql:"pk"`
	Url             string    `json:"url,omitempty"`
	Checks          int       `json:"checks,omitempty"`
	LastCheck       time.Time `json:"last_check,omitempty"`
	LastCheckStatus int       `json:"last_check_status,omitempty"`
}

// The base response object
type BaseResponse struct {
	Message string `json:"message,omitempty"`
	Url     string `json:"url,omitempty"`
}

// The object used to return urls in a response
type UrlsBaseResponse struct {
	Urls []Url `json:"urls,omitempty"`
}

func initOrm() *hood.Hood {
	fmt.Println("Initializing ORM...")
	// Connect to the DB
	hd, err := hood.Open("postgres", getDatabaseUrl())
	// Create the table
	err = hd.CreateTable(&Url{})
	if err != nil {
		panic(err)
	}
	return hd
}

// Insert a URL into the database
func addUrl(u string) (Url, error) {
	var newurl Url
	newurl.Url = u
	tx := hd.Begin()
	err := tx.Save(&newurl)
	if err != nil {
		fmt.Println(err)
		return newurl, err
	}
	url, err := getUrl(int64(newurl.Id))
	if err != nil {
		fmt.Println(err)
	}
	// Do an error check
	return url, err
}

// Delete a url
func deleteUrl(id int64) (Url, error) {
	existurl, err := getUrl(id)
	_, err = orm.Where("id=$1", id).Delete(&existurl)
	return existurl, err
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
