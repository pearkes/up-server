package main

import (
	"fmt"
	"net/http"
	// "runtime"
	"time"
)

var checkChan = make(chan bool)

func checkUrl(url Url) {
	resp, err := http.Get(url.Url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		// Handle HTTP Errors but "Bad" responses
		if resp != nil {
			fmt.Println(url.Url, "check raised an error with status:", resp.StatusCode)
			// Dispatch Notifier for bad response
			// Record date last checked
			// Increment check counter
			checkChan <- false
			return
		}
		fmt.Println(url.Url, "check raised an error with nil response.")
		// Dispatch Notifier for failed lookup, internal error
		// Record date last_checked
		// Increment check counter
		checkChan <- false
		return
	}
	fmt.Println(url.Url, "check was successful with status:", resp.StatusCode)
	checkChan <- true
	return
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
		go recieveChecks()
	}
	return
}

func recieveChecks() {
	var recievedCheck bool
	recievedCheck = <-checkChan
	fmt.Println("Recieved a response from channel:", recievedCheck)
}

func setTimer() {
	c := time.Tick(1 * time.Minute)
	for now := range c {
		go checkUrls()
		fmt.Println("Dispatched URL checks at", now)
	}
}

func initChecks() {
	fmt.Println("Starting checks...")
	// Initalize the checks
	go setTimer()
	return
}
