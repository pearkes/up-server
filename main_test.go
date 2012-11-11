package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

//Get the port for testing
var port = os.Getenv("PORT")

func TestHome(t *testing.T) {

	var route = "/"
	var url = "http://localhost:" + port + route

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("--- WARNING: The server must be running for the tests to pass.")
		t.Errorf("Fail => %v", err, resp)
	}

}
