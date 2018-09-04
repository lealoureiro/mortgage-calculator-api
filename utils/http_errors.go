package utils

import "net/http"

func RespondHTTPError(c int, m string, w http.ResponseWriter) {
	w.WriteHeader(c)
	w.Write([]byte(m))
}
