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
	interestAmountGross := 0.0
	interestAmountNet := 0.0

	totalGrossInterest := 0.0
	totalNetInterest := 0.0

	incomeTax := float64(r.IncomeTax) / 100.0

	for i := 0; i < r.Months && principal > 0; i++ {

		if principal < monthlyRepayment {
			monthlyRepayment = principal
		}

		interestPercentage = r.InterestTiers[0].Interest / 100
		interestAmountGross = (principal * interestPercentage) / 12.0
		interestAmountNet = interestAmountGross - (interestAmountGross * incomeTax)

		totalGrossInterest += interestAmountGross
		totalNetInterest += interestAmountNet

		var payment = model.MonthPayment{}

		payment.Repayment = decimal.NewDecimal64p2FromFloat64(monthlyRepayment)
		payment.InterestAmountGross = decimal.NewDecimal64p2FromFloat64(interestAmountGross)
		payment.Principal = decimal.NewDecimal64p2FromFloat64(principal)
		payment.InterestPercentage = interestPercentage
		payment.TotalGross = decimal.NewDecimal64p2FromFloat64(monthlyRepayment + interestAmountGross)
		payment.TotalNet = decimal.NewDecimal64p2FromFloat64(monthlyRepayment + interestAmountNet)

		result = append(result, payment)

		principal -= monthlyRepayment

	}

	return model.MonthlyPayments{
		Payments:           result,
		TotalGrossInterest: decimal.NewDecimal64p2FromFloat64(totalGrossInterest),
		TotalNetInterest:   decimal.NewDecimal64p2FromFloat64(totalNetInterest)}
}
