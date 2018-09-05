package monthlypayments

import (
	"github.com/lealoureiro/mortgage-calculator-api/model"
	"github.com/strongo/decimal"
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

		interestPercentage = r.InterestTiers[0].Interest / 100
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
