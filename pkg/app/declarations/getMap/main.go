package getMap

import (
	"creatif/pkg/app/declarations/getMap/services"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"errors"
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

	models, err := services.Execute(m.ID, services.CreateStrategy(c.model.Return, c.model.Fields))
	if err != nil {
		var convertedErr appErrors.AppError[struct{}]
		errors.As(err, &convertedErr)
		return LogicModel{}, convertedErr.AddError("GetMap.Get.Logic", nil)
	}

	for key, val := range models {
		if found := sdk.Search(val, "value"); found != nil {
			v := []byte(found.(string))
			var conv interface{}
			if err := json.Unmarshal(v, &conv); err != nil {
				return LogicModel{}, appErrors.NewApplicationError(err).AddError("GetMap.Get.Logic", nil)
			}

			models[key]["value"] = conv
		}

		if found := sdk.Search(val, "groups"); found != nil {
			v := []byte(found.(string))
			var conv interface{}
			if err := json.Unmarshal(v, &conv); err != nil {
				return LogicModel{}, appErrors.NewApplicationError(err).AddError("GetMap.Get.Logic", nil)
			}

			models[key]["groups"] = conv
		}
	}

	return LogicModel{
		nodeMap: m,
		nodes:   models,
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
