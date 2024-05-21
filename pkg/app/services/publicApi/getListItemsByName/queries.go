package getListItemsByName

import (
	"creatif/pkg/app/domain/published"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
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

func getItemSql(locale string, options Options) string {
	localeSql := ""
	if locale != "" {
		localeSql = "AND lv.locale_id = @localeId"
	}

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
    v.project_id = @projectId AND 
    v.name = @versionName AND 
    v.id = lv.version_id AND 
    (lv.name = @structureName OR lv.id = @structureName OR 
    lv.short_id = @structureName) AND 
    lv.variable_name = @variableName %s 
`,
		selectFields,
		(published.PublishedList{}).TableName(),
		(published.Version{}).TableName(),
		localeSql,
	)
}

func getConnectionsSql() string {
	return fmt.Sprintf(`
SELECT 
    v.project_id,
    c.name AS connection_name,
    c.child_type AS connection_type,
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
INNER JOIN %s AS v ON v.project_id = ? AND v.name = ? AND v.id = lv.version_id AND lv.variable_id IN(?)
INNER JOIN %s AS c ON c.project_id = ? AND c.project_id = v.project_id AND v.name = ? AND v.id = c.version_id AND c.child_id IN (?)
`,
		(published.PublishedList{}).TableName(),
		(published.Version{}).TableName(),
		(published.PublishedReference{}).TableName(),
	)
}
