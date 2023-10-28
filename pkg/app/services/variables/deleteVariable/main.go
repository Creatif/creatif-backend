package deleteVariable

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
	c.logBuilder.Add("deleteVariable", "Validating...")
	c.logBuilder.Add("deleteVariable", "Validated")
	return nil
}

func (c Main) Authenticate() error {
	c.logBuilder.Add("deleteVariable", "Validating...")
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

func (c Main) Logic() (interface{}, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("deleteVariable", err.Error())
		return declarations.Variable{}, appErrors.NewNotFoundError(err)
	}

	res := storage.Gorm().Where("name = ? AND project_id = ? AND locale_id = ?", c.model.Name, c.model.ProjectID, localeID).Delete(&declarations.Variable{})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.logBuilder.Add("deleteVariable", res.Error.Error())
		return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteVariable.Logic", nil)
	}

	if res.Error != nil {
		c.logBuilder.Add("deleteVariable", res.Error.Error())
		return nil, appErrors.NewDatabaseError(res.Error).AddError("deleteVariable.Logic", nil)
	}

	return nil, nil
}

func (c Main) Handle() (interface{}, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	if err := c.Authorize(); err != nil {
		return nil, err
	}

	_, err := c.Logic()

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func New(model Model, logBuilder logger.LogBuilder) pkg.Job[Model, interface{}, interface{}] {
	logBuilder.Add("deleteVariable", "Created")
	return Main{model: model, logBuilder: logBuilder}
}
