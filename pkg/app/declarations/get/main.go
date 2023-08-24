package get

import (
	"creatif/pkg/app/declarations/get/services"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"encoding/json"
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

func (c Main) Logic() (NodeWithValueQuery, error) {
	serviceModel, err := services.NewGetService(c.model.ID).GetNode(func(id string) (declarations.Node, error) {
		var node declarations.Node
		if err := storage.Get((&declarations.Node{}).TableName(), c.model.ID, &node, "ID"); err != nil {
			return declarations.Node{}, appErrors.NewDatabaseError(err).AddError("Node.Get.Logic", nil)
		}

		return node, nil
	}, func(name string) (declarations.Node, error) {
		var node declarations.Node

		if err := storage.GetBy((&declarations.Node{}).TableName(), "name", c.model.ID, &node, "ID"); err != nil {
			return declarations.Node{}, appErrors.NewDatabaseError(err).AddError("Node.Get.Logic", nil)
		}

		return node, nil
	})

	var v interface{}
	if serviceModel.Value != nil {
		if err := json.Unmarshal(serviceModel.Value, &v); err != nil {
			return NodeWithValueQuery{}, appErrors.NewDatabaseError(err).AddError("Node.Get.Logic", nil)
		}
	} else {
		v = serviceModel.Value
	}

	return NodeWithValueQuery{
		ID:        serviceModel.ID,
		Name:      serviceModel.Name,
		Type:      serviceModel.Type,
		Behaviour: serviceModel.Behaviour,
		Groups:    serviceModel.Groups,
		Metadata:  serviceModel.Metadata,
		Value:     v,
		CreatedAt: serviceModel.CreatedAt,
		UpdatedAt: serviceModel.UpdatedAt,
	}, err
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

func New(model GetNodeModel) pkg.Job[GetNodeModel, View, NodeWithValueQuery] {
	return Main{model: model}
}
