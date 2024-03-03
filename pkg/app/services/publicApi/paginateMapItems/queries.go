package paginateMapItems

import (
	"creatif/pkg/app/domain/published"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
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

func getItemSql(structureIdentifier string, page int, order, sortBy, search string, lcls, groups []string) (string, map[string]interface{}) {
	offset := (page - 1) * 100
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
ORDER BY %s %s
OFFSET @offset
LIMIT 100
`,
		(published.PublishedMap{}).TableName(),
		(published.Version{}).TableName(),
		searchSql,
		groupsSql,
		localesSql,
		sortBy,
		order,
	), placeholders
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
		(published.PublishedMap{}).TableName(),
		(published.Version{}).TableName(),
		(published.PublishedReference{}).TableName(),
	)
}
