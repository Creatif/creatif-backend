package get

import (
	pkg "creatif/pkg/lib"
)

type Main struct {
	model PaginationModel
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

func (c Main) Logic() (NodeWithValueQuery, error) {
	return NodeWithValueQuery{}, nil
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

func New(model PaginationModel) pkg.Job[PaginationModel, View, NodeWithValueQuery] {
	return Main{model: model}
}
