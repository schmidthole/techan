package techan

import (
	"testing"

	"github.com/sdcoffey/big"
)

func TestBinaryFrequencyIndicator(t *testing.T) {
	ind := NewFixedIndicator(
		0.0, 0.0, 2.0, 2.0, 0.0, 0.0, 0.0, 2.0, 0.0,
	)
	expected := []float64{
		0.0, 0.0, 1.0, 2.0, 2.0, 1.0, 0.0, 1.0, 1.0,
	}

	bin := NewBinaryFrequencyIndicator([]Indicator{ind}, 3, 1.0)

	for i, want := range expected {
		got := bin.Calculate(i)

		if !got.EQ(big.NewDecimal(want)) {
			t.Errorf("binary freq not matching for index %v, want: %v, got: %v", i, want, got.String())
		}
	}
}
