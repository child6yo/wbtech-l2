package memory

import (
	"errors"
	"testing"
	"time"

	"l2.18/internal/repository"
	"l2.18/pkg/models"
)

func TestPut(t *testing.T) {
	now := time.Now()
	eventID := models.EventID("1")
	userID := models.UserID("user1")
	event := models.Event{
		ID:    eventID,
		Date:  now,
		Event: "test_event",
	}

	testCases := []struct {
		name     string
		setup    func(*EventsRepository)
		expected error
	}{
		{
			name: "success - new event",
			setup: func(r *EventsRepository) {
			},
			expected: nil,
		},
		{
			name: "failure - already exists",
			setup: func(r *EventsRepository) {
				_ = r.Put(userID, event)
			},
			expected: repository.ErrAlreadyExist,
		},
		{
			name: "success - different user",
			setup: func(r *EventsRepository) {
				_ = r.Put(models.UserID("user2"), event)
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewEventsRepository()
			tc.setup(repo)
			err := repo.Put(userID, event)
			if !errors.Is(err, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	now := time.Now()
	eventID := models.EventID("1")
	userID := models.UserID("user1")
	event := models.Event{
		ID:    eventID,
		Date:  now,
		Event: "test_event",
	}

	testCases := []struct {
		name     string
		setup    func(*EventsRepository)
		expected error
	}{
		{
			name: "success - event exists",
			setup: func(r *EventsRepository) {
				_ = r.Put(userID, event)
			},
			expected: nil,
		},
		{
			name: "failure - not found",
			setup: func(r *EventsRepository) {
			},
			expected: repository.ErrNotFound,
		},
		{
			name: "success - different user",
			setup: func(r *EventsRepository) {
				_ = r.Put(models.UserID("user2"), event)
			},
			expected: repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewEventsRepository()
			tc.setup(repo)
			_, err := repo.Get(userID, eventID)
			if !errors.Is(err, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	eventID := models.EventID("1")
	userID := models.UserID("user1")
	event := models.Event{
		ID:    eventID,
		Date:  now,
		Event: "test_event",
	}

	testCases := []struct {
		name     string
		setup    func(*EventsRepository)
		input    models.Event
		expected error
	}{
		{
			name: "success - only time",
			setup: func(r *EventsRepository) {
				_ = r.Put(userID, event)
			},
			input: models.Event{
				ID:   eventID,
				Date: time.Now().Add(3 * time.Hour),
			},
			expected: nil,
		},
		{
			name: "success - only description",
			setup: func(r *EventsRepository) {
				_ = r.Put(userID, event)
			},
			input: models.Event{
				ID:    eventID,
				Event: "updated",
			},
			expected: nil,
		},
		{
			name: "success - date&description",
			setup: func(r *EventsRepository) {
				_ = r.Put(userID, event)
			},
			input: models.Event{
				ID:    eventID,
				Date:  time.Now().Add(3 * time.Hour),
				Event: "updated",
			},
			expected: nil,
		},
		{
			name: "failure - not found",
			setup: func(r *EventsRepository) {
				_ = r.Put(userID, event)
			},
			input:    models.Event{},
			expected: repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewEventsRepository()
			tc.setup(repo)
			err := repo.Update(userID, tc.input)
			if !errors.Is(err, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}

			got, _ := repo.Get(userID, eventID)

			expectedEvent := event
			if !tc.input.Date.IsZero() {
				expectedEvent.Date = tc.input.Date
			}
			if tc.input.Event != "" {
				expectedEvent.Event = tc.input.Event
			}

			if got.ID != expectedEvent.ID {
				t.Errorf("ID mismatch: got %v, want %v", got.ID, expectedEvent.ID)
			}
			if !got.Date.Equal(expectedEvent.Date) {
				t.Errorf("Date mismatch: got %v, want %v", got.Date, expectedEvent.Date)
			}
			if got.Event != expectedEvent.Event {
				t.Errorf("Event mismatch: got %q, want %q", got.Event, expectedEvent.Event)
			}

		})
	}
}

func TestDelete(t *testing.T) {
	now := time.Now()
	eventID := models.EventID("1")
	userID := models.UserID("user1")
	event := models.Event{
		ID:    eventID,
		Date:  now,
		Event: "test_event",
	}

	testCases := []struct {
		name     string
		setup    func(*EventsRepository)
		expected error
	}{
		{
			name: "success",
			setup: func(r *EventsRepository) {
				_ = r.Put(userID, event)
			},
			expected: nil,
		},
		{
			name: "failure",
			setup: func(r *EventsRepository) {
			},
			expected: repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewEventsRepository()
			tc.setup(repo)
			err := repo.Delete(userID, eventID)
			if !errors.Is(err, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, err)
			}
		})
	}
}

func TestGetEventsByDateRange(t *testing.T) {
	now := time.Now().Truncate(24 * time.Hour)
	userID := models.UserID("user1")
	otherUserID := models.UserID("user2")

	event1 := models.Event{
		ID:    models.EventID("1"),
		Date:  now.Add(2 * time.Hour),
		Event: "event1",
	}
	event2 := models.Event{
		ID:    models.EventID("2"),
		Date:  now.Add(5 * time.Hour),
		Event: "event2",
	}
	event3 := models.Event{
		ID:    models.EventID("3"),
		Date:  now.Add(24 * time.Hour).Add(3 * time.Hour),
		Event: "event3",
	}
	event4 := models.Event{
		ID:    models.EventID("4"),
		Date:  now.Add(-24 * time.Hour).Add(10 * time.Hour),
		Event: "event4",
	}

	testCases := []struct {
		name     string
		setup    func(*EventsRepository)
		userID   models.UserID
		start    time.Time
		end      time.Time
		expected []models.Event
	}{
		{
			name: "success - events in range today",
			setup: func(r *EventsRepository) {
				_ = r.Put(userID, event1)
				_ = r.Put(userID, event2)
				_ = r.Put(userID, event3)
				_ = r.Put(userID, event4)
			},
			userID:   userID,
			start:    now,
			end:      now.Add(24 * time.Hour),
			expected: []models.Event{event1, event2},
		},
		{
			name: "success - no events in range",
			setup: func(r *EventsRepository) {
				_ = r.Put(userID, event1)
				_ = r.Put(userID, event2)
			},
			userID:   userID,
			start:    now.Add(48 * time.Hour),
			end:      now.Add(72 * time.Hour),
			expected: []models.Event{},
		},
		{
			name: "success - user has no events",
			setup: func(r *EventsRepository) {
			},
			userID:   models.UserID("nonexistent"),
			start:    now,
			end:      now.Add(24 * time.Hour),
			expected: []models.Event{},
		},
		{
			name: "success - other user's events not included",
			setup: func(r *EventsRepository) {
				_ = r.Put(otherUserID, event1)
				_ = r.Put(otherUserID, event2)
			},
			userID:   userID,
			start:    now,
			end:      now.Add(24 * time.Hour),
			expected: []models.Event{},
		},
		{
			name: "success - boundary: start inclusive, end exclusive",
			setup: func(r *EventsRepository) {
				eventAtStart := models.Event{
					ID:    models.EventID("5"),
					Date:  now,
					Event: "at_start",
				}
				eventAtEnd := models.Event{
					ID:    models.EventID("6"),
					Date:  now.Add(24 * time.Hour),
					Event: "at_end",
				}
				_ = r.Put(userID, eventAtStart)
				_ = r.Put(userID, eventAtEnd)
			},
			userID: userID,
			start:  now,
			end:    now.Add(24 * time.Hour),
			expected: []models.Event{
				{ID: models.EventID("5"), Date: now, Event: "at_start"},
			},
		},
		{
			name: "success - multiple events, sorted by date",
			setup: func(r *EventsRepository) {
				_ = r.Put(userID, event2)
				_ = r.Put(userID, event1)
			},
			userID:   userID,
			start:    now,
			end:      now.Add(24 * time.Hour),
			expected: []models.Event{event1, event2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewEventsRepository()
			tc.setup(repo)

			got, err := repo.GetEventsByDateRange(tc.userID, tc.start, tc.end)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tc.expected) {
				t.Fatalf("expected %d events, got %d", len(tc.expected), len(got))
			}

			for i, expectedEvent := range tc.expected {
				if got[i].ID != expectedEvent.ID {
					t.Errorf("event[%d].ID: got %v, want %v", i, got[i].ID, expectedEvent.ID)
				}
				if !got[i].Date.Equal(expectedEvent.Date) {
					t.Errorf("event[%d].Date: got %v, want %v", i, got[i].Date, expectedEvent.Date)
				}
				if got[i].Event != expectedEvent.Event {
					t.Errorf("event[%d].Event: got %q, want %q", i, got[i].Event, expectedEvent.Event)
				}
			}
		})
	}
}
