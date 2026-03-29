package api

import (
	"encoding/json"
	"net/http"
)

// bindJSON decodes the request body into v. Returns false (and writes a 400
// error response) if decoding fails, so the caller can return early.
func bindJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}
