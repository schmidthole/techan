package techan

import (
	"testing"
	"time"

	"github.com/schmidthole/big"
	"github.com/stretchr/testify/assert"
)

func TestAccountHistory_ApplySnapshot(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period}
	pricingSnap := PricingSnapshot{Period: period}

	hist := NewAccountHistory()

	err := hist.ApplySnapshot(&acctSnap, &pricingSnap)
	assert.Nil(t, err)
}

func TestAccountHistory_ApplyBadSnapshot(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	period2 := NewTimePeriod(time.Now().AddDate(1, 0, 0), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period}
	pricingSnap := PricingSnapshot{Period: period2}

	hist := NewAccountHistory()

	err := hist.ApplySnapshot(&acctSnap, &pricingSnap)
	assert.NotNil(t, err)
}

func TestAccountHistory_LastIndex(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period}
	pricingSnap := PricingSnapshot{Period: period}

	hist := NewAccountHistory()

	hist.ApplySnapshot(&acctSnap, &pricingSnap)

	assert.Equal(t, 0, hist.LastIndex())
}

func TestAccountHistory_PriceAtIndex(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	period2 := NewTimePeriod(time.Now().AddDate(1, 0, 0), time.Hour*24)

	prices1 := map[string]big.Decimal{
		MOCK_SECURITY: big.NewDecimal(1.0),
	}
	prices2 := map[string]big.Decimal{
		MOCK_SECURITY: big.NewDecimal(2.0),
	}

	hist := NewAccountHistory()

	acctSnap := AccountSnapshot{Period: period}
	pricingSnap := PricingSnapshot{Period: period, Prices: prices1}
	hist.ApplySnapshot(&acctSnap, &pricingSnap)

	acctSnap2 := AccountSnapshot{Period: period2}
	pricingSnap2 := PricingSnapshot{Period: period2, Prices: prices2}
	hist.ApplySnapshot(&acctSnap2, &pricingSnap2)

	p1, exists := hist.PriceAtIndex(MOCK_SECURITY, 0)
	assert.True(t, exists)
	decimalEquals(t, 1.0, p1)

	p2, exists := hist.PriceAtIndex(MOCK_SECURITY, 1)
	assert.True(t, exists)
	decimalEquals(t, 2.0, p2)

	_, exists = hist.PriceAtIndex("NOT_THERE", 1)
	assert.False(t, exists)
}

func TestAccountHistory_AccountEquityAsIndicator(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	period2 := NewTimePeriod(time.Now().AddDate(1, 0, 0), time.Hour*24)

	prices1 := map[string]big.Decimal{
		MOCK_SECURITY: big.NewDecimal(1.0),
	}
	prices2 := map[string]big.Decimal{
		MOCK_SECURITY: big.NewDecimal(2.0),
	}

	hist := NewAccountHistory()

	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period, Prices: prices1}
	hist.ApplySnapshot(&acctSnap, &pricingSnap)

	acctSnap2 := AccountSnapshot{Period: period2, Equity: big.NewDecimal(2.0)}
	pricingSnap2 := PricingSnapshot{Period: period2, Prices: prices2}
	hist.ApplySnapshot(&acctSnap2, &pricingSnap2)

	ind := hist.AccountEquityAsIndicator()

	decimalEquals(t, 1.0, ind.Calculate(0))
	decimalEquals(t, 2.0, ind.Calculate(1))
}
