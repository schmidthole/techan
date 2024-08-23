# Techan

TechAn is a **tech**nical **an**alysis library for Go! It provides a suite of tools and frameworks to analyze financial data and make trading decisions.

**NOTE:** This is my fork of the original techan library found at [https://github.com/sdcoffey/techan](https://github.com/sdcoffey/techan).
My fork preserves the structure of the timeseries and inidicator calculation while completely overhauling the strategy and trade simulation
to make it more functional for my real use.

## Features 
* Basic and advanced technical analysis indicators
* Strategy building
* Backtesting and account simulation

## Installation
```sh
$ go get github.com/schmidthole/techan
```

## Examples

### Quickstart
```go
series := techan.NewTimeSeries()

// fetch this from your preferred exchange
dataset := [][]string{
	// Timestamp, Open, Close, High, Low, volume
	{"1234567", "1", "2", "3", "5", "6"},
}

for _, datum := range dataset {
	start, _ := strconv.ParseInt(datum[0], 10, 64)
	period := techan.NewTimePeriod(time.Unix(start, 0), time.Hour*24)

	candle := techan.NewCandle(period)
	candle.OpenPrice = big.NewFromString(datum[1])
	candle.ClosePrice = big.NewFromString(datum[2])
	candle.MaxPrice = big.NewFromString(datum[3])
	candle.MinPrice = big.NewFromString(datum[4])

	series.AddCandle(candle)
}

closePrices := techan.NewClosePriceIndicator(series)
movingAverage := techan.NewEMAIndicator(closePrices, 10) // Create an exponential moving average with a window of 10

fmt.Println(movingAverage.Calculate(0).FormattedString(2))
```

### Creating trading strategies
A `Strategy` in Techan is the application of a `Rule` against a particular security/asset. For ease of reference,
the `Strategy` struct contains the original `Timeseries`, all `Indicators` used to calculate the `Rule`, and the
rull itself.

```go
// using the timeseries and indicators from above...

// we simply create a rule that is satisfied if the price is above the moving average.
rule := OverIndicatorRule {
    First: movingAverage,
    Second: closePrices
}

strategy := Strategy{
    Security: "TEST",
    Timeseries: series,
    Indicators: map[string]Indicator{
        "ema10": movingAverage,
    },
    Rule: rule,
}
```

Strategies against individual securities/assets can be combined into a more comprehensive strategy using an 
`Allocator`. An allocator is simply an interface that accepts a list of strategies and outputs a portfolio
allocation. This can be as simple or as complex as needed.

```go
// again using the strategy defined above...

// we create a naive allocator which allows a single position to be 50% of a portfolio and 100% of a portfoliio
// to be allocated.
allocator := NewNaiveAllocator(big.NewDecimal(0.5), big.NewDecimal(1.0))

allocations := allocator.Allocate(0, []Strategy{strategy})
```

Using the outputted portfolio `Allocations`, we can then create a `TradePlan` by combining information from an 
`Account`. The `TradePlan` will use the current security prices, open positions in the account, and the allocations 
provided to identify the orders which can be placed to achieve the desired allocation. This can be used to calculate
orders to execute with a broker based on the strategy results.

```go
account := NewAccount()
account.Deposit(big.NewDecimal(10000.00))

// create pricing data from timeseries candles above.
pricing := Pricing{
    "ONE": big.NewDecimal(2.0)
}

tradePlan, _ := CreateTradePlan(allocations, pricing, account)
```

Putting it all together, we can use all of the components to run a full backtest of a strategy against the historical
`Timeseries` data. Combining multiple `Strategy` objects along with the `Allocator` allow for complex strategies to
be modelled over time.

```go
// we are still using everuyting defined above...

backtest := NewBacktest([]Strategy{strategy}, allocator, account)
history, _ := backtest.Run()
```

Running a backtest will about an `AccountHistory` object, which contains all of the relevant backtest 
snapshots/results over time. The account history can be used to perform analysis of the strategy's 
performance metrics.

```go
profit := history.TotalProfit()
```

### Credits
Techan is heavily influenced by the great [ta4j](https://github.com/ta4j/ta4j). Many of the ideas and frameworks in this library owe their genesis to the great work done over there.

### License

Techan is released under the MIT license. See [LICENSE](./LICENSE) for details.
