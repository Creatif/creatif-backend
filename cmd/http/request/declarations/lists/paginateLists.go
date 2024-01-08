package lists

import (
	"github.com/microcosm-cc/bluemonday"
)

type PaginateLists struct {
	ProjectID      string `param:"projectID"`
	Limit          int    `query:"limit"`
	Page           int    `query:"page"`
	Search         string `query:"search"`
	OrderBy        string `query:"orderBy"`
	OrderDirection string `query:"direction"`
}

func SanitizePaginateLists(model PaginateLists) PaginateLists {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)
	model.OrderBy = p.Sanitize(model.OrderBy)

	return model
}
