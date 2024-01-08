package appendToList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("appendToList", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("appendToList", "Validated.")
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

func (c Main) Logic() (declarations.List, error) {
	list, err := getList(c.model.Name)
	if err != nil {
		return declarations.List{}, err
	}

	assignDefaultGroupsToVariables(c.model.Variables)

	listVariables, err := createListVariables(list.ID, c.model.Variables)
	if err != nil {
		return declarations.List{}, err
	}

	if err := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&listVariables); res.Error != nil {
			c.logBuilder.Add("appendToList", res.Error.Error())
			return res.Error
		}

		return nil
	}); err != nil {
		return declarations.List{}, appErrors.NewApplicationError(err)
	}

	return list, nil
}

func (c Main) Handle() (View, error) {
	if err := c.Validate(); err != nil {
		return View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return View{}, err
	}

	if err := c.Authorize(); err != nil {
		return View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return View{}, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.List] {
	logBuilder.Add("appendToList", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
