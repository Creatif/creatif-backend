package getMap

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"errors"
	"gorm.io/gorm"
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
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (declarations.Map, error) {
	m, err := queryMap(c.model.ProjectID, c.model.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return declarations.Map{}, appErrors.NewNotFoundError(err).AddError("getMap.Logic", nil)
	}

	if err != nil {
		return declarations.Map{}, appErrors.NewDatabaseError(err).AddError("getMap.Logic", nil)
	}

	return m, nil
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, View, declarations.Map] {
	return Main{model: model, auth: auth}
}
