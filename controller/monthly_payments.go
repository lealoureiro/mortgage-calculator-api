package controller

import (
	"encoding/json"
	"github.com/lealoureiro/mortgage-calculator-api/model"
	"github.com/lealoureiro/mortgage-calculator-api/utils"
	"log"
	"net/http"
)

func MonthlyPayments(w http.ResponseWriter, r *http.Request) {

	log.Printf("Calculation montly payments client to: %s", r.RemoteAddr)

	if !isMonthlyPaymentsRequestValid(w, r) {
		return
	}

	var rb model.MonthlyPaymentRequest

	d := json.NewDecoder(r.Body)
	err := d.Decode(&rb)

	if err != nil {
		utils.RespondHTTPError(400, err.Error(), w)
		return
	}

	log.Printf("Calculating monthly payments for %.2f for property with market value %.2f to be paid in %d months.", rb.InitialPrincipal, rb.MarketValue, rb.Months)
}

func isMonthlyPaymentsRequestValid(w http.ResponseWriter, r *http.Request) bool {

	if !utils.IsContentTypeJSON(r) {
		utils.RespondHTTPError(415, "Unsupported Media Type", w)
		return false
	}

	if !utils.IsClientAcceptingJSON(r) {
		utils.RespondHTTPError(406, "Unsupported media!", w)
		return false
	}

	return true
}
