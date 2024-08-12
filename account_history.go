package techan

import (
	"fmt"

	"github.com/sdcoffey/big"
)

// An AccountSnapshot provides the point in time state of an account and its positions.
type AccountSnapshot struct {
	Period    TimePeriod
	Equity    big.Decimal
	Cash      big.Decimal
	Positions []*PositionSnapshot
}

// The PricingSnapshot provides point in time pricing information for all tracked securities,
// including the bechmark if used. The period for each snapshot should match a period supplied
// in an account snapshot.
type PricingSnapshot struct {
	Period TimePeriod
	Prices map[string]big.Decimal
}

// The AccountHistory contains a record of point in time account snapshots as well as a list of all
// of the Securities tracked by the account over time.
type AccountHistory struct {
	Benchmark  string
	Securities []string
	Prices     []*PricingSnapshot
	Snapshots  []*AccountSnapshot
}

// Add a new account and pricing snapshot to the history.
func (ah *AccountHistory) ApplySnapshot(accountSnapshot *AccountSnapshot, pricingSnapshot *PricingSnapshot) error {
	if !accountSnapshot.Period.Start.Equal(pricingSnapshot.Period.Start) {
		return fmt.Errorf(
			"start periods for account and pricing snapshots do not match: %v <-> %v",
			accountSnapshot.Period.Start,
			pricingSnapshot.Period.Start,
		)
	}

	ah.Snapshots = append(ah.Snapshots, accountSnapshot)
	ah.Prices = append(ah.Prices, pricingSnapshot)

	return nil
}

// Helper function to return the last index of snapshot data.
func (ah *AccountHistory) LastIndex() int {
	return len(ah.Snapshots) - 1
}
