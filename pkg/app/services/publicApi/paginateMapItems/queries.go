package paginateMapItems

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/domain/published"
	"creatif/pkg/app/services/publicApi/publicApiError"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Item struct {
	StructureID   string `gorm:"column:structure_id"`
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

type ConnectionItem struct {
	StructureID      string `gorm:"column:structure_id"`
	Path             string `gorm:"column:path"`
	ParentVariableID string `gorm:"column:parent_variable_id"`
	StructureShortID string `gorm:"column:short_id"`
	StructureName    string `gorm:"column:structure_name"`
	ProjectID        string `gorm:"column:project_id"`

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
		return nil, publicApiError.NewError("paginateMapItems", map[string]string{
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
		return nil, publicApiError.NewError("paginateMapItems", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return groupIds, nil
}

func getItem(placeholders map[string]interface{}, defs defaults, subQrs subQueries) ([]Item, error) {
	sql := fmt.Sprintf(`
SELECT 
    v.project_id,
	lv.structure_id,
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
INNER JOIN %s AS v ON v.project_id = @projectId AND v.id = @versionId AND v.id = lv.version_id 
AND (lv.name = @structureIdentifier OR lv.structure_id = @structureIdentifier OR lv.short_id = @structureIdentifier)
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
		subQrs.search,
		subQrs.groups,
		subQrs.locales,
		subQrs.query,
		subQrs.sortBy,
		defs.orderDirections,
		defs.limit,
	)

	var items []Item
	res := storage.Gorm().Raw(sql, placeholders).Scan(&items)
	if res.Error != nil {
		return nil, publicApiError.NewError("paginateMapItems", map[string]string{
			"error": res.Error.Error(),
		}, publicApiError.ApplicationError)
	}

	if res.RowsAffected == 0 {
		return nil, nil
	}

	return items, nil
}

func getConnectionListVariables(versionId, projectId string, parentVariables []string) ([]ConnectionItem, error) {
	selectFields := fmt.Sprintf(`
    v.project_id,
	c.path AS path,
	lv.structure_id,
	lv.short_id AS structure_short_id,
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
`)

	sql := fmt.Sprintf(`
SELECT 
    %s
FROM %s AS lv
INNER JOIN %s AS v ON v.project_id = ? AND lv.version_id = ? AND v.id = lv.version_id
INNER JOIN %s AS c ON c.project_id = ? AND c.version_id = ? AND c.parent_variable_id IN(?) AND c.child_variable_id = lv.variable_id
`,
		selectFields,
		(published.PublishedList{}).TableName(),
		(published.Version{}).TableName(),
		(published.PublishedConnection{}).TableName(),
	)

	var items []ConnectionItem
	res := storage.Gorm().Raw(sql, projectId, versionId, projectId, versionId, parentVariables).Scan(&items)
	if res.Error != nil {
		return nil, publicApiError.NewError("getMapItemByName", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return items, nil
}

func getConnectionMapVariables(versionId, projectId string, parentVariables []string) ([]ConnectionItem, error) {
	selectFields := fmt.Sprintf(`
    v.project_id,
	c.path AS path,
	c.parent_variable_id AS parent_variable_id,
	lv.structure_id,
	lv.short_id AS structure_short_id,
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
`)

	sql := fmt.Sprintf(`
SELECT 
    %s
FROM %s AS lv
INNER JOIN %s AS v ON v.project_id = ? AND lv.version_id = ? AND v.id = lv.version_id
INNER JOIN %s AS c ON c.project_id = ? AND c.version_id = ? AND c.parent_variable_id IN(?) AND c.child_variable_id = lv.variable_id
`,
		selectFields,
		(published.PublishedMap{}).TableName(),
		(published.Version{}).TableName(),
		(published.PublishedConnection{}).TableName(),
	)

	var items []ConnectionItem
	res := storage.Gorm().Raw(sql, projectId, versionId, projectId, versionId, parentVariables).Scan(&items)
	if res.Error != nil {
		return nil, publicApiError.NewError("getMapItemByName", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return items, nil
}

func getConnectionVariables(versionId, projectId string, parentVariables []string) ([]ConnectionItem, error) {
	listVariables, err := getConnectionListVariables(versionId, projectId, parentVariables)
	if err != nil {
		return nil, err
	}

	mapVariables, err := getConnectionMapVariables(versionId, projectId, parentVariables)
	if err != nil {
		return nil, err
	}

	allItems := make([]ConnectionItem, len(listVariables)+len(mapVariables))
	counter := 0
	for _, v := range listVariables {
		allItems[counter] = v
		counter++
	}

	for _, v := range mapVariables {
		allItems[counter] = v
		counter++
	}

	return allItems, nil
}
