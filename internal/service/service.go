package service

import (
	"fmt"
	"nstu/internal/logger"
	"nstu/internal/model"
	"nstu/internal/repository"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

const (
	// Время хранения пользователя в кеше
	userCacheExpiration = 24 * time.Hour
	// Интервал очистки кеша
	userCacheCleanupInterval = time.Hour
)

type Service struct {
	token     string
	repo      repository.Repository
	userCache *gocache.Cache
}

func NewService(repo repository.Repository, token string) *Service {
	return &Service{
		token:     token,
		repo:      repo,
		userCache: gocache.New(userCacheExpiration, userCacheCleanupInterval),
	}
}

// TODO: не обновляющийся кеш
// EnsureUser проверяет наличие пользователя в кеше и создает его в БД при необходимости
func (s *Service) EnsureUser(userID int64, firstName, lastName, userName string) (*model.User, error) {
	// Проверяем наличие пользователя в кеше
	if cachedUser, found := s.userCache.Get(fmt.Sprintf("%d", userID)); found {
		return cachedUser.(*model.User), nil
	}

	user := &model.User{
		BaseModel: model.BaseModel{
			ID: userID,
		},
		FirstName:   firstName,
		LastName:    lastName,
		UserName:    userName,
		OwnedGroups: make([]int64, 0),
	}

	logger.Log.Info().
		Int64("user_id", userID).
		Str("first_name", firstName).
		Str("username", userName).
		Msg("Создаем пользователя в БД")

	if err := s.repo.CreateUserIfNotExists(user); err != nil {
		logger.Log.Error().
			Int64("user_id", userID).
			Err(err).
			Msg("Ошибка создания пользователя в БД")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Получаем актуальные данные из БД
	dbUser, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.Log.Error().
			Int64("user_id", userID).
			Err(err).
			Msg("Ошибка получения пользователя из БД")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Добавляем пользователя в кеш
	s.userCache.Set(fmt.Sprintf("%d", userID), dbUser, gocache.DefaultExpiration)

	return dbUser, nil
}

// AuthInitData проверяет валидность initData и создает пользователя в БД при необходимости
func (s *Service) AuthInitData(initData string) (*model.User, error) {

	user, err := s.repo.GetUserByID(878413772)
	if err != nil {
		return nil, err
	}

	// // Define how long since init data generation date init data is valid.
	// expIn := 1 * time.Hour

	// // Проверяем валидность initData
	// err := initdata.Validate(initData, s.token, expIn)
	// if err != nil {
	// 	return nil, err
	// }

	// // Парсим initData
	// parsedData, err := initdata.Parse(initData)
	// if err != nil {
	// 	return nil, err
	// }

	// // Создаем пользователя
	// user := &model.User{
	// 	BaseModel: model.BaseModel{
	// 		ID: parsedData.User.ID,
	// 	},
	// 	UserName:  parsedData.User.Username,
	// 	FirstName: parsedData.User.FirstName,
	// 	LastName:  parsedData.User.LastName,
	// }

	// // обновляем запись в БД, при необходимости
	// user, err = s.EnsureUser(parsedData.User.ID, parsedData.User.FirstName, parsedData.User.LastName, parsedData.User.Username)
	// if err != nil {
	// 	return nil, err
	// }

	// Возвращаем пользователя
	return user, nil
}

// GetSchedule возвращает расписание на конкретную дату
func (s *Service) GetSchedule(group *model.GroupDetailed, date time.Time) (*[]model.Subject, error) {
	// Получаем день недели (0 = воскресенье, 1 = понедельник, и т.д.)
	weekday := date.Weekday()

	// Если чередование отключено, берем расписание из нечетной недели
	if !group.AlternatingWeeks {
		return getWeekdaySchedule(&group.OddWeek, weekday), nil
	}

	// Для чередования нужна дата нечетного понедельника
	if group.OddMonday == nil {
		return nil, ErrNoOddMonday
	}

	// Вычисляем разницу в неделях между датой нечетного понедельника и запрошенной датой
	weeks := getWeeksBetween(*group.OddMonday, date)

	// Если разница четная - нечетная неделя, если нечетная - четная неделя
	if weeks%2 == 0 {
		return getWeekdaySchedule(&group.OddWeek, weekday), nil
	}

	return getWeekdaySchedule(&group.EvenWeek, weekday), nil
}

// getWeeksBetween возвращает количество недель между двумя датами
// Может быть отрицательным, если date раньше oddMonday
func getWeeksBetween(oddMonday, date time.Time) int {
	// Приводим обе даты к началу дня для корректного сравнения
	oddMonday = time.Date(oddMonday.Year(), oddMonday.Month(), oddMonday.Day(), 0, 0, 0, 0, oddMonday.Location())
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Разница в днях
	days := date.Sub(oddMonday).Hours() / 24

	// Округляем до недель
	return int(days / 7)
}

// getWeekdaySchedule возвращает расписание на конкретный день недели
// В русской системе: 0 = понедельник, 1 = вторник, и т.д.
func getWeekdaySchedule(week *model.Week, weekday time.Weekday) *[]model.Subject {
	// Преобразуем из системы Go (воскресенье = 0) в русскую систему (понедельник = 0)
	// (weekday + 6) % 7 даст нам нужное смещение
	russianWeekday := (weekday + 6) % 7

	switch russianWeekday {
	case 0: // Понедельник
		return &week.Monday
	case 1: // Вторник
		return &week.Tuesday
	case 2: // Среда
		return &week.Wednesday
	case 3: // Четверг
		return &week.Thursday
	case 4: // Пятница
		return &week.Friday
	case 5: // Суббота
		return &week.Saturday
	case 6: // Воскресенье
		return &week.Sunday
	default:
		return &[]model.Subject{}
	}
}
