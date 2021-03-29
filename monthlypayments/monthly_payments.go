package monthlypayments

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/jinzhu/now"
	"github.com/lealoureiro/mortgage-calculator-api/model"
)

// CalculateLinearMonthlyPayments : calculate the monthly payments for a Linear Mortgage
func CalculateLinearMonthlyPayments(r model.MonthlyPaymentsRequest) model.MonthlyPayments {

	result := make([]model.MonthPayment, 0, r.Months)

	var interestSet InterestSet
	if r.AutomaticInterestUpdate {
		interestSet = LoanToValueInterestSet{r.MarketValue.AsFloat(), r.LoanToValueInterestTiers}
	} else {
		interestSet = InterestUpdatesSet{
			r.MarketValue.AsFloat(),
			r.InitialInterestRate.AsFloat(),
			r.InterestTierUpdates,
		}
	}

	monthlyRepayment := r.InitialPrincipal / float64(r.Months)
	principal := r.InitialPrincipal

	interestPercentage := 0.0
	interestGrossAmount := 0.0
	interestNetAmount := 0.0
	totalGrossInterest := 0.0
	totalNetInterest := 0.0

	incomeTax := float64(r.IncomeTax) / 100.0

	currentTime := r.StartDate.AsTime()
	endOfMonth := now.With(currentTime).EndOfMonth()

	remainingDaysInitialMonth := daysBetweenDates(currentTime, endOfMonth)

	initialInterest, marketValue := interestSet.GetInterest(currentTime, principal)
	initialInterestGross := ((principal * initialInterest) / float64(360)) * float64(remainingDaysInitialMonth+1)

	firstDayNextMonth := endOfMonth.Add(time.Nanosecond * time.Duration(1))
	endOfMonth = now.With(firstDayNextMonth).EndOfMonth()

	daysFirstRepaymentMonth := daysBetweenDates(firstDayNextMonth, endOfMonth)
	firstMonthInterestGross := ((principal*initialInterest)/float64(360))*float64(daysFirstRepaymentMonth) + initialInterestGross
	firstMonthInterestNet := firstMonthInterestGross - (firstMonthInterestGross * incomeTax)

	totalGrossInterest += firstMonthInterestGross
	totalNetInterest += firstMonthInterestNet

	paymentDate := firstDayNextMonth.AddDate(0, 1, 0)

	var payment = model.MonthPayment{}

	payment.Month = 1
	payment.PaymentDate = model.NewJSONTime(paymentDate)
	payment.Repayment = model.NewNumber(monthlyRepayment)
	payment.InterestGrossAmount = model.NewNumber(firstMonthInterestGross)
	payment.InterestNetAmount = model.NewNumber(firstMonthInterestNet)
	payment.Principal = model.NewNumber(principal)
	payment.InterestPercentage = model.NewNumber(initialInterest * 100)
	payment.TotalGross = model.NewNumber(monthlyRepayment + firstMonthInterestGross)
	payment.TotalNet = model.NewNumber(monthlyRepayment + firstMonthInterestNet)
	payment.LoanToValueRatio = model.NewNumber(principal / marketValue * 100)
	payment.MarketValue = model.NewNumber(marketValue)

	result = append(result, payment)

	principal -= monthlyRepayment
	currentTime = paymentDate

	repayments := r.Repayments

	monthAveragePrincipal := 0.0

	for i := 2; principal > 0; i++ {

		repayments, monthAveragePrincipal = processExtraRepayments(repayments, &principal, currentTime, now.With(currentTime).EndOfMonth())

		if principal < monthlyRepayment {
			monthlyRepayment = principal
		}

		interestPercentage, marketValue = interestSet.GetInterest(currentTime, principal)
		interestGrossAmount = monthAveragePrincipal * (interestPercentage / 12.0)
		interestNetAmount = interestGrossAmount - (interestGrossAmount * incomeTax)

		totalGrossInterest += interestGrossAmount
		totalNetInterest += interestNetAmount

		currentTime = currentTime.AddDate(0, 1, 0)

		var payment = model.MonthPayment{}

		payment.Month = i
		payment.PaymentDate = model.NewJSONTime(currentTime)
		payment.Repayment = model.NewNumber(monthlyRepayment)
		payment.InterestGrossAmount = model.NewNumber(interestGrossAmount)
		payment.InterestNetAmount = model.NewNumber(interestNetAmount)
		payment.Principal = model.NewNumber(principal)
		payment.InterestPercentage = model.NewNumber(interestPercentage * 100)
		payment.TotalGross = model.NewNumber(monthlyRepayment + interestGrossAmount)
		payment.TotalNet = model.NewNumber(monthlyRepayment + interestNetAmount)
		payment.LoanToValueRatio = model.NewNumber(principal / marketValue * 100)
		payment.MarketValue = model.NewNumber(marketValue)

		result = append(result, payment)

		principal -= monthlyRepayment

	}

	return model.MonthlyPayments{
		Payments:           result,
		TotalGrossInterest: model.NewNumber(totalGrossInterest),
		TotalNetInterest:   model.NewNumber(totalNetInterest)}
}

