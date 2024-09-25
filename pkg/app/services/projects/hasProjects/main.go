package hasProjects

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
)

type Main struct {
	auth auth.Authentication
}

func (c Main) Validate() error {
	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return appErrors.NewAuthenticationError(err)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (bool, error) {
	sql := fmt.Sprintf("SELECT COUNT(id) FROM %s", (app.Project{}).TableName())

	var projectCount int
	res := storage.Gorm().Raw(sql).Scan(&projectCount)
	if res.Error != nil {
		return false, appErrors.NewApplicationError(res.Error)
	}

	if projectCount == 0 {
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

	projectExists, err := c.Logic()

	if err != nil {
		return false, err
	}

	return projectExists, nil
}

func New(auth auth.Authentication) pkg.Job[interface{}, bool, bool] {
	return Main{auth: auth}
}
