package createVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("createVariable", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

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

func (c Main) Logic() (declarations.Variable, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("createVariable", err.Error())
		return declarations.Variable{}, appErrors.NewApplicationError(err)
	}

	var metadata []byte
	var value []byte
	if len(c.model.Metadata) > 0 {
		m, err := sdk.CovertToGeneric(c.model.Metadata)
		if err != nil {
			c.logBuilder.Add("createVariable", err.Error())
			return declarations.Variable{}, appErrors.NewApplicationError(err)
		}

		metadata = m
	}

	if len(c.model.Value) > 0 {
		m, err := sdk.CovertToGeneric(c.model.Value)
		if err != nil {
			c.logBuilder.Add("createVariable", err.Error())
			return declarations.Variable{}, appErrors.NewApplicationError(err)
		}

		value = m
	}

	model := declarations.NewVariable(c.model.ProjectID, localeID, c.model.Name, c.model.Behaviour, c.model.Groups, metadata, value)
	res := storage.Gorm().Model(&model).Clauses(clause.Returning{Columns: []clause.Column{
		{Name: "id"},
		{Name: "name"},
		{Name: "short_id"},
		{Name: "behaviour"},
		{Name: "metadata"},
		{Name: "value"},
		{Name: "groups"},
		{Name: "created_at"},
		{Name: "updated_at"},
	}}).Create(&model)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) || res.RowsAffected == 0 {
		if res.Error != nil {
			c.logBuilder.Add("createVariable", res.Error.Error())
		} else {
			c.logBuilder.Add("createVariable", "No rows returned. Returning 404 not found.")
		}
		return declarations.Variable{}, appErrors.NewNotFoundError(res.Error).AddError("createVariable.Logic", nil)
	}

	if res.Error != nil {
		c.logBuilder.Add("createVariable", res.Error.Error())
		return declarations.Variable{}, appErrors.NewDatabaseError(res.Error).AddError("createVariable.Logic", nil)
	}

	return model, nil
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.Variable] {
	logBuilder.Add("createVariable", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
