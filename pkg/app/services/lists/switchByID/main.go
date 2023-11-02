package switchByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("switchByID", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("replaceListItem", "Validated")

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

func (c Main) Logic() (LogicResult, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("switchByID", err.Error())
		return LogicResult{}, appErrors.NewApplicationError(err)
	}
	source, destination, err := tryUpdates(c.model.ProjectID, localeID, c.model.Name, c.model.ID, c.model.ShortID, c.model.Source, c.model.Destination, 0, 10)
	if err != nil {
		c.logBuilder.Add("switchByID", err.Error())
		return LogicResult{}, appErrors.NewDatabaseError(err)
	}

	return LogicResult{
		To:     source,
		From:   destination,
		Locale: c.model.Locale,
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

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, LogicResult] {
	logBuilder.Add("switchByID", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
