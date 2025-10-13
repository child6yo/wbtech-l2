package events

import (
	"testing"
	"time"

	repomock "l2.18/internal/repository/mock"
	"l2.18/pkg/models"
)

func TestGetEventsForDay(t *testing.T) {
	now := time.Now().Truncate(24 * time.Hour) // начало текущего дня
	userID := models.UserID("user1")

	expectedStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	expectedEnd := expectedStart.AddDate(0, 0, 1)

	testEvents := []models.Event{
		{ID: "1", Date: now.Add(2 * time.Hour), Event: "event1"},
		{ID: "2", Date: now.Add(5 * time.Hour), Event: "event2"},
	}

	mockRepo := &repomock.MockRepository{
		GetEventsByDateRangeFn: func(userID models.UserID, start, end time.Time) ([]models.Event, error) {
			if userID != models.UserID("user1") {
				t.Errorf("expected userID %q, got %q", "user1", userID)
			}
			if !start.Equal(expectedStart) {
				t.Errorf("start: got %v, want %v", start, expectedStart)
			}
			if !end.Equal(expectedEnd) {
				t.Errorf("end: got %v, want %v", end, expectedEnd)
			}
			return testEvents, nil
		},
	}

	svc := &Service{repo: mockRepo}

	got, err := svc.GetEventsForDay(userID, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != len(testEvents) {
		t.Fatalf("expected %d events, got %d", len(testEvents), len(got))
	}
	for i, ev := range testEvents {
		if got[i].ID != ev.ID || got[i].Event != ev.Event {
			t.Errorf("event[%d] mismatch: got %+v, want %+v", i, got[i], ev)
		}
	}
}

func TestGetEventsForWeek(t *testing.T) {
	now := time.Now().Truncate(24 * time.Hour)
	userID := models.UserID("user1")

	expectedStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	expectedEnd := expectedStart.AddDate(0, 0, 7)

	testEvents := []models.Event{
		{ID: "3", Date: now.Add(10 * time.Hour), Event: "week_event1"},
		{ID: "4", Date: now.Add(3 * 24 * time.Hour), Event: "week_event2"},
	}

	mockRepo := &repomock.MockRepository{
		GetEventsByDateRangeFn: func(userID models.UserID, start, end time.Time) ([]models.Event, error) {
			if userID != models.UserID("user1") {
				t.Errorf("expected userID %q, got %q", "user1", userID)
			}
			if !start.Equal(expectedStart) {
				t.Errorf("start: got %v, want %v", start, expectedStart)
			}
			if !end.Equal(expectedEnd) {
				t.Errorf("end: got %v, want %v", end, expectedEnd)
			}
			return testEvents, nil
		},
	}

	svc := &Service{repo: mockRepo}

	got, err := svc.GetEventsForWeek(userID, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != len(testEvents) {
		t.Fatalf("expected %d events, got %d", len(testEvents), len(got))
	}
}

func TestGetEventsForMonth(t *testing.T) {
	// Выберем конкретную дату, чтобы избежать проблем с "31-е число"
	refDate := time.Date(2024, time.March, 15, 14, 30, 0, 0, time.UTC)
	userID := models.UserID("user1")

	expectedStart := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)

	testEvents := []models.Event{
		{ID: "5", Date: time.Date(2024, time.March, 5, 9, 0, 0, 0, time.UTC), Event: "march1"},
		{ID: "6", Date: time.Date(2024, time.March, 20, 18, 0, 0, 0, time.UTC), Event: "march2"},
	}

	mockRepo := &repomock.MockRepository{
		GetEventsByDateRangeFn: func(userID models.UserID, start, end time.Time) ([]models.Event, error) {
			if userID != models.UserID("user1") {
				t.Errorf("expected userID %q, got %q", "user1", userID)
			}
			if !start.Equal(expectedStart) {
				t.Errorf("start: got %v, want %v", start, expectedStart)
			}
			if !end.Equal(expectedEnd) {
				t.Errorf("end: got %v, want %v", end, expectedEnd)
			}
			return testEvents, nil
		},
	}

	svc := &Service{repo: mockRepo}

	got, err := svc.GetEventsForMonth(userID, refDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != len(testEvents) {
		t.Fatalf("expected %d events, got %d", len(testEvents), len(got))
	}
}
