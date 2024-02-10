package addToList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
)

func validateGroupsExist(projectId string, groups []string) error {
	g := make([]declarations.Group, 0)
	if res := storage.Gorm().Raw(fmt.Sprintf("SELECT count(id) FROM %s WHERE project_id = ? AND name IN(?)", (declarations.Group{}).TableName()), projectId, groups).Scan(&g); res.Error != nil {
		return res.Error
	}

	if len(g) != len(groups) {
		return errors.New("Invalid groups. Some of the groups provided do not exist")
	}

	return nil
}
