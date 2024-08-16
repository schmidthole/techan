package techan

import (
	"math"

	"github.com/sdcoffey/big"
)

func CashToShares(assetPrice big.Decimal, cash big.Decimal) big.Decimal {
	if assetPrice.IsZero() || cash.LTE(big.ZERO) {
		return big.ZERO
	}

	return big.NewDecimal(math.Floor(cash.Div(assetPrice).Float()))
}

func SharesToCash(assetPrice big.Decimal, shares big.Decimal) big.Decimal {
	return assetPrice.Mul(shares)
}
