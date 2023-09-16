package getBatchData

import (
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
)

type Main struct {
	model *Model
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

func (c Main) Logic() (map[string]interface{}, error) {
	variables := make([]Variable, 0)
	if len(c.model.variableIds) > 0 {
		n, err := queryVariableValue(c.model.variableIds)
		if err != nil {
			return nil, err
		}

		variables = append(variables, n...)
	}

	maps := make([]QueriesMapVariable, 0)
	if len(c.model.mapIds) > 0 {
		if err := queryMapVariables(c.model.mapIds, &maps); err != nil {
			return nil, err
		}
	}

	mapVariables := make(map[string][]Variable)
	for _, mapVariable := range maps {
		if _, ok := mapVariables[mapVariable.Name]; !ok {
			mapVariables[mapVariable.Name] = make([]Variable, 0)
		}

		mapVariables[mapVariable.Name] = append(mapVariables[mapVariable.Name], Variable{
			ID:        mapVariable.ID,
			Name:      mapVariable.Name,
			Behaviour: mapVariable.Behaviour,
			Groups:    mapVariable.Groups,
			Metadata:  mapVariable.Metadata,
			Value:     mapVariable.Value,
			CreatedAt: mapVariable.CreatedAt,
			UpdatedAt: mapVariable.UpdatedAt,
		})
	}

	return map[string]interface{}{
		"variables": variables,
		"maps":      mapVariables,
	}, nil
}

func (c Main) Handle() (map[string]interface{}, error) {
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

func New(model *Model) pkg.Job[*Model, map[string]interface{}, map[string]interface{}] {
	return Main{model: model}
}
