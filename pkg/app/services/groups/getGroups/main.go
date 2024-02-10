package getGroups

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("getGroups", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	c.logBuilder.Add("getGroups", "Validated.")

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

func (c Main) Logic() ([]declarations.Group, error) {
	var groups []declarations.Group
	if res := storage.Gorm().Raw(fmt.Sprintf("SELECT name FROM %s WHERE project_id = ?", (declarations.Group{}).TableName()), c.model.ProjectID).Scan(&groups); res.Error != nil {
		return []declarations.Group{}, appErrors.NewApplicationError(res.Error)
	}

	return groups, nil
}

func (c Main) Handle() ([]string, error) {
	if err := c.Validate(); err != nil {
		return []string{}, err
	}

	if err := c.Authenticate(); err != nil {
		return []string{}, err
	}

	if err := c.Authorize(); err != nil {
		return []string{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return []string{}, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, []string, []declarations.Group] {
	logBuilder.Add("getGroups", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
