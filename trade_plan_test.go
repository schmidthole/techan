package techan

import (
	"reflect"
	"testing"

	"github.com/schmidthole/big"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name        string
	allocations Allocations
	pricing     Pricing
	account     Account
	shouldError bool
	result      *TradePlan
}{
	{
		name: "missing pricing data",
		allocations: Allocations{
			"ONE": big.NewDecimal(0.5),
			"TWO": big.NewDecimal(0.4),
		},
		pricing: Pricing{
			"ONE": big.NewDecimal(1.0),
		},
		account:     *NewAccount(),
		shouldError: true,
		result:      &TradePlan{},
	},
	{
		name: "add single position",
		allocations: Allocations{
			"ONE": big.NewDecimal(0.5),
		},
		pricing: Pricing{
			"ONE": big.NewDecimal(1.0),
		},
		account: Account{
			Cash:      big.NewDecimal(10.00),
			Positions: map[string]*Position{},
		},
		shouldError: false,
		result: &TradePlan{
			Order{
				Side:     BUY,
				Security: "ONE",
				Amount:   big.NewDecimal(5.0),
				Price:    big.NewDecimal(1.0),
			},
		},
	},
	{
		name:        "sell single position",
		allocations: Allocations{},
		pricing: Pricing{
			"ONE": big.NewDecimal(1.0),
		},
		account: Account{
			Cash: big.NewDecimal(0.00),
			Positions: map[string]*Position{
				"ONE": {
					Security:      "ONE",
					Side:          BUY,
					Amount:        big.NewDecimal(5.0),
					AvgEntryPrice: big.NewDecimal(1.0),
					Price:         big.NewDecimal(1.0),
				},
			},
		},
		shouldError: false,
		result: &TradePlan{
			Order{
				Side:     SELL,
				Security: "ONE",
				Amount:   big.NewDecimal(5.0),
				Price:    big.NewDecimal(1.0),
			},
		},
	},
	{
		name:        "add one and sell one",
		allocations: Allocations{"TWO": big.NewDecimal(0.50)},
		pricing: Pricing{
			"ONE": big.NewDecimal(1.0),
			"TWO": big.NewDecimal(2.0),
		},
		account: Account{
			Cash: big.NewDecimal(5.00),
			Positions: map[string]*Position{
				"ONE": {
					Security:      "ONE",
					Side:          BUY,
					Amount:        big.NewDecimal(5.0),
					AvgEntryPrice: big.NewDecimal(1.0),
					Price:         big.NewDecimal(1.0),
				},
			},
		},
		shouldError: false,
		result: &TradePlan{
			Order{
				Side:     SELL,
				Security: "ONE",
				Amount:   big.NewDecimal(5.0),
				Price:    big.NewDecimal(1.0),
			},
			Order{
				Side:     BUY,
				Security: "TWO",
				Amount:   big.NewDecimal(2.0),
				Price:    big.NewDecimal(2.0),
			},
		},
	},
	{
		name:        "add to open position",
		allocations: Allocations{"ONE": big.NewDecimal(0.60)},
		pricing: Pricing{
			"ONE": big.NewDecimal(1.0),
		},
		account: Account{
			Cash: big.NewDecimal(5.00),
			Positions: map[string]*Position{
				"ONE": {
					Security:      "ONE",
					Side:          BUY,
					Amount:        big.NewDecimal(5.0),
					AvgEntryPrice: big.NewDecimal(1.0),
					Price:         big.NewDecimal(1.0),
				},
			},
		},
		shouldError: false,
		result: &TradePlan{
			Order{
				Side:     BUY,
				Security: "ONE",
				Amount:   big.NewDecimal(1.0),
				Price:    big.NewDecimal(1.0),
			},
		},
	},
	{
		name:        "sell from open position",
		allocations: Allocations{"ONE": big.NewDecimal(0.60)},
		pricing: Pricing{
			"ONE": big.NewDecimal(1.0),
		},
		account: Account{
			Cash: big.NewDecimal(2.00),
			Positions: map[string]*Position{
				"ONE": {
					Security:      "ONE",
					Side:          BUY,
					Amount:        big.NewDecimal(8.0),
					AvgEntryPrice: big.NewDecimal(1.0),
					Price:         big.NewDecimal(1.0),
				},
			},
		},
		shouldError: false,
		result: &TradePlan{
			Order{
				Side:     SELL,
				Security: "ONE",
				Amount:   big.NewDecimal(2.0),
				Price:    big.NewDecimal(1.0),
			},
		},
	},
	{
		name: "add one, sell one, and add to one",
		allocations: Allocations{
			"ONE": big.NewDecimal(0.60),
			"TWO": big.NewDecimal(0.30),
		},
		pricing: Pricing{
			"ONE":   big.NewDecimal(1.0),
			"TWO":   big.NewDecimal(2.0),
			"THREE": big.NewDecimal(3.0),
		},
		account: Account{
			Cash: big.NewDecimal(2.00),
			Positions: map[string]*Position{
				"ONE": {
					Security:      "ONE",
					Side:          BUY,
					Amount:        big.NewDecimal(5.0),
					AvgEntryPrice: big.NewDecimal(1.0),
					Price:         big.NewDecimal(1.0),
				},
				"THREE": {
					Security:      "THREE",
					Side:          BUY,
					Amount:        big.NewDecimal(1.0),
					AvgEntryPrice: big.NewDecimal(3.0),
					Price:         big.NewDecimal(3.0),
				},
			},
		},
		shouldError: false,
		result: &TradePlan{
			Order{
				Side:     SELL,
				Security: "THREE",
				Amount:   big.NewDecimal(1.0),
				Price:    big.NewDecimal(3.0),
			},
			Order{
				Side:     BUY,
				Security: "ONE",
				Amount:   big.NewDecimal(1.0),
				Price:    big.NewDecimal(1.0),
			},
			Order{
				Side:     BUY,
				Security: "TWO",
				Amount:   big.NewDecimal(1.0),
				Price:    big.NewDecimal(2.0),
			},
		},
	},
}

func Test_CreateTradePlan(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := CreateTradePlan(tt.allocations, tt.pricing, &tt.account)

			if tt.shouldError {
				assert.NotNil(t, err)
			} else {
				if !reflect.DeepEqual(tt.result, plan) {
					t.Errorf("got %v, want %v", plan, tt.result)
				}
			}
		})
	}
}
