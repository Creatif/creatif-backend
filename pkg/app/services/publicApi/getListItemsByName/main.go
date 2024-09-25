package getListItemsByName

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/locales"
	connections2 "creatif/pkg/app/services/publicApi/connections"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("getMapItemsByName", errs, publicApiError.ValidationError)
	}

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

	placeholders := make(map[string]interface{})
	placeholders["projectId"] = c.model.ProjectID
	placeholders["versionName"] = version.Name
	placeholders["structureName"] = c.model.StructureName
	placeholders["variableName"] = c.model.Name

	var locale string
	if c.model.Locale != "" {
		l, err := locales.GetIDWithAlpha(c.model.Locale)
		if err != nil {
			return LogicModel{}, publicApiError.NewError("getListItemsByName", map[string]string{
				"invalidLocale": "The locale you provided is invalid and does not exist.",
			}, publicApiError.ValidationError)
		}

		placeholders["localeId"] = l
		locale = l
	}

	var items []Item
	res := storage.Gorm().Raw(
		getItemSql(locale, c.model.Options),
		placeholders,
	).Scan(&items)
	if res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getListItemsByName", map[string]string{
			"error": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res.RowsAffected == 0 {
		return LogicModel{}, nil
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, interface{}, LogicModel] {
	return Main{model: model, auth: auth}
}
