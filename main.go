package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	binance "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

const (
	apiKey    = ""
	secretKey = ""
)

const (
	diff       = 0.01
	symbol     = "BTCUSDT"
	entryPrice = 56601.30
)

func main() {

	// Create a new client
	futuresClient := binance.NewFuturesClient(apiKey, secretKey)

	errHandler := func(err error) {
		fmt.Println(err)
	}

	WsAggTradeHandler := func(event *futures.WsAggTradeEvent) {
		price, err := strconv.ParseFloat(event.Price, 64)
		if err != nil {
			log.Println(err)
		}

		priceDiff := getPercentageDiff(entryPrice, price)
		if priceDiff <= diff {
			fmt.Println("Price diff: ", priceDiff, " %", " Price: ", price)
		}
	}

	doneC, _, _ := futures.WsAggTradeServe(symbol, WsAggTradeHandler, errHandler)
	<-doneC

	res, err := futuresClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		panic(err)
	}
	futuresClient.NewChangeLeverageService().Symbol("ETHUSDT").Leverage(55).Do(context.Background())

	for _, v := range res.Symbols {
		if v.Status == "TRADING" && v.ContractType == "PERPETUAL" {
			// res, _ := client.NewGetLeverageBracketService().Symbol(v.Symbol).Do(context.Background())
			// data, _ := json.Marshal(res)
			fmt.Println(v.Symbol)
			// fmt.Println(string(data))
		}
	}

	// client.NewListPricesService().Symbol("BTCUSDT").Do(context.Background())

	// data, _ := json.Marshal(order)
	// fmt.Println(string(data), err)

	// orders, err := futuresClient.NewListOrdersService().Symbol("BTCUSDT").Do(context.Background())
	// if err != nil {
	// 	log.Println(err)
	// }

	// for _, order := range orders {
	// 	if order.Price == "56601.30" {
	// 		data, _ := json.Marshal(order)
	// 		fmt.Println(string(data))
	// 	}
	// }
}

func getPercentageDiff(v1, v2 float64) float64 {
	return (math.Abs(v1-v2) / ((v1 + v2) / 2)) * 100
}
