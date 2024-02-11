package mapCreate

import (
	"creatif/pkg/app/auth"
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
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("mapCreate", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	c.logBuilder.Add("getMap", "Validated.")
	return nil
}

func (c Main) Authenticate() error {
	// user check by project id should be gotten here, with authentication cookie
	var project app.Project
	if err := storage.Get((app.Project{}).TableName(), c.model.ProjectID, &project); err != nil {
		return appErrors.NewAuthenticationError(err).AddError("createVariable.Authenticate", nil)
	}

	return nil
}

func (c Main) Authorize() error {

	return nil
}

func (c Main) Logic() (LogicResult, error) {
	newMap := declarations.NewMap(c.model.ProjectID, c.model.Name)
	names := make([]ViewVariable, 0)
	if err := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&newMap); res.Error != nil {
			return res.Error
		}

		domainEntries := make([]declarations.MapVariable, len(c.model.Variables))
		variables := c.model.Variables
		for i, variable := range variables {
			localeId, _ := locales.GetIDWithAlpha(variable.Locale)

			domainEntries[i] = declarations.NewMapVariable(
				newMap.ID,
				localeId,
				variable.Name,
				variable.Behaviour,
				variable.Metadata,
				variable.Value,
			)
		}

		if len(domainEntries) != 0 {
			if res := tx.Create(&domainEntries); res.Error != nil {
				return res.Error
			}

			for _, d := range domainEntries {
				if d.ID != "" {
					locale, _ := locales.GetAlphaWithID(d.LocaleID)

					names = append(names, ViewVariable{
						ID:      d.ID,
						Locale:  locale,
						ShortID: d.ShortID,
						Name:    d.Name,
					})
				}
			}
		}

		return nil
	}); err != nil {
		c.logBuilder.Add("mapCreate", err.Error())
		return LogicResult{}, appErrors.NewDatabaseError(err).AddError("mapCreate.Logic", nil)
	}

	return LogicResult{
		ID:        newMap.ID,
		ShortID:   newMap.ShortID,
		ProjectID: newMap.ProjectID,
		Name:      newMap.Name,
		Variables: names,
	}, nil
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, LogicResult] {
	logBuilder.Add("mapCreate", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
