package getMap

import (
	"creatif/pkg/app/declarations/getMap/services"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"encoding/json"
	"errors"
	"fmt"
)

type Main struct {
	model GetMapModel
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
	m, err := services.GetMap(c.model.ID)
	if err != nil {
		return LogicModel{}, err
	}

	strategy := services.CreateStrategy(c.model.Return, c.model.Fields)
	models, err := services.Execute(m.ID, strategy)
	if err != nil {
		var convertedErr appErrors.AppError[struct{}]
		errors.As(err, &convertedErr)
		return LogicModel{}, convertedErr.AddError("GetMap.Get.Logic", nil)
	}

	nodes := make([]FullNode, 0)
	for _, model := range models {
		node := FullNode{
			ID:        model.ID,
			Name:      model.Name,
			Type:      model.Type,
			Behaviour: model.Behaviour,
			Groups:    model.Groups,
			Metadata:  model.Metadata,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		}

		if model.Value != nil {
			var conv interface{}
			if err := json.Unmarshal(model.Value, &conv); err != nil {
				fmt.Println("VALUE ERROR")
				return LogicModel{}, appErrors.NewApplicationError(err).AddError("GetMap.Get.Logic", nil)
			}

			node.Value = conv
		}

		nodes = append(nodes, node)
	}

	return LogicModel{
		nodeMap:  m,
		nodes:    nodes,
		strategy: strategy.Name(),
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

func New(model GetMapModel) pkg.Job[GetMapModel, View, LogicModel] {
	return Main{model: model}
}
