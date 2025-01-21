package service

import (
	"nstu/internal/model"
	"time"
)

// HandleSchedule обрабатывает запрос на получение расписания
func (s *Service) HandleSchedule(userID int64, date string) (*[]model.Subject, error) {
	// Получаем пользователя
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Проверяем, что пользователь состоит в группе
	if user.MemberOf == nil {
		return nil, ErrNoGroup
	}

	// Проверяем, что id группы не равен 0
	if *user.MemberOf <= 0 {
		return nil, ErrInvalidGroupID
	}

	// Парсим дату. Дата может быть пустой, тогда используем текущую дату
	var dateTime time.Time
	if date == "" {
		dateTime = time.Now()
	} else {
		dateTime, err = time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}
	}

	// Получаем полную информацию о группе
	group, err := s.repo.GetGroupByID(*user.MemberOf)
	if err != nil {
		return nil, err
	}

	// Получаем расписание на дату
	subjects, err := s.GetSchedule(group, dateTime)
	if err != nil {
		return nil, err
	}

	return subjects, nil
}

// HandleLeave обрабатывает запрос на выход из группы
func (s *Service) HandleLeave(userID int64) error {
	// Получаем пользователя
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Проверяем, что пользователь состоит в группе
	if user.MemberOf == nil {
		return nil
	}

	// Проверяем, что id группы не равен 0
	if *user.MemberOf <= 0 {
		return nil
	}

	// Удаляем пользователя из группы
	user.MemberOf = nil

	return s.repo.UpdateUser(user)
}

// HandleGetGroups обрабатывает запрос на получение списка групп
func (s *Service) HandleGetGroups(userID int64, search string, owned bool, offset, limit int) ([]model.Group, error) {
	// Ограничиваем длину поиска до 256 символов
	if len(search) > 256 {
		search = search[:256]
	}

	// Проверяем offset
	if offset < 0 {
		offset = 0
	}

	// Проверяем limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	return s.repo.ListGroups(search, owned, offset, limit, userID)
}

// HandleCreateGroup обрабатывает запрос на создание группы
func (s *Service) HandleCreateGroup(userID int64, group model.GroupDetailed) error {
	// Получаем пользователя
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if len(user.OwnedGroups) >= 100 {
		return ErrTooManyGroups
	}

	return s.repo.CreateGroup(&group)
}

// HandleGetGroup обрабатывает запрос на получение информации о группе
func (s *Service) HandleGetGroup(groupID int64) (*model.GroupDetailed, error) {
	return s.repo.GetGroupByID(groupID)
}

// HandleUpdateGroup обрабатывает запрос на обновление информации о группе
func (s *Service) HandleUpdateGroup(userID int64, groupID int64, group model.GroupDetailed) error {
	// Получаем пользователя
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	for _, id := range user.OwnedGroups {
		if id == groupID {
			return s.repo.UpdateGroup(&group)
		}
	}

	return ErrNoOwnedGroup
}

// HandleDeleteGroup обрабатывает запрос на удаление группы
func (s *Service) HandleDeleteGroup(userID int64, groupID int64) error {
	// Получаем пользователя
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	for _, id := range user.OwnedGroups {
		if id == groupID {
			return s.repo.DeleteGroup(groupID)
		}
	}

	return ErrNoOwnedGroup
}

// HandleJoinGroup обрабатывает запрос на вступление в группу
func (s *Service) HandleJoinGroup(userID int64, groupID int64) error {
	// Получаем пользователя
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.MemberOf = &groupID

	return s.repo.UpdateUser(user)
}
