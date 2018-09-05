package utils

import (
	"encoding/json"
	"net/http"
)

func RespondHTTPError(c int, m string, w http.ResponseWriter) {

	data := map[string]string{"errorMessage": m}
	jsonBody, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)

	w.Write([]byte(jsonBody))
}
