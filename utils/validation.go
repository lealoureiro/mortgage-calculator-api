package utils

import "net/http"

// IsClientAcceptingJSON : validate the http client is requestion JSON
func IsClientAcceptingJSON(r *http.Request) bool {
	return r.Header.Get("accept") == "application/json"
}

// IsContentTypeJSON : validate if the content of request is JSON
func IsContentTypeJSON(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}
