package orderbook

import (
	"github.com/shopspring/decimal"
)

type OrderKind int

const BUY OrderKind = 0
const SELL OrderKind = 1

type Order struct {
	id         int64
	kind       OrderKind // either BUY or SELL
	time       string
	instrument string // maybe we do not need it
	quantity   int64
	price      decimal.Decimal
}

func OrderComparator(a, b interface{}) int {
	return a.(decimal.Decimal).Cmp(b.(decimal.Decimal))
}

func CreateOrder(id int64, kind OrderKind, time string, instrument string, quantity int64, price decimal.Decimal) *Order {
	return &Order{
		id:         id,
		kind:       kind,
		time:       time,
		instrument: instrument,
		quantity:   quantity,
		price:      price,
	}
}
