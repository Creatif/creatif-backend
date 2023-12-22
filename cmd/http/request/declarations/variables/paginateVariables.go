package variables

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type PaginateVariables struct {
	ProjectID      string `param:"projectID"`
	Locales        string `query:"locales"`
	Limit          int    `query:"limit"`
	Page           int    `query:"page"`
	Filters        string `query:"filters"`
	Groups         string `query:"groups"`
	Behaviour      string `query:"behaviour"`
	Search         string `query:"search"`
	OrderBy        string `query:"orderBy"`
	OrderDirection string `query:"direction"`
	In             string `query:"in"`

	SanitizedGroups  []string
	SanitizedLocales []string
}

func SanitizePaginateVariables(model PaginateVariables) PaginateVariables {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)
	model.Behaviour = p.Sanitize(model.Behaviour)
	model.In = p.Sanitize(model.In)
	model.OrderBy = p.Sanitize(model.OrderBy)

	if model.Groups != "" {
		model.SanitizedGroups = sdk.Map(strings.Split(model.Groups, ","), func(idx int, value string) string {
			return p.Sanitize(strings.TrimSpace(value))
		})
	} else {
		model.SanitizedGroups = make([]string, 0)
	}

	if model.Locales != "" {
		model.SanitizedLocales = sdk.Map(strings.Split(model.Locales, ","), func(idx int, value string) string {
			return p.Sanitize(strings.TrimSpace(value))
		})
	} else {
		model.SanitizedLocales = make([]string, 0)
	}

	return model
}
