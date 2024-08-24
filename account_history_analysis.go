package techan

import (
	"math"

	"github.com/sdcoffey/big"
)

// Total profit of the account history.
func (ah *AccountHistory) TotalProfit() big.Decimal {
	endEquity := ah.Snapshots[ah.LastIndex()].Equity
	startEquity := ah.Snapshots[0].Equity

	return endEquity.Sub(startEquity)
}

// Total percent gain of the account history.
func (ah *AccountHistory) PercentGain() big.Decimal {
	if ah.Snapshots[0].Equity.EQ(big.ZERO) {
		return big.ZERO
	}

	return ah.TotalProfit().Div(ah.Snapshots[0].Equity).Mul(big.NewDecimal(100.00))
}

// Get the annualized return of the account equity
func (ah *AccountHistory) AnnualizedReturn() big.Decimal {
	startTimestamp := ah.Snapshots[0].Period.Start
	endTimestamp := ah.Snapshots[ah.LastIndex()].Period.Start

	days := big.NewDecimal(endTimestamp.Sub(startTimestamp).Hours()).Div(big.NewDecimal(24.0))
	years := days.Div(big.NewDecimal(365.00))

	startEquity := ah.Snapshots[0].Equity
	endEquity := ah.Snapshots[ah.LastIndex()].Equity

	if years.IsZero() || startEquity.IsZero() {
		return big.ZERO
	}

	base := endEquity.Sub(startEquity).Div(startEquity).Add(big.ONE).Float()
	exponent := big.ONE.Div(years).Float()

	return big.NewDecimal(math.Pow(base, exponent)).Sub(big.ONE).Mul(big.NewDecimal(100.00))
}
