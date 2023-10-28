package getVariable

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
}

func (c Main) Validate() error {
	c.logBuilder.Add("getVariable", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	c.logBuilder.Add("getVariable", "Validated")
	return nil
}

func (c Main) Authenticate() error {
	// user check by project id should be gotten here, with authentication cookie
	var project app.Project
	if err := storage.Get((app.Project{}).TableName(), c.model.ProjectID, &project); err != nil {
		return appErrors.NewAuthenticationError(err).AddError("createVariable.Authenticate", nil)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (declarations.Variable, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.LocaleAlpha)
	if err != nil {
		c.logBuilder.Add("getVariable", err.Error())
		return declarations.Variable{}, appErrors.NewNotFoundError(err)
	}

	variable, err := queryValue(c.model.ProjectID, localeID, c.model.Name, c.model.Fields)
	if err != nil {
		c.logBuilder.Add("getVariable", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return declarations.Variable{}, appErrors.NewNotFoundError(err)
		}

		return declarations.Variable{}, appErrors.NewDatabaseError(err)
	}

	return variable, nil
}

func (c Main) Handle() (map[string]interface{}, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	if err := c.Authorize(); err != nil {
		return nil, err
	}

	model, err := c.Logic()

	if err != nil {
		return nil, err
	}

	return newView(model, c.model.Fields, c.model.LocaleAlpha), nil
}

func New(model Model, logBuilder logger.LogBuilder) pkg.Job[Model, map[string]interface{}, declarations.Variable] {
	logBuilder.Add("getVariable", "Created.")
	return Main{model: model, logBuilder: logBuilder}
}
