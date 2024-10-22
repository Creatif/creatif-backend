package createAdmin

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
)

func adminExists() (bool, error) {
	sql := fmt.Sprintf("SELECT id FROM %s WHERE is_admin = true", (app.User{}).TableName())

	var user app.User
	res := storage.Gorm().Raw(sql).Scan(&user)
	if res.Error != nil {
		return false, appErrors.NewAuthorizationError(res.Error)
	}

	if res.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
