package techan

import (
	"fmt"

	"github.com/sdcoffey/big"
)

// Positions holds iformation about an open position
type Position struct {
	Security      string
	Side          OrderSide
	Amount        big.Decimal
	AvgEntryPrice big.Decimal
	Price         big.Decimal
}

// NewPosition returns a new Position with the passed-in order as the open order
func NewPosition(order *Order) *Position {
	pos := new(Position)
	pos.Security = order.Security
	pos.Side = order.Side
	pos.Amount = order.Amount
	pos.AvgEntryPrice = order.Price
	pos.Price = order.Price
	
	return pos
}

// ExecuteOrder takes a new order to apply to the position and inceases/deducts the amount from
// the current position. If a buy order is placed, it will also recalculate the average entry
// price for the position. In all cases, if the order results in a difference to the cash value
// of the account placing the order a positive (for sells) or negative (for buys) big.Decimal is
// returned to denote how the account's cash should be modified.
//
// Only long positions are supported currently.
func (p *Position) ExecuteOrder(order *Order) error {
	if (p.Side == BUY) && (order.Side == BUY) {
		newTotalValue := p.AvgEntryPrice.Mul(p.Amount).Add(order.Price.Mul(order.Amount))
		newAmount := p.Amount.Add(order.Amount)

		p.AvgEntryPrice = newTotalValue.Div(newAmount)
		p.Amount = newAmount

		return nil
	} else if (p.Side == BUY) && (order.Side == SELL) {
		intermediate := p.Amount.Sub(order.Amount)
		if intermediate.LT(big.ZERO) {
			return fmt.Errorf(
				"invalid long sell on position: %v. tried to sell %v when position has %v",
				p.Security,
				order.Amount.String(),
				p.Amount.String(),
			)
		}

		p.Amount = intermediate

		return nil
	} else {
		return fmt.Errorf("unsupported order operation %v on %v position", order.Side, p.Side)
	}
}

// Returns if the position is closed, meaning there is currently a zero amount.
func (p *Position) IsClosed() bool {
	return p.Amount.EQ(big.ZERO)
}

// CostBasis returns the price to enter this order
func (p *Position) CostBasis() big.Decimal {
	return p.Amount.Mul(p.AvgEntryPrice)
}

// Update the current price to reflect realtime values and to calculate unrealized gains/equity
func (p *Position) UpdatePrice(newPrice big.Decimal) {
	p.Price = newPrice
}

// Calculate the unrealized equity of an open position
func (p *Position) UnrealizedEquity() big.Decimal {
	return p.Amount.Mul(p.Price)
}

// Calculate the unrealized gains of an open position
func (p *Position) UnrealizedGains() big.Decimal {
	return p.Price.Sub(p.AvgEntryPrice).Mul(p.Amount)
}
