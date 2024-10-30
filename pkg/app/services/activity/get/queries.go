package get

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/storage"
	"fmt"
)

func getActivities(projectId string) ([]LogicModel, error) {
	sql := fmt.Sprintf("SELECT id, data, created_at FROM %s WHERE project_id = ?", (app.Activity{}).TableName())

	var models []LogicModel
	if res := storage.Gorm().Raw(sql, projectId).Scan(&models); res.Error != nil {
		return nil, res.Error
	}

	return models, nil
}
