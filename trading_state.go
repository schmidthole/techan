package techan

// TradingState is used in the MarketSnapshot to convey whether the market for each
// Security is open, closed, or halted
type TradingState int

// Possible trading states
const (
	OPEN int = iota
	CLOSED
	HALTED
)
