package createList

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
}

func (c Main) Validate() error {
	c.logBuilder.Add("createList", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("createList", "Validated")
	return nil
}

func (c Main) Authenticate() error {
	// user check by project id should be gotten here, with authentication cookie
	var project app.Project
	if err := storage.Get((app.Project{}).TableName(), c.model.ProjectID, &project); err != nil {
		c.logBuilder.Add("appendToList", err.Error())
		return appErrors.NewAuthenticationError(err).AddError("createVariable.Authenticate", nil)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (declarations.List, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("createList", err.Error())
		return declarations.List{}, appErrors.NewApplicationError(err).AddError("createList.Logic", nil)
	}

	list := declarations.NewList(c.model.ProjectID, c.model.Name, localeID)
	if err := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&list); res.Error != nil {
			return res.Error
		}

		listVariables := make([]declarations.ListVariable, len(c.model.Variables))
		for i := 0; i < len(c.model.Variables); i++ {
			v := c.model.Variables[i]
			listVariables[i] = declarations.NewListVariable(list.ID, localeID, v.Name, v.Behaviour, v.Metadata, v.Groups, v.Value)
		}

		if res := tx.Create(&listVariables); res.Error != nil {
			return res.Error
		}

		return nil
	}); err != nil {
		c.logBuilder.Add("createList", err.Error())
		return declarations.List{}, appErrors.NewDatabaseError(err).AddError("createList.Logic", nil)
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

	return newView(model, c.model.Locale), nil
}

func New(model Model, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.List] {
	logBuilder.Add("createList", "Created")
	return Main{model: model, logBuilder: logBuilder}
}
