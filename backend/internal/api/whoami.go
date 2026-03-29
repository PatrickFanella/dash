package api

import (
	"net/http"

	"github.com/patrickfanella/dash/backend/internal/identity"
)

func handleWhoami(w http.ResponseWriter, r *http.Request) {
	id := identity.FromContext(r.Context())
	if id == nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	writeJSON(w, http.StatusOK, id)
}
