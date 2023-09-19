package getVariable

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"errors"
	"gorm.io/gorm"
)

type Main struct {
	model Model
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

func (c Main) Logic() (declarations.Variable, error) {
	variable, err := queryValue(c.model.Name, c.model.Fields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return declarations.Variable{}, appErrors.NewNotFoundError(err)
		}

		return declarations.Variable{}, appErrors.NewDatabaseError(err)
	}

	return variable, nil
}

func (c Main) Handle() (map[string]interface{}, error) {
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

	return newView(model, c.model.Fields), nil
}

func New(model Model) pkg.Job[Model, map[string]interface{}, declarations.Variable] {
	return Main{model: model}
}
