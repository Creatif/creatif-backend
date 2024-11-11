package paginateMapItems

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/domain/published"
	"creatif/pkg/app/services/publicApi/publicApiError"
	"creatif/pkg/app/services/shared/queryProcessor"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Item struct {
	ID            string `gorm:"column:id"`
	ShortID       string `gorm:"column:short_id"`
	StructureName string `gorm:"column:structure_name"`
	ProjectID     string `gorm:"column:project_id"`

	ItemName    string         `gorm:"column:variable_name"`
	ItemID      string         `gorm:"column:variable_id"`
	ItemShortID string         `gorm:"column:variable_short_id"`
	Value       datatypes.JSON `gorm:"type:jsonb"`
	Behaviour   string
	Locale      string         `gorm:"column:locale_id"`
	Index       float64        `gorm:"type:float"`
	Groups      pq.StringArray `gorm:"type:text[];not_null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func getItemSql(structureIdentifier string, page, limit int, order, sortBy, search string, lcls, groups []string, query []queryProcessor.Query) (string, map[string]interface{}, error) {
	offset := (page - 1) * limit
	placeholders := make(map[string]interface{})
	placeholders["offset"] = offset
	placeholders["structureIdentifier"] = structureIdentifier

	var searchSql string
	if search != "" {
		searchSql = fmt.Sprintf("AND (%s ILIKE @searchOne OR %s ILIKE @searchTwo OR %s ILIKE @searchThree OR %s ILIKE @searchFour)", "lv.variable_name", "lv.variable_name", "lv.variable_name", "lv.variable_name")
		placeholders["searchOne"] = fmt.Sprintf("%%%s", search)
		placeholders["searchTwo"] = fmt.Sprintf("%s%%", search)
		placeholders["searchThree"] = fmt.Sprintf("%%%s%%", search)
		placeholders["searchFour"] = search
	}

	var groupsSql string
	if len(groups) > 0 {
		groupsSql = fmt.Sprintf("AND'{%s}'::text[] && lv.groups ", strings.Join(groups, ","))
	}

	var localesSql string
	if len(lcls) > 0 {
		placeholders["locales"] = lcls
		localesSql = fmt.Sprintf("AND lv.locale_id IN (@locales)")
	}

	var querySql string
	if len(query) != 0 {
		s, err := queryProcessor.CreateSql(query)
		if err != nil {
			return "", nil, err
		}

		querySql = fmt.Sprintf("AND %s", s)
	}

	return fmt.Sprintf(`
SELECT 
    v.project_id,
	lv.id,
	lv.short_id,
	lv.name AS structure_name,
	lv.variable_name AS variable_name,
	lv.variable_id AS variable_id,
	lv.variable_short_id AS variable_short_id,
	lv.value,
	lv.behaviour,
	lv.locale_id,
	lv.index,
	lv.created_at,
	lv.updated_at,
	lv.groups
FROM %s AS lv
INNER JOIN %s AS v ON v.project_id = @projectId AND v.name = @versionName AND v.id = lv.version_id 
AND (lv.name = @structureIdentifier OR lv.id = @structureIdentifier OR lv.short_id = @structureIdentifier)
%s
%s
%s
%s
ORDER BY %s %s
OFFSET @offset
LIMIT %d
`,
		(published.PublishedMap{}).TableName(),
		(published.Version{}).TableName(),
		searchSql,
		groupsSql,
		localesSql,
		querySql,
		sortBy,
		order,
		limit,
	), placeholders, nil
}

func getVersion(projectId, versionName string) (published.Version, error) {
	var version published.Version
	var res *gorm.DB
	if versionName == "" {
		res = storage.Gorm().Raw(
			fmt.Sprintf("SELECT id, name FROM %s WHERE project_id = ? ORDER BY created_at DESC LIMIT 1",
				(published.Version{}).TableName()),
			projectId).
			Scan(&version)
	} else {
		res = storage.Gorm().Raw(
			fmt.Sprintf("SELECT * FROM %s WHERE project_id = ? AND name = ?", (published.Version{}).TableName()), projectId, versionName).Scan(&version)
	}

	if res.Error != nil {
		return published.Version{}, publicApiError.NewError("paginateMapItems", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res.RowsAffected == 0 {
		return published.Version{}, publicApiError.NewError("paginateMapItems", map[string]string{
			"notFound": "This version does not exist",
		}, publicApiError.NotFoundError)
	}

	return version, nil
}

func getGroups(ids []string) (map[string][]string, error) {
	sql := fmt.Sprintf(
		"SELECT g.name, vg.variable_id FROM %s AS g INNER JOIN %s AS vg ON vg.variable_id IN(?) AND g.id = ANY(vg.groups) GROUP BY g.name, vg.variable_id",
		(declarations.Group{}).TableName(),
		(declarations.VariableGroup{}).TableName(),
	)

	type Group struct {
		Name       string `gorm:"column:name"`
		VariableID string `gorm:"column:variable_id"`
	}

	var groups []Group
	res := storage.Gorm().Raw(sql, ids).Scan(&groups)
	if res.Error != nil {
		return nil, publicApiError.NewError("getListItemById", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	results := make(map[string][]string, 0)
	for _, v := range groups {
		if _, ok := results[v.VariableID]; !ok {
			results[v.VariableID] = make([]string, 0)
		}

		results[v.VariableID] = append(results[v.VariableID], v.Name)
	}

	return results, nil
}

func getGroupIdsByName(projectId string, groups []string) ([]string, error) {
	sql := fmt.Sprintf("SELECT id FROM %s WHERE name IN(?) AND project_id = ?", (declarations.Group{}).TableName())

	var groupIds []string
	if res := storage.Gorm().Raw(sql, groups, projectId).Scan(&groupIds); res.Error != nil {
		return nil, publicApiError.NewError("getListItemById", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return groupIds, nil
}
