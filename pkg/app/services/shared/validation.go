package shared

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
)

func ValidateGroupsExist(projectId string, groups []string) (int, error) {
	var count int
	if res := storage.Gorm().Raw(fmt.Sprintf("SELECT count(id) FROM %s WHERE project_id = ? AND id IN(?)", (declarations.Group{}).TableName()), projectId, groups).Scan(&count); res.Error != nil {
		return 0, res.Error
	}

	if count != len(groups) {
		return 0, errors.New("Invalid groups. Some of the groups provided do not exist")
	}

	return count, nil
}
