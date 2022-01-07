package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/yashptel/binance-bot/pkg/models/db"
)

type Order struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	Status      string `json:"status"`
}

type OrderRepositoryImpl struct {
	db *firestore.Client
}

type OrderRepository interface {
	Create(order *Order) error
}

func NewOrderModel() (OrderRepository, error) {
	db, err := db.Connect()
	if err != nil {
		return nil, err
	}
	return &OrderRepositoryImpl{db: db}, nil
	// return OrderRepositoryImpl{}, nil
	// return OrderRepository{}, err
}

func (o OrderRepositoryImpl) Create(order *Order) error {
	res, err := o.db.Collection("orders").Doc(order.ID).Set(context.Background(), order)
	if err != nil {
		return err
	}
	fmt.Printf("Created a new order: %v\n", res)
	return nil
}
