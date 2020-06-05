package model

import "github.com/strongo/decimal"

// MonthlyPaymentsRequest : calculate monthly payment request object
type MonthlyPaymentsRequest struct {
	InitialPrincipal         float64
	MarketValue              float64
	Months                   int
	IncomeTax                int
	AutomaticInterestUpdate  bool
	LoanToValueInterestTiers []LoanToValueInterestTier
	InterestTierUpdates      []InterestTierUpdate
	Repayments               []Repayment
}

// InterestTier : Represents a Interest Tier aka ‘loan-to-value ratio’
type LoanToValueInterestTier struct {
	Percentage float64
	Interest   float64
}

// InterestTierUpdate : Represents an interest update for a certain month
type InterestTierUpdate struct {
	Month    int
	Interest float64
}

// Repayment : a extra repayment during the mortgage period
type Repayment struct {
	Month  int
	Amount float64
}

// MonthPayment : a Mortgage monthly payment
type MonthPayment struct {
	Month               int                 `json:"month"`
	Repayment           decimal.Decimal64p2 `json:"repayment"`
	InterestGrossAmount decimal.Decimal64p2 `json:"interestGrossAmount"`
	InterestNetAmount   decimal.Decimal64p2 `json:"interestNetAmount"`
	InterestPercentage  decimal.Decimal64p2 `json:"interestPercentage"`
	Principal           decimal.Decimal64p2 `json:"principal"`
	TotalGross          decimal.Decimal64p2 `json:"totalGross"`
	TotalNet            decimal.Decimal64p2 `json:"totalNet"`
}

// MonthlyPayments : the response model of Mortgage monthly payments operation
type MonthlyPayments struct {
	Payments           []MonthPayment      `json:"payments"`
	TotalGrossInterest decimal.Decimal64p2 `json:"totalGrossInterest"`
	TotalNetInterest   decimal.Decimal64p2 `json:"totalNetInterest"`
}
