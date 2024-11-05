package getMapGroups

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return appErrors.NewAuthenticationError(err)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() ([]declarations2.Group, error) {
	sql := fmt.Sprintf(`
	SELECT g.name, g.id FROM %s AS g
	INNER JOIN %s AS vg ON vg.variable_id = ? AND g.project_id = ? AND g.id = ANY(vg.groups)
`, (declarations2.Group{}).TableName(), (declarations2.VariableGroup{}).TableName())

	var groups []declarations2.Group
	res := storage.Gorm().Raw(sql, c.model.ItemID, c.model.ProjectID).Scan(&groups)
	if res.Error != nil {
		return []declarations2.Group{}, appErrors.NewApplicationError(res.Error)
	}

	return groups, nil
}

func (c Main) Handle() ([]View, error) {
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, []View, []declarations2.Group] {
	return Main{model: model, auth: auth}
}
