package api

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickfanella/dash/backend/internal/importer"
)

func handleImport(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, err := importer.Parse(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid YAML: "+err.Error())
			return
		}

		result, err := importer.Run(r.Context(), pool, cfg)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "import failed: "+err.Error())
			return
		}

		writeJSON(w, http.StatusOK, result)
	}
}
