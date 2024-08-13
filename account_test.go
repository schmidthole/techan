package techan

import (
	"testing"

	"github.com/sdcoffey/big"
	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acct := NewAccount()
	decimalEquals(t, 0.0, acct.Cash)
	assert.Equal(t, 0, len(acct.Positions))
	assert.Equal(t, 0, len(acct.TradeRecord))
}

func TestAccount_OpenPosition(t *testing.T) {
	acct := NewAccount()
	acct.Positions[MOCK_SECURITY] = mockPosition()

	pos, exists := acct.OpenPosition(MOCK_SECURITY)
	assert.True(t, exists)
	assert.Equal(t, MOCK_SECURITY, pos.Security)
	assert.Equal(t, BUY, pos.Side)
	decimalEquals(t, 2.0, pos.Amount)
	decimalEquals(t, 2.0, pos.Price)

	notPos, exists := acct.OpenPosition("NOTTHERE")
	assert.False(t, exists)
	assert.Nil(t, notPos)
}

func TestAccount_Deposit(t *testing.T) {
	acct := NewAccount()
	acct.Deposit(big.NewDecimal(1.0))
	decimalEquals(t, 1.0, acct.Cash)
	acct.Deposit(big.NewDecimal(1001.15))
	decimalEquals(t, 1002.15, acct.Cash)
}

func TestAccount_Withdraw(t *testing.T) {
	acct := NewAccount()
	acct.Deposit(big.NewDecimal(1001.15))

	err := acct.Withdraw(big.NewDecimal(1.15))
	assert.Nil(t, err)
	decimalEquals(t, 1000.00, acct.Cash)

	err = acct.Withdraw(big.NewDecimal(999.00))
	assert.Nil(t, err)
	decimalEquals(t, 1.0, acct.Cash)

	err = acct.Withdraw(big.NewDecimal(1.01))
	assert.NotNil(t, err)

	err = acct.Withdraw(big.NewDecimal(1.0))
	assert.Nil(t, err)
	decimalEquals(t, 0.0, acct.Cash)
}

func TestAccount_UpdatePrices(t *testing.T) {
	acct := NewAccount()
	acct.Positions[MOCK_SECURITY] = mockPosition()

	acct.UpdatePrices(map[string]big.Decimal{MOCK_SECURITY: big.NewDecimal(3.0)})
	pos, _ := acct.OpenPosition(MOCK_SECURITY)

	decimalEquals(t, 3.0, pos.Price)
}

func TestAccount_HasSufficientFunds(t *testing.T) {
	acct := NewAccount()

	res := acct.HasSufficientFunds(&mockOrder)
	assert.False(t, res)

	acct.Deposit(big.NewDecimal(3.0))
	res = acct.HasSufficientFunds(&mockOrder)
	assert.False(t, res)

	acct.Deposit(big.NewDecimal(1.0))
	res = acct.HasSufficientFunds(&mockOrder)
	assert.True(t, res)

	acct.Deposit(big.NewDecimal(100.0))
	res = acct.HasSufficientFunds(&mockOrder)
	assert.True(t, res)
}

func TestAccount_Equity(t *testing.T) {
	acct := NewAccount()

	decimalEquals(t, 0.0, acct.Equity())

	acct.Deposit(big.NewDecimal(4.0))
	decimalEquals(t, 4.0, acct.Equity())

	acct.Positions[MOCK_SECURITY] = mockPosition()
	decimalEquals(t, 8.0, acct.Equity())

	acct.UpdatePrices(map[string]big.Decimal{MOCK_SECURITY: big.NewDecimal(3.0)})
	decimalEquals(t, 10.0, acct.Equity())
}
