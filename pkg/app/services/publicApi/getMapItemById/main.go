package getMapItemById

import (
	"creatif/pkg/app/auth"
	connections2 "creatif/pkg/app/services/publicApi/connections"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("getVersions", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("getMapItemById", errs, publicApiError.ValidationError)
	}

	c.logBuilder.Add("getVersions", "Validated")
	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return publicApiError.NewError("getMapItemById", map[string]string{
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

	var mapItem Item
	res := storage.Gorm().Raw(getItemSql(c.model.Options), c.model.ProjectID, version.Name, c.model.ItemID).Scan(&mapItem)
	if res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getMapItemById", map[string]string{
			"data": res.Error.Error(),
		}, publicApiError.ApplicationError)
	}

	if res.RowsAffected == 0 {
		return LogicModel{}, publicApiError.NewError("getMapItemById", map[string]string{
			"data": "This item does not exist",
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
		Item:        mapItem,
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, interface{}, LogicModel] {
	logBuilder.Add("getVersions", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
