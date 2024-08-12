package techan

import "github.com/sdcoffey/big"

type AccountSnapshot struct {
	Period    TimePeriod
	Equity    big.Decimal
	Cash      big.Decimal
	Positions []*PositionSnapshot
}

type AccountHistory struct {
	snapshots []*AccountSnapshot
}
