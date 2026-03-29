package identity

import (
	"net/http"
	"strings"
)

// Identity represents an authenticated user parsed from Authelia forwarded headers.
type Identity struct {
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Email       string   `json:"email"`
	Groups      []string `json:"groups"`
}

// ParseHeaders reads Authelia forwarded headers from the request and returns
// an Identity. Returns nil if no identity headers are present.
func ParseHeaders(r *http.Request) *Identity {
	if r == nil {
		return nil
	}

	username := r.Header.Get("Remote-User")
	displayName := r.Header.Get("Remote-Name")
	email := r.Header.Get("Remote-Email")
	rawGroups := r.Header.Get("Remote-Groups")

	if username == "" && displayName == "" && email == "" && rawGroups == "" {
		return nil
	}

	var groups []string
	if rawGroups != "" {
		for _, group := range strings.Split(rawGroups, ",") {
			group = strings.TrimSpace(group)
			if group != "" {
				groups = append(groups, group)
			}
		}
	}

	return &Identity{
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		Groups:      groups,
	}
}
