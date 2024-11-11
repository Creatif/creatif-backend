package paginateMapItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("paginateMapItems", errs, publicApiError.ValidationError)
	}
	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return publicApiError.NewError("paginateMapItems", map[string]string{
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

	if c.model.Limit == 0 {
		c.model.Limit = 100
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

	groups, err := getGroupIdsByName(c.model.ProjectID, c.model.Groups)
	if err != nil {
		return LogicModel{}, err
	}

	itemsSql, placeholders, err := getItemSql(c.model.StructureName, c.model.Page, c.model.Limit, order, sortBy, c.model.Search, lcls, groups, c.model.Query)
	if err != nil {
		return LogicModel{}, publicApiError.NewError("paginateListItems", map[string]string{
			"error": err.Error(),
		}, publicApiError.ApplicationError)
	}

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
			Items: []Item{},
		}, nil
	}

	normalizedGroups, err := getGroups(sdk.Map(items, func(idx int, value Item) string {
		return value.ItemID
	}))
	if err != nil {
		return LogicModel{}, err
	}

	for i, item := range items {
		if _, ok := normalizedGroups[item.ItemID]; ok {
			item.Groups = normalizedGroups[item.ItemID]
			items[i] = item
		}
	}

	return LogicModel{
		Items:   items,
		Options: c.model.Options,
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, interface{}, LogicModel] {
	return Main{model: model, auth: auth}
}
