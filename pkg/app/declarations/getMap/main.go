package create

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
)

type Main struct {
	model GetNodeModel
}

func (c Main) Validate() error {
	return nil
}

func (c Main) Authenticate() error {
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() ([]declarations.Node, error) {
	return nil, nil
}

func (c Main) Handle() (map[string]View, error) {
	if err := c.Validate(); err != nil {
		return map[string]View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return map[string]View{}, err
	}

	if err := c.Authorize(); err != nil {
		return map[string]View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return map[string]View{}, err
	}

	return newView(model), nil
}

func New(model GetNodeModel) pkg.Job[GetNodeModel, map[string]View, []declarations.Node] {
	return Main{model: model}
}
