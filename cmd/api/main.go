package main

import (
	"fmt"
	"math/rand"
	"nstu/internal/api/handler"
	"nstu/internal/api/router"
	"nstu/internal/api/server"
	"nstu/internal/config"
	"nstu/internal/logger"
	"nstu/internal/model"
	"nstu/internal/repository/postgres"
	"nstu/internal/service"
	"nstu/internal/tg"
	"os"
	"time"
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

	// Загрузка конфигурации с путем к .env
	cfg, err := config.Load(envPath)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Ошибка загрузки конфигурации")
	}

	// Подключение к БД
	db, err := postgres.NewPostgresDB(cfg.DB)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Ошибка подключения к базе данных")
	}
	defer db.Close()

	// Иницилизация структуры для работы с БД
	repo := postgres.NewRepository(db)

	// Иницилизация структуры бизнес логики
	service := service.NewService(repo, cfg.TG.Token)

	// Иницилизация бота
	tg.InitBot(&cfg.TG, service)

	handler := handler.NewHandler(repo, service)

	router := router.NewRouter(handler, service)

	server := server.NewServer(router, cfg.Api.Addr, cfg.Api.Port, service)

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

	// Генерируем 50 пользователей
	logger.Log.Info().Msg("Генерация пользователей...")
	for i := int64(1); i <= 50; i++ {
		user := &model.User{
			BaseModel: model.BaseModel{
				ID: i,
			},
			FirstName:   fmt.Sprintf("Имя%d", i),
			LastName:    fmt.Sprintf("Фамилия%d", i),
			UserName:    fmt.Sprintf("user%d", i),
			MemberOf:    nil,
			OwnedGroups: []int64{},
		}
		if err := repo.UserRepo.CreateUser(user); err != nil {
			return fmt.Errorf("ошибка создания пользователя: %w", err)
		}
	}

	// Генерируем 250 групп
	logger.Log.Info().Msg("Генерация групп...")
	for i := 1; i <= 250; i++ {
		// Случайный владелец из 50 пользователей
		ownerID := rand.Int63n(50) + 1

		group := &model.GroupDetailed{
			Group: model.Group{
				OwnerID:          ownerID,
				Name:             fmt.Sprintf("ПМИ-%d", i),
				Title:            fmt.Sprintf("Прикладная математика и информатика %d", i),
				Participants:     []int64{},
				AlternatingWeeks: rand.Int31n(2) == 1,
			},
			OddWeek: model.Week{
				Monday: []model.Subject{
					{Name: "Математика"},
					{Name: "Физика"},
				},
				Tuesday: []model.Subject{
					{Name: "Информатика"},
					{Name: "Английский"},
				},
				Wednesday: []model.Subject{
					{Name: "Программирование"},
					{Name: "База данных"},
				},
				Thursday: []model.Subject{
					{Name: "Алгоритмы"},
					{Name: "Сети"},
				},
				Friday: []model.Subject{
					{Name: "Веб-разработка"},
					{Name: "Защита информации"},
				},
				Saturday: []model.Subject{
					{Name: "Практика"},
				},
				Sunday: []model.Subject{},
			},
			EvenWeek: model.Week{
				Monday: []model.Subject{
					{Name: "Физика"},
					{Name: "Математика"},
				},
				Tuesday: []model.Subject{
					{Name: "Английский"},
					{Name: "Информатика"},
				},
				Wednesday: []model.Subject{
					{Name: "База данных"},
					{Name: "Программирование"},
				},
				Thursday: []model.Subject{
					{Name: "Сети"},
					{Name: "Алгоритмы"},
				},
				Friday: []model.Subject{
					{Name: "Защита информации"},
					{Name: "Веб-разработка"},
				},
				Saturday: []model.Subject{
					{Name: "Практика"},
				},
				Sunday: []model.Subject{},
			},
		}

		if err := repo.GroupRepo.CreateGroup(group); err != nil {
			return fmt.Errorf("ошибка создания группы: %w", err)
		}

		// Добавляем группу в owned_groups владельца
		owner, err := repo.UserRepo.GetUserByID(ownerID)
		if err != nil {
			return fmt.Errorf("ошибка получения владельца: %w", err)
		}
		owner.OwnedGroups = append(owner.OwnedGroups, int64(i))
		if err := repo.UserRepo.UpdateUser(owner); err != nil {
			return fmt.Errorf("ошибка обновления владельца: %w", err)
		}
	}

	logger.Log.Info().Msg("Тестовые данные успешно сгенерированы")
	return nil
}
