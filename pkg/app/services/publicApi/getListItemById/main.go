package getListItemById

import (
	"creatif/pkg/app/auth"
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
	c.logBuilder.Add("getListItemById", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("getListItemById", errs, publicApiError.ValidationError)
	}

	c.logBuilder.Add("getListItemById", "Validated")
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

	var connections []ConnectionItem
	if !c.model.Options.ValueOnly {
		res = storage.Gorm().Raw(getConnectionsMapSql(), c.model.ProjectID, version.ID, c.model.ItemID, c.model.ProjectID, version.Name, c.model.ItemID).Scan(&connections)
		if res.Error != nil {
			return LogicModel{}, publicApiError.NewError("getListItemById", map[string]string{
				"notFound": res.Error.Error(),
			}, publicApiError.DatabaseError)
		}
	}

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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, interface{}, LogicModel] {
	logBuilder.Add("getVersions", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
