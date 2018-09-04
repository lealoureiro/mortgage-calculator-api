package monthlypayments

import "github.com/lealoureiro/mortgage-calculator-api/model"

func CalculateLinearMonthlyPayments(r model.MonthlyPaymentRequest) model.MonthlyPayments {

	result := make([]model.MonthPayment, 0, r.Months)

	monthlyRepayment := r.InitialPrincipal / float64(r.Months)
	remainingAmount := r.InitialPrincipal

	totalGrossInterest := 0.0
	totalNetInterest := 0.0

	incomeTax := float64(r.IncomeTax) / 100.0

	for i := 0; i < r.Months && remainingAmount > 0; i++ {

		if remainingAmount < monthlyRepayment {
			monthlyRepayment = remainingAmount
		}

		var payment = model.MonthPayment{}

		interestPercentage := r.InterestTiers[0].Interest / 100

		payment.GrossAmount = monthlyRepayment
		payment.InterestPercentage = interestPercentage
		payment.InterestAmountGross = (remainingAmount * interestPercentage) / 12.0

		interestAmountNet := payment.InterestAmountGross - (payment.InterestAmountGross * incomeTax)

		payment.TotalGross = payment.InterestAmountGross + monthlyRepayment
		payment.TotalNet = monthlyRepayment + interestAmountNet

		totalGrossInterest += payment.InterestAmountGross
		totalNetInterest += interestAmountNet

		payment.RemainingAmount = remainingAmount

		result = append(result, payment)

		remainingAmount -= monthlyRepayment

	}

	return model.MonthlyPayments{Payments: result, TotalGrossInterest: totalGrossInterest, TotalNetInterest: totalNetInterest}
}
