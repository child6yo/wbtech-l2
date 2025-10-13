package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"l2.18/pkg/models"
)

type eventsService interface {
	AddEvent(userID models.UserID, event models.Event) error
	UpdateEvent(userID models.UserID, event models.Event) error
	RemoveEvent(userID models.UserID, eventID models.EventID) error
	GetEventsForDay(userID models.UserID, day time.Time) ([]models.Event, error)
	GetEventsForWeek(userID models.UserID, weekStart time.Time) ([]models.Event, error)
	GetEventsForMonth(userID models.UserID, month time.Time) ([]models.Event, error)
}

// EventsHandler обрабатывает CRUD событий.
type EventsHandler struct {
	service eventsService
}

// NewEventsHandler создает новый EventsHandler.
func NewEventsHandler(service eventsService) *EventsHandler {
	return &EventsHandler{service: service}
}

type eventRequest struct {
	UserID models.UserID `json:"user_id"`
	Event  struct {
		ID    string `json:"id"`
		Date  string `json:"date"`
		Event string `json:"event"`
	}
}

type eventResponse struct {
	Result []models.Event `json:"result"`
}

// CreateEvent обрабатывает POST /create_event.
func (eh *EventsHandler) CreateEvent(w http.ResponseWriter, r *http.Request) error {
	data, err := io.ReadAll(r.Body)
	if err != nil || len(data) == 0 {
		return fmt.Errorf("%w: %v", errInvalidData, err)
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Println("body was not closed: ", err)
		}
	}()

	var req eventRequest

	err = json.Unmarshal(data, &req)
	if err != nil {
		return fmt.Errorf("%w: %v", errInvalidData, err)
	}

	parsedDate, err := time.Parse("2006-01-02", req.Event.Date)
	if err != nil {
		return fmt.Errorf("%w: %v", errInvalidData, err)
	}

	event := models.Event{Date: parsedDate, Event: req.Event.Event}

	err = eh.service.AddEvent(req.UserID, event)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// UpdateEvent обрабатывает POST /update_event.
//
// ! В тз описано POST, желательно заменить на PATCH.
func (eh *EventsHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) error {
	data, err := io.ReadAll(r.Body)
	if err != nil || len(data) == 0 {
		return fmt.Errorf("%w: %v", errInvalidData, err)
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Println("body was not closed: ", err)
		}
	}()

	var req eventRequest

	err = json.Unmarshal(data, &req)
	if err != nil {
		return fmt.Errorf("%w: %v", errInvalidData, err)
	}

	parsedDate, err := time.Parse("2006-01-02", req.Event.Date)
	if err != nil {
		return fmt.Errorf("%w: %v", errInvalidData, err)
	}

	event := models.Event{ID: models.EventID(req.Event.ID), Date: parsedDate, Event: req.Event.Event}

	err = eh.service.UpdateEvent(req.UserID, event)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// DeleteEvent обрабатывает POST /delete_event.
//
// ! В тз описано POST, желательно заменить на DELETE.
func (eh *EventsHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) error {
	userID := r.FormValue("user_id")
	if userID == "" {
		return errInvalidData
	}

	eventID := r.FormValue("id")
	if eventID == "" {
		return errInvalidData
	}

	err := eh.service.RemoveEvent(models.UserID(userID), models.EventID(eventID))
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// EventsForDay обрабатывает GET /events_for_day.
func (eh *EventsHandler) EventsForDay(w http.ResponseWriter, r *http.Request) error {
	userID := r.FormValue("user_id")
	if userID == "" {
		return errInvalidData
	}

	day := r.FormValue("date")
	if day == "" {
		return errInvalidData
	}

	t, err := time.Parse("2006-01-02", day)
	if err != nil {
		return fmt.Errorf("%w: %v", errInvalidData, err)
	}

	res, err := eh.service.GetEventsForDay(models.UserID(userID), t)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(eventResponse{Result: res}); err != nil {
		return err
	}

	return nil
}

// EventsForWeek обрабатывает GET /events_for_week.
func (eh *EventsHandler) EventsForWeek(w http.ResponseWriter, r *http.Request) error {
	userID := r.FormValue("user_id")
	if userID == "" {
		return errInvalidData
	}

	weekStart := r.FormValue("date")
	if weekStart == "" {
		return errInvalidData
	}

	t, err := time.Parse("2006-01-02", weekStart)
	if err != nil {
		return fmt.Errorf("%w: %v", errInvalidData, err)
	}

	res, err := eh.service.GetEventsForWeek(models.UserID(userID), t)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(eventResponse{Result: res}); err != nil {
		return err
	}

	return nil
}

// EventsForMonth обрабатывает GET /events_for_month.
func (eh *EventsHandler) EventsForMonth(w http.ResponseWriter, r *http.Request) error {
	userID := r.FormValue("user_id")
	if userID == "" {
		return errInvalidData
	}

	month := r.FormValue("date")
	if month == "" {
		return errInvalidData
	}

	t, err := time.Parse("2006-01-02", month)
	if err != nil {
		return fmt.Errorf("%w: %v", errInvalidData, err)
	}

	res, err := eh.service.GetEventsForMonth(models.UserID(userID), t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(eventResponse{Result: res}); err != nil {
		return err
	}

	return nil
}
