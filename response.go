package main

import (
	"fmt"
)

// Responses //

func HomeResponse() BaseResponse {
	resp := BaseResponse{
		Message: "I will keep you up.",
		Url:     "github.com/pearkes/up-server",
	}
	return resp
}

// Builds up the Urls in a response object
func UrlsResponse() UrlsBaseResponse {
	urls, err := getUrls()
	if err != nil {
		fmt.Println(err)
	}
	resp := UrlsBaseResponse{
		Urls: urls,
	}
	return resp
}

// Builds up a Url response object
func UrlResponse(id int64) Url {
	// temp for testing, will be real urls from db
	url, err := getUrl(id)
	if err != nil {
		fmt.Println(err)
	}
	return url
}

// Builds up a Url response object
func DeleteUrlResponse(id int64) (Url, error) {
	// temp for testing, will be real urls from db
	url, err := deleteUrl(id)
	if err != nil {
		fmt.Println(err)
	}
	return url, err
}

// Builds up a Add Url response object
func AddUrlResponse(u string) (Url, error) {
	// temp for testing, will be real urls from db
	url, err := addUrl(u)
	if err != nil {
		fmt.Println(err)
	}
	return url, err
}
