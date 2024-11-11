package paginateListItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return publicApiError.NewError("paginateListItems", errs, publicApiError.ValidationError)
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

	groups, err := getGroupIdsByName(c.model.ProjectID, c.model.Groups)
	if err != nil {
		return LogicModel{}, err
	}

	defs := createDefaults(c.model.Page, c.model.Limit, c.model.Order)
	placeholders := createPlaceholders(c.model.ProjectID, version.ID, defs.page, defs.limit, c.model.StructureName, c.model.Locales, c.model.Search)
	subQrs, err := createSubQueries(c.model.SortBy, c.model.Search, groups, placeholders["locales"].([]string), c.model.Query)
	if err != nil {
		return LogicModel{}, publicApiError.NewError("paginateListItems", map[string]string{
			"error": err.Error(),
		}, publicApiError.ApplicationError)
	}

	items, err := getItem(placeholders, defs, subQrs)
	if err != nil {
		return LogicModel{}, err
	}

	items, err = placeGroupsInItems(items)
	if err != nil {
		return LogicModel{}, err
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
