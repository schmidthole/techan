package techan

import (
	"testing"

	"github.com/sdcoffey/big"
)

func Test_CashToShare(t *testing.T) {
	zeroShares := CashToShares(big.ZERO, big.NewDecimal(1.0))
	decimalEquals(t, 0.0, zeroShares)

	moreZeroShares := CashToShares(big.NewDecimal(100.00), big.NewDecimal(-1.0))
	decimalEquals(t, 0.0, moreZeroShares)

	tenShares := CashToShares(big.NewDecimal(10.00), big.NewDecimal(100.00))
	decimalEquals(t, 10.0, tenShares)

	stillTenShares := CashToShares(big.NewDecimal(10.00), big.NewDecimal(109.99))
	decimalEquals(t, 10.0, stillTenShares)

	nineShares := CashToShares(big.NewDecimal(10.01), big.NewDecimal(100.00))
	decimalEquals(t, 9.0, nineShares)
}
