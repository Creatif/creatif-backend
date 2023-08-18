package create

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
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

func (c Main) Logic() (declarations.Node, error) {
	var node declarations.Node
	if sdk.IsValidUuid(c.model.ID) {
		if err := storage.Get(node.TableName(), c.model.ID, &node); err != nil {
			return declarations.Node{}, appErrors.NewDatabaseError(err).AddError("Node.Get.Logic", nil)
		}
	} else {
		if err := storage.GetBy(node.TableName(), "name", c.model.ID, &node); err != nil {
			return declarations.Node{}, appErrors.NewDatabaseError(err).AddError("Node.Get.Logic", nil)
		}
	}

	return node, nil
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

func New(model GetNodeModel) pkg.Job[GetNodeModel, View, declarations.Node] {
	return Main{model: model}
}
