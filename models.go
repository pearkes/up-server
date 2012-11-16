package main

// Objects //

// URL Object
type Url struct {
	Id        int    `PK`
	Url       string `json:"url,omitempty"`
	Checks    int    `json:"checks,omitempty"`
	LastCheck string `json:"last_check,omitempty"`
}

// The base response object
type BaseResponse struct {
	Message string `json:"message,omitempty"`
	Error   bool   `json:"error,omitempty"`
	Url     Url    `json:"url,omitempty"`
	Urls    []Url  `json:"urls,omitempty"`
}
