package main

import (
	"./model"
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
	router.HandleFunc("/info", showInfo).Methods("GET")
	router.HandleFunc("/monthly-payments", monthlyPayments).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func monthlyPayments(w http.ResponseWriter, r *http.Request) {

	log.Printf("Calculation montly payments client to: %s", r.RemoteAddr)

	if !isMonthlyPaymentRequestValid(w, r) {
		return
	}

	var rb model.MonthlyPaymentRequest

	d := json.NewDecoder(r.Body)
	err := d.Decode(&rb)

	if err != nil {
		respondHTTPError(400, err.Error(), w)
		return
	}

	log.Printf("Calculating monthly payments for %.2f for property with market value %.2f to be paid in %d months.", rb.InitialPrincipal, rb.MarketValue, rb.Months)
}

func showInfo(w http.ResponseWriter, r *http.Request) {

	log.Printf("Show application information to: %s", r.RemoteAddr)

	if !isClientAcceptingJSON(r) {
		respondHTTPError(406, "Unsupported media!", w)
		return
	}

	m := map[string]string{"applicationName": "Mortgage Calculator", "applicationVersion": "0.0.1"}
	j, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(j))
}

func isMonthlyPaymentRequestValid(w http.ResponseWriter, r *http.Request) bool {

	if !isContentTypeJSON(r) {
		respondHTTPError(415, "Unsupported Media Type", w)
		return false
	}

	if !isClientAcceptingJSON(r) {
		respondHTTPError(406, "Unsupported media!", w)
		return false
	}

	return true
}

func isClientAcceptingJSON(r *http.Request) bool {
	return r.Header.Get("accept") == "application/json"
}

func isContentTypeJSON(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}

func respondHTTPError(c int, m string, w http.ResponseWriter) {
	w.WriteHeader(c)
	w.Write([]byte(m))
}
