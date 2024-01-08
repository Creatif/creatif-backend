package createList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func validateListWithNameExists(name, projectId string) error {
	var variable declarations.List
	res := storage.Gorm().Where("name = ? AND project_id = ?", name, projectId).Select("ID").First(&variable)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if res.Error != nil {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Record with name '%s' already exists", name),
		})
	}

	if variable.ID != "" {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Record with name '%s' already exists", name),
		})
	}

	return nil
}
