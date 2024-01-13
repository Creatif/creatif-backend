package updateVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/constants"
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
	c.logBuilder.Add("updateVariable", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	// check if the variable to be updated with id/locale_id exists
	var existing declarations.Variable
	res := storage.Gorm().Where("(id = ? OR short_id = ?) AND project_id = ?", c.model.ID, c.model.ID, c.model.ProjectID).Select("id", "name").First(&existing)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return appErrors.NewValidationError(map[string]string{
			"notExists": fmt.Sprintf("Variable with ID '%s'.", c.model.ID),
		})
	} else if res.Error != nil {
		return appErrors.NewValidationError(map[string]string{
			"notExists": fmt.Sprintf("Variable with ID '%s'.", c.model.ID),
		})
	}

	// check if variable with name and locale already exists
	if sdk.Includes(c.model.Fields, "locale") || sdk.Includes(c.model.Fields, "name") {
		name := c.model.Values.Name
		updatingLocaleId, _ := locales.GetIDWithAlpha(c.model.Values.Locale)

		var variable declarations.Variable
		res := storage.Gorm().Where("name = ? AND project_id = ? AND locale_id = ? AND id != ?", name, c.model.ProjectID, updatingLocaleId, existing.ID).Select("id").First(&existing)
		if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return appErrors.NewValidationError(map[string]string{
				"exists": fmt.Sprintf("Variable with name '%s' and locale '%s' already exists.", c.model.Values.Name, c.model.Values.Locale),
			})
		}

		if variable.ID != "" {
			return appErrors.NewValidationError(map[string]string{
				"exists": fmt.Sprintf("Variable with name '%s' and locale '%s' already exists.", c.model.Values.Name, c.model.Values.Locale),
			})
		}
	}

	// check that groups number is correct
	type GroupBehaviourCheck struct {
		Count     int    `gorm:"column:count"`
		Behaviour string `gorm:"column:behaviour"`
	}

	var check GroupBehaviourCheck
	res = storage.Gorm().Raw(fmt.Sprintf("SELECT cardinality(groups) AS count, behaviour FROM %s WHERE project_id = ? AND (id = ? OR short_id = ?)", (declarations.Variable{}).TableName()), c.model.ProjectID, c.model.ID, c.model.ID).Scan(&check)
	if res.Error != nil || res.RowsAffected == 0 {
		if res.Error != nil {
			c.logBuilder.Add("updateVariable", res.Error.Error())
		} else {
			c.logBuilder.Add("updateVariable", "No rows returned in group check")
		}
		return appErrors.NewValidationError(map[string]string{
			"groups": fmt.Sprintf("Invalid number of groups for variable with ID '%s'. Maximum number of groups per variable is 20.", c.model.ID),
		})
	}

	if check.Count+len(c.model.Values.Groups) > 20 {
		return appErrors.NewValidationError(map[string]string{
			"groups": fmt.Sprintf("Invalid number of groups for variable with ID '%s'. Maximum number of groups per variable is 20.", c.model.ID),
		})
	}

	if check.Behaviour == constants.ReadonlyBehaviour {
		return appErrors.NewValidationError(map[string]string{
			"behaviourReadonly": fmt.Sprintf("List item with ID '%s' is readonly and cannot be updated.", c.model.ID),
		})
	}

	c.logBuilder.Add("updateVariable", "Validated.")

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

func (c Main) Logic() (declarations.Variable, error) {
	var existing declarations.Variable
	if res := storage.Gorm().Where(fmt.Sprintf("(id = ? OR short_id = ?) AND project_id = ?"), c.model.ID, c.model.ID, c.model.ProjectID).First(&existing); res.Error != nil {
		c.logBuilder.Add("updateVariable", res.Error.Error())
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.Variable{}, appErrors.NewNotFoundError(res.Error).AddError("updateVariable.Logic", nil)
		}

		return declarations.Variable{}, appErrors.NewDatabaseError(res.Error).AddError("updateVariable.Logic", nil)
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
			updatingLocale, _ := locales.GetIDWithAlpha(c.model.Values.Locale)
			existing.LocaleID = updatingLocale
		}
	}

	var updated declarations.Variable
	if res := storage.Gorm().Model(&updated).Clauses(clause.Returning{Columns: []clause.Column{
		{Name: "id"},
		{Name: "project_id"},
		{Name: "short_id"},
		{Name: "name"},
		{Name: "locale_id"},
		{Name: "behaviour"},
		{Name: "metadata"},
		{Name: "value"},
		{Name: "groups"},
		{Name: "created_at"},
		{Name: "updated_at"},
	}}).Where("id = ?", existing.ID).Select(sdk.Replace(c.model.Fields, "locale", "locale_id")).Updates(existing); res.Error != nil {
		c.logBuilder.Add("updateVariable", res.Error.Error())
		return declarations.Variable{}, appErrors.NewApplicationError(res.Error).AddError("updateVariable.Logic", nil)
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.Variable] {
	logBuilder.Add("updateVariable", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
