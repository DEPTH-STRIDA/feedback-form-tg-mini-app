package config

import (
	"fmt"
	"net/url"
	"nstu/internal/logger"
	"os"
)

const (
	errEnvNotFound = "переменная окружения %s не найдена"
)

var (
	errDBHost    = fmt.Errorf(errEnvNotFound, "DBHOST")
	errDBPort    = fmt.Errorf(errEnvNotFound, "DBPORT")
	errDBUser    = fmt.Errorf(errEnvNotFound, "DBUSER")
	errDBPass    = fmt.Errorf(errEnvNotFound, "DBPASS")
	errDBName    = fmt.Errorf(errEnvNotFound, "DBNAME")
	errDBSSLMode = fmt.Errorf(errEnvNotFound, "DBSSLMODE")
)

type DBConfig struct {
	Host    string
	Port    string
	User    string
	Pass    string
	DBName  string
	SSLMode string
}

func (c DBConfig) URL() string {
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

func LoadDB() (*DBConfig, error) {
	host, ok := os.LookupEnv("DBHOST")
	if !ok {
		return nil, errDBHost
	}
	logger.Log.Info().Str("host", host).Msg("Загружен хост БД")

	port, ok := os.LookupEnv("DBPORT")
	if !ok {
		return nil, errDBPort
	}
	logger.Log.Info().Str("port", port).Msg("Загружен порт БД")

	user, ok := os.LookupEnv("DBUSER")
	if !ok {
		return nil, errDBUser
	}
	logger.Log.Info().Str("user", user).Msg("Загружен пользователь БД")

	pass, ok := os.LookupEnv("DBPASS")
	if !ok {
		return nil, errDBPass
	}
	logger.Log.Info().Str("password", "**").Msg("Загружен пароль БД")

	dbName, ok := os.LookupEnv("DBNAME")
	if !ok {
		return nil, errDBName
	}
	logger.Log.Info().Str("database", dbName).Msg("Загружено имя БД")

	sslMode, ok := os.LookupEnv("DBSSLMODE")
	if !ok {
		return nil, errDBSSLMode
	}
	logger.Log.Info().Str("sslmode", sslMode).Msg("Загружен режим SSL")

	return &DBConfig{
		Host:    host,
		Port:    port,
		User:    user,
		Pass:    pass,
		DBName:  dbName,
		SSLMode: sslMode,
	}, nil
}
