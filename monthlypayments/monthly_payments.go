package monthlypayments

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/jinzhu/now"
	"github.com/lealoureiro/mortgage-calculator-api/model"
	"github.com/strongo/decimal"
)

// CalculateLinearMonthlyPayments : calculate the monthly payments for a Linear Mortgage
func CalculateLinearMonthlyPayments(r model.MonthlyPaymentsRequest) model.MonthlyPayments {

	result := make([]model.MonthPayment, 0, r.Months)

	var interestSet InterestSet
	if r.AutomaticInterestUpdate {
		interestSet = LoanToValueInterestSet{r.MarketValue.AsFloat64(), r.LoanToValueInterestTiers}
	} else {
		interestSet = InterestUpdatesSet{
			r.MarketValue.AsFloat64(),
			r.InitialInterestRate.AsFloat64(),
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

	initialInterest, marketValue := interestSet.GetInterest(1, principal)
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
	payment.Repayment = decimal.NewDecimal64p2FromFloat64(monthlyRepayment)
	payment.InterestGrossAmount = decimal.NewDecimal64p2FromFloat64(firstMonthInterestGross)
	payment.InterestNetAmount = decimal.NewDecimal64p2FromFloat64(firstMonthInterestNet)
	payment.Principal = decimal.NewDecimal64p2FromFloat64(principal)
	payment.InterestPercentage = decimal.NewDecimal64p2FromFloat64(initialInterest * 100)
	payment.TotalGross = decimal.NewDecimal64p2FromFloat64(monthlyRepayment + firstMonthInterestGross)
	payment.TotalNet = decimal.NewDecimal64p2FromFloat64(monthlyRepayment + firstMonthInterestNet)
	payment.LoanToValueRatio = decimal.NewDecimal64p2FromFloat64(principal / marketValue * 100)
	payment.MarketValue = decimal.NewDecimal64p2FromFloat64(marketValue)

	result = append(result, payment)

	principal -= monthlyRepayment
	currentTime = paymentDate

	for i := 2; principal > 0; i++ {

		currentTime = currentTime.AddDate(0, 1, 0)

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
		payment.PaymentDate = model.NewJSONTime(currentTime)
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
		initialRatio := r.InitialPrincipal / r.MarketValue.AsFloat64()

		if initialRatio > initialTierPercentage {
			return false, fmt.Sprintf("No interest tier found for initial percentage of %.2f %%", initialRatio*100)
		}

	} else {

		if r.InitialInterestRate == nil {
			return false, "Missing initial interest rate!"
		}

		for _, u := range r.InterestTierUpdates {

			if u.Month < 1 || u.Month > r.Months {
				return false, fmt.Sprintf("Interest update month %d outside of range!", u.Month)
			}

			if u.Interest == nil && u.MarketValue == nil {
				return false, fmt.Sprintf("Manually update for month %d should contain at least Market Value or Interest Rate!", u.Month)
			}

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

func daysBetweenDates(t1, t2 time.Time) int32 {
	duration := t2.Sub(t1).Hours() / 24
	return int32(math.Round(duration))
}
