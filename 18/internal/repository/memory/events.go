package memory

import (
	"sort"
	"sync"
	"time"

	"l2.18/internal/repository"
	"l2.18/pkg/models"
)

// EventsRepository хранит в оперативной памяти информацию о событиях.
type EventsRepository struct {
	sync.RWMutex

	events    map[models.UserID]map[models.EventID]*models.Event
	dateIndex map[models.UserID][]*models.Event
}

// NewEventsRepository создает новый EventsRepository.
func NewEventsRepository() *EventsRepository {
	return &EventsRepository{
		events:     make(map[models.UserID]map[models.EventID]*models.Event),
		dateIndex:  make(map[models.UserID][]*models.Event),
	}
}

// Put добавляет новое событие. Если событие уже существует - вернет ошибку.
func (er *EventsRepository) Put(userID models.UserID, event models.Event) error {
	er.Lock()
	defer er.Unlock()

	if er.events[userID] == nil {
		er.events[userID] = make(map[models.EventID]*models.Event)
		er.dateIndex[userID] = make([]*models.Event, 0)
	}

	if _, exists := er.events[userID][event.ID]; exists {
		return repository.ErrAlreadyExist
	}


	er.events[userID][event.ID] = &event
	er.dateIndex[userID] = insertSorted(er.dateIndex[userID], &event)
	return nil
}

// Get возвращает событие пользователя по его айди. Если события нет - вернет ошибку.
func (er *EventsRepository) Get(userID models.UserID, eventID models.EventID) (*models.Event, error) {
	er.RLock()
	defer er.RUnlock()

	if _, ok := er.events[userID][eventID]; !ok {
		return nil, repository.ErrNotFound
	}

	return er.events[userID][eventID], nil
}

// Update обновляет событие пользователя, заменяя существующие поля,
// полями переданными в функцию в event.
func (er *EventsRepository) Update(userID models.UserID, event models.Event) error {
	er.Lock()
	defer er.Unlock()

	eventPtr, exists := er.events[userID][event.ID]
	if !exists {
		return repository.ErrNotFound
	}

	oldDate := eventPtr.Date
	if !event.Date.IsZero() {
		eventPtr.Date = event.Date
	}
	if event.Event != "" {
		eventPtr.Event = event.Event
	}

	if !event.Date.IsZero() && !event.Date.Equal(oldDate) {
		er.dateIndex[userID] = nil
		for _, e := range er.events[userID] {
			er.dateIndex[userID] = insertSorted(er.dateIndex[userID], e)
		}
	}
	return nil
}

// Delete удаляет события пользователя по айди.
func (er *EventsRepository) Delete(userID models.UserID, eventID models.EventID) error {
	er.Lock()
	defer er.Unlock()

	eventPtr, exists := er.events[userID][eventID]
	if !exists {
		return repository.ErrNotFound
	}

	delete(er.events[userID], eventID)

	idx := -1
	for i, e := range er.dateIndex[userID] {
		if e == eventPtr {
			idx = i
			break
		}
	}
	if idx != -1 {
		er.dateIndex[userID] = append(er.dateIndex[userID][:idx], er.dateIndex[userID][idx+1:]...)
	}

	return nil
}

// GetEventsByDateRange возвращает все события пользователя в диапазоне [start, end).
func (er *EventsRepository) GetEventsByDateRange(
	userID models.UserID,
	start, end time.Time,
) ([]models.Event, error) {
	er.RLock()
	defer er.RUnlock()

	index, ok := er.dateIndex[userID]
	if !ok || len(index) == 0 {
		return []models.Event{}, nil
	}

	left := sort.Search(len(index), func(i int) bool {
		return !index[i].Date.Before(start)
	})

	right := sort.Search(len(index), func(i int) bool {
		return !index[i].Date.Before(end)
	})

	if left >= right {
		return []models.Event{}, nil
	}

	result := make([]models.Event, right-left)
	for i, j := left, 0; i < right; i, j = i+1, j+1 {
		result[j] = *index[i]
	}

	return result, nil
}

func insertSorted(slice []*models.Event, event *models.Event) []*models.Event {
	i := sort.Search(len(slice), func(i int) bool {
		return !slice[i].Date.Before(event.Date)
	})
	slice = append(slice, nil)
	copy(slice[i+1:], slice[i:])
	slice[i] = event
	return slice
}
