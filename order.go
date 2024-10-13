package techan

import (
	"time"

	"github.com/sdcoffey/big"
)

// OrderSide is a simple enumeration representing the side of an Order (buy or sell)
type OrderSide string

// BUY and SELL enumerations
const (
	BUY  OrderSide = "BUY"
	SELL OrderSide = "SELL"
)

// OrderType defines common order types accepted by brokers
type OrderType string

// OrderType enumerations, we only support the basics at this time
// there are more complex types such as "stop limit" "trail" etc. that may be added later
const (
	MARKET OrderType = "MKT"
	LIMIT  OrderType = "LMT"
	STOP   OrderType = "STP"
)

// TimeInForce defines how long an order is good for before it is auto cancelled
type TimeInForce string

// TimeInForce option enumerations
const (
	GTC TimeInForce = "GTC"
	OPG TimeInForce = "OPG"
	DAY TimeInForce = "DAY"
	IOC TimeInForce = "IOC"
)

// Order represents a trade execution (buy or sell) with associated metadata.
type Order struct {
	Side          OrderSide
	Security      string
	Price         big.Decimal
	Type          string
	Amount        big.Decimal
	TimeInForce   TimeInForce
	ExecutionTime time.Time
}

// Return the total cost to execute the order.
func (o *Order) CostBasis() big.Decimal {
	return o.Amount.Mul(o.Price)
}
