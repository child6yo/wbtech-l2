package repomock

import (
	"time"

	"l2.18/pkg/models"
)

// MockRepository - repository mock.
type MockRepository struct {
	GetEventsByDateRangeFn func(userID models.UserID, start, end time.Time) ([]models.Event, error)
}

// Put mock.
func (m *MockRepository) Put(userID models.UserID, event models.Event) error {
	panic("not implemented")
}

// Get mock.
func (m *MockRepository) Get(userID models.UserID, eventID models.EventID) (*models.Event, error) {
	panic("not implemented")
}

// Update mock.
func (m *MockRepository) Update(userID models.UserID, event models.Event) error {
	panic("not implemented")
}

// Delete mock.
func (m *MockRepository) Delete(userID models.UserID, eventID models.EventID) error {
	panic("not implemented")
}

// GetEventsByDateRange mock.
func (m *MockRepository) GetEventsByDateRange(userID models.UserID, start, end time.Time) ([]models.Event, error) {
	if m.GetEventsByDateRangeFn != nil {
		return m.GetEventsByDateRangeFn(userID, start, end)
	}
	return nil, nil
}
