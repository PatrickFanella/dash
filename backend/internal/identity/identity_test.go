package identity

import (
	"net/http"
	"reflect"
	"testing"
)

func TestParseHeaders_NoHeadersReturnsNil(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	got := ParseHeaders(req)
	if got != nil {
		t.Fatalf("ParseHeaders() = %+v, want nil", got)
	}
}

func TestParseHeaders_AllHeaders(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Remote-User", "alice")
	req.Header.Set("Remote-Name", "Alice Smith")
	req.Header.Set("Remote-Email", "alice@example.com")
	req.Header.Set("Remote-Groups", "admins, devs,team-a")

	got := ParseHeaders(req)
	if got == nil {
		t.Fatal("ParseHeaders() = nil, want identity")
	}

	if got.Username != "alice" {
		t.Errorf("Username = %q, want %q", got.Username, "alice")
	}
	if got.DisplayName != "Alice Smith" {
		t.Errorf("DisplayName = %q, want %q", got.DisplayName, "Alice Smith")
	}
	if got.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", got.Email, "alice@example.com")
	}

	wantGroups := []string{"admins", "devs", "team-a"}
	if !reflect.DeepEqual(got.Groups, wantGroups) {
		t.Errorf("Groups = %#v, want %#v", got.Groups, wantGroups)
	}
}

func TestParseHeaders_PartialHeaders(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Remote-Email", "bob@example.com")

	got := ParseHeaders(req)
	if got == nil {
		t.Fatal("ParseHeaders() = nil, want identity")
	}

	if got.Username != "" {
		t.Errorf("Username = %q, want empty", got.Username)
	}
	if got.DisplayName != "" {
		t.Errorf("DisplayName = %q, want empty", got.DisplayName)
	}
	if got.Email != "bob@example.com" {
		t.Errorf("Email = %q, want %q", got.Email, "bob@example.com")
	}
	if got.Groups != nil {
		t.Errorf("Groups = %#v, want nil", got.Groups)
	}
}

func TestParseHeaders_GroupsFiltersEmptyValues(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Remote-Groups", " admins, ,devs,,  ")

	got := ParseHeaders(req)
	if got == nil {
		t.Fatal("ParseHeaders() = nil, want identity")
	}

	wantGroups := []string{"admins", "devs"}
	if !reflect.DeepEqual(got.Groups, wantGroups) {
		t.Errorf("Groups = %#v, want %#v", got.Groups, wantGroups)
	}
}
