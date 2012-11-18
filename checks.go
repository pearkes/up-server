package main

import (
	"fmt"
	"net/http"
	"time"
)

func checkUrl(url Url) (bool, error) {
	resp, err := http.Get(url.Url)
	if err != nil {
		fmt.Println(err)
		// Handle HTTP Errors but "Bad" responses
		if resp != nil {
			fmt.Println(url.Url, "check raised an error with status:", resp.StatusCode)
			// Dispatch Notifier for bad response
			// Record date last checked
			// Increment check counter
			return false, err
		}
		fmt.Println(url.Url, "check raised an error with nil response.")
		// Dispatch Notifier for failed lookup, internal error
		// Record date last_checked
		// Increment check counter
		return false, err
	}
	fmt.Println(url.Url, "check was successful with status:", resp.StatusCode)
	return true, err
	// Record date last_checked
	// Increment check counter
}

func checkUrls() {
	urls, err := getUrls()
	if err != nil {
		fmt.Println("Could not fetch urls for checks")
	}
	for _, url := range urls {
		go checkUrl(url)
	}
}

func initChecks() {
	// Initalize the checks
	c := time.Tick(1 * time.Minute)
	for now := range c {
		go checkUrls()
		fmt.Println("Dispatched URL checks at", now)
	}
}
