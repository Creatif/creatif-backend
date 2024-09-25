package getListItemById

import (
	"creatif/pkg/app/auth"
	connections2 "creatif/pkg/app/services/publicApi/connections"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/storage"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("getListItemById", errs, publicApiError.ValidationError)
	}

	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return publicApiError.NewError("getListItemById", map[string]string{
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

	var item Item
	res := storage.Gorm().Raw(getListItemSql(c.model.Options), c.model.ProjectID, version.ID, c.model.ItemID).Scan(&item)
	if res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getListItemById", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res.RowsAffected == 0 {
		return LogicModel{}, publicApiError.NewError("getListItemById", map[string]string{
			"notFound": "This list item does not exist",
		}, publicApiError.NotFoundError)
	}

	connections := newConnections()
	parents := make([]string, 0)
	children := make([]string, 0)
	models, err := connections2.GetConnections(version.ID, c.model.ProjectID, c.model.ItemID)
	if err != nil {
		return LogicModel{}, err
	}

	for _, model := range models {
		if model.Parent == c.model.ItemID {
			children = append(children, model.Child)
		}

		if model.Child == c.model.ItemID {
			parents = append(parents, model.Parent)
		}
	}

	connections.parents = parents
	connections.children = children

	return LogicModel{
		Item:        item,
		Connections: connections,
		Options:     c.model.Options,
	}, nil
}

func (c Main) Handle() (interface{}, error) {
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, interface{}, LogicModel] {
	return Main{model: model, auth: auth}
}
