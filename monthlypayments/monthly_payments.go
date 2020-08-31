package monthlypayments

import (
	"fmt"
	"sort"

	"github.com/lealoureiro/mortgage-calculator-api/model"
	"github.com/strongo/decimal"
)

// CalculateLinearMonthlyPayments : calculate the monthly payments for a Linear Mortgage
func CalculateLinearMonthlyPayments(r model.MonthlyPaymentsRequest) model.MonthlyPayments {

	result := make([]model.MonthPayment, 0, r.Months)

	var interestSet InterestSet
	if r.AutomaticInterestUpdate {
		interestSet = LoanToValueInterestSet{r.MarketValue, r.LoanToValueInterestTiers}
	} else {
		interestSet = InterestUpdatesSet{r.MarketValue, r.InterestTierUpdates}
	}

	monthlyRepayment := r.InitialPrincipal / float64(r.Months)
	principal := r.InitialPrincipal

	interestPercentage := 0.0
	interestGrossAmount := 0.0
	interestNetAmount := 0.0
	totalGrossInterest := 0.0
	totalNetInterest := 0.0
	marketValue := 0.0

	incomeTax := float64(r.IncomeTax) / 100.0

	for i := 1; i <= r.Months && principal > 0; i++ {

		if principal < monthlyRepayment {
			monthlyRepayment = principal
		}

		interestPercentage, marketValue = interestSet.GetInterest(i, principal)
		interestGrossAmount = (principal * interestPercentage) / 12.0
		interestNetAmount = interestGrossAmount - (interestGrossAmount * incomeTax)

		totalGrossInterest += interestGrossAmount
		totalNetInterest += interestNetAmount

		var payment = model.MonthPayment{}

		payment.Month = i
		payment.Repayment = decimal.NewDecimal64p2FromFloat64(monthlyRepayment)
		payment.InterestGrossAmount = decimal.NewDecimal64p2FromFloat64(interestGrossAmount)
		payment.InterestNetAmount = decimal.NewDecimal64p2FromFloat64(interestNetAmount)
		payment.Principal = decimal.NewDecimal64p2FromFloat64(principal)
		payment.InterestPercentage = decimal.NewDecimal64p2FromFloat64(interestPercentage * 100)
		payment.TotalGross = decimal.NewDecimal64p2FromFloat64(monthlyRepayment + interestGrossAmount)
		payment.TotalNet = decimal.NewDecimal64p2FromFloat64(monthlyRepayment + interestNetAmount)
		payment.LoanToValueRatio = decimal.NewDecimal64p2FromFloat64(principal / marketValue * 100)
		payment.MarketValue = decimal.NewDecimal64p2FromFloat64(marketValue)

		result = append(result, payment)

		principal -= monthlyRepayment

		processExtraRepayments(r.Repayments, i, &principal)

	}

	return model.MonthlyPayments{
		Payments:           result,
		TotalGrossInterest: decimal.NewDecimal64p2FromFloat64(totalGrossInterest),
		TotalNetInterest:   decimal.NewDecimal64p2FromFloat64(totalNetInterest)}
}

// ValidateInputData : validate the input data request to calculate Mortgage Monthly payments
func ValidateInputData(r model.MonthlyPaymentsRequest) (bool, string) {

	if r.AutomaticInterestUpdate && len(r.LoanToValueInterestTiers) == 0 {
		return false, "No loan to value interest tiers provided!"
	}

	if !r.AutomaticInterestUpdate && len(r.InterestTierUpdates) == 0 {
		return false, "No interest tiers month updates provided!"
	}

	if r.AutomaticInterestUpdate {

		sort.Slice(r.LoanToValueInterestTiers[:], func(i, j int) bool {
			return r.LoanToValueInterestTiers[i].Percentage < r.LoanToValueInterestTiers[j].Percentage
		})

		initialTierPercentage := r.LoanToValueInterestTiers[len(r.LoanToValueInterestTiers)-1].Percentage / 100
		initialRatio := r.InitialPrincipal / r.MarketValue

		if initialRatio > initialTierPercentage {
			return false, fmt.Sprintf("No interest tier found for initial percentage of %.2f %%", initialRatio*100)
		}

	} else {

		sort.Slice(r.InterestTierUpdates[:], func(i, j int) bool {
			return r.InterestTierUpdates[i].Month < r.InterestTierUpdates[j].Month
		})

		if r.InterestTierUpdates[0].Month != 1 {
			return false, "Interest Rate updates does not contain interest for 1st month!"
		}

	}

	if r.IncomeTax < 0 || r.IncomeTax > 100 {
		return false, "Income tax should be between 0% and 100%!"
	}

	return true, ""
}

func processExtraRepayments(rp []model.Repayment, m int, p *float64) {

	if rp == nil || len(rp) == 0 {
		return
	}

	for _, e := range rp {
		if e.Month == m {
			*p -= e.Amount
		}
	}

}
