package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	log.Printf("Starting application on port: %s", port)

	router := mux.NewRouter()
	router.HandleFunc("/info", showInfo)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func showInfo(w http.ResponseWriter, r *http.Request) {

	log.Printf("Show application information to: %s", r.RemoteAddr)

	if r.Header.Get("accept") != "application/json" {
		respondHTTPError(406, "Unsupported Media!", w)
		return
	}

	m := map[string]string {"applicationName": "Mortgage Calculator", "applicationVersion": "0.0.1"}
	j, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(j))

}

func respondHTTPError(c int, m string, w http.ResponseWriter) {
	w.WriteHeader(c)
	w.Write([]byte(m))
}




