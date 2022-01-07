package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

var once sync.Once
var conn *firestore.Client

func Connect() (*firestore.Client, error) {

	var err error
	once.Do(func() {
		var app *firebase.App

		dir, _ := os.Getwd()
		keyPath := dir + "/firebase-key.json"

		opt := option.WithCredentialsFile(keyPath)

		app, err = firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			fmt.Printf("error initializing app: %v\n", err)
			return
		}
		conn, err = app.Firestore(context.Background())
	})
	return conn, err
}

func Disconnect() {
	if conn != nil {
		conn.Close()
	}
}
