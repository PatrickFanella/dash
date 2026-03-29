package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickfanella/dash/backend/internal/domain"
	"github.com/patrickfanella/dash/backend/internal/models"
)

type ServiceService struct {
	queries *models.Queries
	pool    *pgxpool.Pool
}

func NewServiceService(q *models.Queries, pool *pgxpool.Pool) *ServiceService {
	return &ServiceService{queries: q, pool: pool}
}

func (s *ServiceService) List(ctx context.Context) ([]domain.Service, error) {
	svcs, err := s.queries.ListServices(ctx)
	if err != nil {
		return nil, err
	}
	return domain.ServicesFromModels(svcs), nil
}

func (s *ServiceService) Get(ctx context.Context, id string) (domain.Service, error) {
	uid, err := domain.ParseUUID(id)
	if err != nil {
		return domain.Service{}, domain.ValidationErr(err.Error())
	}
	svc, err := s.queries.GetService(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Service{}, domain.NotFoundErr("service", id)
		}
		return domain.Service{}, err
	}
	return domain.ServiceFromModel(svc), nil
}

func (s *ServiceService) GetSectionIDs(ctx context.Context, serviceID string) ([]string, error) {
	uid, err := domain.ParseUUID(serviceID)
	if err != nil {
		return nil, err
	}
	ids, err := s.queries.ListSectionIDsByService(ctx, uid)
	if err != nil {
		return nil, err
	}
	return domain.UUIDsToStrings(ids), nil
}

type CreateServiceInput struct {
	Title          string
	URL            string
	Description    string
	Icon           string
	StatusCheck    *bool
	StatusCheckURL *string
	SortOrder      int32
	SectionIDs     []string
}

func (s *ServiceService) Create(ctx context.Context, input CreateServiceInput) (domain.Service, error) {
	if input.Title == "" || input.URL == "" {
		return domain.Service{}, domain.ValidationErr("title and url are required")
	}

	statusCheck := true
	if input.StatusCheck != nil {
		statusCheck = *input.StatusCheck
	}

	var statusCheckURL pgtype.Text
	if input.StatusCheckURL != nil {
		statusCheckURL = pgtype.Text{String: *input.StatusCheckURL, Valid: true}
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.Service{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)
	svc, err := qtx.CreateService(ctx, models.CreateServiceParams{
		Title:          input.Title,
		Url:            input.URL,
		Description:    input.Description,
		Icon:           input.Icon,
		StatusCheck:    statusCheck,
		StatusCheckUrl: statusCheckURL,
		SortOrder:      input.SortOrder,
	})
	if err != nil {
		return domain.Service{}, err
	}

	for i, sid := range input.SectionIDs {
		secID, err := domain.ParseUUID(sid)
		if err != nil {
			return domain.Service{}, domain.ValidationErr(err.Error())
		}
		_, err = qtx.AddServiceToSection(ctx, models.AddServiceToSectionParams{
			ServiceID: svc.ID,
			SectionID: secID,
			SortOrder: int32(i),
		})
		if err != nil {
			if isFKViolation(err) {
				return domain.Service{}, domain.ValidationErr("section not found: " + sid)
			}
			return domain.Service{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Service{}, err
	}
	return domain.ServiceFromModel(svc), nil
}

func (s *ServiceService) Update(ctx context.Context, id string, input CreateServiceInput) (domain.Service, error) {
	uid, err := domain.ParseUUID(id)
	if err != nil {
		return domain.Service{}, domain.ValidationErr(err.Error())
	}

	statusCheck := true
	if input.StatusCheck != nil {
		statusCheck = *input.StatusCheck
	}

	var statusCheckURL pgtype.Text
	if input.StatusCheckURL != nil {
		statusCheckURL = pgtype.Text{String: *input.StatusCheckURL, Valid: true}
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.Service{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)
	svc, err := qtx.UpdateService(ctx, models.UpdateServiceParams{
		ID:             uid,
		Title:          input.Title,
		Url:            input.URL,
		Description:    input.Description,
		Icon:           input.Icon,
		StatusCheck:    statusCheck,
		StatusCheckUrl: statusCheckURL,
		SortOrder:      input.SortOrder,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Service{}, domain.NotFoundErr("service", id)
		}
		return domain.Service{}, err
	}

	if input.SectionIDs != nil {
		if err := qtx.DeleteMappingsByService(ctx, svc.ID); err != nil {
			return domain.Service{}, err
		}
		for i, sid := range input.SectionIDs {
			secID, err := domain.ParseUUID(sid)
			if err != nil {
				return domain.Service{}, domain.ValidationErr(err.Error())
			}
			_, err = qtx.AddServiceToSection(ctx, models.AddServiceToSectionParams{
				ServiceID: svc.ID,
				SectionID: secID,
				SortOrder: int32(i),
			})
			if err != nil {
				if isFKViolation(err) {
					return domain.Service{}, domain.ValidationErr("section not found: " + sid)
				}
				return domain.Service{}, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Service{}, err
	}
	return domain.ServiceFromModel(svc), nil
}

func (s *ServiceService) Delete(ctx context.Context, id string) error {
	uid, err := domain.ParseUUID(id)
	if err != nil {
		return domain.ValidationErr(err.Error())
	}
	return s.queries.DeleteService(ctx, uid)
}

func isFKViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
