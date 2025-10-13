package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"l2.18/internal/service"
)

type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// Middleware выполняет промежуточные операции (логирование, обработка ошибок).
type Middleware struct {
	log logger
}

// NewMiddleware создает новый Middleware.
func NewMiddleware(log logger) *Middleware {
	return &Middleware{log: log}
}

// Logging выполняет логирование запроса и обработку ошибок.
func (m *Middleware) Logging(h ErrHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log.Info("new request", "method", r.Method, "path", r.URL.Path)

		var statusCode int

		if err := h(w, r); err != nil {
			switch {
			case errors.Is(err, service.ErrAlreadyExist):
				statusCode = http.StatusServiceUnavailable
			case errors.Is(err, service.ErrNotFound):
				statusCode = http.StatusServiceUnavailable
			case errors.Is(err, errInvalidData):
				statusCode = http.StatusBadRequest
			default:
				statusCode = http.StatusInternalServerError
			}

			m.log.Error("request failed",
				"method", r.Method,
				"path", r.URL.Path,
				"error", err.Error(),
				"status", statusCode)

			w.WriteHeader(statusCode)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		}
	})
}
