package publicApi

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type PaginateListItems struct {
	ProjectID      string `param:"projectID"`
	ListName       string `param:"name"`
	Locales        string `query:"locales"`
	Page           int    `query:"page"`
	Groups         string `query:"groups"`
	Search         string `query:"search"`
	OrderBy        string `query:"orderBy"`
	OrderDirection string `query:"direction"`

	SanitizedGroups  []string
	SanitizedLocales []string
	SanitizedFields  []string
}

func SanitizePaginateListItems(model PaginateListItems) PaginateListItems {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)
	model.OrderBy = p.Sanitize(model.OrderBy)

	if model.Groups != "" {
		model.SanitizedGroups = sdk.Map(strings.Split(model.Groups, ","), func(idx int, value string) string {
			return p.Sanitize(strings.TrimSpace(value))
		})
	}

	if model.Locales != "" {
		model.SanitizedLocales = sdk.Map(strings.Split(model.Locales, ","), func(idx int, value string) string {
			return p.Sanitize(strings.TrimSpace(value))
		})
	}

	return model
}
