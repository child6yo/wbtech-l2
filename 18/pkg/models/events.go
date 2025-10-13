package models

import "time"

// EventID определяет модель айди события
type EventID string

// Event определяет модель события.
type Event struct {
	ID    EventID   `json:"id"`
	Date  time.Time `json:"date"`
	Event string    `json:"event"`
}
