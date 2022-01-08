package config

type Binance struct {
	APIKey    string `envconfig:"BINANCE_API_KEY"`
	APISecret string `envconfig:"BINANCE_API_SECRET"`
}
