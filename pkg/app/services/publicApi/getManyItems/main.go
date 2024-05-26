package getManyItems

import (
	"creatif/pkg/app/auth"
	connections2 "creatif/pkg/app/services/publicApi/connections"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("getListItemsByName", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("getMapItemsByName", errs, publicApiError.ValidationError)
	}

	c.logBuilder.Add("getListItemsByName", "Validated")
	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return publicApiError.NewError("getListItemsByName", map[string]string{
			"unauthorized": "You are unauthorized to use this route",
		}, 403)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (LogicModel, error) {
	version, err := getVersion(c.model.ProjectID, c.model.VersionName)
	if err != nil {
		return LogicModel{}, err
	}

	items, err := getItemsConcurrently(c.model.ProjectID, version.ID, c.model.IDs, c.model.Options)
	if err != nil {
		return LogicModel{}, err
	}

	childIds := sdk.Map(items, func(idx int, value Item) string {
		return value.ItemID
	})

	mappedConnections := make(map[string]connections)
	connections := newConnections()
	models, err := connections2.GetManyConnections(version.ID, c.model.ProjectID, childIds)
	if err != nil {
		return LogicModel{}, err
	}

	for _, item := range items {
		parents := make([]string, 0)
		children := make([]string, 0)
		for _, model := range models {
			if model.Parent == item.ItemID {
				children = append(children, model.Child)
			}

			if model.Child == item.ItemID {
				parents = append(parents, model.Parent)
			}
		}

		connections.parents = parents
		connections.children = children
		mappedConnections[item.ItemID] = connections
	}

	return LogicModel{
		Items:       items,
		Connections: mappedConnections,
		Options:     c.model.Options,
	}, nil
}

func (c Main) Handle() (interface{}, error) {
	if err := c.Validate(); err != nil {
		return []View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return []View{}, err
	}

	if err := c.Authorize(); err != nil {
		return []View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return []View{}, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, interface{}, LogicModel] {
	logBuilder.Add("getListItemsByName", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
