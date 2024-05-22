package paginateListItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/published"
	"creatif/pkg/app/services/locales"
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
		return publicApiError.NewError("paginateListItems", errs, 422)
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
	var version published.Version
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT * FROM %s WHERE project_id = ? AND is_production_version = true", (published.Version{}).TableName()), c.model.ProjectID).Scan(&version)
	if res.Error != nil {
		return LogicModel{}, publicApiError.NewError("paginateListItems", map[string]string{
			"internalError": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res.RowsAffected == 0 {
		return LogicModel{}, publicApiError.NewError("paginateListItems", map[string]string{
			"error": "Production version has not been found",
		}, publicApiError.NotFoundError)
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
	res = storage.Gorm().Raw(itemsSql, placeholders).Scan(&items)
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

	var connections []ConnectionItem
	res = storage.Gorm().Raw(getConnectionsSql(), c.model.ProjectID, version.Name, childIds, c.model.ProjectID, version.Name, childIds).Scan(&connections)
	if res.Error != nil {
		return LogicModel{}, publicApiError.NewError("paginateListItems", map[string]string{
			"error": res.Error.Error(),
		}, publicApiError.ApplicationError)
	}

	mappedConnections := make(map[string][]ConnectionItem)
	if len(connections) > 0 {
		for _, conn := range connections {
			if _, ok := mappedConnections[conn.ItemID]; !ok {
				mappedConnections[conn.ItemID] = make([]ConnectionItem, 0)
			}

			mappedConnections[conn.ItemID] = append(mappedConnections[conn.ItemID], conn)
		}
	}

	return LogicModel{
		Items:       items,
		Connections: mappedConnections,
	}, nil
}

func (c Main) Handle() ([]View, error) {
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, []View, LogicModel] {
	logBuilder.Add("paginateListItems", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
