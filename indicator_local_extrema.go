package techan

import "github.com/sdcoffey/big"

type localExtremaIndicator struct {
	extrema []int
}

// Defines an indicator to determine local extrema (minima/maxima). This indicator is really
// useful for post analysis and calculating things like max drawdown. It is *NOT* good for
// live trading as it uses a lookback to find the local extrema when an roc switches direction.
//
// This indicator can be useful for identification of "boxes" or "bases" in an asset/security
// where it is not critical that the extrema be identified immediately.
func NewLocalExtremaIndicator(timeseries *TimeSeries, window int) Indicator {
	closes := NewClosePriceIndicator(timeseries)
	roc := NewRateOfChangeIndicator(closes, window)

	extrema := make([]int, timeseries.LastIndex()+1)
	for i := range extrema {
		extrema[i] = 0
	}

	for i := 0; i <= timeseries.LastIndex(); i++ {
		roc0 := roc.Calculate(i)
		roc1 := roc.Calculate(i - 1)

		startIndex := i - window
		if roc0.GT(big.ZERO) && roc1.LTE(big.ZERO) {
			minimaIndex := findMinIndexInWindow(closes, startIndex, i)
			extrema[minimaIndex] = -1
		} else if roc0.LTE(big.ZERO) && roc1.GT(big.ZERO) {
			maximaIndex := findMaxIndexInWindow(closes, startIndex, i)
			extrema[maximaIndex] = 1
		}
	}

	return localExtremaIndicator{extrema: extrema}
}

func (ext localExtremaIndicator) Calculate(index int) big.Decimal {
	return big.NewFromInt(ext.extrema[index])
}

func findMaxIndexInWindow(ind Indicator, startIndex int, endIndex int) int {
	maxIndex := startIndex
	for i := startIndex; i <= endIndex; i++ {
		if ind.Calculate(i).GTE(ind.Calculate(maxIndex)) {
			maxIndex = i
		}
	}

	return maxIndex
}

func findMinIndexInWindow(ind Indicator, startIndex int, endIndex int) int {
	minIndex := startIndex
	for i := startIndex; i <= endIndex; i++ {
		if ind.Calculate(i).LTE(ind.Calculate(minIndex)) {
			minIndex = i
		}
	}

	return minIndex
}
