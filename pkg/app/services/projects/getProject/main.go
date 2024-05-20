package getProject

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
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
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("projectService", "Validating...")

	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("projectService", "Validated")

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

func (c Main) Logic() (app.Project, error) {
	id := c.model.ProjectID

	var project app.Project
	if res := storage.Gorm().Where("id = ? OR name = ?", id, id).First(&project); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return app.Project{}, appErrors.NewNotFoundError(res.Error)
		}

		c.logBuilder.Add("getProject.databaseError", res.Error.Error())
		return app.Project{}, appErrors.NewDatabaseError(res.Error)
	}

	return project, nil
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

func New(model Model, auth auth.Authentication, builder logger.LogBuilder) pkg.Job[Model, View, app.Project] {
	builder.Add("projectService", "Get project")
	return Main{model: model, logBuilder: builder, auth: auth}
}
