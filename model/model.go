package model

type MonthlyPaymentRequest struct {
	InitialPrincipal float64
	MarketValue      float64
	Months           int
	IncomeTax        int
	InterestTiers    []InterestTier
}

type InterestTier struct {
	Percentage int
	Interest   float64
}

type MonthPayment struct {
	GrossAmount         float64
	InterestAmountGross float64
	InterestPercentage  float64
	RemainingAmount     float64
	TotalGross          float64
	TotalNet            float64
}

type MonthlyPayments struct {
	Payments           []MonthPayment
	TotalGrossInterest float64
	TotalNetInterest   float64
}
