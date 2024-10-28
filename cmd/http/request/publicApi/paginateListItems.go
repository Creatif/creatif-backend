package publicApi

import (
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type PaginateListItems struct {
	ProjectID      string `param:"projectID"`
	ListName       string `param:"name"`
	Locales        string `query:"locales"`
	Page           int    `query:"page"`
	Limit          int    `query:"limit"`
	Groups         string `query:"groups"`
	Search         string `query:"search"`
	OrderBy        string `query:"orderBy"`
	OrderDirection string `query:"direction"`
	Options        string `query:"options"`
	VersionName    string
	Query          string `query:"query"`

	SanitizedGroups  []string
	SanitizedLocales []string
	SanitizedFields  []string
	ResolvedOptions  GetListItemByIDOptions
	SanitizedQuery   []Query
}

func SanitizePaginateListItems(model PaginateListItems) (PaginateListItems, error) {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.VersionName = p.Sanitize(model.VersionName)

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

	if model.Options != "" {
		model.ResolvedOptions = resolveListOptions(model.Options)
	}

	if model.Query != "" {
		var q []Query
		if err := json.Unmarshal([]byte(model.Query), &q); err != nil {
			return model, err
		}

		model.SanitizedQuery = sdk.Map(q, func(idx int, value Query) Query {
			return Query{
				Column:   p.Sanitize(value.Column),
				Value:    p.Sanitize(value.Value),
				Operator: p.Sanitize(value.Operator),
			}
		})
	}

	return model, nil
}
