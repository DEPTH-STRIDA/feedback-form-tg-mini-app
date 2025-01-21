package config

import (
	"fmt"
	"os"
)

const (
	errApiAddr = "переменная окружения API_ADDR не найдена"
	errApiPort = "переменная окружения API_PORT не найдена"
)

type ApiConfig struct {
	Addr string
	Port string
}

func LoadApi() (*ApiConfig, error) {
	addr, ok := os.LookupEnv("API_ADDR")
	if !ok {
		return nil, fmt.Errorf(errApiAddr)
	}

	port, ok := os.LookupEnv("API_PORT")
	if !ok {
		return nil, fmt.Errorf(errApiPort)
	}

	return &ApiConfig{
		Addr: addr,
		Port: port,
	}, nil
}
