package publish

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/domain/published"
	"fmt"
)

func getSelectListSql() string {
	return fmt.Sprintf(`
SELECT 
lv.id AS variableId,
lv.name AS variableName,
lv.behaviour AS behaviour,
lv.value AS value,
lv.short_id AS variableShortId,
lv.locale_id AS locale,
lv.index AS index,
l.id AS ID,
l.short_id AS shortId,
l.name AS name,
(SELECT g.groups FROM %s AS g WHERE lv.id = g.variable_id LIMIT 1) AS groups
FROM %s AS l 
INNER JOIN %s AS lv ON l.project_id = ? AND lv.list_id = l.id
`,
		(declarations.VariableGroup{}).TableName(),
		(declarations.List{}).TableName(),
		(declarations.ListVariable{}).TableName())
}

func getSelectMapSql() string {
	return fmt.Sprintf(`
SELECT 
lv.id AS variableId,
lv.name AS variableName,
lv.behaviour AS behaviour,
lv.value AS value,
lv.short_id AS variableShortId,
lv.locale_id AS locale,
lv.index AS index,
l.id AS ID,
l.short_id AS shortId,
l.name AS name,
(SELECT g.groups FROM %s AS g WHERE lv.id = g.variable_id LIMIT 1) AS groups
FROM %s AS l
INNER JOIN %s AS lv ON l.project_id = ? AND lv.map_id = l.id
`,
		(declarations.VariableGroup{}).TableName(),
		(declarations.Map{}).TableName(),
		(declarations.MapVariable{}).TableName())
}

func getReferencesSql() string {
	return fmt.Sprintf(`
SELECT 
lv.id AS id,
lv.project_id AS projectId,
lv.name AS name,
lv.parent_type AS parentType,
lv.child_type AS childType,
lv.parent_structure_id AS parentStructureId,
lv.child_structure_id AS childStructureId,
lv.parent_id AS parentId,
lv.child_id AS childId
FROM %s AS lv WHERE project_id = ?
`,
		(declarations.Reference{}).TableName())
}

func getReferenceMergeSql(versionId, selectSql string) string {
	return fmt.Sprintf(`
	MERGE INTO %s AS p
	USING (%s) AS t
	ON p.project_id = t.projectId
	WHEN NOT MATCHED THEN
        INSERT (
			id, 
			project_id, 
			version_id, 
			name, 
			parent_type, 
			child_type, 
			parent_structure_id, 
			child_structure_id, 
			parent_id, 
			child_id
		) VALUES (
			t.id, 
			t.projectId,
			'%s', 
			t.name, 
			t.parentType, 
			t.childType, 
			t.parentStructureId, 
			t.childStructureId, 
			t.parentId, 
			t.childId
		)
`,
		(published.PublishedReference{}).TableName(),
		selectSql,
		versionId,
	)
}

func getMergeSql(versionId, tableName, selectSql string) string {
	return fmt.Sprintf(`
	MERGE INTO %s AS p
	USING (%s) AS t
	ON p.variable_id != t.variableId
	WHEN NOT MATCHED THEN
        INSERT (
			id, 
			short_id, 
			version_id, 
			name, 
			variable_name, 
			variable_id, 
			variable_short_id, 
			index, 
			behaviour, 
			value, 
			locale_id, 
			groups
		) VALUES (
			t.ID, 
			t.shortId, 
			'%s', 
			t.name, 
			t.variableName, 
			t.variableId, 
			t.variableShortId, 
			t.index, 
			t.behaviour, 
			t.value, 
			t.locale, 
			t.groups
		)
`,
		tableName,
		selectSql,
		versionId,
	)
}
