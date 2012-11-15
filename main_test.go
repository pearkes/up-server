package up

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

// Get the port for testing
var port = os.Getenv("PORT")

// Various error handling for ease of development
func checks() {
	fmt.Println("--- WARNING: The server must be running for the tests to pass.")
}

func TestHome(t *testing.T) {

	var route = "/"
	var url = "http://localhost:" + port + route

	resp, err := http.Get(url)
	if err != nil {
		checks()
		t.Errorf("Fail => %v", err, resp)
	}

}

func TestUrls(t *testing.T) {

	var route = "/urls"
	var url = "http://localhost:" + port + route

	resp, err := http.Get(url)
	if err != nil {
		checks()
		t.Errorf("Fail => %v", err, resp)
	}
}
