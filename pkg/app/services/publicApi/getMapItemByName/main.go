package getMapItemByName

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/locales"
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
	c.logBuilder.Add("getMapItemByName", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("getMapItemByName", errs, publicApiError.ValidationError)
	}

	c.logBuilder.Add("getMapItemByName", "Validated")
	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return publicApiError.NewError("getMapItemByName", map[string]string{
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

	placeholders := make(map[string]interface{})
	placeholders["projectId"] = c.model.ProjectID
	placeholders["versionName"] = version.Name
	placeholders["structureName"] = c.model.StructureName
	placeholders["variableName"] = c.model.Name

	var locale string
	if c.model.Locale != "" {
		l, _ := locales.GetIDWithAlpha(c.model.Locale)
		locale = l
		placeholders["localeId"] = l
	}

	var mapItem Item
	res := storage.Gorm().Raw(getItemSql(locale), placeholders).Scan(&mapItem)
	if res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getMapItemByName", map[string]string{
			"error": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res.RowsAffected == 0 {
		return LogicModel{}, publicApiError.NewError("getMapItemByName", map[string]string{
			"notFound": "Item has not been found.",
		}, publicApiError.NotFoundError)
	}

	connections := newConnections()
	parents := make([]string, 0)
	children := make([]string, 0)
	models, err := connections2.GetConnections(version.ID, c.model.ProjectID, mapItem.ItemID)
	if err != nil {
		return LogicModel{}, err
	}

	for _, model := range models {
		if model.Parent == mapItem.ItemID {
			children = append(children, model.Child)
		}

		if model.Child == mapItem.ItemID {
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, interface{}, LogicModel] {
	logBuilder.Add("getMapItemByName", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
