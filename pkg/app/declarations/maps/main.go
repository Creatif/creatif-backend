package maps

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"gorm.io/gorm"
)

type Main struct {
	model CreateMapModel
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

func (c Main) Logic() (LogicResult, error) {
	var nodes []declarations.Node

	if res := storage.Gorm().Select("ID").Where("ID IN (?)", c.model.Nodes).Find(&nodes); res.Error != nil {
		return LogicResult{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Create.Logic", nil)
	}

	if len(nodes) != len(c.model.Nodes) {
		return LogicResult{}, appErrors.NewValidationError(map[string]string{
			"validNum": "Found invalid number of nodes. Some of the nodes you provided do not exist.",
		})
	}

	if !sdk.Every(nodes, func(idx int, value declarations.Node) bool {
		return sdk.Includes(c.model.Nodes, value.ID)
	}) {
		return LogicResult{}, appErrors.NewValidationError(map[string]string{
			"validNum": "Found invalid number of nodes. Some of the nodes you provided do not exist.",
		})
	}

	mapNodes := sdk.Map(nodes, func(idx int, value declarations.Node) *declarations.MapNode {
		m := declarations.NewMapNode(value.ID)
		return &m
	})

	m := declarations.NewMap(c.model.Name)
	if err := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&m); res.Error != nil {
			return res.Error
		}

		for _, mapNode := range mapNodes {
			mapNode.MapID = m.ID
		}

		if res := tx.Create(&mapNodes); res.Error != nil {
			return res.Error
		}

		assignmentNode := assignments.NewMapNode(m.Name, m.ID)
		if res := tx.Create(&assignmentNode); res.Error != nil {
			return res.Error
		}

		assignmentNodeValue := assignments.NewMapValueNode(assignmentNode.ID, nil)
		if res := tx.Create(&assignmentNodeValue); res.Error != nil {
			return res.Error
		}

		return nil
	}); err != nil {
		return LogicResult{}, appErrors.NewDatabaseError(err).AddError("Node.Create.Logic", nil)
	}

	return LogicResult{
		ID: m.ID,
		Nodes: sdk.Map(mapNodes, func(idx int, value *declarations.MapNode) string {
			return value.ID
		}),
		Name: m.Name,
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

	return newView(model), nil
}

func New(model CreateMapModel) pkg.Job[CreateMapModel, View, LogicResult] {
	return Main{model: model}
}
