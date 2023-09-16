package getValue

import (
	"creatif/pkg/lib/storage"
)

func queryValue(name string) (Variable, error) {
	var variable Variable
	if res := storage.Gorm().Raw(`SELECT n.value FROM declarations.variables AS n WHERE n.name = ?`, name).Scan(&variable); res.Error != nil {
		return Variable{}, res.Error
	}

	return variable, nil
}
