package publish

import (
	"creatif/pkg/app/domain/published"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
)

func validateVersionNameExists(projectId, name string) error {
	var version published.Version
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT id FROM %s WHERE project_id = ? AND name = ?", (published.Version{}).TableName()), projectId, name).Scan(&version)

	if res.Error != nil {
		return appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected != 0 {
		return appErrors.NewValidationError(map[string]string{
			"versionExists": fmt.Sprintf("Version with name '%s' already exists.", name),
		})
	}

	return nil
}
