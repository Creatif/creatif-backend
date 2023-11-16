package app

import (
	"github.com/microcosm-cc/bluemonday"
)

type PaginateProjects struct {
	Limit          int    `query:"limit"`
	Page           int    `query:"page"`
	Search         string `query:"search"`
	OrderBy        string `query:"orderBy"`
	OrderDirection string `query:"direction"`
}

func SanitizePaginateProjects(model PaginateProjects) PaginateProjects {
	p := bluemonday.StrictPolicy()
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)

	return model
}
