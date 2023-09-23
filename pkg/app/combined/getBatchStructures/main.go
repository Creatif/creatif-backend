package getBatchStructures

import (
	"creatif/pkg/app/domain/app"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
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
	// user check by project id should be gotten here, with authentication cookie
	var project app.Project
	if err := storage.Get((app.Project{}).TableName(), c.model.ProjectID, &project); err != nil {
		return appErrors.NewAuthenticationError(err).AddError("createVariable.Authenticate", nil)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (map[string]interface{}, error) {
	variables := make([]Variable, 0)
	if len(c.model.variableIds) > 0 {
		n, err := queryVariables(c.model.variableIds)
		if err != nil {
			return nil, err
		}

		variables = append(variables, n...)
	}

	queriedMaps := make([]QueriesMapVariable, 0)
	if len(c.model.mapIds) > 0 {
		if err := queryMaps(c.model.mapIds, &queriedMaps); err != nil {
			return nil, err
		}
	}

	maps := make(map[string][]Variable)
	for _, mapVariable := range queriedMaps {
		mapName := mapVariable.MapName
		if _, ok := maps[mapName]; !ok {
			maps[mapName] = make([]Variable, 0)
		}

		maps[mapName] = append(maps[mapName], Variable{
			ID:        mapVariable.ID,
			ProjectID: mapVariable.ProjectID,
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
		"maps":      maps,
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
