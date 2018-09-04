package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lealoureiro/mortgage-calculator-api/model"
	"log"
	"net/http"
	"os"
)

type CORSEnabledRouter struct {
	r *mux.Router
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	log.Printf("Starting application on port: %s", port)

	r := mux.NewRouter()
	r.HandleFunc("/info", showInfo).Methods("GET")
	r.HandleFunc("/monthly-payments", monthlyPayments).Methods("POST")

	http.Handle("/", &CORSEnabledRouter{r})
	http.ListenAndServe(":"+port, nil)
}

func (s *CORSEnabledRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	if r.Method == "OPTIONS" {
		return
	}

	s.r.ServeHTTP(w, r)
}

func monthlyPayments(w http.ResponseWriter, r *http.Request) {

	log.Printf("Calculation montly payments client to: %s", r.RemoteAddr)

	if !isMonthlyPaymentsRequestValid(w, r) {
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

func isMonthlyPaymentsRequestValid(w http.ResponseWriter, r *http.Request) bool {

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
