package techan

import "github.com/sdcoffey/big"

// An allocator is a functional component which provides a portfolio allocation across many Securities.
// A list of strategy structures is passed in to calculate the allocations based on rule analysis or
// timeseries data. All outputted allocation fractions will add up to 1.0 (big.ONE).
type Allocator interface {
	Allocate(index int, strategies []*Strategy) map[string]big.Decimal
}

// The naive allocator will proportionally allocate to all strategies which are triggered. A maximum
// position and total allocation can be provided in order to limit single position size and how much
// of a given portfolio is allocated.
type NaiveAllocator struct {
	maxSinglePositionFraction big.Decimal
	maxTotalPositionFraction  big.Decimal
}

// Sets up a new naive allocator. if any of the maximum fractions exceed 1.0, then 1.0 will be used.
// If the max single fraction exceeds the total fraction, then it will be set at the total fraction.
func NewNaiveAllocator(maxSinglePositionFraction big.Decimal, maxTotalPositionFraction big.Decimal) *NaiveAllocator {
	maxSingle := maxSinglePositionFraction
	if maxSinglePositionFraction.GT(big.ONE) {
		maxSingle = big.ONE
	}

	maxTotal := maxTotalPositionFraction
	if maxTotalPositionFraction.GT(big.ONE) {
		maxTotal = big.ONE
	}

	if maxSingle.GT(maxTotal) {
		maxSingle = maxTotal
	}

	return &NaiveAllocator{
		maxSinglePositionFraction: maxSingle,
		maxTotalPositionFraction:  maxTotal,
	}
}

func (na *NaiveAllocator) Allocate(index int, strategies []Strategy) map[string]big.Decimal {
	triggers := make([]string, 0)
	allocations := make(map[string]big.Decimal, 0)

	for _, s := range strategies {
		if s.Rule.IsSatisfied(index) {
			triggers = append(triggers, s.Security)
		}
	}

	if len(triggers) == 0 {
		return allocations
	}

	allocationFraction := na.maxTotalPositionFraction.Div(big.NewFromInt(len(triggers)))
	if allocationFraction.GT(na.maxSinglePositionFraction) {
		allocationFraction = na.maxTotalPositionFraction
	}

	for _, t := range triggers {
		allocations[t] = allocationFraction
	}

	return allocations
}
