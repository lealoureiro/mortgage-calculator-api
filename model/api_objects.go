package model

// Marshaler : interface to allow a type to marshalled to JSON
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// MonthlyPaymentsRequest : calculate monthly payment request object
// swagger:model
type MonthlyPaymentsRequest struct {

	// Initial Principal lended from bank
	//
	// required: true
	// example: 200000
	InitialPrincipal float64 `json:"initialPrincipal"`

	// Initial Market Value of the property
	//
	// required: true
	// example: 210000
	MarketValue *Number `json:"marketValue"`

	// Initial Interest Rate of the Mortgage
	//
	// required: true
	// example: 2.0
	InitialInterestRate *Number `json:"initialInterestRate"`

	// Number of months to pay back the Mortgage
	//
	// required: true
	// example: 360
	Months int `json:"months"`

	// Start date of the Mortgage
	//
	// required: true
	// example: 2020-01-20
	StartDate *JSONTime `json:"startDate"`

	// Income Tax of the Mortgage payer in %, used to calculate interest tax benefit
	//
	// required: true
	// example: 40
	IncomeTax int `json:"incomeTax"`

	// Indication if Bank updates the interest based on LoanToValue ration
	//
	// required: true
	// example: false
	AutomaticInterestUpdate bool `json:"automaticInterestUpdate"`

	// Loan To Value Interest Rate Tiers, need to be provided in case **automaticInterestUpdate** is **true**
	LoanToValueInterestTiers []LoanToValueInterestTier `json:"loanToValueInterestTiers"`

	// List of Interest Tier updates during the Mortgage period
	InterestTierUpdates []InterestTierUpdate `json:"interestTierUpdates"`

	// List of extra payments during the Mortgage period
	Repayments []Repayment `json:"repayments"`
}

// LoanToValueInterestTier : Represents a Interest Tier aka ‘loan-to-value ratio’
// swagger:model
type LoanToValueInterestTier struct {

	// Loan To Value ratio percentage in %
	//
	// required: true
	// example: 90
	Percentage float64 `json:"percentage"`

	// Interest rate for this Loan To Value ratio
	//
	// required: true
	// example: 1.95
	Interest float64 `json:"interest"`
}

// InterestTierUpdate : Represents an interest update for a certain month
// swagger:model
type InterestTierUpdate struct {

	// Date when the financial instituion register the new update
	//
	// required: true
	// example: 2020-11-10
	UpdateDate JSONTime `json:"updateDate"`

	// Market value in the moment of the update
	//
	// required: true
	// example: 225000
	MarketValue *Number `json:"marketValue"`

	// Current interest rate taking in account new LoanToValue ratio after update
	//
	// required: true
	// example: 1.70
	Interest *Number `json:"interest"`
}

// Repayment : a extra repayment during the mortgage period
// swagger:model
type Repayment struct {

	// Date of the extra repayment
	//
	// required: true
	// example: 2023-01-20
	Date JSONTime `json:"date"`

	// Amount of the extra repayment
	//
	// required: true
	// example: 1000
	Amount float64 `json:"amount"`
}

// MonthPayment : a Mortgage monthly payment
type MonthPayment struct {

	// number of the month
	//
	// required: true
	// example: 1
	Month int `json:"month"`

	// the date when the amount will paid/debited
	//
	// required: true
	// example: 2021-03-01
	PaymentDate JSONTime `json:"paymentDate"`

	// amount of the reapayment for this monthly payment
	//
	// required: true
	// example: 500.00
	Repayment Number `json:"repayment"`

	// gross interest amount
	//
	// required: true
	// example: 300.00
	InterestGrossAmount Number `json:"interestGrossAmount"`

	// Net interest amount
	//
	// required: true
	// example: 150.00
	InterestNetAmount Number `json:"interestNetAmount"`

	// interest percentage used to calculate this monthly payment
	//
	// required: true
	// example: 1.89
	InterestPercentage Number `json:"interestPercentage"`

	// remaining principal of the mortgage
	//
	// required: true
	// example: 189000
	Principal Number `json:"principal"`

	// total gross amount of the monthly payment
	//
	// required: true
	// example: 800.00
	TotalGross Number `json:"totalGross"`

	// total net amount of the montly payment
	//
	// required: true
	// example: 650.00
	TotalNet Number `json:"totalNet"`

	// loan-to-value after this montly payment
	//
	// required: true
	// example 0.87
	LoanToValueRatio Number `json:"loanToValueRatio"`

	// current market value of the property
	//
	// required: true
	// example: 245000
	MarketValue Number `json:"marketValue"`
}

// MonthlyPayments : the response model of Mortgage monthly payments operation
type MonthlyPayments struct {

	// List of monthly payments of your proposed mortgage
	//
	// required: true
	Payments []MonthPayment `json:"payments"`

	// Total amount of Gross Interest paid during the whole Mortgage
	//
	// required: true
	TotalGrossInterest Number `json:"totalGrossInterest"`

	// Total amount of Net Interest paid during the whole Mortgage
	//
	// required: true
	TotalNetInterest Number `json:"totalNetInterest"`
}

// BadRequest : the response model to hold response for a bad request
type BadRequest struct {
	ErrorMessage string `json:"errorMessage"`
}

// Info : the response model to show application info
type Info struct {

	// application name
	//
	// required: true
	// example: MortgageCalculatorAPI
	ApplicationName string `json:"applicationName"`

	// application version
	//
	// required: true
	// example: v0.0.1-1232131
	ApplicationVersion string `json:"applicationVersion"`
}
