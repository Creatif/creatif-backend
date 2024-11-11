package getListItemsByName

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
	ID            string `gorm:"column:id"`
	ShortID       string `gorm:"column:short_id"`
	StructureName string `gorm:"column:structure_name"`

	ConnectionName string `gorm:"column:connection_name"`
	ConnectionType string `gorm:"column:connection_type"`

	ProjectID string `gorm:"column:project_id"`

	Name        string         `gorm:"column:variable_name"`
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

func getItem(placeholders map[string]interface{}, options Options) ([]Item, error) {
	localeSql := ""
	if placeholders["localeId"].(string) != "" {
		localeSql = "AND lv.locale_id = @localeId"
	}

	selectFields := fmt.Sprintf(`
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
`)

	if options.ValueOnly {
		selectFields = fmt.Sprintf(`
	lv.value
`)
	}

	sql := fmt.Sprintf(`
SELECT 
    %s
FROM %s AS lv
INNER JOIN %s AS v ON 
    v.project_id = @projectId AND 
    v.id = @versionId AND 
    v.id = lv.version_id AND 
    (lv.name = @structureName OR lv.structure_id = @structureName OR 
    lv.short_id = @structureName) AND 
    lv.variable_name = @variableName %s 
`,
		selectFields,
		(published.PublishedList{}).TableName(),
		(published.Version{}).TableName(),
		localeSql,
	)

	var items []Item
	res := storage.Gorm().Raw(
		sql,
		placeholders,
	).Scan(&items)
	if res.Error != nil {
		return nil, publicApiError.NewError("getListItemsByName", map[string]string{
			"error": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res.RowsAffected == 0 {
		return nil, nil
	}

	return items, nil
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
			fmt.Sprintf("SELECT id, name FROM %s WHERE project_id = ? AND name = ?", (published.Version{}).TableName()), projectId, versionName).Scan(&version)
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
