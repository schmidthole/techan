package techan

import (
	"testing"

	"github.com/schmidthole/big"
	"github.com/stretchr/testify/assert"
)

var MOCK_SECURITY = "FAKE"

var mockOrder = Order{
	Security: MOCK_SECURITY,
	Side:     BUY,
	Amount:   big.NewFromString("2"),
	Price:    big.NewFromString("2"),
}

func mockPosition() *Position {
	return NewPosition(&mockOrder)
}

func TestPosition_NewPosition(t *testing.T) {
	position := mockPosition()

	assert.Equal(t, position.Side, mockOrder.Side)
	assert.Equal(t, position.Amount, mockOrder.Amount)
	decimalEquals(t, 2.0, position.Price)
	decimalEquals(t, 2.0, position.AvgEntryPrice)
}

func TestPosition_ExecuteOrder_Add(t *testing.T) {
	position := mockPosition()

	orderAdd := Order{
		Side:   BUY,
		Amount: big.NewFromString("2"),
		Price:  big.NewFromString("3"),
	}

	err := position.ExecuteOrder(&orderAdd)
	assert.Nil(t, err)

	decimalEquals(t, 4.0, position.Amount)
	decimalAlmostEquals(t, big.NewDecimal(2.5), position.AvgEntryPrice, 0.2)
	decimalEquals(t, 3.0, position.Price)
}

func TestPosition_ExecuteOrder_Sell(t *testing.T) {
	position := mockPosition()

	orderSell := Order{
		Side:   SELL,
		Amount: big.ONE,
		Price:  big.NewFromString("3"),
	}

	err := position.ExecuteOrder(&orderSell)
	assert.Nil(t, err)

	decimalEquals(t, 1.0, position.Amount)
	decimalEquals(t, 3.0, position.Price)
	decimalEquals(t, 2.0, position.AvgEntryPrice)
	assert.False(t, position.IsClosed())
}

func TestPosition_ExecuteOrder_SellAll_Close(t *testing.T) {
	position := mockPosition()

	orderSell := Order{
		Side:   SELL,
		Amount: big.NewFromString("2"),
		Price:  big.NewFromString("3"),
	}

	err := position.ExecuteOrder(&orderSell)
	assert.Nil(t, err)

	decimalEquals(t, 0.0, position.Amount)
	decimalEquals(t, 3.0, position.Price)
	decimalEquals(t, 2.0, position.AvgEntryPrice)
	assert.True(t, position.IsClosed())
}

func TestPosition_ExecuteOrder_SellTooMany(t *testing.T) {
	position := mockPosition()

	orderSell := Order{
		Side:   SELL,
		Amount: big.NewFromString("3"),
		Price:  big.NewFromString("3"),
	}

	err := position.ExecuteOrder(&orderSell)
	assert.NotNil(t, err)

	decimalEquals(t, 2.0, position.Amount)
	decimalEquals(t, 2.0, position.Price)
}

func TestPosition_ExecuteOrder_TryToSellShort(t *testing.T) {
	position := mockPosition()

	orderSell := Order{
		Side:   SELL,
		Amount: big.NewFromString("3"),
		Price:  big.NewFromString("3"),
	}

	position.Side = SELL

	err := position.ExecuteOrder(&orderSell)
	assert.NotNil(t, err)
}

func TestPosition_UpdatePrice(t *testing.T) {
	position := mockPosition()

	decimalEquals(t, 2.0, position.Price)
	position.UpdatePrice(big.NewDecimal(11.0))
	decimalEquals(t, 11.0, position.Price)
}

func TestPosition_UnrealizedEquity(t *testing.T) {
	position := mockPosition()
	equity := position.UnrealizedEquity()
	decimalEquals(t, 4.0, equity)
}
