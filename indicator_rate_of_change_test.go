package techan

import (
	"testing"

	"github.com/sdcoffey/big"
)

func Test_RateOfChangeIndicator(t *testing.T) {
	ind := NewFixedIndicator(
		0.0, 0.5, 1.0, 0.5, 0.0, 0.5, 1.0, 0.5, 0.0,
	)
	expected := []float64{
		0.0, 0.0, 0.33, 0.0, -0.33, 0.0, 0.33, 0.00, -0.33,
	}

	roc := NewRateOfChangeIndicator(ind, 3)

	for i, want := range expected {
		got := roc.Calculate(i)

		top := got.Add(big.NewDecimal(0.1))
		bottom := got.Sub(big.NewDecimal(0.1))

		wantDec := big.NewDecimal(want)

		if wantDec.GT(top) || wantDec.LT(bottom) {
			t.Errorf("roc indicator failed at index %v. want: %v, got: %v.", i, want, got.String())
		}
	}
}
