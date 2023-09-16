package getValue

import (
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"gorm.io/datatypes"
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

func (c Main) Logic() (Variable, error) {
	return queryValue(c.model.Name)
}

func (c Main) Handle() (datatypes.JSON, error) {
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

func New(model Model) pkg.Job[Model, datatypes.JSON, Variable] {
	return Main{model: model}
}
