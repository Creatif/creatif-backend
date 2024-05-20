package createAndDiff

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
)

func getProject(projectId string) (app.Project, error) {
	sql := fmt.Sprintf("SELECT id FROM %s WHERE id = ?", (app.Project{}).TableName())
	var project app.Project
	res := storage.Gorm().Raw(sql, projectId).Scan(&project)
	if res.Error != nil {
		return app.Project{}, appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		return app.Project{}, appErrors.NewApplicationError(res.Error)
	}

	return project, nil
}

func getStructures(projectId string, tableName string, model interface{}) error {
	sql := fmt.Sprintf("SELECT id, short_id, name FROM %s WHERE project_id = ?", tableName)
	res := storage.Gorm().Raw(sql, projectId).Scan(model)
	if res.Error != nil {
		return appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		return nil
	}

	return nil
}

func getProjectMetadata(projectId string) ([]MetadataModel, error) {
	var logicModels []MetadataModel
	res := storage.Gorm().Raw(fmt.Sprintf(`
SELECT 
p.id,
p.name,
p.state,
p.user_id,
m.name AS map_name,
m.id AS map_id,
m.short_id AS map_short_id,
l.name AS list_name,
l.id AS list_id,
l.short_id AS list_short_id
FROM %s AS p
LEFT JOIN %s AS m ON m.project_id = p.id AND m.project_id = ?
LEFT JOIN %s AS l ON l.project_id = p.id AND l.project_id = ?
WHERE p.id = ?
`,
		(app.Project{}).TableName(),
		(declarations.Map{}).TableName(),
		(declarations.List{}).TableName(),
	),
		projectId,
		projectId,
		projectId,
	).Scan(&logicModels)

	if res.Error != nil {
		return nil, appErrors.NewApplicationError(res.Error)
	}

	return logicModels, nil
}