func processExtraRepayments(rp []model.Repayment, p *float64, s, e time.Time) ([]model.Repayment, float64) {

	if len(rp) == 0 {
		return rp, *p
	}

	days := 0
	totalPrincipal := 0.0

	for d := s; d.Before(e); d = d.AddDate(0, 0, 1) {

		rp = findRepaymentsForDate(rp, d, p)

		days++
		totalPrincipal += *p

	}

	monthAveragePrincipal := totalPrincipal / float64(days)

	return rp, monthAveragePrincipal

}

func findRepaymentsForDate(rp []model.Repayment, d time.Time, p *float64) []model.Repayment {

	if len(rp) == 0 {
		return rp
	}

	remaining := make([]model.Repayment, 0, len(rp))

	for _, v := range rp {

		if d.Equal(v.Date.AsTime()) {
			*p -= v.Amount
		} else {
			remaining = append(remaining, v)
		}

	}

	return remaining
}

// ValidateInputData : validate the input data request to calculate Mortgage Monthly payments
func ValidateInputData(r model.MonthlyPaymentsRequest) (bool, string) {

	if r.AutomaticInterestUpdate && len(r.LoanToValueInterestTiers) == 0 {
		return false, "No loan to value interest tiers provided!"
	}

	if !r.AutomaticInterestUpdate && len(r.InterestTierUpdates) == 0 {
		return false, "No interest tiers month updates provided!"
	}

	if r.MarketValue == nil {
		return false, "Missing initial Market Value!"
	}

	if r.StartDate == nil {
		return false, "Missing start date!"
	}

	if r.AutomaticInterestUpdate {

		sort.Slice(r.LoanToValueInterestTiers[:], func(i, j int) bool {
			return r.LoanToValueInterestTiers[i].Percentage < r.LoanToValueInterestTiers[j].Percentage
		})

		initialTierPercentage := r.LoanToValueInterestTiers[len(r.LoanToValueInterestTiers)-1].Percentage / 100
		initialRatio := r.InitialPrincipal / r.MarketValue.AsFloat()

		if initialRatio > initialTierPercentage {
			return false, fmt.Sprintf("No interest tier found for initial percentage of %.2f %%", initialRatio*100)
		}

	} else {

		if r.InitialInterestRate == nil {
			return false, "Missing initial interest rate!"
		}

		for _, u := range r.InterestTierUpdates {

			if u.UpdateDate.AsTime().Before(r.StartDate.AsTime()) {
				return false, fmt.Sprintf("Interest update date %d before mortgage start date!", u.UpdateDate)
			}

			if u.Interest == nil && u.MarketValue == nil {
				return false, fmt.Sprintf("Manually update for month %d should contain at least Market Value or Interest Rate!", u.UpdateDate)
			}

		}

	}

	if r.IncomeTax < 0 || r.IncomeTax > 100 {
		return false, "Income tax should be between 0% and 100%!"
	}

	return true, ""
}

func daysBetweenDates(t1, t2 time.Time) int32 {
	duration := t2.Sub(t1).Hours() / 24
	return int32(math.Round(duration))
}
