package techan

import (
	"os"
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

func TestAccountHistory_AnnualizedVolatility(t *testing.T) {
	tests := []struct {
		name     string
		equities []float64
		expected float64
	}{
		{
			name:     "basic",
			equities: []float64{100.0, 105.0, 102.0, 107.0},
			expected: 11.36,
		},
		{
			name:     "no variance",
			equities: []float64{100.0, 100.0, 100.0, 100.0},
			expected: 0.0,
		},
		{
			name:     "single",
			equities: []float64{100.0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ah := NewAccountHistory()
			now := time.Now()

			for i, e := range tt.equities {
				period := NewTimePeriod(now.Add(time.Hour*24*time.Duration(i)), time.Hour*24)
				snap := AccountSnapshot{Period: period, Equity: big.NewDecimal(e)}
				pricing := PricingSnapshot{Period: period}

				ah.ApplySnapshot(&snap, &pricing)
			}

			aVol := ah.AnnualizedVolatility()

			if tt.expected == 0.0 {
				if !aVol.Zero() {
					t.Errorf("result should be 0.0, but got %v", aVol.String())
				}
			} else {
				decimalAlmostEquals(t, big.NewDecimal(tt.expected), aVol, 0.1)
			}
		})
	}
}

func TestAccountHistory_MonthlyPercentGains(t *testing.T) {
	tests := []struct {
		name     string
		periods  []time.Time
		equities []float64
		expected []ReturnPeriod
	}{
		{
			name:     "zero days",
			periods:  []time.Time{},
			equities: []float64{},
			expected: []ReturnPeriod{},
		},
		{
			name:     "one day",
			periods:  []time.Time{time.Now()},
			equities: []float64{1.0},
			expected: []ReturnPeriod{},
		},
		{
			name: "two days same month",
			periods: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
			equities: []float64{1.0, 2.0},
			expected: []ReturnPeriod{
				{
					Period:      TimePeriod{},
					PercentGain: big.NewDecimal(100.00),
					TotalProfit: big.NewDecimal(1.0),
				},
			},
		},
		{
			name: "two periods",
			periods: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.February, 28, 0, 0, 0, 0, time.UTC),
			},
			equities: []float64{1.0, 2.0, 2.0, 4.0},
			expected: []ReturnPeriod{
				{
					Period:      TimePeriod{},
					PercentGain: big.NewDecimal(100.00),
					TotalProfit: big.NewDecimal(1.0),
				},
				{
					Period:      TimePeriod{},
					PercentGain: big.NewDecimal(100.00),
					TotalProfit: big.NewDecimal(2.0),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ah := NewAccountHistory()

			for i, e := range tt.equities {
				period := NewTimePeriod(tt.periods[i], time.Hour*24)
				snap := AccountSnapshot{Period: period, Equity: big.NewDecimal(e)}
				pricing := PricingSnapshot{Period: period}

				ah.ApplySnapshot(&snap, &pricing)
			}

			monthlyReturns := ah.MonthlyPercentGains()

			if len(monthlyReturns) != len(tt.expected) {
				t.Errorf("expected %v monthly returns, got %v", len(tt.expected), len(monthlyReturns))
			}

			for i, r := range monthlyReturns {
				decimalAlmostEquals(t, tt.expected[i].TotalProfit, r.TotalProfit, 1.0)
				decimalAlmostEquals(t, tt.expected[i].PercentGain, r.PercentGain, 1.0)
			}
		})
	}
}

func TestAccountHistory_ExportSnapshotsYaml(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2, Equity: big.NewDecimal(2.0)}
	pricingSnap2 := PricingSnapshot{Period: period2}

	ah := NewAccountHistory()
	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	filepath := "snapshots.yaml"
	err := ah.ExportSnapshotsYaml(filepath)
	if err != nil {
		t.Errorf("failed to export snapshots to yaml: %v", err)
	}

	os.Remove(filepath)
}

func TestAccountHistory_ExportAnalysisSummaryYaml(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24*365*2), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2, Equity: big.NewDecimal(2.0)}
	pricingSnap2 := PricingSnapshot{Period: period2}

	ah := NewAccountHistory()
	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	filepath := "analysis_summary.yaml"
	err := ah.ExportAnalysisSummaryYaml(filepath)
	if err != nil {
		t.Errorf("failed to export analysis summary to yaml: %v", err)
	}

	os.Remove(filepath)
}

func TestAccountHistory_ExportMonthlyGainsYaml(t *testing.T) {
	period := NewTimePeriod(time.Now(), time.Hour*24)
	acctSnap := AccountSnapshot{Period: period, Equity: big.NewDecimal(1.0)}
	pricingSnap := PricingSnapshot{Period: period}

	period2 := NewTimePeriod(time.Now().Add(time.Hour*24*30), time.Hour*24)
	acctSnap2 := AccountSnapshot{Period: period2, Equity: big.NewDecimal(2.0)}
	pricingSnap2 := PricingSnapshot{Period: period2}

	ah := NewAccountHistory()
	ah.ApplySnapshot(&acctSnap, &pricingSnap)
	ah.ApplySnapshot(&acctSnap2, &pricingSnap2)

	filepath := "monthly_gains.yaml"
	err := ah.ExportMonthlyGainsYaml(filepath)
	if err != nil {
		t.Errorf("failed to export monthly gains to yaml: %v", err)
	}

	os.Remove(filepath)
}
