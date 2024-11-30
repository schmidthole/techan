package techan

import (
	"testing"

	"github.com/schmidthole/big"
	"github.com/stretchr/testify/assert"
)

func TestAllocator_NewNaiveAllocator(t *testing.T) {
	alc := NewNaiveAllocator(big.NewDecimal(0.2), big.NewDecimal(1.0))
	decimalEquals(t, 0.2, alc.maxSinglePositionFraction)
	decimalEquals(t, 1.0, alc.maxTotalPositionFraction)

	alc = NewNaiveAllocator(big.NewDecimal(0.9), big.NewDecimal(0.8))
	decimalEquals(t, 0.8, alc.maxSinglePositionFraction)
	decimalEquals(t, 0.8, alc.maxTotalPositionFraction)

	alc = NewNaiveAllocator(big.NewDecimal(1.1), big.NewDecimal(1.0))
	decimalEquals(t, 1.0, alc.maxSinglePositionFraction)
	decimalEquals(t, 1.0, alc.maxTotalPositionFraction)

	alc = NewNaiveAllocator(big.NewDecimal(0.3), big.NewDecimal(1.8))
	decimalEquals(t, 0.3, alc.maxSinglePositionFraction)
	decimalEquals(t, 1.0, alc.maxTotalPositionFraction)
}

func TestAllocator_NaiveAllocatorAllocate(t *testing.T) {
	alc := NewNaiveAllocator(big.NewDecimal(0.4), big.NewDecimal(1.0))

	rule1 := mockRule{[]bool{true, true, true, false}}
	rule2 := mockRule{[]bool{false, true, true, false}}
	rule3 := mockRule{[]bool{false, true, false, false}}

	strat1 := Strategy{Security: "ONE", Rule: rule1}
	strat2 := Strategy{Security: "TWO", Rule: rule2}
	strat3 := Strategy{Security: "THREE", Rule: rule3}

	strats := []Strategy{strat1, strat2, strat3}

	alc0 := alc.Allocate(0, strats)
	assert.Equal(t, 1, len(alc0))
	decimalAlmostEquals(t, big.NewDecimal(0.4), alc0["ONE"], 0.1)

	alc1 := alc.Allocate(1, strats)
	assert.Equal(t, 3, len(alc1))
	decimalAlmostEquals(t, big.NewDecimal(0.33), alc1["ONE"], 0.1)
	decimalAlmostEquals(t, big.NewDecimal(0.33), alc1["TWO"], 0.1)
	decimalAlmostEquals(t, big.NewDecimal(0.33), alc1["THREE"], 0.1)

	alc2 := alc.Allocate(2, strats)
	assert.Equal(t, 2, len(alc2))
	decimalAlmostEquals(t, big.NewDecimal(0.4), alc2["ONE"], 0.1)
	decimalAlmostEquals(t, big.NewDecimal(0.4), alc2["TWO"], 0.1)

	alc3 := alc.Allocate(3, strats)
	assert.Equal(t, 0, len(alc3))
}
