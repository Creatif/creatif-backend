package getProjectMetadata

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
)

func getVariablesMetadata(projectId, userId string) ([]LogicModel, error) {
	var logicModels []LogicModel
	res := storage.Gorm().Raw(fmt.Sprintf(`
SELECT 
p.id,
p.name,
p.state,
p.user_id,
v.name AS variable_name,
v.id AS variable_id,
v.short_id AS variable_short_id,
m.name AS map_name,
m.id AS map_id,
m.short_id AS map_short_id,
l.name AS list_name,
l.id AS list_id,
l.short_id AS list_short_id,
v.locale_id AS variable_locale
FROM %s AS p
LEFT JOIN %s AS v ON p.id = ? AND p.user_id = ? AND v.project_id = p.id AND v.project_id = ?
LEFT JOIN %s AS m ON m.project_id = p.id AND m.project_id = ?
LEFT JOIN %s AS l ON l.project_id = p.id AND l.project_id = ?
WHERE p.id = ? AND p.user_id = ?
`,
		(app.Project{}).TableName(),
		(declarations.Variable{}).TableName(),
		(declarations.Map{}).TableName(),
		(declarations.List{}).TableName(),
	),
		projectId,
		userId,
		projectId,
		projectId,
		projectId,
		projectId,
		userId,
	).Scan(&logicModels)

	if res.Error != nil {
		return nil, appErrors.NewNotFoundError(res.Error)
	}

	return logicModels, nil
}
