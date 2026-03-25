package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/patrickfanella/dash/backend/internal/models"
)

type SectionService struct {
	queries *models.Queries
}

func NewSectionService(q *models.Queries) *SectionService {
	return &SectionService{queries: q}
}

func (s *SectionService) List(ctx context.Context) ([]models.Section, error) {
	return s.queries.ListSections(ctx)
}

func (s *SectionService) Get(ctx context.Context, id pgtype.UUID) (models.Section, error) {
	return s.queries.GetSection(ctx, id)
}

func (s *SectionService) Create(ctx context.Context, params models.CreateSectionParams) (models.Section, error) {
	return s.queries.CreateSection(ctx, params)
}

func (s *SectionService) Update(ctx context.Context, params models.UpdateSectionParams) (models.Section, error) {
	return s.queries.UpdateSection(ctx, params)
}

func (s *SectionService) Delete(ctx context.Context, id pgtype.UUID) error {
	return s.queries.DeleteSection(ctx, id)
}

func ParseUUID(s string) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(s); err != nil {
		return id, fmt.Errorf("invalid UUID: %s", s)
	}
	return id, nil
}
