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

func TestAccountHistory_BenchmarkBuyHoldTotalProfitStartDNE(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period, Prices: map[string]big.Decimal{}}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2}
	pricingSnap2 := PricingSnapshot{Period: period2, Prices: map[string]big.Decimal{"BENCH": big.NewDecimal(2.0)}}

	ah := NewAccountHistory()
	ah.Benchmark = "BENCH"

	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	pg := ah.BenchmarkBuyHoldTotalProfit()
	decimalEquals(t, 0.0, pg)
}

func TestAccountHistory_BenchmarkBuyHoldTotalProfitEndDNE(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period, Prices: map[string]big.Decimal{"BENCH": big.NewDecimal(2.0)}}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2}
	pricingSnap2 := PricingSnapshot{Period: period2, Prices: map[string]big.Decimal{}}

	ah := NewAccountHistory()
	ah.Benchmark = "BENCH"

	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	pg := ah.BenchmarkBuyHoldTotalProfit()
	decimalEquals(t, 0.0, pg)
}

func TestAccountHistory_BenchmarkBuyHoldTotalProfit(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period, Prices: map[string]big.Decimal{"BENCH": big.NewDecimal(1.0)}}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2}
	pricingSnap2 := PricingSnapshot{Period: period2, Prices: map[string]big.Decimal{"BENCH": big.NewDecimal(2.0)}}

	period3 := NewTimePeriod(time.Now().Add(time.Hour*48), time.Hour*24)
	acctSnap3 := AccountSnapshot{Period: period3}
	pricingSnap3 := PricingSnapshot{Period: period3, Prices: map[string]big.Decimal{"BENCH": big.NewDecimal(0.0)}}

	ah := NewAccountHistory()
	ah.Benchmark = "BENCH"

	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	pg := ah.BenchmarkBuyHoldTotalProfit()
	decimalEquals(t, 1.00, pg)

	ah.ApplySnapshot(&acctSnap3, &pricingSnap3)

	pg = ah.BenchmarkBuyHoldTotalProfit()
	decimalEquals(t, -1.0, pg)
}

func TestAccountHistory_BenchmarkBuyHoldPercentGainStartDNE(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period, Prices: map[string]big.Decimal{}}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2}
	pricingSnap2 := PricingSnapshot{Period: period2, Prices: map[string]big.Decimal{"BENCH": big.NewDecimal(2.0)}}

	ah := NewAccountHistory()
	ah.Benchmark = "BENCH"

	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	pg := ah.BenchmarkBuyHoldPercentGain()
	decimalEquals(t, 0.0, pg)
}

func TestAccountHistory_BenchmarkBuyHoldPercentGainEndDNE(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period, Prices: map[string]big.Decimal{"BENCH": big.NewDecimal(2.0)}}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2}
	pricingSnap2 := PricingSnapshot{Period: period2, Prices: map[string]big.Decimal{}}

	ah := NewAccountHistory()
	ah.Benchmark = "BENCH"

	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	pg := ah.BenchmarkBuyHoldPercentGain()
	decimalEquals(t, 0.0, pg)
}

func TestAccountHistory_BenchmarkBuyHoldPercentGain(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period}
	pricingSnap := PricingSnapshot{Period: period, Prices: map[string]big.Decimal{"BENCH": big.NewDecimal(1.0)}}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2}
	pricingSnap2 := PricingSnapshot{Period: period2, Prices: map[string]big.Decimal{"BENCH": big.NewDecimal(2.0)}}

	period3 := NewTimePeriod(time.Now().Add(time.Hour*48), time.Hour*24)
	acctSnap3 := AccountSnapshot{Period: period3}
	pricingSnap3 := PricingSnapshot{Period: period3, Prices: map[string]big.Decimal{"BENCH": big.NewDecimal(0.0)}}

	ah := NewAccountHistory()
	ah.Benchmark = "BENCH"

	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	pg := ah.BenchmarkBuyHoldPercentGain()
	decimalEquals(t, 100.00, pg)

	ah.ApplySnapshot(&acctSnap3, &pricingSnap3)

	pg = ah.BenchmarkBuyHoldPercentGain()
	decimalEquals(t, -100.0, pg)
}
