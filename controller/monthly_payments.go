package controller

import (
	"encoding/json"
	"net/http"

	"github.com/lealoureiro/mortgage-calculator-api/model"
	"github.com/lealoureiro/mortgage-calculator-api/monthlypayments"
	"github.com/lealoureiro/mortgage-calculator-api/utils"
	log "github.com/sirupsen/logrus"
)

// MonthlyPayments : REST resource to calculate Mortgage Monthly Payments
func MonthlyPayments(w http.ResponseWriter, r *http.Request) {

	log.Printf("Calculation montly payments client to: %s", r.RemoteAddr)

	if !isMonthlyPaymentsRequestValid(w, r) {
		return
	}

	var rb model.MonthlyPaymentsRequest

	d := json.NewDecoder(r.Body)
	err := d.Decode(&rb)

	if err != nil {
		utils.RespondHTTPError(400, err.Error(), w)
		return
	}

	valid, err2 := monthlypayments.ValidateInputData(rb)

	if !valid {
		log.Printf("Failed to validate input data, reason: %s", err2)
		utils.RespondHTTPError(400, err2, w)
		return
	}

	log.Printf("Calculating monthly payments for %.2f for property with initial market value %.2f to be paid in %d months with automatic updates %v and Extra Repayments %v", rb.InitialPrincipal, rb.MarketValue.AsFloat64(), rb.Months, rb.AutomaticInterestUpdate, rb.Repayments)

	mps := monthlypayments.CalculateLinearMonthlyPayments(rb)

	mpsJSON, err3 := json.Marshal(mps)

	if err3 != nil {
		log.Printf("Error serializing the response when calculating monthly payments, reason: %s", err3)
		utils.RespondHTTPError(500, err3.Error(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(mpsJSON)

}

func isMonthlyPaymentsRequestValid(w http.ResponseWriter, r *http.Request) bool {

	if !utils.IsContentTypeJSON(r) {
		utils.RespondHTTPError(415, "Unsupported Media Type", w)
		return false
	}

	return true
}
