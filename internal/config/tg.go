package config

import (
	"fmt"
	"nstu/internal/logger"
	"os"
	"strconv"
	"time"
)

var (
	errTGToken           = fmt.Errorf("TG_TOKEN не найден")
	errTGExpiration      = fmt.Errorf("TG_EXPIRATION_HOURS не найден")
	errTGCleanupInterval = fmt.Errorf("TG_CLEANUP_INTERVAL_MINUTES не найден")
	errTGExpirationParse = fmt.Errorf("не удалось преобразовать TG_EXPIRATION_HOURS в число")
	errTGCleanupParse    = fmt.Errorf("не удалось преобразовать TG_CLEANUP_INTERVAL_MINUTES в число")
)

type TGConfig struct {
	Token           string
	Expiration      time.Duration
	CleanupInterval time.Duration
}

func LoadTG() (*TGConfig, error) {
	token, ok := os.LookupEnv("TG_TOKEN")
	if !ok {
		return nil, errTGToken
	}
	// Маскируем токен для логов
	maskedToken := "**:**"
	if len(token) > 10 {
		maskedToken = token[:6] + "..." + token[len(token)-4:]
	}
	logger.Log.Info().Str("token", maskedToken).Msg("Загружен токен Telegram")

	expirationStr, ok := os.LookupEnv("TG_EXPIRATION_HOURS")
	if !ok {
		return nil, errTGExpiration
	}
	logger.Log.Info().Str("expiration_hours", expirationStr).Msg("Загружено время хранения состояний")

	expirationHours, err := strconv.ParseInt(expirationStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errTGExpirationParse, err)
	}
	expiration := time.Duration(expirationHours) * time.Hour

	cleanupStr, ok := os.LookupEnv("TG_CLEANUP_INTERVAL_MINUTES")
	if !ok {
		return nil, errTGCleanupInterval
	}
	logger.Log.Info().Str("cleanup_interval", cleanupStr).Msg("Загружен интервал очистки")

	cleanupMinutes, err := strconv.ParseInt(cleanupStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errTGCleanupParse, err)
	}
	cleanupInterval := time.Duration(cleanupMinutes) * time.Minute

	return &TGConfig{
		Token:           token,
		Expiration:      expiration,
		CleanupInterval: cleanupInterval,
	}, nil
}
