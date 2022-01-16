package db

import (
	"context"
	"fmt"
	"log"
	"path"
	"sync"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/yashptel/binance-bot/pkg/config"

	"google.golang.org/api/option"
)

var once sync.Once
var conn *firestore.Client

func Connect() (*firestore.Client, error) {

	var err error
	once.Do(func() {
		var app *firebase.App

		cfg := config.GetConfig()
		root, err := config.GetRootDir()
		if err != nil {
			log.Println(err)
			return
		}

		keyPath := path.Join(root, cfg.Firebase.KeyPath)

		opt := option.WithCredentialsFile(keyPath)

		app, err = firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			fmt.Printf("error initializing app: %v\n", err)
			return
		}
		conn, err = app.Firestore(context.Background())
		if err != nil {
			fmt.Printf("error initializing firestore: %v\n", err)
			return
		}
	})
	return conn, err
}

func Disconnect() {
	if conn != nil {
		conn.Close()
	}
}
