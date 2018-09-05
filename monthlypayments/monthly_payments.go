package monthlypayments

import (
	"fmt"
	"github.com/lealoureiro/mortgage-calculator-api/model"
	"github.com/strongo/decimal"
	"sort"
)

func CalculateLinearMonthlyPayments(r model.MonthlyPaymentRequest) model.MonthlyPayments {

	result := make([]model.MonthPayment, 0, r.Months)

	monthlyRepayment := r.InitialPrincipal / float64(r.Months)
	principal := r.InitialPrincipal

	interestPercentage := 0.0
	interestGrossAmount := 0.0
	interestNetAmount := 0.0

	totalGrossInterest := 0.0
	totalNetInterest := 0.0

	incomeTax := float64(r.IncomeTax) / 100.0

	for i := 0; i < r.Months && principal > 0; i++ {

		if principal < monthlyRepayment {
			monthlyRepayment = principal
		}

		interestPercentage = getInterestTierPercentage(r.MarketValue, principal, r.InterestTiers) / 100
		interestGrossAmount = (principal * interestPercentage) / 12.0
		interestNetAmount = interestGrossAmount - (interestGrossAmount * incomeTax)

		totalGrossInterest += interestGrossAmount
		totalNetInterest += interestNetAmount

		var payment = model.MonthPayment{}

		payment.Repayment = decimal.NewDecimal64p2FromFloat64(monthlyRepayment)
		payment.InterestGrossAmount = decimal.NewDecimal64p2FromFloat64(interestGrossAmount)
		payment.InterestNetAmount = decimal.NewDecimal64p2FromFloat64(interestNetAmount)
		payment.Principal = decimal.NewDecimal64p2FromFloat64(principal)
		payment.InterestPercentage = interestPercentage * 100
		payment.TotalGross = decimal.NewDecimal64p2FromFloat64(monthlyRepayment + interestGrossAmount)
		payment.TotalNet = decimal.NewDecimal64p2FromFloat64(monthlyRepayment + interestNetAmount)

		result = append(result, payment)

		principal -= monthlyRepayment

	}

	return model.MonthlyPayments{
		Payments:           result,
		TotalGrossInterest: decimal.NewDecimal64p2FromFloat64(totalGrossInterest),
		TotalNetInterest:   decimal.NewDecimal64p2FromFloat64(totalNetInterest)}
}

func ValidateInputData(r model.MonthlyPaymentRequest) (bool, string) {

	sort.Slice(r.InterestTiers[:], func(i, j int) bool {
		return r.InterestTiers[i].Percentage < r.InterestTiers[j].Percentage
	})

	initialTierPercentage := r.InterestTiers[len(r.InterestTiers)-1].Percentage / 100
	initialRatio := r.InitialPrincipal / r.MarketValue

	if initialRatio > initialTierPercentage {
		return false, fmt.Sprintf("No interest tier found for initial percentage of %.2f %", initialRatio*100)
	}

	return true, ""
}

func getInterestTierPercentage(m float64, p float64, l []model.InterestTier) float64 {

	ratio := p / m * 100

	for _, e := range l {
		if ratio <= e.Percentage {
			return e.Interest
		}
	}

	return 0.0
}
