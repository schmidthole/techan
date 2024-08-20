package techan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type truthRule struct{}

func (tr truthRule) IsSatisfied(index int) bool {
	return true
}

type falseRule struct{}

func (fr falseRule) IsSatisfied(index int) bool {
	return false
}

type mockRule struct{ values []bool }

func (mr mockRule) IsSatisfied(index int) bool {
	return mr.values[index]
}

func TestAndRule(t *testing.T) {
	t.Run("both truthy", func(t *testing.T) {
		rule := And(truthRule{}, truthRule{})

		assert.True(t, rule.IsSatisfied(0))
	})

	t.Run("both falsey", func(t *testing.T) {
		rule := And(falseRule{}, falseRule{})

		assert.False(t, rule.IsSatisfied(0))
	})

	t.Run("one of each", func(t *testing.T) {
		rule := And(truthRule{}, falseRule{})

		assert.False(t, rule.IsSatisfied(0))
	})
}

func TestOrRule(t *testing.T) {
	t.Run("both truthy", func(t *testing.T) {
		rule := Or(truthRule{}, truthRule{})

		assert.True(t, rule.IsSatisfied(0))
	})

	t.Run("both falsey", func(t *testing.T) {
		rule := Or(falseRule{}, falseRule{})

		assert.False(t, rule.IsSatisfied(0))
	})

	t.Run("one of each", func(t *testing.T) {
		rule := Or(truthRule{}, falseRule{})

		assert.True(t, rule.IsSatisfied(0))
	})
}

func TestOverIndicatorRule(t *testing.T) {
	highIndicator := NewConstantIndicator(1)
	lowIndicator := NewConstantIndicator(0)

	t.Run("returns true when first indicator is over second indicator", func(t *testing.T) {
		rule := OverIndicatorRule{
			First:  highIndicator,
			Second: lowIndicator,
		}

		assert.True(t, rule.IsSatisfied(0))
	})

	t.Run("returns false when first indicator is under second indicator", func(t *testing.T) {
		rule := OverIndicatorRule{
			First:  lowIndicator,
			Second: highIndicator,
		}

		assert.False(t, rule.IsSatisfied(0))
	})
}

func TestUnderIndicatorRule(t *testing.T) {
	highIndicator := NewConstantIndicator(1)
	lowIndicator := NewConstantIndicator(0)

	t.Run("returns true when first indicator is under second indicator", func(t *testing.T) {
		rule := UnderIndicatorRule{
			First:  lowIndicator,
			Second: highIndicator,
		}

		assert.True(t, rule.IsSatisfied(0))
	})

	t.Run("returns false when first indicator is over second indicator", func(t *testing.T) {
		rule := UnderIndicatorRule{
			First:  highIndicator,
			Second: lowIndicator,
		}

		assert.False(t, rule.IsSatisfied(0))
	})
}

func TestPercentChangeRule(t *testing.T) {
	t.Run("returns false when percent change is less than the amount", func(t *testing.T) {
		series := mockTimeSeries("1", "1.1")
		rule := NewPercentChangeRule(NewClosePriceIndicator(series), 0.25)

		assert.False(t, rule.IsSatisfied(1))
	})

	t.Run("returns true when percent change is greater than the amount", func(t *testing.T) {
		series := mockTimeSeries("1", "1.11")
		rule := NewPercentChangeRule(NewClosePriceIndicator(series), 0.1)

		assert.True(t, rule.IsSatisfied(1))
	})
}
