package api

import (
	"io/fs"
	"net/http"
	"strings"
)

func FileServer(fsys fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(fsys))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// API routes are handled by the router, not the file server.
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Try to serve the static file. If it doesn't exist, serve index.html
		// so React Router can handle client-side routes.
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else {
			path = strings.TrimPrefix(path, "/")
		}

		if _, err := fs.Stat(fsys, path); err != nil {
			// File doesn't exist — serve index.html for SPA routing.
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(w, r)
	})
}
