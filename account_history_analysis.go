package techan

import (
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

	return ah.TotalProfit().Div(ah.Snapshots[0].Equity)
}

// Get the total profit for a buy and hold of the benchmark security.
func (ah *AccountHistory) BenchmarkBuyHoldTotalProfit() big.Decimal {
	benchStartPrice, exists := ah.PriceAtIndex(ah.Benchmark, 0)
	if !exists || benchStartPrice.EQ(big.ZERO) {
		return big.ZERO
	}
	benchEndPrice, exists := ah.PriceAtIndex(ah.Benchmark, ah.LastIndex())
	if !exists || benchEndPrice.EQ(big.ZERO) {
		return big.ZERO
	}

	benchChange := benchEndPrice.Sub(benchStartPrice)

	return ah.Snapshots[0].Equity.Mul(benchChange)
}

func (ah *AccountHistory) BenchmarkBuyHoldPercentGain() big.Decimal {
	benchStartPrice, exists := ah.PriceAtIndex(ah.Benchmark, 0)
	if !exists || benchStartPrice.EQ(big.ZERO) {
		return big.ZERO
	}
	benchEndPrice, exists := ah.PriceAtIndex(ah.Benchmark, ah.LastIndex())
	if !exists || benchEndPrice.EQ(big.ZERO) {
		return big.ZERO
	}

	return benchEndPrice.Sub(benchStartPrice).Div(benchStartPrice)
}
