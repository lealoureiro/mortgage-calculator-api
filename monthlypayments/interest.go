package monthlypayments

import "github.com/lealoureiro/mortgage-calculator-api/model"

// InterestSet : represents a set of interest
type InterestSet interface {
	GetInterest(month int, principal float64) (float64, float64)
}

// LoanToValueInterestSet : list of interest rates based on Loan to Value Ratio
type LoanToValueInterestSet struct {
	marketValue   float64
	interestTiers []model.LoanToValueInterestTier
}

// InterestUpdatesSet : list of interest rates based on updates (changed manually LoanToRation or possible market value change)
type InterestUpdatesSet struct {
	marketValue   float64
	interestTiers []model.InterestTierUpdate
}

// GetInterest : get current interest rate based on Loan to Value ratio
func (s LoanToValueInterestSet) GetInterest(_ int, principal float64) (float64, float64) {

	ratio := principal / s.marketValue * 100

	for _, e := range s.interestTiers {
		if ratio <= e.Percentage {
			return e.Interest / 100, s.marketValue
		}
	}

	return 0.0, s.marketValue

}

// GetInterest : get current interest based on manually interest rate updates
func (s InterestUpdatesSet) GetInterest(month int, _ float64) (float64, float64) {

	interest := 0.0
	marketValue := 0.0

	for _, e := range s.interestTiers {

		if e.Month > month {
			return interest / 100, marketValue
		}

		interest = e.Interest
		marketValue = e.MarketValue
	}

	return interest / 100, marketValue
}
