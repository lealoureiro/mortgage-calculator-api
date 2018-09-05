package controller

import (
	"encoding/json"
	"github.com/lealoureiro/mortgage-calculator-api/model"
	"github.com/lealoureiro/mortgage-calculator-api/monthlypayments"
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

	log.Printf("Calculating monthly payments for %.2f for property with market value %.2f to be paid in %d months with Interest tiers %v", rb.InitialPrincipal, rb.MarketValue, rb.Months, rb.InterestTiers)

	mps := monthlypayments.CalculateLinearMonthlyPayments(rb)

	mpsJson, _ := json.Marshal(mps)

	w.Header().Set("Content-Type", "application/json")

	w.Write(mpsJson)

}

func isMonthlyPaymentsRequestValid(w http.ResponseWriter, r *http.Request) bool {

	if !utils.IsContentTypeJSON(r) {
		utils.RespondHTTPError(415, "Unsupported Media Type", w)
		return false
	}

	return true
}
