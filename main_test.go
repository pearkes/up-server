package main

import (
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
		t.Errorf("Fail => %v", err, resp)
	}

}
