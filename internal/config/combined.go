package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

var errLoadEnvFile = fmt.Errorf("ошибка загрузки .env файла")

type Config struct {
	DB  DBConfig
	TG  TGConfig
	Api ApiConfig
}

func Load(envPath string) (*Config, error) {
	if envPath != "" {
		err := LoadEnv(envPath)
		if err != nil {
			return nil, err
		}
	}
	tgConfig, err := LoadTG()
	if err != nil {
		return nil, err
	}
	dbConfig, err := LoadDB()
	if err != nil {
		return nil, err
	}
	apiConfig, err := LoadApi()
	if err != nil {
		return nil, err
	}
	return &Config{
		TG:  *tgConfig,
		DB:  *dbConfig,
		Api: *apiConfig,
	}, nil
}

func LoadEnv(envPath string) error {
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("%w: %v", errLoadEnvFile, err)
	}
	return nil
}
