// Package config содержит конфигурацию приложения
package config

import (
	"fmt"
	"net/url"
	"time"
)

// Api конфигурация API
type Api struct {
	Addr         string `envconfig:"API_ADDR" required:"true"`
	Port         string `envconfig:"API_PORT" required:"true"`
	LimiterRate  int    `envconfig:"LIMITER_RATE" default:"5"`
	LimiterBurst int    `envconfig:"LIMITER_BURST" default:"10"`
}

func (c *Api) URL() string {
	return fmt.Sprintf("%s:%s", c.Addr, c.Port)
}
func (c *Api) GetLimiterRate() int {
	return c.LimiterRate
}
func (c *Api) GetLimiterBurst() int {
	return c.LimiterBurst
}

// Database конфигурация базы данных

type Database struct {
	Host    string `envconfig:"DB_HOST" required:"true"`
	Port    string `envconfig:"DB_PORT" required:"true"`
	User    string `envconfig:"DB_USER" required:"true"`
	Pass    string `envconfig:"DB_PASS" required:"true"`
	DBName  string `envconfig:"DB_NAME" required:"true"`
	SSLMode string `envconfig:"DB_SSLMODE" required:"true"`
}

func (c *Database) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		url.QueryEscape(c.Pass),
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

type Telegram struct {
	Token              string `envconfig:"TG_TOKEN" required:"true"`
	ExpirationRow      int    `envconfig:"TG_EXPIRATION_HOURS" required:"true"`
	CleanupIntervalRow int    `envconfig:"TG_CLEANUP_INTERVAL_MINUTES" required:"true"`

	Expiration      time.Duration `ignored:"true"`
	CleanupInterval time.Duration `ignored:"true"`
}

func (c *Telegram) GetToken() string {
	return c.Token
}

func (c *Telegram) GetExpiration() time.Duration {
	return time.Duration(c.ExpirationRow) * time.Hour
}
func (c *Telegram) GetCleanupInterval() time.Duration {
	return time.Duration(c.CleanupIntervalRow) * time.Minute
}
