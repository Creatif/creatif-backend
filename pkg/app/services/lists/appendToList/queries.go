package appendToList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func getHighestIndex(listId string) (float64, error) {
	var idx float64
	res := storage.Gorm().Raw(fmt.Sprintf(`
	SELECT index FROM %s WHERE list_id = ? ORDER BY index DESC LIMIT 1
`, (declarations.ListVariable{}).TableName()),
		listId,
	).Scan(&idx)

	if res.Error != nil {
		return idx, res.Error
	}

	return idx, nil
}

func getList(name string) (declarations.List, error) {
	var list declarations.List
	if res := storage.Gorm().Where("id = ? OR short_id = ?", name, name).Select("id", "serial", "project_id", "name", "created_at", "updated_at").First(&list); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.List{}, appErrors.NewNotFoundError(res.Error).AddError("appendToList.Logic", nil)
		}

		return declarations.List{}, appErrors.NewDatabaseError(res.Error).AddError("appendToList.Logic", nil)
	}

	return list, nil
}
