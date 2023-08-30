package getBatchNodes

import (
	pkg "creatif/pkg/lib"
)

type Main struct {
	model GetBatchedNodesModel
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

func (c Main) Logic() ([]NodeWithValueQuery, error) {
	return []NodeWithValueQuery{}, nil
}

func (c Main) Handle() (map[string]View, error) {
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

func New(model GetBatchedNodesModel) pkg.Job[GetBatchedNodesModel, map[string]View, []NodeWithValueQuery] {
	return Main{model: model}
}
