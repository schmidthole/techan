package techan

//lint:file-ignore S1038 prefer Fprintln

import (
	"github.com/sdcoffey/big"
)

// Analysis is an interface that describes a methodology for taking a AccountHistory as input,
// and giving back some float value that describes it's performance with respect to that methodology.
type Analysis interface {
	Analyze(*AccountHistory) big.Decimal
}

// TotalProfitAnalysis analyzes the trading record for total profit.
type TotalProfitAnalysis struct{}

// Analyze analyzes the trading record for total profit.
func (tps TotalProfitAnalysis) Analyze(hist *AccountHistory) big.Decimal {
	endEquity := hist.Snapshots[hist.LastIndex()].Equity
	startEquity := hist.Snapshots[0].Equity

	return endEquity.Sub(startEquity)
}

// PercentGainAnalysis analyzes the trading record for the percentage profit gained relative to start
type PercentGainAnalysis struct{}

// Analyze analyzes the trading record for the percentage profit gained relative to start
func (pga PercentGainAnalysis) Analyze(hist *AccountHistory) big.Decimal {
	endEquity := hist.Snapshots[hist.LastIndex()].Equity
	startEquity := hist.Snapshots[0].Equity

	if startEquity.EQ(big.ZERO) {
		return big.ZERO
	}

	return endEquity.Sub(startEquity).Div(startEquity)
}

// // Analyze returns the profit based on a simple buy and hold strategy
// func (baha BuyAndHoldAnalysis) Analyze(hist *AccountHistory) big.Decimal {

// }
