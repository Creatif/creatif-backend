package dashboard

import "C"
import (
	"creatif/pkg/app/auth"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
)

type Results struct {
	Errors []error
}

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
	lists, err := getLists()
	if err != nil {
		return nil, appErrors.NewApplicationError(err)
	}

	maps, err := getMaps()
	if err != nil {
		return nil, appErrors.NewApplicationError(err)
	}

	lists = append(lists, maps...)

	return lists, nil
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
