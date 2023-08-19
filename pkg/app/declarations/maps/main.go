package create

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
)

type Main struct {
	model CreateMapModel
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

func (c Main) Logic() ([]string, error) {
	var nodes []declarations.Node

	if res := storage.Gorm().Select("ID").Where("ID IN (?)", c.model.Names).Find(&nodes); res.Error != nil {
		return []string{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Create.Logic", nil)
	}

	if len(nodes) != len(c.model.Names) {
		return []string{}, appErrors.NewValidationError(map[string]string{
			"validNum": "Found invalid number of nodes. Some of the nodes you provided do not exist.",
		})
	}

	if !sdk.Every(nodes, func(idx int, value declarations.Node) bool {
		return sdk.Includes(c.model.Names, value.ID)
	}) {
		return []string{}, appErrors.NewValidationError(map[string]string{
			"validNum": "Found invalid number of nodes. Some of the nodes you provided do not exist.",
		})
	}

	maps := sdk.Map(nodes, func(idx int, value declarations.Node) declarations.Map {
		return declarations.NewMap(value.ID)
	})

	if err := storage.Create((declarations.Map{}).TableName(), &maps, false); err != nil {
		return []string{}, appErrors.NewDatabaseError(err).AddError("Node.Create.Logic", nil)
	}

	return c.model.Names, nil
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

func New(model CreateMapModel) pkg.Job[CreateMapModel, View, []string] {
	return Main{model: model}
}
