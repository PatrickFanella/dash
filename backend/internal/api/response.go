package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/patrickfanella/dash/backend/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeServiceError(w http.ResponseWriter, err error) {
	var domErr *domain.Error
	if errors.As(err, &domErr) {
		switch domErr.Kind {
		case domain.ErrNotFound:
			writeError(w, http.StatusNotFound, domErr.Message)
		case domain.ErrValidation:
			writeError(w, http.StatusBadRequest, domErr.Message)
		case domain.ErrConflict:
			writeError(w, http.StatusConflict, domErr.Message)
		default:
			writeError(w, http.StatusInternalServerError, domErr.Message)
		}
		return
	}
	writeError(w, http.StatusInternalServerError, "internal error")
}
