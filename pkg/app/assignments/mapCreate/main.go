package mapCreate

import (
	pkg "creatif/pkg/lib"
)

type Create struct {
	model *CreateNodeModel
}

func (c Create) Validate() error {
	return nil
}

func (c Create) Authenticate() error {
	return nil
}

func (c Create) Authorize() error {
	return nil
}

func (c Create) Logic() (interface{}, error) {
	return nil, nil
}

func (c Create) Handle() (View, error) {
	if err := c.Validate(); err != nil {
		return View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return View{}, err
	}

	if err := c.Authorize(); err != nil {
		return View{}, err
	}

	_, err := c.Logic()

	if err != nil {
		return View{}, err
	}

	return newView(), nil
}

func New(model *CreateNodeModel) pkg.Job[*CreateNodeModel, View, interface{}] {
	return Create{model: model}
}
