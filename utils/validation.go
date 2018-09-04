package utils

import "net/http"

func IsClientAcceptingJSON(r *http.Request) bool {
	return r.Header.Get("accept") == "application/json"
}

func IsContentTypeJSON(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}
