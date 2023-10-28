package getValue

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

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

func (c Main) Logic() (Variable, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		return Variable{}, appErrors.NewApplicationError(err)
	}

	value, err := queryValue(c.model.ProjectID, c.model.Name, localeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Variable{}, appErrors.NewNotFoundError(err)
		}

		return Variable{}, appErrors.NewDatabaseError(err)
	}

	return value, nil
}

func (c Main) Handle() (datatypes.JSON, error) {
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

	return newView(model), nil
}

func New(model Model, logBuilder logger.LogBuilder) pkg.Job[Model, datatypes.JSON, Variable] {
	return Main{model: model, logBuilder: logBuilder}
}
