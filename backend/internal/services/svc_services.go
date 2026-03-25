package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickfanella/dash/backend/internal/models"
)

type ServiceService struct {
	queries *models.Queries
	pool    *pgxpool.Pool
}

func NewServiceService(q *models.Queries, pool *pgxpool.Pool) *ServiceService {
	return &ServiceService{queries: q, pool: pool}
}

func (s *ServiceService) List(ctx context.Context) ([]models.Service, error) {
	return s.queries.ListServices(ctx)
}

func (s *ServiceService) Get(ctx context.Context, id pgtype.UUID) (models.Service, error) {
	return s.queries.GetService(ctx, id)
}

func (s *ServiceService) GetSectionIDs(ctx context.Context, serviceID pgtype.UUID) ([]pgtype.UUID, error) {
	return s.queries.ListSectionIDsByService(ctx, serviceID)
}

func (s *ServiceService) Create(ctx context.Context, params models.CreateServiceParams, sectionIDs []pgtype.UUID) (models.Service, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.Service{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)
	svc, err := qtx.CreateService(ctx, params)
	if err != nil {
		return models.Service{}, err
	}

	for i, secID := range sectionIDs {
		_, err := qtx.AddServiceToSection(ctx, models.AddServiceToSectionParams{
			ServiceID: svc.ID,
			SectionID: secID,
			SortOrder: int32(i),
		})
		if err != nil {
			return models.Service{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Service{}, err
	}
	return svc, nil
}

func (s *ServiceService) Update(ctx context.Context, params models.UpdateServiceParams, sectionIDs []pgtype.UUID) (models.Service, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return models.Service{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)
	svc, err := qtx.UpdateService(ctx, params)
	if err != nil {
		return models.Service{}, err
	}

	if sectionIDs != nil {
		if err := qtx.DeleteMappingsByService(ctx, svc.ID); err != nil {
			return models.Service{}, err
		}
		for i, secID := range sectionIDs {
			_, err := qtx.AddServiceToSection(ctx, models.AddServiceToSectionParams{
				ServiceID: svc.ID,
				SectionID: secID,
				SortOrder: int32(i),
			})
			if err != nil {
				return models.Service{}, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Service{}, err
	}
	return svc, nil
}

func (s *ServiceService) Delete(ctx context.Context, id pgtype.UUID) error {
	return s.queries.DeleteService(ctx, id)
}
