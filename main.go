package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2/futures"
)

const (
	apiKey    = ""
	secretKey = ""
)

func main() {
	// futures.UseTestnet = true

	futuresClient := futures.NewClient(apiKey, secretKey) // USDT-M Futures
	// order, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").Do(context.Background())
	// if err != nil {
	// 	log.Println(err)
	// }

	order, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").Type("LIMIT").Side("BUY").TimeInForce("GTC").Quantity("0.001").Price("0.000001").Do(context.Background())
	if err != nil {
		log.Println(err)
	}

	// order, err := futuresClient.OrderTypeStopMarket().OrderList([]*futures.CreateOrderService{}).Type(futures.OrderTypeStopMarket).Symbol("ADAUSDT").Side(futures.SideTypeBuy).PositionSide(futures.PositionSideTypeLong).Quantity("10").StopPrice("2.0").Do(context.Background())

	data, _ := json.Marshal(order)
	fmt.Println(string(data), err)

	orders, err := futuresClient.NewListOrdersService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		log.Println(err)
	}

	for _, order := range orders {
		if order.Price == "56601.30" {
			data, _ := json.Marshal(order)
			fmt.Println(string(data))
		}
	}
	println("Hello World")
}
