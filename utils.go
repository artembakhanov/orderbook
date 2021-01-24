package orderbook

import "github.com/shopspring/decimal"

func less(a, b decimal.Decimal) bool {
	return a.LessThan(b)
}

func greater(a, b decimal.Decimal) bool {
	return a.GreaterThan(b)
}
