package main

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
	resp := HomeResponse()
	success200(r)
	fmt.Fprintf(w, encodeJson(resp))
}

func UrlsHandler(w http.ResponseWriter, r *http.Request) {
	resp := UrlsResponse()
	success200(r)
	fmt.Fprintf(w, encodeJson(resp))
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
	success200(r)
	writeJson(w, urlresponse)
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
	newurl, err := AddUrlResponse(url.Url)
	if err != nil {
		errresp := BaseResponse{
			Message: err.Error(),
		}
		writeJson(w, errresp)
		abort400(w, r)
		return
	}
	success200(r)
	writeJson(w, newurl)
}

func DeleteUrlHandler(w http.ResponseWriter, r *http.Request) {
	// Type conversion
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	// If something other then an int is sent, abort
	if err != nil {
		fmt.Println("Failed to convert to useable id: " + string(id))
	}
	url, err := DeleteUrlResponse(id)
	if err != nil {
		fmt.Println("Failed to delete url: " + err.Error())
		abort400(w, r)
		return
	}
	if url.Id == 0 {
		abort404(w, r)
		return
	}
	resp := BaseResponse{
		Message: "Successfully deleted url: " + url.Url,
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
	s.HandleFunc("/url/{id:[0-9]+}", UrlHandler).Methods("GET")
	s.HandleFunc("/url/{id:[0-9]+}", DeleteUrlHandler).Methods("DELETE")
	s.HandleFunc("/url", AddUrlHandler).Methods("POST")

	// Pass the base router to the http handler
	http.Handle("/", r)
	// Get the port from the environment and serve the people
	port := getPort()
	fmt.Println("Starting web service...")
	http.ListenAndServe(port, nil)
}
