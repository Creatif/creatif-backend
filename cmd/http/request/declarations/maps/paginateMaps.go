package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type PaginateMaps struct {
	ProjectID      string `param:"projectID"`
	Limit          int    `query:"limit"`
	Page           int    `query:"page"`
	Search         string `query:"search"`
	OrderBy        string `query:"orderBy"`
	OrderDirection string `query:"direction"`
}

func SanitizePaginateMaps(model PaginateMaps) PaginateMaps {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)
	model.OrderBy = p.Sanitize(model.OrderBy)

	return model
}
