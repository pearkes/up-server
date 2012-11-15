package up

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Routes //

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	resp := BaseResponse{
		Message: "I will keep you up.",
		Url:     Url{Url: "github.com/pearkes/up"},
	}
	// Handle 404's
	if r.URL.Path != "/" {
		abort404(w, r)
	}
	success200(r)
	writeJson(w, resp)
}

func UrlsHandler(w http.ResponseWriter, r *http.Request) {
	resp := BaseResponse{
		Urls: UrlsResponse(),
	}
	// Handle 404's
	if r.URL.Path != "/urls" {
		abort404(w, r)
		return
	}
	success200(r)
	writeJson(w, resp)
}

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	// Type conversion
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	// If something other then an int is sent, abort
	if err != nil {
		fmt.Println("Failed to convert to useable id: " + string(id))
	}
	urlresponse := UrlResponse(id)
	// If the object doesn't exist.
	if urlresponse.Id == 0 {
		abort404(w, r)
		return
	}
	resp := BaseResponse{
		Url: urlresponse,
	}
	success200(r)
	writeJson(w, resp)
}

func AddUrlHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println("Failed to parse URL from request body: " + err.Error())
		abort400(w, r)
		return
	}
	url := &Url{}
	_ = json.Unmarshal(body, &url)
	// TODO move this into a response function
	newurl, err := addUrl(url.Url)
	if err != nil {
		fmt.Println("Failed to add url: " + err.Error())
		abort400(w, r)
		return
	}
	resp := BaseResponse{
		Url: newurl,
	}
	success200(r)
	writeJson(w, resp)
}

// Initalize the web server
func initServer() {
	// The Base Router
	r := mux.NewRouter()
	// Unauthenticated index page
	r.HandleFunc("/", HomeHandler).Methods("GET")

	// Setup an authenticated subrouter
	s := r.Headers("Content-Type", "application/json").
		Headers("X-Up-Auth", getSecret()).
		Subrouter()

	// The authenticated routes
	s.HandleFunc("/urls", UrlsHandler).Methods("GET")
	s.HandleFunc("/url/{id:[0-9]+}", UrlHandler).Methods("GET", "DELETE")
	s.HandleFunc("/url", AddUrlHandler).Methods("POST")

	// Pass the base router to the http handler
	http.Handle("/", r)
	// Get the port from the environment and serve the people
	port := getPort()
	fmt.Println("Starting web service...")
	http.ListenAndServe(port, nil)
}
