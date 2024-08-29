package techan

import (
	"testing"

	"github.com/sdcoffey/big"
)

func TestLocalExtremaIndicator(t *testing.T) {
	ind := NewFixedIndicator(
		1.0, 2.0, 3.0, 2.0, 1.0, 0.5, 1.0, 1.5, 2.0, 3.0, 4.0, 2.0, 1.0, 0.5,
	)
	expected := []float64{
		-1.0, 0.0, 1.0, 0.0, 0.0, -1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0,
	}

	ext := NewLocalExtremaIndicator(ind, 3, len(expected))

	for i, want := range expected {
		got := ext.Calculate(i)

		if !got.EQ(big.NewDecimal(want)) {
			t.Errorf("local extrema fail at index %v. want: %v, got: %v", i, want, got.String())
		}
	}
}
