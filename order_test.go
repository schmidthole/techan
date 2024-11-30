package techan

import (
	"testing"
	"time"

	"github.com/schmidthole/big"
)

func TestOrderCostBasis(t *testing.T) {
	order := Order{
		Side:          BUY,
		Security:      "FAKE",
		Price:         big.NewDecimal(10.00),
		Amount:        big.NewDecimal(10.00),
		ExecutionTime: time.Now(),
	}

	costBasis := order.CostBasis()

	decimalEquals(t, 100.0, costBasis)
}
