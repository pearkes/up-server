package up

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Helpers //
//Helper to encode JSON responses and catch encoding errors
func encodeJson(r BaseResponse) string {
	j, err := json.MarshalIndent(r, "", "  ")

	// Catch JSON encoding errors
	if err != nil {
		fmt.Println("error encoding json:", err)
	}

	return string(j)
}

// 200 HTTP
func success200(r *http.Request) {
	fmt.Println("200 " + r.URL.Path)
}

// 404 HTTP
func abort404(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
	fmt.Println("404 " + r.URL.Path)
	return
}

// 400 HTTP
func abort400(w http.ResponseWriter, r *http.Request) {
	status_code := http.StatusBadRequest
	http.Error(w, http.StatusText(status_code), status_code)
	fmt.Println("400 " + r.URL.Path)
	return
}

// Return JSON response to the ResponseWriter

func writeJson(w http.ResponseWriter, resp BaseResponse) {
	fmt.Fprintf(w, encodeJson(resp))
}
