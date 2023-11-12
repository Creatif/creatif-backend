package createProject

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
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
	model := app.NewProject(c.model.Name, c.auth.User().ID)

	if err := storage.Create((app.Project{}).TableName(), &model, false); err != nil {
		c.logBuilder.Add("error", err.Error())

		return app.Project{}, appErrors.NewApplicationError(err).AddError("createProject", nil)
	}

	c.logBuilder.Add("projectService", fmt.Sprintf("Project %s created", c.model.Name))
	return model, nil
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
	builder.Add("projectService", "Created")
	return Main{model: model, logBuilder: builder, auth: auth}
}
