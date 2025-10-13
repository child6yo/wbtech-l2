package events

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"l2.18/internal/repository"
	"l2.18/internal/service"
	"l2.18/pkg/models"
)

type eventsRepository interface {
	Put(userID models.UserID, event models.Event) error
	Update(userID models.UserID, event models.Event) error
	Delete(userID models.UserID, eventID models.EventID) error
	GetEventsByDateRange(userID models.UserID, start, end time.Time) ([]models.Event, error)
}

// Service реализует сервис работы с событиями.
type Service struct {
	repo eventsRepository
}

// New создает новый Service.
func New(repo eventsRepository) *Service {
	return &Service{repo: repo}
}

// AddEvent добавляет новое событие.
func (s *Service) AddEvent(userID models.UserID, event models.Event) error {
	eventUID := uuid.NewString()
	event.ID = models.EventID(eventUID)

	err := s.repo.Put(userID, event)
	if errors.Is(err, repository.ErrAlreadyExist) {
		return service.ErrAlreadyExist
	} else if err != nil {
		return err
	}

	return nil
}

// UpdateEvent обновляет событие.
func (s *Service) UpdateEvent(userID models.UserID, event models.Event) error {
	err := s.repo.Update(userID, event)
	if errors.Is(err, repository.ErrNotFound) {
		return service.ErrNotFound
	} else if err != nil {
		return err
	}

	return nil
}

// RemoveEvent удаляет событие.
func (s *Service) RemoveEvent(userID models.UserID, eventID models.EventID) error {
	err := s.repo.Delete(userID, eventID)
	if errors.Is(err, repository.ErrNotFound) {
		return service.ErrNotFound
	} else if err != nil {
		return err
	}

	return nil
}

// GetEventsForDay возвращает все события пользователя на указанный день.
func (s *Service) GetEventsForDay(userID models.UserID, day time.Time) ([]models.Event, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.AddDate(0, 0, 1)

	return s.repo.GetEventsByDateRange(userID, start, end)
}

// GetEventsForWeek возвращает события на неделю (понедельник–воскресенье или просто 7 дней от даты).
func (s *Service) GetEventsForWeek(userID models.UserID, weekStart time.Time) ([]models.Event, error) {
	start := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
	end := start.AddDate(0, 0, 7)

	return s.repo.GetEventsByDateRange(userID, start, end)
}

// GetEventsForMonth возвращает события на месяц.
func (s *Service) GetEventsForMonth(userID models.UserID, month time.Time) ([]models.Event, error) {
	start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	end := start.AddDate(0, 1, 0)

	return s.repo.GetEventsByDateRange(userID, start, end)
}
