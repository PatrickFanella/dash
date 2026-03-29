package identity

import (
	"context"
	"net/http"
)

type contextKey struct{}

// Middleware is a Chi middleware that parses Authelia identity headers and
// attaches the identity to the request context. It does NOT block requests
// without identity — Authelia handles authentication upstream.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := ParseHeaders(r)
		if id != nil {
			ctx := context.WithValue(r.Context(), contextKey{}, id)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// FromContext retrieves the Identity from the request context.
// Returns nil if no identity is present.
func FromContext(ctx context.Context) *Identity {
	id, _ := ctx.Value(contextKey{}).(*Identity)
	return id
}
