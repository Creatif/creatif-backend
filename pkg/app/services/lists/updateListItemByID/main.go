package updateListItemByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
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
	c.logBuilder.Add("updateListItemByID", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("updateListItemByID", err.Error())
		return appErrors.NewApplicationError(err).AddError("updateVariable.Logic", nil)
	}

	type GroupBehaviourCheck struct {
		Count     int    `gorm:"column:count"`
		Behaviour string `gorm:"column:behaviour"`
	}

	var check GroupBehaviourCheck
	res := storage.Gorm().Raw(fmt.Sprintf(`
SELECT cardinality(lv.groups) AS count, behaviour
FROM %s AS lv 
INNER JOIN %s AS l ON l.name = ? AND l.project_id = ? AND l.locale_id = ? AND l.id = lv.list_id AND lv.id = ?`,
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName()),
		c.model.ListName,
		c.model.ProjectID,
		localeID,
		c.model.ItemID,
	).Scan(&check)

	if res.Error != nil || res.RowsAffected == 0 {
		if res.Error != nil {
			c.logBuilder.Add("updateListItemByID", res.Error.Error())
		}
		return appErrors.NewValidationError(map[string]string{
			"groups": fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", c.model.ItemID),
		})
	}

	if check.Count+len(c.model.Values.Groups) > 20 {
		return appErrors.NewValidationError(map[string]string{
			"groups": fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", c.model.ItemID),
		})
	}

	if check.Behaviour == constants.ReadonlyBehaviour {
		return appErrors.NewValidationError(map[string]string{
			"behaviour": fmt.Sprintf("Cannot update a readonly list item with ID '%s'", c.model.ItemID),
		})
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

func (c Main) Logic() (declarations.ListVariable, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("updateListItemByID", err.Error())
		return declarations.ListVariable{}, appErrors.NewNotFoundError(err).AddError("updateListItemByID.Logic", nil)
	}

	listId, listVal := shared.DetermineID("", c.model.ListName, c.model.ListID, c.model.ListShortID)
	var list declarations.List
	if res := storage.Gorm().Where(fmt.Sprintf("%s AND project_id = ? AND locale_id = ?", listId), listVal, c.model.ProjectID, localeID).Select("id").First(&list); res.Error != nil {
		c.logBuilder.Add("updateListItemByID", res.Error.Error())

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.ListVariable{}, appErrors.NewNotFoundError(res.Error).AddError("updateListItemByID.Logic", nil)
		}

		return declarations.ListVariable{}, appErrors.NewDatabaseError(res.Error).AddError("updateListItemByID.Logic", nil)
	}

	varId, varVal := shared.DetermineID("", "", c.model.ItemID, c.model.ItemShortID)
	var existing declarations.ListVariable
	if res := storage.Gorm().Where(fmt.Sprintf("%s AND list_id = ? AND locale_id = ?", varId), varVal, list.ID, localeID).First(&existing); res.Error != nil {
		c.logBuilder.Add("updateListItemByID", res.Error.Error())

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.ListVariable{}, appErrors.NewNotFoundError(res.Error).AddError("updateListItemByID.Logic", nil)
		}

		return declarations.ListVariable{}, appErrors.NewDatabaseError(res.Error).AddError("updateListItemByID.Logic", nil)
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

	var updated declarations.ListVariable
	if res := storage.Gorm().Model(&updated).Clauses(clause.Returning{Columns: []clause.Column{
		{Name: "id"},
		{Name: "name"},
		{Name: "behaviour"},
		{Name: "metadata"},
		{Name: "value"},
		{Name: "groups"},
		{Name: "created_at"},
		{Name: "updated_at"},
	}}).Where("id = ?", existing.ID).Select(c.model.Fields).Updates(existing); res.Error != nil {
		c.logBuilder.Add("updateListItemByID", res.Error.Error())

		return declarations.ListVariable{}, appErrors.NewApplicationError(res.Error).AddError("updateListItemByID.Logic", nil)
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.ListVariable] {
	logBuilder.Add("updateListItemByID", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
