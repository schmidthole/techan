package techan

import (
	"github.com/sdcoffey/big"
)

type windowedPercentChangeIndicator struct {
	indicator Indicator
	window    int
}

// The window percent change indicator calculates the change percentage over a fixed
// lookback window.
func NewWindowedPercentChangeIndicator(indicator Indicator, window int) Indicator {
	return windowedPercentChangeIndicator{
		indicator: indicator,
		window:    window,
	}
}

func (s windowedPercentChangeIndicator) Calculate(index int) big.Decimal {
	if index < s.window {
		return big.ZERO
	}

	end := s.indicator.Calculate(index)
	start := s.indicator.Calculate(index - s.window)

	if start.EQ(big.ZERO) {
		return big.ZERO
	}

	return end.Sub(start).Div(start)
}
