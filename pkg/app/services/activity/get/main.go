package get

import (
	"creatif/pkg/app/auth"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
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

func (c Main) Logic() ([]LogicModel, error) {
	activities, err := getActivities(c.model.ProjectID)
	if err != nil {
		return nil, appErrors.NewApplicationError(err)
	}

	return activities, nil
}

func (c Main) Handle() ([]View, error) {
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, []View, []LogicModel] {
	return Main{model: model, auth: auth}
}
