package lists

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type PaginateListItems struct {
	ProjectID      string `param:"projectID"`
	Locales        string `query:"locales"`
	ListName       string `param:"name"`
	Limit          int    `query:"limit"`
	Page           int    `query:"page"`
	Filters        string `query:"filters"`
	Behaviour      string `query:"behaviour"`
	Groups         string `query:"groups"`
	Search         string `query:"search"`
	OrderBy        string `query:"orderBy"`
	OrderDirection string `query:"direction"`
	In             string `query:"in"`

	SanitizedGroups  []string
	SanitizedLocales []string
}

func SanitizePaginateListItems(model PaginateListItems) PaginateListItems {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)
	model.In = p.Sanitize(model.In)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Behaviour = p.Sanitize(model.Behaviour)

	if model.Groups != "" {
		model.SanitizedGroups = sdk.Map(strings.Split(model.Groups, ","), func(idx int, value string) string {
			return p.Sanitize(strings.TrimSpace(value))
		})
	}

	if model.Locales != "" {
		model.SanitizedLocales = sdk.Map(strings.Split(model.Locales, ","), func(idx int, value string) string {
			sanitized := p.Sanitize(strings.TrimSpace(value))
			locale, _ := locales.GetIDWithAlpha(sanitized)

			return locale
		})
	}

	return model
}
