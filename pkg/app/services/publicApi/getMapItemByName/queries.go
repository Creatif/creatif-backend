package getMapItemByName

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

func getItem(placeholders map[string]interface{}) (Item, error) {
	localeSql := ""
	if placeholders["localeId"].(string) != "" {
		localeSql = "AND lv.locale_id = @localeId"
	}

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
INNER JOIN %s AS v ON v.project_id = @projectId AND v.id = @versionId AND v.id = lv.version_id AND lv.name = @structureName AND lv.variable_name = @variableName %s  
`,
		(published.PublishedMap{}).TableName(),
		(published.Version{}).TableName(),
		localeSql,
	)

	var item Item
	res := storage.Gorm().Raw(sql, placeholders).Scan(&item)
	if res.Error != nil {
		return Item{}, publicApiError.NewError("getMapItemByName", map[string]string{
			"error": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res.RowsAffected == 0 {
		return Item{}, publicApiError.NewError("getMapItemByName", map[string]string{
			"notFound": "Item has not been found.",
		}, publicApiError.NotFoundError)
	}

	return item, nil
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
		return published.Version{}, publicApiError.NewError("getListItemById", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res.RowsAffected == 0 {
		return published.Version{}, publicApiError.NewError("getListItemById", map[string]string{
			"notFound": "This list item does not exist",
		}, publicApiError.NotFoundError)
	}

	return version, nil
}

func getGroups(itemId string) ([]string, error) {
	sql := fmt.Sprintf(
		"SELECT g.name FROM %s AS g INNER JOIN %s AS vg ON vg.variable_id = ? AND g.id = ANY(vg.groups)",
		(declarations.Group{}).TableName(),
		(declarations.VariableGroup{}).TableName(),
	)

	var groups []string
	res := storage.Gorm().Raw(sql, itemId).Scan(&groups)
	if res.Error != nil {
		return nil, publicApiError.NewError("getListItemById", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return groups, nil
}

func getConnectionListVariables(versionId, projectId, parentVariableId string) ([]ConnectionItem, error) {
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
INNER JOIN %s AS c ON c.project_id = ? AND c.version_id = ? AND c.parent_variable_id = ? AND c.child_variable_id = lv.variable_id
`,
		selectFields,
		(published.PublishedList{}).TableName(),
		(published.Version{}).TableName(),
		(published.PublishedConnection{}).TableName(),
	)

	var items []ConnectionItem
	res := storage.Gorm().Raw(sql, projectId, versionId, projectId, versionId, parentVariableId).Scan(&items)
	if res.Error != nil {
		return nil, publicApiError.NewError("getMapItemByName", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return items, nil
}

func getConnectionMapVariables(versionId, projectId, parentVariableId string) ([]ConnectionItem, error) {
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
INNER JOIN %s AS c ON c.project_id = ? AND c.version_id = ? AND c.parent_variable_id = ? AND c.child_variable_id = lv.variable_id
`,
		selectFields,
		(published.PublishedMap{}).TableName(),
		(published.Version{}).TableName(),
		(published.PublishedConnection{}).TableName(),
	)

	var items []ConnectionItem
	res := storage.Gorm().Raw(sql, projectId, versionId, projectId, versionId, parentVariableId).Scan(&items)
	if res.Error != nil {
		return nil, publicApiError.NewError("getMapItemByName", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return items, nil
}

func getConnectionVariables(versionId, projectId, parentVariableId string) ([]ConnectionItem, error) {
	listVariables, err := getConnectionListVariables(versionId, projectId, parentVariableId)
	if err != nil {
		return nil, err
	}

	mapVariables, err := getConnectionMapVariables(versionId, projectId, parentVariableId)
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
