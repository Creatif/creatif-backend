package get

import (
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
)

type Main struct {
	model GetNodeModel
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

func (c Main) Logic() (Node, error) {
	return queryValue(c.model.ID, c.model.Fields)
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

func New(model GetNodeModel) pkg.Job[GetNodeModel, map[string]interface{}, Node] {
	return Main{model: model}
}
