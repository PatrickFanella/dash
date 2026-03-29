package identity

import (
	"net/http"
	"strings"
)

type Identity struct {
	Username    string
	DisplayName string
	Email       string
	Groups      []string
}

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
