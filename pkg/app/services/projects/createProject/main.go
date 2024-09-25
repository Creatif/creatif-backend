package createProject

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

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
		return app.Project{}, appErrors.NewApplicationError(err).AddError("createProject", nil)
	}

	if err := createProjectInPublicDirectory(model.ID); err != nil {
		return app.Project{}, err
	}

	if err := createProjectInAssetsDirectory(model.ID); err != nil {
		return app.Project{}, err
	}

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

func New(model Model, auth auth.Authentication) pkg.Job[Model, View, app.Project] {
	return Main{model: model, auth: auth}
}
