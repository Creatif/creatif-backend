package getVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func queryValue(projectId, localeID, name string, fields []string) (declarations.Variable, error) {
	resolvedFields := strings.Join(sdk.Map(fields, func(idx int, value string) string {
		return fmt.Sprintf("n.%s", value)
	}), ",")

	var variable declarations.Variable
	res := storage.Gorm().Raw(fmt.Sprintf(`SELECT n.id, n.name, n.project_id, n.locale_id, n.short_id, %s FROM declarations.variables AS n WHERE (n.name = ? OR n.short_id = ? OR n.id = ?) AND n.project_id = ? AND locale_id = ?`, resolvedFields), name, name, name, projectId, localeID).Scan(&variable)

	if res.RowsAffected == 0 {
		return declarations.Variable{}, gorm.ErrRecordNotFound
	}

	if res.Error != nil {
		return declarations.Variable{}, res.Error
	}

	return variable, nil
}
