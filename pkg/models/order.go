package models

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/yashptel/binance-bot/pkg/models/db"
	"google.golang.org/api/iterator"
)

type Order struct {
	ID           string                   `json:"id"`
	Symbol       string                   `json:"symbol"`
	Price        float64                  `json:"price"`
	Qty          string                   `json:"qty"`
	Side         futures.SideType         `json:"side"`
	PositionSide futures.PositionSideType `json:"position_side"`
	StopPrice    float64                  `json:"stop_price"`
}

func (o *Order) SetStopPrice(percentage float64) {

	sl := 1.00

	if o.PositionSide == futures.PositionSideTypeLong {
		sl = (1 - (percentage / 100))
	} else if o.PositionSide == futures.PositionSideTypeShort {
		sl = (1 + (percentage / 100))
	}
	o.Price = o.Price * sl
}

// futures.Order{
// 	Symbol:           symbol,
// 	OrderID:          0,
// 	ClientOrderID:    "",
// 	Price:            "",
// 	ReduceOnly:       false,
// 	OrigQuantity:     "",
// 	ExecutedQuantity: "",
// 	CumQuantity:      "",
// 	CumQuote:         "",
// 	Status:           "",
// 	TimeInForce:      "",
// 	Type:             "",
// 	Side:             "",
// 	StopPrice:        "",
// 	Time:             0,
// 	UpdateTime:       0,
// 	WorkingType:      "",
// 	ActivatePrice:    "",
// 	PriceRate:        "",
// 	AvgPrice:         "",
// 	OrigType:         "",
// 	PositionSide:     "",
// 	PriceProtect:     false,
// 	ClosePosition:    false,
// }

type OrderRepositoryImpl struct {
	db *firestore.Client
}

type OrderRepository interface {
	GetAll() ([]*Order, error)
	Create(order *Order) (*Order, error)
	Get(id string) (*Order, error)
	Delete(id string) error
}

func NewOrderModel() (OrderRepository, error) {
	db, err := db.Connect()
	if err != nil {
		return nil, err
	}
	return &OrderRepositoryImpl{db: db}, nil
}

func (o OrderRepositoryImpl) Create(order *Order) (*Order, error) {

	ref := o.db.Collection("orders").NewDoc()
	order.ID = ref.ID

	_, err := ref.Set(context.Background(), order)
	if err != nil {
		return nil, err
	}
	return o.Get(ref.ID)
}

func (o OrderRepositoryImpl) GetAll() ([]*Order, error) {
	var orders []*Order
	iter := o.db.Collection("orders").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var order Order
		doc.DataTo(&order)
		orders = append(orders, &order)
	}
	return orders, nil
}

func (o OrderRepositoryImpl) Get(id string) (*Order, error) {
	doc, err := o.db.Collection("orders").Doc(id).Get(context.Background())
	if err != nil {
		return nil, err
	}
	var order Order
	doc.DataTo(&order)
	return &order, nil
}

func (o OrderRepositoryImpl) Delete(id string) error {
	_, err := o.db.Collection("orders").Doc(id).Delete(context.Background())
	return err
}
