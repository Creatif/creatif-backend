package variables

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type PaginateVariables struct {
	ProjectID      string `param:"projectID"`
	Locale         string `param:"locale"`
	Limit          int    `query:"limit"`
	Page           int    `query:"page"`
	Filters        string `query:"filters"`
	Groups         string `query:"groups"`
	Search         string `query:"search"`
	OrderBy        string `query:"orderBy"`
	OrderDirection string `query:"orderDirection"`
	In             string `query:"in"`

	SanitizedGroups []string
}

func SanitizePaginateVariables(model PaginateVariables) PaginateVariables {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)
	model.In = p.Sanitize(model.In)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Locale = p.Sanitize(model.Locale)

	if model.Groups != "" {
		newGroups := sdk.Map(strings.Split(model.Groups, ","), func(idx int, value string) string {
			return p.Sanitize(strings.TrimSpace(value))
		})

		model.SanitizedGroups = newGroups
	}

	return model
}
