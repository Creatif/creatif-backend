package paginateListItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/locales"
	connections2 "creatif/pkg/app/services/publicApi/connections"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("paginateListItems", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("paginateListItems", errs, publicApiError.ValidationError)
	}

	c.logBuilder.Add("paginateListItems", "Validated")
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

	var items []Item
	sortBy := "lv.index"
	if c.model.SortBy != "" {
		sortBy = fmt.Sprintf("lv.%s", c.model.SortBy)
	}

	order := "desc"
	if c.model.Order != "" {
		order = c.model.Order
	}

	lcls := make([]string, len(c.model.Locales))
	for i, l := range c.model.Locales {
		alpha, _ := locales.GetIDWithAlpha(l)
		lcls[i] = alpha
	}

	itemsSql, placeholders := getItemSql(c.model.StructureName, c.model.Page, order, sortBy, c.model.Search, lcls, c.model.Groups)
	placeholders["projectId"] = c.model.ProjectID
	placeholders["versionName"] = version.Name
	res := storage.Gorm().Raw(itemsSql, placeholders).Scan(&items)
	if res.Error != nil {
		return LogicModel{}, publicApiError.NewError("paginateListItems", map[string]string{
			"error": res.Error.Error(),
		}, publicApiError.ApplicationError)
	}

	if res.RowsAffected == 0 {
		return LogicModel{
			Items:       []Item{},
			Connections: nil,
		}, nil
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
	logBuilder.Add("paginateListItems", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
