package techan

import "github.com/sdcoffey/big"

type maximumDrawdownIndicator struct {
	drawdowns []big.Decimal
}

// Revamped maximum drawdown indicator which uses a local extrema calculation
// to find the difference between almost absolute peaks and valleys over time.
//
// Each "peak" and "valley" are paired to introduce a single "drawdown"
func NewMaximumDrawdownIndicator(timeseries *TimeSeries, window int) Indicator {
	closes := NewClosePriceIndicator(timeseries)
	extrema := NewLocalExtremaIndicator(timeseries, window)

	drawdowns := []big.Decimal{}
	maxDrawdown := big.ZERO
	lastPeak := big.ZERO
	inDrawdown := false

	for i := 0; i < timeseries.LastIndex(); i++ {
		isExtrema := extrema.Calculate(i)

		if isExtrema.GT(big.ZERO) {
			lastPeak = closes.Calculate(i)
			inDrawdown = true
		} else if isExtrema.LT(big.ZERO) && inDrawdown {
			drawdown := lastPeak.Sub(closes.Calculate(i))

			if drawdown.GT(maxDrawdown) {
				maxDrawdown = drawdown
			}

			inDrawdown = false
		}

		drawdowns = append(drawdowns, maxDrawdown)
	}

	return maximumDrawdownIndicator{drawdowns: drawdowns}
}

func (mdi maximumDrawdownIndicator) Calculate(index int) big.Decimal {
	return mdi.drawdowns[index]
}
