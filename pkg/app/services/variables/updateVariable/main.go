package updateVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/logger"
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

	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("updateVariable", err.Error())
		return appErrors.NewApplicationError(err).AddError("updateVariable.Logic", nil)
	}

	type GroupBehaviourCheck struct {
		Count     int    `gorm:"column:count"`
		Behaviour string `gorm:"column:behaviour"`
	}

	var check GroupBehaviourCheck
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT cardinality(groups) AS count, behaviour FROM %s WHERE name = ? AND project_id = ? AND locale_id = ?", (declarations.Variable{}).TableName()), c.model.Name, c.model.ProjectID, localeID).Scan(&check)
	if res.Error != nil || res.RowsAffected == 0 {
		if res.Error != nil {
			c.logBuilder.Add("updateVariable", res.Error.Error())
		} else {
			c.logBuilder.Add("updateVariable", "No rows returned in group check")
		}
		return appErrors.NewValidationError(map[string]string{
			"groups": fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", c.model.Name),
		})
	}

	if check.Count+len(c.model.Values.Groups) > 20 {
		return appErrors.NewValidationError(map[string]string{
			"groups": fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", c.model.Name),
		})
	}

	if check.Behaviour == constants.ReadonlyBehaviour {
		return appErrors.NewValidationError(map[string]string{
			"behaviour": fmt.Sprintf("Cannot update a readonly variable '%s'", c.model.Name),
		})
	}

	c.logBuilder.Add("updateVariable", "Validated.")

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
	localeID, _ := locales.GetIDWithAlpha(c.model.Locale)

	var existing declarations.Variable
	if res := storage.Gorm().Where("name = ? AND project_id = ? AND locale_id = ?", c.model.Name, c.model.ProjectID, localeID).First(&existing); res.Error != nil {
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
	}

	var updated declarations.Variable
	if res := storage.Gorm().Model(&updated).Clauses(clause.Returning{Columns: []clause.Column{
		{Name: "id"},
		{Name: "project_id"},
		{Name: "name"},
		{Name: "behaviour"},
		{Name: "metadata"},
		{Name: "value"},
		{Name: "groups"},
		{Name: "created_at"},
		{Name: "updated_at"},
	}}).Where("id = ?", existing.ID).Select(c.model.Fields).Updates(existing); res.Error != nil {
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

	return newView(model, c.model.Locale), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.Variable] {
	logBuilder.Add("updateVariable", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
