package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/patrickfanella/dash/backend/internal/domain"
	"github.com/patrickfanella/dash/backend/internal/models"
)

type SectionService struct {
	queries *models.Queries
}

func NewSectionService(q *models.Queries) *SectionService {
	return &SectionService{queries: q}
}

func (s *SectionService) List(ctx context.Context) ([]domain.Section, error) {
	sections, err := s.queries.ListSections(ctx)
	if err != nil {
		return nil, err
	}
	return domain.SectionsFromModels(sections), nil
}

func (s *SectionService) ListNested(ctx context.Context) ([]domain.NestedSection, error) {
	sections, err := s.queries.ListSections(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.NestedSection, len(sections))
	for i, sec := range sections {
		svcs, err := s.queries.ListServicesBySection(ctx, sec.ID)
		if err != nil {
			svcs = []models.Service{}
		}
		result[i] = domain.NestedSection{
			Section:  domain.SectionFromModel(sec),
			Services: domain.ServicesFromModels(svcs),
		}
	}
	return result, nil
}

func (s *SectionService) Get(ctx context.Context, id string) (domain.Section, error) {
	uid, err := domain.ParseUUID(id)
	if err != nil {
		return domain.Section{}, domain.ValidationErr(err.Error())
	}
	sec, err := s.queries.GetSection(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Section{}, domain.NotFoundErr("section", id)
		}
		return domain.Section{}, err
	}
	return domain.SectionFromModel(sec), nil
}

type CreateSectionInput struct {
	Name        string
	Icon        string
	Cols        int32
	Collapsed   bool
	SortOrder   int32
	SectionType string
}

func (s *SectionService) Create(ctx context.Context, input CreateSectionInput) (domain.Section, error) {
	if input.Name == "" {
		return domain.Section{}, domain.ValidationErr("name is required")
	}
	if input.SectionType == "" {
		input.SectionType = "services"
	}
	if input.Cols == 0 {
		input.Cols = 3
	}

	sec, err := s.queries.CreateSection(ctx, models.CreateSectionParams{
		Name:        input.Name,
		Icon:        input.Icon,
		Cols:        input.Cols,
		Collapsed:   input.Collapsed,
		SortOrder:   input.SortOrder,
		SectionType: input.SectionType,
	})
	if err != nil {
		return domain.Section{}, err
	}
	return domain.SectionFromModel(sec), nil
}

func (s *SectionService) Update(ctx context.Context, id string, input CreateSectionInput) (domain.Section, error) {
	uid, err := domain.ParseUUID(id)
	if err != nil {
		return domain.Section{}, domain.ValidationErr(err.Error())
	}
	if input.SectionType == "" {
		input.SectionType = "services"
	}
	if input.Cols == 0 {
		input.Cols = 3
	}

	sec, err := s.queries.UpdateSection(ctx, models.UpdateSectionParams{
		ID:          uid,
		Name:        input.Name,
		Icon:        input.Icon,
		Cols:        input.Cols,
		Collapsed:   input.Collapsed,
		SortOrder:   input.SortOrder,
		SectionType: input.SectionType,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Section{}, domain.NotFoundErr("section", id)
		}
		return domain.Section{}, err
	}
	return domain.SectionFromModel(sec), nil
}

func (s *SectionService) Delete(ctx context.Context, id string) error {
	uid, err := domain.ParseUUID(id)
	if err != nil {
		return domain.ValidationErr(err.Error())
	}
	return s.queries.DeleteSection(ctx, uid)
}
