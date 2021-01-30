package main

import (
	ob "./orderbook"
	"encoding/csv"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

)

type trade struct {
}

func main() {
	closeTime := map[string]string{
		"USD000000TOD": "180000000000",//"174500000000",
		"USD000UTSTOM": "235000000000",
		"EUR_RUB__TOD": "150000000000",
		"EUR_RUB__TOM": "235000000000",
		"EURUSD000TOM": "235000000000",
		"EURUSD000TOD": "150000000000",
	}

	files, err := ioutil.ReadDir("./test")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		orderBooks := map[string]*ob.OrderBook{
			"USD000000TOD": ob.NewOrderBook(),
			"USD000UTSTOM": ob.NewOrderBook(),
			"EUR_RUB__TOD": ob.NewOrderBook(),
			"EUR_RUB__TOM": ob.NewOrderBook(),
			"EURUSD000TOM": ob.NewOrderBook(),
			"EURUSD000TOD": ob.NewOrderBook(),
		}

		fmt.Println(f.Name())
		csvfile, err := os.Open("./data/" + f.Name())
		if err != nil {
			log.Fatalln("Couldn't open the csv file", err)
		}

		// Parse the file
		r := csv.NewReader(csvfile)

		// Iterate through the records
		for i := 0; ; i++ {
			// Read each record from csv
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if i == 0 || record[3] >= closeTime[record[1]] {
				continue
			}

			if record[5] == "1" {
				quantity, _ := decimal.NewFromString(record[7])
				price, _ := decimal.NewFromString(record[6])
				if record[2] == "S" {
					if price.IsZero() {
						_, _, _, _, _ = orderBooks[record[1]].ProcessMarketOrder(ob.Sell, quantity)
					} else {
						_, _, _, _ = orderBooks[record[1]].ProcessLimitOrder(ob.Sell, record[4], quantity, price)
					}
				} else {
					if price.IsZero() {
						_, _, _, _, _ = orderBooks[record[1]].ProcessMarketOrder(ob.Buy, quantity)
					} else {
						_, _, _, _ = orderBooks[record[1]].ProcessLimitOrder(ob.Buy, record[4], quantity, price)
					}
				}
			} else if record[5] == "0" {
				orderBooks[record[1]].CancelOrder(record[4])
			} else {

			}

			var SpectrumValues [20]int
			j := 0
			sides := [] *OrderSide {orderBooks[record[1]].Bids, orderBooks[record[1]].Asks}

			for _, side :=  range sides {
			distance := side.MaxPriceQueue().price - side.MinPriceQueue().price
			step := distance / 10
			volume := side.Volume

			for price:= side.MinPriceQueue().price; price<side.MaxPriceQueue().price; price+=step{
				SpectrumValues[j] = side.GreaterThan(price).Len() - side.GreaterThan(price+step).Len() / volume * 100
				j++
			}

			}
		}

		for name, orderBook := range orderBooks {
			//f, err := os.Create("./out/" + strings.Split(f.Name(), ".")[0] + "_" + name + ".json")
			f, err := os.Create("./stakan_test/" + strings.Split(f.Name(), ".")[0] + "_" + name + ".json")
			if err != nil {
				fmt.Printf("There was a problem with %s %s", f.Name(), name)
			}

			//jjson, _ := orderBook.MarshalJSON()
			_, _ = f.WriteString(orderBook.String())
			_ = f.Close()

		}
	}
}
