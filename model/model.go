package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/strongo/decimal"
)

// Marshaler : interface to allow a type to marshalled to JSON
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// JSONTime : wrapper time to hold a specific format for a date
//swagger:strfmt date
type JSONTime time.Time

// AsTime : unwrap primitive time type
func (t JSONTime) AsTime() time.Time {
	return time.Time(t)
}

// NewJSONTime : wrap primitive time type to JSONTime
func NewJSONTime(time time.Time) JSONTime {
	return JSONTime(time)
}

// MarshalJSON : convert a time type to string with a specific format YYYY-MM-DD
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
}

// UnmarshalJSON : convert date string with format YYYY-MM-DD to time.Time
func (t *JSONTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	time, err := time.Parse("2006-01-02", s)
	*t = JSONTime(time)
	return
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
	MarketValue *decimal.Decimal64p2 `json:"marketValue"`

	// Initial Interest Rate of the Mortgage
	//
	// required: true
	// example: 2.0
	InitialInterestRate *decimal.Decimal64p2 `json:"initialInterestRate"`

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

	// Income Tax of the Mortgage payer in %
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
	MarketValue *decimal.Decimal64p2 `json:"marketValue"`

	// Current interest rate taking in account new LoanToValue ratio after update
	//
	// required: true
	// example: 1.70
	Interest *decimal.Decimal64p2 `json:"interest"`
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
	Month               int                 `json:"month"`
	PaymentDate         JSONTime            `json:"paymentDate"`
	Repayment           decimal.Decimal64p2 `json:"repayment"`
	InterestGrossAmount decimal.Decimal64p2 `json:"interestGrossAmount"`
	InterestNetAmount   decimal.Decimal64p2 `json:"interestNetAmount"`
	InterestPercentage  decimal.Decimal64p2 `json:"interestPercentage"`
	Principal           decimal.Decimal64p2 `json:"principal"`
	TotalGross          decimal.Decimal64p2 `json:"totalGross"`
	TotalNet            decimal.Decimal64p2 `json:"totalNet"`
	LoanToValueRatio    decimal.Decimal64p2 `json:"loanToValueRatio"`
	MarketValue         decimal.Decimal64p2 `json:"marketValue"`
}

// MonthlyPayments : the response model of Mortgage monthly payments operation
type MonthlyPayments struct {
	Payments           []MonthPayment      `json:"payments"`
	TotalGrossInterest decimal.Decimal64p2 `json:"totalGrossInterest"`
	TotalNetInterest   decimal.Decimal64p2 `json:"totalNetInterest"`
}

// BadRequest : the response model to hold response for a bad request
type BadRequest struct {
	ErrorMessage string `json:"errorMessage"`
}

// Info : the response model to show application info
type Info struct {
	ApplicationName    string `json:"applicationName"`
	ApplicationVersion string `json:"applicationVersion"`
}
