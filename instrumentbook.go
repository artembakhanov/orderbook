package orderbook

import (
	"container/list"
	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/shopspring/decimal"
)

type InstrumentBook struct {
	instrument   string
	buys         *redblacktree.Tree
	sells        *redblacktree.Tree
	orders       map[int64]*list.Element
	queues       map[decimal.Decimal]*Queue
	quantitySell int64
	quantityBuy  int64
	bid          decimal.Decimal
	ask          decimal.Decimal
}

func CreateInstrumentBook(instrument string) *InstrumentBook {
	var instrumentBook = &InstrumentBook{
		instrument:   instrument,
		buys:         redblacktree.NewWith(OrderComparator),
		sells:        redblacktree.NewWith(OrderComparator),
		orders:       make(map[int64]*list.Element, 0),
		queues:       make(map[decimal.Decimal]*Queue, 0),
		quantitySell: 0,
		quantityBuy:  0,
		bid:          decimal.NewFromInt(-9223372036854775807),
		ask:          decimal.NewFromInt(9223372036854775807), // max value
	}

	return instrumentBook
}

func (b *InstrumentBook) ProcessOrder(id int64, kind OrderKind, price decimal.Decimal, quantity int64) *Order {
	// here we should process the order with that kind
	// if price == 0 then we need to process the order with market price
	// what if there not enough items with price == 0?¿¿?¿?!!
	// if price is not 0

	order := CreateOrder(id, kind, "", "", quantity, price)
	if kind == BUY {
		// log successful orders (trades) and then compare them to trades in the trade log to find out errors

		if price.IsZero() {
			// process market prices
			if b.quantitySell < quantity {
				// todo: log error since such order should not exist
				return nil
			}
			// here we need to choose all orders with overall quantity = quantity with order (price, time) first

		} else {
			b.processBuy(order)
		}
	} else {

		if price.IsZero() {
			// process market prices

			if b.quantityBuy < quantity {
				// todo: log error since such order should not exist
				return nil
			}
		} else {
			if price.LessThan(b.bid) {
				// todo: log this is a mistake, no need to put such a price
				return nil
			}
		}
	}

	return nil
}

func (b *InstrumentBook) RemoveOrder(id int64) *Order {
	e, ok := b.orders[id]
	if !ok {
		// here we should count an error because the order does not exist
		return nil
	}

	delete(b.orders, id)

	order := e.Value.(*Order)
	queue := b.queues[order.price]

	if order.kind == BUY {
		b.quantityBuy -= order.quantity
	} else {
		b.quantitySell -= order.quantity
	}

	queue.Remove(e)
	if queue.orders.Len() == 0 {
		// we can remove the price from the tree then
		delete(b.queues, order.price)
		if order.kind == BUY {
			b.buys.Remove(order.price)
		} else {
			b.buys.Remove(order.price)
		}
	}

	return order
}

func (b *InstrumentBook) postOrder(order *Order) bool {

	var tree *redblacktree.Tree
	var quantity *int64
	var bidask *decimal.Decimal
	var comp func(a, b decimal.Decimal) bool

	// kostyl?
	if order.kind == BUY {
		tree = b.buys
		quantity = &b.quantityBuy
		bidask = &b.bid
		comp = greater
	} else {
		tree = b.sells
		quantity = &b.quantitySell
		bidask = &b.ask
		comp = less
	}

	price := order.price
	queue, ok := b.queues[price]

	if !ok {
		// add price to the tree if does not exist
		// create queue (if price did not exist)
		queue = CreateQueue(price)
		b.queues[price] = queue
		tree.Put(price, queue)
	}

	// add order to the queue and get that element
	e := queue.Add(order)
	// add the element to b.orders
	b.orders[order.id] = e
	// update bid or ask
	if comp(price, *bidask) {
		*bidask = price
	}

	// update quantity
	*quantity += order.quantity

	return true
}

func (b *InstrumentBook) processBuy(order *Order) bool {
	if order.price.GreaterThan(b.ask) {
		// todo: log this is a mistake, no need to put such a price
		return false
	}

	queue, ok := b.queues[order.price]

	if !ok {
		// no such price, post the order

	}

	if queue.quantity < order.quantity {
		// not enough; remove queue and price from the tree
		// place the order
	} else if queue.quantity == order.quantity {
		// exactly what we need; remove queue and price from the tree
		// do not place the order
	} else {
		// enough; no removal
		// do not place the order
	}

	return true
}
