package publish

import (
	"context"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/domain/published"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SingleItem struct {
	ID      string `gorm:"type:text;column:id"`
	ShortID string `gorm:"type:text;column:short_id"`
	Name    string `gorm:"type:text;column:name"`

	VariableID      string         `gorm:"type:text;column:variable_id"`
	VariableName    string         `gorm:"type:text;column:variable_name"`
	Behaviour       string         `gorm:"type:text;column:behaviour"`
	Value           datatypes.JSON `gorm:"type:text"`
	Groups          pq.StringArray `gorm:"type:[]text"`
	VariableShortID string         `gorm:"type:text;column:variable_short_id"`
	Locale          string         `gorm:"type:text;column:locale"`
	Index           float64        `gorm:"type:text;column:index"`
}

func publishMaps(tx *gorm.DB, projectId, versionId string, ctx context.Context) error {
	sql := fmt.Sprintf(`
INSERT INTO %s (
    version_id, 
    variable_id, 
    variable_name, 
    behaviour, 
    structure_id,
    value, 
    variable_short_id, 
    locale_id, 
    "index", 
    ID, 
    short_id, 
    name, 
    created_at,
	updated_at
)
SELECT
    '%s' AS version_id,
    lv.id AS variable_id,
    lv.name AS variable_name,
    lv.behaviour AS behaviour,
    lv.map_id AS structure_id,
    COALESCE(lv.value::jsonb, '{}'::jsonb) AS value,
    lv.short_id AS variable_short_id,
    lv.locale_id AS locale_id,
    lv."index" AS "index",
    l.id AS ID,
    l.short_id AS short_id,
    l.name AS name,
    lv.created_at,
    lv.updated_at
FROM %s AS l
INNER JOIN %s AS lv ON l.project_id = ? AND lv.map_id = l.id`,
		(published.PublishedMap{}).TableName(),
		versionId,
		(declarations.Map{}).TableName(),
		(declarations.MapVariable{}).TableName(),
	)

	if res := tx.WithContext(ctx).Exec(sql, projectId); res.Error != nil {
		return res.Error
	}

	return nil
}

func publishLists(tx *gorm.DB, projectId, versionId string, ctx context.Context) error {
	sql := fmt.Sprintf(`
INSERT INTO %s (
    version_id, 
    variable_id, 
    variable_name, 
    behaviour, 
    value, 
    variable_short_id, 
    locale_id, 
    "index", 
    ID, 
    short_id, 
    name, 
    created_at,
    updated_at
)
SELECT
    '%s' AS version_id,
    lv.id AS variable_id,
    lv.name AS variable_name,
    lv.behaviour AS behaviour,
    COALESCE(lv.value::jsonb, '{}'::jsonb) AS value,
    lv.short_id AS variable_short_id,
    lv.locale_id AS locale_id,
    lv."index" AS "index",
    l.id AS ID,
    l.short_id AS short_id,
    l.name AS name,
    lv.created_at,
    lv.updated_at
FROM %s AS l
INNER JOIN %s AS lv ON l.project_id = ? AND lv.list_id = l.id`,
		(published.PublishedList{}).TableName(),
		versionId,
		(declarations.List{}).TableName(),
		(declarations.ListVariable{}).TableName(),
	)

	if res := tx.WithContext(ctx).Exec(sql, projectId); res.Error != nil {
		return res.Error
	}

	return nil
}

func publishFiles(tx *gorm.DB, projectId, versionId string, ctx context.Context) error {
	sql := fmt.Sprintf(`
INSERT INTO %s (
    id, 
    version_id,
    project_id, 
    name, 
    file_name, 
    mime_type, 
    created_at,
    updated_at
)
SELECT
    id,
    '%s' AS version_id,
    project_id,
    name,
    file_name,
    mime_type,
    created_at,
    updated_at
FROM %s AS l WHERE l.project_id = ?`,
		(published.PublishedFile{}).TableName(),
		versionId,
		(declarations.File{}).TableName(),
	)

	if res := tx.WithContext(ctx).Exec(sql, projectId); res.Error != nil {
		return res.Error
	}

	return nil
}
