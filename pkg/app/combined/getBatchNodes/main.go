package getBatchNodes

import (
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
)

type Main struct {
	model *GetBatchedNodesModel
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
	nodes := make([]Node, 0)
	if len(c.model.nodeIds) > 0 {
		n, err := queryNodesValue(c.model.nodeIds)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, n...)
	}

	maps := make([]MapNode, 0)
	if len(c.model.mapIds) > 0 {
		n, err := queryMapValues(c.model.mapIds)
		if err != nil {
			return nil, err
		}

		maps = append(maps, n...)
	}

	mapNodes := make(map[string][]Node)
	for _, mapNode := range maps {
		if _, ok := mapNodes[mapNode.Name]; !ok {
			mapNodes[mapNode.Name] = make([]Node, 0)
		}

		mapNodes[mapNode.Name] = append(mapNodes[mapNode.Name], Node{
			ID:        mapNode.ID,
			Name:      mapNode.Name,
			Behaviour: mapNode.Behaviour,
			Groups:    mapNode.Groups,
			Metadata:  mapNode.Metadata,
			Value:     mapNode.Value,
			CreatedAt: mapNode.CreatedAt,
			UpdatedAt: mapNode.UpdatedAt,
		})
	}

	return map[string]interface{}{
		"nodes": nodes,
		"maps":  mapNodes,
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

func New(model *GetBatchedNodesModel) pkg.Job[*GetBatchedNodesModel, map[string]interface{}, map[string]interface{}] {
	return Main{model: model}
}
