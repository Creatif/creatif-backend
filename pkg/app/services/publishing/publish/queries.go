package publish

import (
	"creatif/pkg/app/domain/declarations"
	"fmt"
)

func getSelectListSql(projectId string) string {
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
INNER JOIN %s AS lv ON l.project_id = ?
`,
		(declarations.VariableGroup{}).TableName(),
		(declarations.List{}).TableName(),
		(declarations.ListVariable{}).TableName())
}

func getSelectMapSql(projectId string) string {
	return fmt.Sprintf(`
SELECT 
lv.id AS variableId,
lv.name AS variableName,
lv.behaviour AS behaviour,
lv.value AS value,
lv.shortId AS variableShortId,
lv.locale_id AS locale,
lv.index AS index,
l.id AS ID,
l.name AS name,
(SELECT g.groups FROM %s AS g WHERE lv.id = g.variable_id LIMIT 1) AS groups
FROM %s AS l
INNER JOIN %s AS lv ON l.project_id = ?
`,
		(declarations.VariableGroup{}).TableName(),
		(declarations.Map{}).TableName(),
		(declarations.MapVariable{}).TableName())
}
