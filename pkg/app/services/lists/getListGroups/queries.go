package getListGroups

import (
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
	"strings"
)

func getGroups(name, itemId, projectId string) ([]LogicModel, error) {
	sql := fmt.Sprintf(`
SELECT groups FROM %s AS lv 
    INNER JOIN %s AS l ON l.project_id = ? AND lv.list_id = l.id AND (l.id = ? OR l.short_id = ?) AND (lv.id = ? OR lv.short_id = ?)
`, (declarations2.ListVariable{}).TableName(), (declarations2.List{}).TableName())

	var duplicatedModel []LogicModel
	res := storage.Gorm().Raw(sql, projectId, name, name, itemId, itemId).Scan(&duplicatedModel)

	if res.Error != nil && res.RowsAffected == 0 {
		return nil, appErrors.NewNotFoundError(res.Error)
	} else if res.Error != nil && strings.Contains(res.Error.Error(), "cannot accumulate empty arrays") {
		return []LogicModel{}, nil
	} else if res.Error != nil {
		return nil, appErrors.NewApplicationError(res.Error)
	}

	return duplicatedModel, nil
}
