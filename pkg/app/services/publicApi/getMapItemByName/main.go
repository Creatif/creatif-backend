package getMapItemByName

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
		return publicApiError.NewError("getMapItemByName", errs, publicApiError.ValidationError)
	}

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

	placeholders, err := createPlaceholders(c.model.ProjectID, version.ID, c.model.StructureName, c.model.Name, c.model.Locale)
	if err != nil {
		return LogicModel{}, err
	}

	item, err := getItem(placeholders)
	if err != nil {
		return LogicModel{}, err
	}

	groups, err := getGroups(item.ItemID)
	if err != nil {
		return LogicModel{}, err
	}

	item.Groups = groups

	return LogicModel{
		Item:    item,
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
