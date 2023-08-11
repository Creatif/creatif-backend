package create

import (
	"creatif/pkg/app/domain"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
)

type Create struct {
	model CreateNodeModel
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

func (c Create) Logic() (domain.Node, error) {
	model := domain.NewNode(c.model.Name, c.model.Type, c.model.Group, c.model.Behaviour)

	if err := storage.Create(model.TableName(), &model); err != nil {
		return domain.Node{}, appErrors.NewDatabaseError(err).AddError("Node.Create.Logic", nil)
	}

	return model, nil
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

	model, err := c.Logic()

	if err != nil {
		return View{}, err
	}

	return newView(model), nil
}

func New(model CreateNodeModel) pkg.Job[CreateNodeModel, View, domain.Node] {
	return Create{model: model}
}
