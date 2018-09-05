package main

import (
	"github.com/gorilla/mux"
	"github.com/lealoureiro/mortgage-calculator-api/controller"
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
	r.HandleFunc("/info", controller.ShowInfo).Methods("GET")
	r.HandleFunc("/monthly-payments", controller.MonthlyPayments).Methods("POST")

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

	w.Header().Set("Content-Type", "application/json")

	s.r.ServeHTTP(w, r)

}
