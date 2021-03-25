package model

// HTTP Status 200 with application information
// swagger:response infoResponse
type swaggerInfoResponse struct {
	// in:body
	Body Info
}

// HTTP Status 400 when something wrong with request data
// swagger:response badRequestResponse
type swaggerBadRequest struct {
	// in:body
	Body BadRequest
}

// HTTP Status 200 with list of monthly payments
// swagger:response monthlyPaymentsResponse
type swaggerMonthlyPaymentsResponse struct {
	// in:body
	Body MonthlyPayments
}

// HTTP Request body to calculate mortgate monthly payments
// swagger:parameters monthlyPaymentsResquest
type swaggerMonthlyPaymentsRequest struct {
	// in:body
	Body MonthlyPaymentsRequest
}
