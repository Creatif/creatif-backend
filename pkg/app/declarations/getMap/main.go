package getMap

import (
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
)

type Main struct {
	model Model
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

func (c Main) Logic() (LogicModel, error) {
	m, err := queryMap(c.model.Name)
	if err != nil {
		return LogicModel{}, appErrors.NewNotFoundError(err).AddError("getMap.Logic", nil)
	}

	var nodes []Node
	if err := queryNodes(m.ID, c.model.Fields, &nodes); err != nil {
		return LogicModel{}, err
	}

	return LogicModel{
		nodeMap: m,
		nodes:   nodes,
	}, nil
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

	return newView(model, c.model.Fields), nil
}

func New(model Model) pkg.Job[Model, View, LogicModel] {
	return Main{model: model}
}
