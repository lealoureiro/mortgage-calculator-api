package model

import "github.com/strongo/decimal"

type MonthlyPaymentRequest struct {
	InitialPrincipal float64
	MarketValue      float64
	Months           int
	IncomeTax        int
	InterestTiers    []InterestTier
	Repayments       []Repayment
}

type InterestTier struct {
	Percentage float64
	Interest   float64
}

type Repayment struct {
	Month  int
	Amount float64
}

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

type MonthlyPayments struct {
	Payments           []MonthPayment      `json:"payments"`
	TotalGrossInterest decimal.Decimal64p2 `json:"totalGrossInterest"`
	TotalNetInterest   decimal.Decimal64p2 `json:"totalNetInterest"`
}
