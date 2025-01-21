package main

import (
	"nstu/internal/api/handler"
	"nstu/internal/api/router"
	"nstu/internal/api/server"
	"nstu/internal/logger"
	"nstu/internal/repository/postgres"
	"nstu/internal/service"
	"nstu/internal/tg"
	"os"
	"time"

	cnfModel "nstu/internal/config"
	cnfLoad "nstu/pkg/config"
)

func init() {
	// Устанавливаем московское время как локальное
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Ошибка установки московского времени")
	}
	time.Local = loc
}

func main() {
	// Получение пути к .env файлу из аргументов
	var envPath string
	if len(os.Args) > 1 {
		envPath = os.Args[1]
		logger.Log.Info().Str("env_path", envPath).Msg("Загружаем конфигурацию из .env файла")

	}

	dbConf := &cnfModel.Database{}
	tgConf := &cnfModel.Telegram{}
	apiConf := &cnfModel.Api{}

	// Загрузка конфигурации с путем к .env
	err := cnfLoad.Load(envPath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Ошибка загрузки конфигурации")
	}

	// Подключение к БД
	db, err := postgres.NewPostgresDB(dbConf)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Ошибка подключения к базе данных")
	}
	defer db.Close()

	// Иницилизация структуры для работы с БД
	repo := postgres.NewRepository(db)

	// Иницилизация структуры бизнес логики
	srv := service.NewService(repo)

	// Иницилизация бота
	tg.InitBot(tgConf)

	handler := handler.NewHandler(srv)

	router := router.NewRouter(handler, apiConf.LimiterRate, apiConf.LimiterBurst)

	server := server.NewServer(router, apiConf.URL())

	// Запускаем сервер и ждем сигнала для завершения
	if err := server.RunUntilSignal(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Ошибка работы сервера")
	}
}

func generateTestData(repo *postgres.PostgresRepository) error {
	// Проверяем текущее количество пользователей
	users, err := repo.UserRepo.GetUserByID(1)
	if err == nil && users != nil {
		logger.Log.Info().Msg("Тестовые данные уже существуют")
		return nil
	}

	return nil
}
