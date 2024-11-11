package paginateListItems

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared/queryProcessor"
	"fmt"
	"strings"
)

type defaults struct {
	page            int
	limit           int
	orderDirections string
}

type subQueries struct {
	sortBy  string
	search  string
	groups  string
	locales string
	query   string
}

func createPlaceholders(
	projectId string,
	versionId string,
	page,
	limit int,
	structureName string,
	providedLocales []string,
	search string,
) map[string]interface{} {
	placeholders := make(map[string]interface{})
	placeholders["projectId"] = projectId
	placeholders["versionId"] = versionId

	offset := (page - 1) * limit
	placeholders["offset"] = offset
	placeholders["structureIdentifier"] = structureName

	lcls := make([]string, len(providedLocales))
	for i, l := range providedLocales {
		alpha, _ := locales.GetIDWithAlpha(l)
		lcls[i] = alpha
	}
	placeholders["locales"] = lcls

	if search != "" {
		placeholders["searchOne"] = fmt.Sprintf("%%%s", search)
		placeholders["searchTwo"] = fmt.Sprintf("%s%%", search)
		placeholders["searchThree"] = fmt.Sprintf("%%%s%%", search)
		placeholders["searchFour"] = search
	}

	return placeholders
}

func createDefaults(page, limit int, orderDirection string) defaults {
	var def defaults
	def.limit = 100
	def.page = page

	if page == 0 || page < 1 {
		def.page = 1
	}

	if limit == 0 || limit < 0 {
		def.limit = 100
	}

	if orderDirection == "" {
		orderDirection = "ASC"
	}

	def.orderDirections = strings.ToUpper(orderDirection)

	return def
}

func createSubQueries(
	sortBy,
	search string,
	groups []string,
	lcls []string,
	query []queryProcessor.Query,
) (subQueries, error) {
	var sq subQueries

	sortByDefault := "lv.index"
	if sortBy != "" {
		sortByDefault = fmt.Sprintf("lv.%s", sortBy)
	}
	sq.sortBy = sortByDefault

	var searchSql string
	if search != "" {
		searchSql = fmt.Sprintf("AND (%s ILIKE @searchOne OR %s ILIKE @searchTwo OR %s ILIKE @searchThree OR %s ILIKE @searchFour)", "lv.variable_name", "lv.variable_name", "lv.variable_name", "lv.variable_name")
	}

	var groupsSql string
	if len(groups) > 0 {
		groupsSql = fmt.Sprintf("INNER JOIN %s AS g ON lv.id = g.variable_id AND '{%s}'::text[] && g.groups", (declarations.VariableGroup{}).TableName(), strings.Join(groups, ","))
	}

	var localesSql string
	if len(lcls) > 0 {
		localesSql = fmt.Sprintf("AND lv.locale_id IN (@locales)")
	}

	var querySql string
	if len(query) != 0 {
		s, err := queryProcessor.CreateSql(query)
		if err != nil {
			return subQueries{}, err
		}

		querySql = fmt.Sprintf("AND %s", s)
	}

	sq.search = searchSql
	sq.groups = groupsSql
	sq.locales = localesSql
	sq.query = querySql

	return sq, nil
}
