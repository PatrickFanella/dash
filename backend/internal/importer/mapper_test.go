package importer

import (
	"testing"
)

func TestMapSection_ServiceType(t *testing.T) {
	ds := DashySection{
		Name:  "MAGI-01 // Media",
		Items: []DashyItem{{Title: "Plex", URL: "https://plex.subcult.tv"}},
	}
	params := MapSection(ds, 0)
	if params.SectionType != "services" {
		t.Errorf("expected services, got %s", params.SectionType)
	}
}

func TestMapSection_MetricsType(t *testing.T) {
	ds := DashySection{
		Name:    "SYNAPSE // Overview",
		Widgets: []DashyWidget{{Type: "clock"}},
	}
	params := MapSection(ds, 0)
	if params.SectionType != "metrics" {
		t.Errorf("expected metrics, got %s", params.SectionType)
	}
}

func TestMapSection_DefaultCols(t *testing.T) {
	ds := DashySection{Name: "Test", DisplayData: DashyDisplayData{Cols: 0}}
	params := MapSection(ds, 5)
	if params.Cols != 3 {
		t.Errorf("expected default cols 3, got %d", params.Cols)
	}
	if params.SortOrder != 5 {
		t.Errorf("expected sort_order 5, got %d", params.SortOrder)
	}
}

func TestMapSection_CollapsedPreserved(t *testing.T) {
	ds := DashySection{
		Name:        "External // Links",
		DisplayData: DashyDisplayData{Collapsed: true, Cols: 4},
	}
	params := MapSection(ds, 10)
	if !params.Collapsed {
		t.Error("expected collapsed true")
	}
	if params.Cols != 4 {
		t.Errorf("expected cols 4, got %d", params.Cols)
	}
}

func TestMapItem_ExplicitStatusCheckTrue(t *testing.T) {
	b := true
	di := DashyItem{Title: "Plex", URL: "https://plex.subcult.tv", StatusCheck: &b}
	params := MapItem(di, 0, true)
	if !params.StatusCheck {
		t.Error("expected status_check true")
	}
}

func TestMapItem_ExplicitStatusCheckFalse(t *testing.T) {
	b := false
	di := DashyItem{Title: "Test", URL: "https://example.com", StatusCheck: &b}
	params := MapItem(di, 0, true)
	if params.StatusCheck {
		t.Error("expected status_check false")
	}
}

func TestMapItem_NilStatusCheck_ServiceSection(t *testing.T) {
	di := DashyItem{Title: "SomeService", URL: "https://example.com", StatusCheck: nil}
	params := MapItem(di, 0, true)
	if !params.StatusCheck {
		t.Error("expected status_check true (service section default)")
	}
}

func TestMapItem_NilStatusCheck_ExternalSection(t *testing.T) {
	di := DashyItem{Title: "GitHub", URL: "https://github.com", StatusCheck: nil}
	params := MapItem(di, 0, false)
	if params.StatusCheck {
		t.Error("expected status_check false (external section default)")
	}
}

func TestMapItem_StatusCheckURL_Present(t *testing.T) {
	b := true
	di := DashyItem{
		Title:          "Plex",
		URL:            "https://plex.subcult.tv",
		StatusCheck:    &b,
		StatusCheckURL: "http://10.0.0.200:32400/identity",
	}
	params := MapItem(di, 0, true)
	if !params.StatusCheckUrl.Valid {
		t.Fatal("expected status_check_url to be valid")
	}
	if params.StatusCheckUrl.String != "http://10.0.0.200:32400/identity" {
		t.Errorf("got %q", params.StatusCheckUrl.String)
	}
}

func TestMapItem_StatusCheckURL_Absent(t *testing.T) {
	b := true
	di := DashyItem{Title: "Sonarr", URL: "https://sonarr.subcult.tv", StatusCheck: &b}
	params := MapItem(di, 0, true)
	if params.StatusCheckUrl.Valid {
		t.Error("expected status_check_url to be null")
	}
}

func TestMapItem_SortOrder(t *testing.T) {
	di := DashyItem{Title: "Test", URL: "https://example.com"}
	params := MapItem(di, 7, true)
	if params.SortOrder != 7 {
		t.Errorf("expected sort_order 7, got %d", params.SortOrder)
	}
}

func TestMapItem_IconPreserved(t *testing.T) {
	di := DashyItem{Title: "Test", URL: "https://example.com", Icon: "hl-plex"}
	params := MapItem(di, 0, true)
	if params.Icon != "hl-plex" {
		t.Errorf("expected icon hl-plex, got %s", params.Icon)
	}
}

func TestSectionDefaultStatusCheck_AllExplicit(t *testing.T) {
	b := true
	items := []DashyItem{{StatusCheck: &b}, {StatusCheck: &b}}
	if !sectionDefaultStatusCheck(items) {
		t.Error("expected true when items have explicit statusCheck")
	}
}

func TestSectionDefaultStatusCheck_AllNil(t *testing.T) {
	items := []DashyItem{{StatusCheck: nil}, {StatusCheck: nil}}
	if sectionDefaultStatusCheck(items) {
		t.Error("expected false when no items have explicit statusCheck")
	}
}

func TestSectionDefaultStatusCheck_Empty(t *testing.T) {
	if sectionDefaultStatusCheck(nil) {
		t.Error("expected false for empty items")
	}
}
