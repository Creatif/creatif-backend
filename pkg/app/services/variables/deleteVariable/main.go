package deleteVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("deleteVariable", "Validating...")
	c.logBuilder.Add("deleteVariable", "Validated")
	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return err
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

	res := storage.Gorm().Where(fmt.Sprintf("project_id = ? AND locale_id = ? AND (name = ? OR id = ? OR short_id = ?)"), c.model.ProjectID, localeID, c.model.Name, c.model.Name, c.model.Name).Delete(&declarations.Variable{})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.logBuilder.Add("deleteVariable", res.Error.Error())
		return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteVariable.Logic", nil)
	}

	if res.RowsAffected == 0 {
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, interface{}, interface{}] {
	logBuilder.Add("deleteVariable", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
