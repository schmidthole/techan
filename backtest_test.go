package techan

import (
	"testing"

	"github.com/schmidthole/big"
	"github.com/stretchr/testify/assert"
)

func Test_NewBacktest(t *testing.T) {
	ts := mockTimeSeriesFl(1.0, 2.0, 3.0, 4.0, 5.0)
	strat := Strategy{
		Security:   "ONE",
		Timeseries: *ts,
		Rule:       truthRule{},
		Indicators: map[string]Indicator{},
	}
	alloc := NewNaiveAllocator(big.NewDecimal(1.0), big.NewDecimal(1.0))
	acct := NewAccount()

	bt := NewBacktest([]Strategy{strat}, alloc, acct)

	assert.Equal(t, 0, bt.tick)
	assert.Equal(t, 1, len(bt.history.Snapshots))
}

func Test_BacktestRun(t *testing.T) {
	ts := mockTimeSeriesFl(1.0, 2.0, 3.0, 4.0, 5.0)
	strat := Strategy{
		Security:   "ONE",
		Timeseries: *ts,
		Rule:       truthRule{},
		Indicators: map[string]Indicator{},
	}
	alloc := NewNaiveAllocator(big.NewDecimal(1.0), big.NewDecimal(1.0))
	acct := NewAccount()

	bt := NewBacktest([]Strategy{strat}, alloc, acct)

	hist, err := bt.Run()

	assert.Nil(t, err)
	assert.NotNil(t, hist)
	assert.Equal(t, 4, bt.tick)
	assert.Equal(t, 6, len(hist.Snapshots))
}
