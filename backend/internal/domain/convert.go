package domain

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/patrickfanella/dash/backend/internal/models"
)

// SectionFromModel converts a sqlc-generated Section to a domain Section.
func SectionFromModel(m models.Section) Section {
	return Section{
		ID:          FormatUUID(m.ID),
		Name:        m.Name,
		Icon:        m.Icon,
		Cols:        m.Cols,
		Collapsed:   m.Collapsed,
		SortOrder:   m.SortOrder,
		SectionType: m.SectionType,
		CreatedAt:   m.CreatedAt.Time,
		UpdatedAt:   m.UpdatedAt.Time,
	}
}

// SectionsFromModels converts a slice of sqlc Sections to domain Sections.
func SectionsFromModels(ms []models.Section) []Section {
	out := make([]Section, len(ms))
	for i, m := range ms {
		out[i] = SectionFromModel(m)
	}
	return out
}

// ServiceFromModel converts a sqlc-generated Service to a domain Service.
func ServiceFromModel(m models.Service) Service {
	s := Service{
		ID:          FormatUUID(m.ID),
		Title:       m.Title,
		URL:         m.Url,
		Description: m.Description,
		Icon:        m.Icon,
		StatusCheck: m.StatusCheck,
		SortOrder:   m.SortOrder,
		CreatedAt:   m.CreatedAt.Time,
		UpdatedAt:   m.UpdatedAt.Time,
	}
	if m.StatusCheckUrl.Valid {
		s.StatusCheckURL = &m.StatusCheckUrl.String
	}
	return s
}

// ServicesFromModels converts a slice of sqlc Services to domain Services.
func ServicesFromModels(ms []models.Service) []Service {
	out := make([]Service, len(ms))
	for i, m := range ms {
		out[i] = ServiceFromModel(m)
	}
	return out
}

// ParseUUID converts a string UUID to a pgtype.UUID for database queries.
func ParseUUID(s string) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(s); err != nil {
		return id, fmt.Errorf("invalid UUID: %s", s)
	}
	return id, nil
}

// FormatUUID converts a pgtype.UUID to its string representation.
func FormatUUID(id pgtype.UUID) string {
	b := id.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// UUIDsToStrings converts a slice of pgtype.UUIDs to string UUIDs.
func UUIDsToStrings(ids []pgtype.UUID) []string {
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = FormatUUID(id)
	}
	return strs
}
