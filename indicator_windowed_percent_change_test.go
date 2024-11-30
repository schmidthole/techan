package techan

import (
	"testing"

	"github.com/schmidthole/big"
)

func TestWindowedPercentChange(t *testing.T) {
	ind := NewFixedIndicator(
		1.0, 1.5, 2.0, 1.5, 1.0, 1.5, 2.0, 1.5, 1.0,
	)
	expected := []float64{
		0.0, 50.0, 100.00, 0.0, -50.0, 0.0, 100.00, 0.0, -50.00,
	}

	pct := NewWindowedPercentChangeIndicator(ind, 3)

	for i, want := range expected {
		got := pct.Calculate(i)

		if !got.EQ(big.NewDecimal(want)) {
			t.Errorf("error in window pct change index %v. want: %v, got %v", i, want, got.String())
		}
	}
}
