package techan

import "testing"

func TestNewSupertrendIndicator(t *testing.T) {
	ts := mockTimeSeriesOCHL(
		[]float64{1.0, 2.0, 2.5, 0.5},
		[]float64{2.0, 3.0, 3.5, 1.5},
		[]float64{3.0, 4.0, 4.5, 2.5},
		[]float64{4.0, 5.0, 5.5, 3.5},
	)

	indicator := NewSupertrendIndicator(ts, 3, 2)
	expectedValues := []float64{1.5, 2.5, 3.5, 8.5}

	indicatorEquals(t, expectedValues, indicator)
}
