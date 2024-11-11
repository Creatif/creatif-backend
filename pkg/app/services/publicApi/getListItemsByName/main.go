package getListItemsByName

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/sdk"
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

	placeholders, err := createPlaceholders(c.model.ProjectID, version.ID, c.model.StructureName, c.model.Name, c.model.Locale)
	if err != nil {
		return LogicModel{}, err
	}

	items, err := getItem(placeholders, c.model.Options)
	if err != nil {
		return LogicModel{}, err
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
