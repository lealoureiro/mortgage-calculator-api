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
type MonthlyPaymentsRequest struct {
	InitialPrincipal         float64
	MarketValue              *decimal.Decimal64p2 `json:"marketValue"`
	InitialInterestRate      *decimal.Decimal64p2 `json:"initialInterestRate"`
	Months                   int
	StartDate                *JSONTime `json:"startDate"`
	IncomeTax                int
	AutomaticInterestUpdate  bool
	LoanToValueInterestTiers []LoanToValueInterestTier
	InterestTierUpdates      []InterestTierUpdate
	Repayments               []Repayment
}

// LoanToValueInterestTier : Represents a Interest Tier aka ‘loan-to-value ratio’
type LoanToValueInterestTier struct {
	Percentage float64
	Interest   float64
}

// InterestTierUpdate : Represents an interest update for a certain month
type InterestTierUpdate struct {
	UpdateDate  JSONTime
	MarketValue *decimal.Decimal64p2 `json:"marketValue"`
	Interest    *decimal.Decimal64p2 `json:"interest"`
}

// Repayment : a extra repayment during the mortgage period
type Repayment struct {
	Date   JSONTime
	Amount float64
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
