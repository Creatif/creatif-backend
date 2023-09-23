package getValue

import (
	"creatif/pkg/lib/storage"
	"gorm.io/gorm"
)

func queryValue(projectId, name string) (Variable, error) {
	var variable Variable
	res := storage.Gorm().Raw(`SELECT n.value FROM declarations.variables AS n WHERE n.name = ? AND n.project_id = ?`, name, projectId).Scan(&variable)
	if res.RowsAffected == 0 {
		return Variable{}, gorm.ErrRecordNotFound
	}

	if res.Error != nil {
		return Variable{}, res.Error
	}

	return variable, nil
}
