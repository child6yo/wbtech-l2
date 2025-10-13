package handler

import "net/http"

// ErrHandlerFunc кастомная функция хендлера, для лучшей реализации миддлвеера.
// Характеризуется ошибкой в возврате.
type ErrHandlerFunc func(w http.ResponseWriter, r *http.Request) error
