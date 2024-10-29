package app

import "github.com/microcosm-cc/bluemonday"

type GetDashboardStats struct {
	ProjectID string `param:"projectId"`
}

func SanitizeGetDashboardStats(model GetDashboardStats) GetDashboardStats {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
