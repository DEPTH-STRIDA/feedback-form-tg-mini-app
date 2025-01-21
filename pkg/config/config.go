// Package config предоставляет утилиты для загрузки конфигурации из переменных окружения.
// Поддерживает загрузку из .env файла и переменных окружения системы.
// Использует теги envconfig для маппинга переменных окружения на поля структур.
//
// Пример использования:
//
//	type Config struct {
//	    Host string `envconfig:"HOST" required:"true"`
//	    Port int    `envconfig:"PORT" default:"8080"`
//	}
//
//	var cfg Config
//	err := config.Load(".env", &cfg)
package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// LoadConfig загружает конфигурацию из переменных окружения в указанные структуры
// envPath может быть пустым - тогда .env файл не будет загружен
// configs - список указателей на структуры с тегами env
func Load(envPath string, configs ...interface{}) error {
	// Загружаем .env файл только если путь указан
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			return fmt.Errorf("ошибка загрузки .env файла: %v", err)
		}
	}

	// Загружаем конфигурацию в каждую структуру
	for _, cfg := range configs {
		if err := envconfig.Process("", cfg); err != nil {
			return fmt.Errorf("ошибка загрузки конфигурации: %v", err)
		}
	}

	return nil
}
