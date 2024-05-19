package adminExists

import (
	"creatif/pkg/app/domain/app"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
)

type Main struct {
	logBuilder logger.LogBuilder
}

func (c Main) Validate() error {
	return nil
}

func (c Main) Authenticate() error {
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (bool, error) {
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

func (c Main) Handle() (bool, error) {
	if err := c.Validate(); err != nil {
		return false, err
	}

	if err := c.Authenticate(); err != nil {
		return false, err
	}

	if err := c.Authorize(); err != nil {
		return false, err
	}

	adminExists, err := c.Logic()

	if err != nil {
		return false, err
	}

	return adminExists, nil
}

func New(logBuilder logger.LogBuilder) pkg.Job[interface{}, bool, bool] {
	logBuilder.Add("createAdmin", "Created")
	return Main{logBuilder: logBuilder}
}
