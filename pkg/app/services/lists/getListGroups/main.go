package getListGroups

import (
	"creatif/pkg/app/auth"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("getMapGroups", "Validating...")

	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("queryMapVariable", "Validated")

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

func (c Main) Logic() ([]string, error) {
	duplicatedModels, err := getGroups(c.model.Name, c.model.ProjectID)
	if err != nil {
		return []string{}, err
	}

	distinctModels := make([]string, 0)
	for _, v := range duplicatedModels {
		groups := v.Groups

		for _, g := range groups {
			if !sdk.Includes(distinctModels, g) {
				distinctModels = append(distinctModels, g)
			}
		}
	}

	return distinctModels, nil
}

func (c Main) Handle() ([]string, error) {
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, []string, []string] {
	logBuilder.Add("queryMapVariable", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
