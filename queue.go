package orderbook

import (
	"container/list"
	"github.com/shopspring/decimal"
)

type Queue struct {
	price    decimal.Decimal
	quantity int64
	orders   *list.List
}

func CreateQueue(price decimal.Decimal) *Queue {
	return &Queue{
		price:    price,
		quantity: 0,
		orders:   list.New(),
	}
}

// return order acc. to FIFO order
func (q *Queue) First() *list.Element {

	return q.orders.Front()
}

// add order to the queue
func (q *Queue) Add(order *Order) *list.Element {
	q.quantity += order.quantity
	return q.orders.PushBack(order)
}

// remove order from the queue
func (q *Queue) Remove(element *list.Element) *Order {
	q.quantity -= element.Value.(*Order).quantity
	return q.orders.Remove(element).(*Order) // todo: if found ??
}
