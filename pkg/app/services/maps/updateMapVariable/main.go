package updateMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	if err := validateGroups(c.model.MapName, c.model.ProjectID, c.model.VariableName, c.model.Values.Groups, c.logBuilder); err != nil {
		return err
	}

	if sdk.Includes(c.model.Fields, "name") {
		return validateUniqueName(c.model.MapName, c.model.VariableName, c.model.Values.Name, c.model.ProjectID)
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

func (c Main) Logic() (declarations.MapVariable, error) {
	var m declarations.Map
	if res := storage.Gorm().Where(
		fmt.Sprintf("(name = ? OR id = ? OR short_id = ?) AND project_id = ?"),
		c.model.MapName,
		c.model.MapName,
		c.model.MapName,
		c.model.ProjectID).
		Select("id", "name").First(&m); res.Error != nil {
		c.logBuilder.Add("updateMapVariable", res.Error.Error())

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.MapVariable{}, appErrors.NewNotFoundError(res.Error).AddError("updateMapVariable.Logic", nil)
		}

		return declarations.MapVariable{}, appErrors.NewDatabaseError(res.Error).AddError("updateMapVariable.Logic", nil)
	}

	var existing declarations.MapVariable
	if res := storage.Gorm().Where(fmt.Sprintf("(id = ? OR short_id = ?) AND map_id = ?"),
		c.model.VariableName,
		c.model.VariableName,
		m.ID).
		First(&existing); res.Error != nil {
		c.logBuilder.Add("updateMapVariable", res.Error.Error())

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.MapVariable{}, appErrors.NewNotFoundError(res.Error).AddError("updateMapVariable.Logic", nil)
		}

		return declarations.MapVariable{}, appErrors.NewDatabaseError(res.Error).AddError("updateMapVariable.Logic", nil)
	}

	for _, f := range c.model.Fields {
		if f == "name" {
			existing.Name = c.model.Values.Name
		}

		if f == "metadata" {
			existing.Metadata = c.model.Values.Metadata
		}

		if f == "value" {
			existing.Value = c.model.Values.Value
		}

		if f == "groups" {
			existing.Groups = c.model.Values.Groups
		}

		if f == "behaviour" {
			existing.Behaviour = c.model.Values.Behaviour
		}

		if f == "locale" {
			localeID, _ := locales.GetIDWithAlpha(c.model.Values.Locale)
			existing.LocaleID = localeID
		}
	}

	var updated declarations.MapVariable
	if err := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Model(&updated).Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
			{Name: "name"},
			{Name: "behaviour"},
			{Name: "metadata"},
			{Name: "locale_id"},
			{Name: "value"},
			{Name: "groups"},
			{Name: "created_at"},
			{Name: "updated_at"},
		}}).Where("id = ?", existing.ID).Updates(existing); res.Error != nil {
			c.logBuilder.Add("updateMapVariable", res.Error.Error())

			return res.Error
		}

		if err := shared.UpdateReferences(c.model.References, m.ID, updated.ID, c.model.ProjectID, tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return declarations.MapVariable{}, appErrors.NewApplicationError(err).AddError("updateMapVariable.Logic", nil)
	}

	return updated, nil
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.MapVariable] {
	logBuilder.Add("updateMapVariable", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
