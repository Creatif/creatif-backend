package getVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func queryValue(name string, fields []string) (declarations.Variable, error) {
	resolvedFields := strings.Join(sdk.Map(fields, func(idx int, value string) string {
		return fmt.Sprintf("n.%s", value)
	}), ",")

	var variable declarations.Variable
	if res := storage.Gorm().Raw(fmt.Sprintf(`SELECT n.id, n.name, %s FROM declarations.variables AS n WHERE n.name = ?`, resolvedFields), name).Scan(&variable); errors.Is(res.Error, gorm.ErrRecordNotFound) || res.RowsAffected == 0 {
		return declarations.Variable{}, gorm.ErrRecordNotFound
	}

	return variable, nil
}
