package getMap

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"errors"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("getMap", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	c.logBuilder.Add("getMap", "Validated.")
	return nil
}

func (c Main) Authenticate() error {
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (LogicModel, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("getMap", err.Error())
		return LogicModel{}, appErrors.NewApplicationError(err).AddError("getMap.Logic", nil)
	}

	m, err := queryMap(c.model.ProjectID, c.model.Name, localeID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.logBuilder.Add("getMap", err.Error())
		return LogicModel{}, appErrors.NewNotFoundError(err).AddError("getMap.Logic", nil)
	}

	if err != nil {
		c.logBuilder.Add("getMap", err.Error())
		return LogicModel{}, appErrors.NewDatabaseError(err).AddError("getMap.Logic", nil)
	}

	var variables []Variable
	if err := queryVariables(m.ID, localeID, c.model.Fields, c.model.Groups, &variables); err != nil {
		c.logBuilder.Add("getMap", err.Error())
		return LogicModel{}, err
	}

	return LogicModel{
		variableMap: m,
		variables:   variables,
	}, nil
}

func (c Main) Handle() (View, error) {
	if err := c.Validate(); err != nil {
		return View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return View{}, err
	}

	if err := c.Authorize(); err != nil {
		return View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return View{}, err
	}

	return newView(model, c.model.Fields, c.model.Locale), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, LogicModel] {
	logBuilder.Add("getMap", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
