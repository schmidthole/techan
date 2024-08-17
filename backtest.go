package techan

// The Backtest is a holder struct to run a simulated trading strategy against an account.
type Backtest struct {
	Tick       int
	Strategies []Strategy
	Allocator  Allocator
	Account    *Account
	History    *AccountHistory
}

// Create a new backtest with the provided strategies, allocator, and starting account.
// All strategies are assumed to be normalized and have the same length, indexes, and periods.
func NewBacktest(strategies []Strategy, allocator Allocator, account *Account) *Backtest {
	backtest := Backtest{
		Tick:       0,
		Strategies: strategies,
		Allocator:  allocator,
		Account:    account,
		History:    NewAccountHistory(),
	}

	// setup the initial account state to be one period before all data.
	// this will serve as a starting point before orders are executed.
	initialPeriod := strategies[0].Timeseries.Candles[0].Period.Advance(-1)
	backtest.History.ApplySnapshot(
		backtest.Account.ExportSnapshot(initialPeriod),
		&PricingSnapshot{Period: initialPeriod, Prices: Pricing{}},
	)

	return &backtest
}

// Execute the strategy for the given tick.
func (b *Backtest) ExecuteTick() error {
	prices := Pricing{}
	for _, strat := range b.Strategies {
		prices[strat.Security] = strat.Timeseries.Candles[b.Tick].ClosePrice
	}

	allocations := b.Allocator.Allocate(b.Tick, b.Strategies)

	tradePlan, err := CreateTradePlan(allocations, prices, b.Account)
	if err != nil {
		return err
	}

	for _, order := range *tradePlan {
		b.Account.ExecuteOrder(&order)
	}

	period := b.Strategies[0].Timeseries.Candles[b.Tick].Period
	b.History.ApplySnapshot(
		b.Account.ExportSnapshot(period),
		&PricingSnapshot{Period: period, Prices: prices},
	)

	return nil
}

// Advance the backtest's index tick.
func (b *Backtest) AdvanceTick() {
	b.Tick = b.Tick + 1
}
