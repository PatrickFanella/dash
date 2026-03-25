package importer

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickfanella/dash/backend/internal/models"
)

type ImportResult struct {
	SectionsCreated int `json:"sections_created"`
	SectionsUpdated int `json:"sections_updated"`
	ServicesCreated int `json:"services_created"`
	ServicesUpdated int `json:"services_updated"`
}

func Run(ctx context.Context, pool *pgxpool.Pool, cfg *DashyConfig) (*ImportResult, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := models.New(tx)
	result := &ImportResult{}

	for sectionIdx, ds := range cfg.Sections {
		sectionParams := MapSection(ds, sectionIdx)
		section, created, err := upsertSection(ctx, q, sectionParams)
		if err != nil {
			return nil, fmt.Errorf("upsert section %q: %w", ds.Name, err)
		}
		if created {
			result.SectionsCreated++
		} else {
			result.SectionsUpdated++
		}

		if len(ds.Items) == 0 {
			continue
		}

		defaultStatus := sectionDefaultStatusCheck(ds.Items)
		for itemIdx, di := range ds.Items {
			serviceParams := MapItem(di, itemIdx, defaultStatus)
			svc, svcCreated, err := upsertService(ctx, q, serviceParams)
			if err != nil {
				return nil, fmt.Errorf("upsert service %q: %w", di.Title, err)
			}
			if svcCreated {
				result.ServicesCreated++
			} else {
				result.ServicesUpdated++
			}

			_, err = q.AddServiceToSection(ctx, models.AddServiceToSectionParams{
				ServiceID: svc.ID,
				SectionID: section.ID,
				SortOrder: int32(itemIdx),
			})
			if err != nil {
				if !errors.Is(err, pgx.ErrNoRows) {
					return nil, fmt.Errorf("link %q to %q: %w", di.Title, ds.Name, err)
				}
				err = q.UpdateMappingSortOrder(ctx, models.UpdateMappingSortOrderParams{
					ServiceID: svc.ID,
					SectionID: section.ID,
					SortOrder: int32(itemIdx),
				})
				if err != nil {
					return nil, fmt.Errorf("update mapping sort order: %w", err)
				}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}

func upsertSection(ctx context.Context, q *models.Queries, params models.CreateSectionParams) (models.Section, bool, error) {
	existing, err := q.GetSectionByName(ctx, params.Name)
	if err == nil {
		updated, err := q.UpdateSection(ctx, models.UpdateSectionParams{
			ID:          existing.ID,
			Name:        params.Name,
			Icon:        params.Icon,
			Cols:        params.Cols,
			Collapsed:   params.Collapsed,
			SortOrder:   params.SortOrder,
			SectionType: params.SectionType,
		})
		return updated, false, err
	}
	created, err := q.CreateSection(ctx, params)
	return created, true, err
}

func upsertService(ctx context.Context, q *models.Queries, params models.CreateServiceParams) (models.Service, bool, error) {
	existing, err := q.GetServiceByURL(ctx, params.Url)
	if err == nil {
		updated, err := q.UpdateService(ctx, models.UpdateServiceParams{
			ID:             existing.ID,
			Title:          params.Title,
			Url:            params.Url,
			Description:    params.Description,
			Icon:           params.Icon,
			StatusCheck:    params.StatusCheck,
			StatusCheckUrl: params.StatusCheckUrl,
			SortOrder:      params.SortOrder,
		})
		return updated, false, err
	}
	created, err := q.CreateService(ctx, params)
	return created, true, err
}
