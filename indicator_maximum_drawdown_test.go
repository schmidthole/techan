package techan

import (
	"testing"

	"github.com/schmidthole/big"
)

func Test_MaximumDrawdownIndicator(t *testing.T) {
	ind := NewFixedIndicator(
		1.0, 2.0, 3.0, 2.0, 1.0, 0.5, 1.0, 1.5, 2.0, 3.0, 4.0, 2.0, 1.0, 0.5,
	)

	expected := []float64{
		0.0, 0.0, 0.0, 33.33, 66.66, 83.33, 83.33, 83.33, 83.33, 83.33, 83.33, 83.33, 83.33, 87.5,
	}

	md := NewMaximumDrawdownIndicator(ind, 3, 14)

	for i, want := range expected {
		got := md.Calculate(i)

		upper := big.NewDecimal(want + 1.0)
		lower := big.NewDecimal(want - 1.0)

		if got.GT(upper) || got.LT(lower) {
			t.Errorf("max drawdown fail at index %v. want: %v, got: %v", i, want, got.String())
		}
	}
}
