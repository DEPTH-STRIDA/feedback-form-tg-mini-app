package main

import (
	"database/sql"
	"embed"
	cnfModel "nstu/internal/config"
	"nstu/internal/logger"
	"nstu/internal/migrator"
	cnfLoad "nstu/pkg/config"
	"os"
)

const migrationsDir = "migrations"

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	var envPath string
	if len(os.Args) > 1 {
		envPath = os.Args[1]
		logger.Log.Info().Str("env_path", envPath).Msg("Загружаем конфигурацию из .env файла")
	}
	dbCfg := &cnfModel.Database{}
	// Загрузка конфигурации для БД
	err := cnfLoad.Load(envPath, dbCfg)
	if err != nil {
		logger.Log.Err(err).Msg("Ошибка загрузки конфигурации")
	}

	migrator := migrator.MustGetNewMigrator(migrationsFS, migrationsDir)

	conn, err := sql.Open("postgres", dbCfg.URL())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = migrator.ApplyMigrations(conn)
	if err != nil {
		panic(err)
	}
	logger.Log.Info().Msg("Миграции применены")
}
