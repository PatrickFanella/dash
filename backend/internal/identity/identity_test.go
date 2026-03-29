package identity

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParseHeadersAllFields(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Remote-User", "patrick")
	r.Header.Set("Remote-Name", "Patrick Fanella")
	r.Header.Set("Remote-Email", "patrick@example.com")
	r.Header.Set("Remote-Groups", "admins, users")

	id := ParseHeaders(r)
	if id == nil {
		t.Fatal("expected identity, got nil")
	}
	if id.Username != "patrick" {
		t.Errorf("expected username patrick, got %s", id.Username)
	}
	if id.DisplayName != "Patrick Fanella" {
		t.Errorf("expected display_name Patrick Fanella, got %s", id.DisplayName)
	}
	if id.Email != "patrick@example.com" {
		t.Errorf("expected email patrick@example.com, got %s", id.Email)
	}

	wantGroups := []string{"admins", "users"}
	if !reflect.DeepEqual(id.Groups, wantGroups) {
		t.Errorf("expected groups %v, got %v", wantGroups, id.Groups)
	}
}

func TestParseHeadersNoHeaders(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	id := ParseHeaders(r)
	if id != nil {
		t.Errorf("expected nil for no headers, got %+v", id)
	}
}

func TestParseHeadersNilRequest(t *testing.T) {
	id := ParseHeaders(nil)
	if id != nil {
		t.Errorf("expected nil for nil request, got %+v", id)
	}
}

func TestParseHeadersPartialHeaders(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Remote-Email", "bob@example.com")

	id := ParseHeaders(r)
	if id == nil {
		t.Fatal("expected identity, got nil")
	}
	if id.Username != "" {
		t.Errorf("expected empty username, got %s", id.Username)
	}
	if id.DisplayName != "" {
		t.Errorf("expected empty display_name, got %s", id.DisplayName)
	}
	if id.Email != "bob@example.com" {
		t.Errorf("expected email bob@example.com, got %s", id.Email)
	}
	if id.Groups != nil {
		t.Errorf("expected nil groups, got %v", id.Groups)
	}
}

func TestParseHeadersSingleGroup(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Remote-User", "patrick")
	r.Header.Set("Remote-Groups", "admins")

	id := ParseHeaders(r)
	if id == nil {
		t.Fatal("expected identity, got nil")
	}
	if len(id.Groups) != 1 || id.Groups[0] != "admins" {
		t.Errorf("expected groups [admins], got %v", id.Groups)
	}
}

func TestParseHeadersGroupsFiltersEmptyValues(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Remote-Groups", " admins, ,devs,,  ")

	id := ParseHeaders(r)
	if id == nil {
		t.Fatal("expected identity, got nil")
	}

	wantGroups := []string{"admins", "devs"}
	if !reflect.DeepEqual(id.Groups, wantGroups) {
		t.Errorf("expected groups %v, got %v", wantGroups, id.Groups)
	}
}

func TestMiddlewareSetsContext(t *testing.T) {
	var captured *Identity
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = FromContext(r.Context())
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Remote-User", "patrick")
	r.Header.Set("Remote-Name", "Patrick Fanella")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if captured == nil {
		t.Fatal("expected identity in context, got nil")
	}
	if captured.Username != "patrick" {
		t.Errorf("expected username patrick, got %s", captured.Username)
	}
}

func TestMiddlewareNoHeadersPassesThrough(t *testing.T) {
	called := false
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if FromContext(r.Context()) != nil {
			t.Error("expected nil identity for unauthenticated request")
		}
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if !called {
		t.Error("expected next handler to be called")
	}
}

func TestFromContextEmptyContext(t *testing.T) {
	id := FromContext(context.Background())
	if id != nil {
		t.Errorf("expected nil from empty context, got %+v", id)
	}
}
