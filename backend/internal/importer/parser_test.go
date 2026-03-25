package importer

import (
	"strings"
	"testing"
)

func TestParseFile_SectionCount(t *testing.T) {
	cfg, err := ParseFile("../../testdata/dashy_conf.yml")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(cfg.Sections) != 11 {
		t.Fatalf("expected 11 sections, got %d", len(cfg.Sections))
	}
}

func TestParseFile_SectionNames(t *testing.T) {
	cfg, _ := ParseFile("../../testdata/dashy_conf.yml")
	expected := []string{
		"SYNAPSE // Overview",
		"SYNAPSE // CPU & Memory",
		"SYNAPSE // Network & Temp",
		"SYNAPSE // Service Status",
		"MAGI-01 // Media",
		"MAGI-02 // Cloud",
		"MAGI-03 // Operations",
		"Terminal // Dev",
		"Archive // Knowledge",
		"Vitals // Life",
		"External // Links",
	}
	for i, name := range expected {
		if cfg.Sections[i].Name != name {
			t.Errorf("section[%d]: expected %q, got %q", i, name, cfg.Sections[i].Name)
		}
	}
}

func TestParseFile_WidgetSectionsHaveNoItems(t *testing.T) {
	cfg, _ := ParseFile("../../testdata/dashy_conf.yml")
	for i := 0; i < 4; i++ {
		sec := cfg.Sections[i]
		if len(sec.Widgets) == 0 {
			t.Errorf("section[%d] %q: expected widgets", i, sec.Name)
		}
		if len(sec.Items) != 0 {
			t.Errorf("section[%d] %q: expected no items, got %d", i, sec.Name, len(sec.Items))
		}
	}
}

func TestParseFile_ServiceSectionItemCounts(t *testing.T) {
	cfg, _ := ParseFile("../../testdata/dashy_conf.yml")
	// Count all items to verify
	total := 0
	for _, sec := range cfg.Sections {
		total += len(sec.Items)
	}
	if total < 40 {
		t.Errorf("expected at least 40 total items, got %d", total)
	}
}

func TestParseFile_PlexStatusCheckURL(t *testing.T) {
	cfg, _ := ParseFile("../../testdata/dashy_conf.yml")
	media := cfg.Sections[4]
	plex := media.Items[0]
	if plex.StatusCheckURL != "http://10.0.0.200:32400/identity" {
		t.Errorf("Plex statusCheckUrl: got %q", plex.StatusCheckURL)
	}
}

func TestParseFile_ExternalItemsHaveNilStatusCheck(t *testing.T) {
	cfg, _ := ParseFile("../../testdata/dashy_conf.yml")
	external := cfg.Sections[10]
	for _, item := range external.Items {
		if item.StatusCheck != nil {
			t.Errorf("External item %q: expected nil StatusCheck, got %v", item.Title, *item.StatusCheck)
		}
	}
}

func TestParseFile_IconPrefixes(t *testing.T) {
	cfg, _ := ParseFile("../../testdata/dashy_conf.yml")
	media := cfg.Sections[4]
	if !strings.HasPrefix(media.Items[0].Icon, "hl-") {
		t.Errorf("expected hl- prefix for Plex icon, got %q", media.Items[0].Icon)
	}
	if !strings.HasPrefix(media.Icon, "fas ") {
		t.Errorf("expected fas prefix for Media section icon, got %q", media.Icon)
	}
}

func TestParseFile_CollapsedSection(t *testing.T) {
	cfg, _ := ParseFile("../../testdata/dashy_conf.yml")
	if !cfg.Sections[10].DisplayData.Collapsed {
		t.Error("External section should be collapsed")
	}
	if cfg.Sections[4].DisplayData.Collapsed {
		t.Error("Media section should not be collapsed")
	}
}

func TestParse_EmptyInput(t *testing.T) {
	_, err := Parse(strings.NewReader(""))
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestParse_NoSections(t *testing.T) {
	_, err := Parse(strings.NewReader("pageInfo:\n  title: test\n"))
	if err == nil {
		t.Error("expected error for config with no sections")
	}
}
