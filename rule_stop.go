package talib4g

type StopLossRule struct {
	Indicator
	tolerance Decimal
}

// Returns a new stop loss rule based on a timeseries and a loss tolerance
// The loss tolerance should be a number between -1 and 1, where negative
// values represent a loss and vice versa.
func NewStopLossRule(series *TimeSeries, lossTolerance float64) Rule {
	return StopLossRule{
		Indicator: NewClosePriceIndicator(series),
		tolerance: NewDecimal(lossTolerance),
	}
}

func (slr StopLossRule) IsSatisfied(index int, record *TradingRecord) bool {
	if !record.CurrentTrade().IsOpen() {
		return false
	}

	openPrice := record.CurrentTrade().CostBasis()
	loss := slr.Indicator.Calculate(index).Div(openPrice).Sub(ONE)
	return loss.LTE(slr.tolerance)
}