package techan

import (
	"testing"
	"time"

	"github.com/sdcoffey/big"
)

func TestAccountHistory_TotalProfit(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2, Equity: big.NewDecimal(2.0)}
	pricingSnap2 := PricingSnapshot{Period: period2}

	period3 := NewTimePeriod(time.Now().Add(time.Hour*48), time.Hour*24)
	acctSnap3 := AccountSnapshot{Period: period3, Equity: big.NewDecimal(0.0)}
	pricingSnap3 := PricingSnapshot{Period: period3}

	ah := NewAccountHistory()
	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	profit := ah.TotalProfit()
	decimalEquals(t, 1.0, profit)

	ah.ApplySnapshot(&acctSnap3, &pricingSnap3)

	profit = ah.TotalProfit()
	decimalEquals(t, -1.0, profit)
}

func TestAccountHistory_PercentGainZeroStart(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(0.0)}
	pricingSnap := PricingSnapshot{Period: period}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2, Equity: big.NewDecimal(2.0)}
	pricingSnap2 := PricingSnapshot{Period: period2}

	ah := NewAccountHistory()
	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	pg := ah.PercentGain()
	decimalEquals(t, 0.0, pg)
}

func TestAccountHistory_PercentGain(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2, Equity: big.NewDecimal(2.0)}
	pricingSnap2 := PricingSnapshot{Period: period2}

	period3 := NewTimePeriod(time.Now().Add(time.Hour*48), time.Hour*24)
	acctSnap3 := AccountSnapshot{Period: period3, Equity: big.NewDecimal(0.0)}
	pricingSnap3 := PricingSnapshot{Period: period3}

	ah := NewAccountHistory()
	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	pg := ah.PercentGain()
	decimalEquals(t, 100.00, pg)

	ah.ApplySnapshot(&acctSnap3, &pricingSnap3)

	pg = ah.PercentGain()
	decimalEquals(t, -100.0, pg)
}

func TestAccountHistory_AnnualizedReturn(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24*365*2), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2, Equity: big.NewDecimal(2.0)}
	pricingSnap2 := PricingSnapshot{Period: period2}

	ah := NewAccountHistory()
	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	arFl := ah.AnnualizedReturn().Float()
	if (arFl > 60.0) || (arFl < 40.0) {
		t.Errorf("annualized return out of bounds. expected %v, got %v", 50.00, arFl)
	}
}
