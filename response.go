package up

import (
	"fmt"
)

// Responses //

// Builds up the Urls in a response object
func UrlsResponse() []Url {
	urls, err := getUrls()
	if err != nil {
		fmt.Println(err)
	}
	return urls
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

// Builds up a Add Url response object
func AddUrlResponse(u string) Url {
	// temp for testing, will be real urls from db
	url, err := addUrl(u)
	if err != nil {
		fmt.Println(err)
	}
	return url
}
