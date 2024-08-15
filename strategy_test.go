package techan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrategyLastIndex(t *testing.T) {
	strat := Strategy{
		Timeseries: *mockedTimeSeries,
	}

	assert.Equal(t, 11, strat.LastIndex())
}
