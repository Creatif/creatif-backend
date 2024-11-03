package getManyItems

import (
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

func getItemsSql(table string, options Options) string {
	selectFields := fmt.Sprintf(`
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
`)

	if options.ValueOnly {
		selectFields = fmt.Sprintf(`
	lv.value
`)
	}

	return fmt.Sprintf(`
SELECT 
    %s
FROM %s AS lv
INNER JOIN %s AS v ON 
    v.project_id = ? AND 
    v.id = ? AND 
	lv.variable_id IN(?)
`,
		selectFields,
		table,
		(published.Version{}).TableName(),
	)
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

func getItemsConcurrently(projectId, versionId string, ids []string, options Options) ([]Item, error) {
	type result struct {
		items []Item
		error error
	}
	ch := make(chan result)
	tables := []string{(published.PublishedList{}).TableName(), (published.PublishedMap{}).TableName()}

	for _, t := range tables {
		go func(table string) {
			var items []Item
			res := storage.Gorm().Raw(
				getItemsSql(table, options),
				projectId, versionId, ids,
			).Scan(&items)
			if res.Error != nil {
				ch <- result{
					items: nil,
					error: publicApiError.NewError("getManyItems", map[string]string{
						"error": res.Error.Error(),
					}, publicApiError.DatabaseError),
				}

				return
			}

			ch <- result{
				items: items,
				error: nil,
			}
		}(t)
	}

	results := []result{<-ch, <-ch}
	items := make([]Item, 0)
	for _, res := range results {
		if res.error != nil {
			return nil, res.error
		}

		for _, t := range res.items {
			items = append(items, t)
		}
	}

	close(ch)

	return items, nil
}
