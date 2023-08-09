package sdk

import (
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type PaginationView[T any] struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Data  []T   `json:"data"`
}

type LogicView[T any] struct {
	Total int64
	Data  []T
}

type Pagination struct {
	Limit          int    `query:"limit"`
	Page           int    `query:"page"`
	Filters        string `query:"filters"`
	Search         string `query:"search"`
	OrderBy        string `query:"orderBy"`
	OrderDirection string `query:"orderDirection"`
	In             string `query:"in"`
}

func SanitizePagination(model Pagination) Pagination {
	p := bluemonday.StrictPolicy()
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)
	model.In = p.Sanitize(model.In)

	return model
}

func ParseIn(items string) []string {
	if items == "" {
		return []string{}
	}

	return strings.Split(strings.ReplaceAll(items, " ", ""), ",")
}

func ParseFilters(filters string) map[string]string {
	sanitized := strings.ReplaceAll(filters, " ", "")

	parts := strings.Split(sanitized, ",")
	filtersMap := make(map[string]string)

	for _, part := range parts {
		split := strings.Split(part, "|")
		if len(split) == 2 {
			filtersMap[split[0]] = split[1]
		}
	}

	return filtersMap
}
