package main

import (
	"fmt"
	"net/http"
	"time"
)

var checkChan = make(chan Url)

func failCheck(url Url, status int) {
	url.LastCheckStatus = status
	url.LastCheck = time.Now().UTC()
	url.Checks = url.Checks + 1
	// Save Check
	orm.Save(&url)
	checkChan <- url
}

func passCheck(url Url, status int) {
	url.LastCheckStatus = status
	url.LastCheck = time.Now().UTC()
	url.Checks = url.Checks + 1
	// Save Check
	orm.Save(&url)
	checkChan <- url
}

func checkUrl(url Url) {
	resp, err := http.Get(url.Url)
	if err != nil {
		fmt.Println("Check Failed:", err)
		failCheck(url, 0)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			passCheck(url, resp.StatusCode)
			return
		} else {
			failCheck(url, resp.StatusCode)
			return
		}
	} else {
		failCheck(url, 0)
		return
	}
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
	checkedUrl := <-checkChan
	if checkedUrl.LastCheckStatus != 200 {
		// Send a notifier
		fmt.Println("Triggering notifier for", checkedUrl.Url)
		go sendNotifier(checkedUrl)
	}
}

func setTimer() {
	c := time.Tick(10 * time.Second)
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
