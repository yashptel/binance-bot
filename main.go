package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"

	binance "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/yashptel/binance-bot/pkg/models"
)

const (
	apiKey    = ""
	secretKey = ""
)

const (
	diff       = 10.01
	symbol     = "BTCUSDT"
	entryPrice = 42960.30
)

var mu sync.Mutex
var trades int

func takeTrade() {
	mu.Lock()
	defer mu.Unlock()
	// fmt.Println("Trying to take trade")
	if trades == 0 {
		fmt.Println("Taking trade")

		// Create a new client
		futuresClient := binance.NewFuturesClient(apiKey, secretKey)

		// Create a new order
		mainOrder := futuresClient.NewCreateOrderService()
		mainOrder.Symbol(symbol)
		mainOrder.Side(futures.SideTypeBuy)
		mainOrder.PositionSide(futures.PositionSideTypeLong)
		mainOrder.Type(futures.OrderTypeMarket)
		mainOrder.Quantity("0.1")

		// Create a new order
		stopLossOrder := futuresClient.NewCreateOrderService()
		stopLossOrder.Symbol(symbol)
		stopLossOrder.Side(futures.SideTypeSell)
		stopLossOrder.PositionSide(futures.PositionSideTypeLong)
		stopLossOrder.Type(futures.OrderTypeStopMarket)
		stopLossOrder.Quantity("0.1")
		stopLossOrder.StopPrice("39000.30")
		stopLossOrder.TimeInForce(futures.TimeInForceTypeGTC)

		// res, err := order.Do(context.Background())
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// data, _ := json.Marshal(res)
		// fmt.Println(data)

		batchOrdersService := futuresClient.NewCreateBatchOrdersService()
		batchOrdersService.OrderList([]*futures.CreateOrderService{mainOrder, stopLossOrder})
		res, err := batchOrdersService.Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		data, _ := json.Marshal(res)
		fmt.Println(string(data))

		trades++
	}
}

func main() {

	orderRepo, err := models.NewOrderModel()
	if err != nil {
		log.Fatal(err)
	}

	err = orderRepo.Create(&models.Order{
		ID:          "1",
		UserID:      "1",
		ProductID:   "1",
		ProductName: "1",
		Quantity:    1,
		Price:       1,
		Status:      "1",
	})
	if err != nil {
		log.Fatal(err)
	}
	return

	futures.UseTestnet = true

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
			// fmt.Println("Price diff: ", priceDiff, " %", " Price: ", price)
			takeTrade()
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
