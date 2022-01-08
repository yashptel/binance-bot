package config

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const projectDirName = "binance-bot"

type AppConfig struct {
	IsDevelopment bool   `envconfig:"IS_DEVELOPMENT"`
	Debug         bool   `envconfig:"DEBUG"`
	Env           string `envconfig:"APP_ENV"`
	Host          string `envconfig:"API_HOST"`
	Port          string `envconfig:"API_PORT"`
	Firebase      Firebase
	Binance       Binance
}

var _app AppConfig
var _once sync.Once

func GetConfig() AppConfig {

	_once.Do(func() {

		_ = godotenv.Load()
		err := envconfig.Process("", &_app)
		if err != nil {
			panic(err)
		}
	})
	return _app
}

func LoadTestEnv() {
	rootDir, err := GetRootDir()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load(path.Join(rootDir, "./.env.testing"))
	if err != nil {
		panic(err)
	}

	GetConfig()
}

func IsTestEnv() bool {
	return _app.Env == "testing"
}

func IsProd() bool {
	if _app.Env == "testing" {
		return false
	}
	return !_app.IsDevelopment
}

func GetRequiredEnv(envName string) string {
	if envVar := os.Getenv(envName); envVar == "" {
		panic(fmt.Sprintf("Env variable '%s' not provided", envName))
	} else {
		return envVar
	}
}

func UpdateTestConfig(updateFunc func(config *AppConfig)) {
	if !IsTestEnv() {
		panic("Config updates allowed only in test env")
	}
	updateFunc(&_app)
}

func recoverPanic() {
	err := recover()
	if err != nil {
		panic(err)
	}
}

func GetRootDir() (string, error) {
	re := regexp.MustCompile("^(.*" + projectDirName + ")")
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	rootPath := re.Find([]byte(cwd))
	return string(rootPath), nil
}
