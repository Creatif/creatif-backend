package getValue

import (
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
)

func queryValue(projectId, id, value, localeID string) (Variable, error) {
	var variable Variable
	res := storage.Gorm().Raw(fmt.Sprintf(`SELECT n.value FROM declarations.variables AS n WHERE %s AND n.project_id = ? AND locale_id = ?`, id), value, projectId, localeID).Scan(&variable)
	if res.RowsAffected == 0 {
		return Variable{}, gorm.ErrRecordNotFound
	}

	if res.Error != nil {
		return Variable{}, res.Error
	}

	return variable, nil
}
