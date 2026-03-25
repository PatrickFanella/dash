package importer

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/patrickfanella/dash/backend/internal/models"
)

func MapSection(ds DashySection, sectionIndex int) models.CreateSectionParams {
	sectionType := "services"
	if len(ds.Items) == 0 && len(ds.Widgets) > 0 {
		sectionType = "metrics"
	}

	cols := ds.DisplayData.Cols
	if cols == 0 {
		cols = 3
	}

	return models.CreateSectionParams{
		Name:        ds.Name,
		Icon:        ds.Icon,
		Cols:        cols,
		Collapsed:   ds.DisplayData.Collapsed,
		SortOrder:   int32(sectionIndex),
		SectionType: sectionType,
	}
}

func MapItem(di DashyItem, itemIndex int, sectionDefault bool) models.CreateServiceParams {
	statusCheck := resolveStatusCheck(di.StatusCheck, sectionDefault)

	var statusCheckURL pgtype.Text
	if di.StatusCheckURL != "" {
		statusCheckURL = pgtype.Text{String: di.StatusCheckURL, Valid: true}
	}

	return models.CreateServiceParams{
		Title:          di.Title,
		Url:            di.URL,
		Description:    di.Description,
		Icon:           di.Icon,
		StatusCheck:    statusCheck,
		StatusCheckUrl: statusCheckURL,
		SortOrder:      int32(itemIndex),
	}
}

func resolveStatusCheck(explicit *bool, sectionDefault bool) bool {
	if explicit != nil {
		return *explicit
	}
	return sectionDefault
}

func sectionDefaultStatusCheck(items []DashyItem) bool {
	for _, item := range items {
		if item.StatusCheck != nil {
			return true
		}
	}
	return false
}
