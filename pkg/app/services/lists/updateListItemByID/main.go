package updateListItemByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
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
	"strings"
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

	if len(c.model.Values.Groups) > 0 {
		count, err := shared.ValidateGroupsExist(c.model.ProjectID, c.model.Values.Groups)
		if err != nil {
			return appErrors.NewValidationError(map[string]string{
				"groupsExist": err.Error(),
			})
		}

		if count+len(c.model.Values.Groups) > 20 {
			return appErrors.NewValidationError(map[string]string{
				"maximumGroups": fmt.Sprintf("You are trying to add %d more groups but you already have %d assigned to this item. Maximum number of groups per item is 20", len(c.model.Values.Groups), count),
			})
		}
	}

	type GroupBehaviourCheck struct {
		Count     int    `gorm:"column:count"`
		Behaviour string `gorm:"column:behaviour"`
	}

	var check GroupBehaviourCheck
	res := storage.Gorm().Raw(fmt.Sprintf(`
SELECT behaviour
FROM %s AS lv 
INNER JOIN %s AS l ON (l.id = ? OR l.short_id = ?) AND l.project_id = ? AND l.id = lv.list_id AND (lv.id = ? OR lv.short_id = ?)`,
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName()),
		c.model.ListName,
		c.model.ListName,
		c.model.ProjectID,
		c.model.ItemID,
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

	if check.Behaviour == constants.ReadonlyBehaviour {
		return appErrors.NewValidationError(map[string]string{
			"behaviourReadonly": fmt.Sprintf("List item with ID '%s' is readonly and cannot be updated.", c.model.ItemID),
		})
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

func (c Main) Logic() (LogicResult, error) {
	var list declarations.List
	if res := storage.Gorm().Where(
		fmt.Sprintf("(id = ? OR short_id = ?) AND project_id = ?"),
		c.model.ListName,
		c.model.ListName,
		c.model.ProjectID).
		Select("id").First(&list); res.Error != nil {
		c.logBuilder.Add("updateListItemByID", res.Error.Error())

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return LogicResult{}, appErrors.NewNotFoundError(res.Error).AddError("updateListItemByID.Logic", nil)
		}

		return LogicResult{}, appErrors.NewDatabaseError(res.Error).AddError("updateListItemByID.Logic", nil)
	}

	var existing declarations.ListVariable
	if res := storage.Gorm().Where(fmt.Sprintf("(id = ? OR short_id = ?) AND list_id = ?"),
		c.model.ItemID,
		c.model.ItemID,
		list.ID).
		First(&existing); res.Error != nil {
		c.logBuilder.Add("updateListItemByID", res.Error.Error())

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return LogicResult{}, appErrors.NewNotFoundError(res.Error).AddError("updateListItemByID.Logic", nil)
		}

		return LogicResult{}, appErrors.NewDatabaseError(res.Error).AddError("updateListItemByID.Logic", nil)
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

		if f == "behaviour" {
			existing.Behaviour = c.model.Values.Behaviour
		}

		if f == "locale" {
			localeID, _ := locales.GetIDWithAlpha(c.model.Values.Locale)
			existing.LocaleID = localeID
		}
	}

	var updated declarations.ListVariable
	if transactionErr := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Model(&updated).Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
			{Name: "name"},
			{Name: "behaviour"},
			{Name: "metadata"},
			{Name: "locale_id"},
			{Name: "value"},
			{Name: "created_at"},
			{Name: "updated_at"},
		}}).Where("id = ?", existing.ID).Updates(existing); res.Error != nil {
			c.logBuilder.Add("updateListItemByID", res.Error.Error())

			return appErrors.NewApplicationError(res.Error).AddError("updateListItemByID.Logic", nil)
		}

		if sdk.Includes(c.model.Fields, "groups") {
			if res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE variable_id = ?", (declarations.VariableGroup{}).TableName()), c.model.ItemID); res.Error != nil {
				return res.Error
			}

			if c.model.Values.Groups != nil && len(c.model.Values.Groups) > 0 {
				variablesGroups := make([]declarations.VariableGroup, 0)
				for _, g := range c.model.Values.Groups {
					variablesGroups = append(variablesGroups, declarations.NewVariableGroup(g, c.model.ItemID, c.model.Values.Groups))
				}

				if res := tx.Create(&variablesGroups); res.Error != nil {
					return res.Error
				}
			}
		}

		if sdk.Includes(c.model.Fields, "references") {
			if err := shared.UpdateReferences(c.model.References, list.ID, updated.ID, c.model.ProjectID, tx); err != nil {
				return err
			}
		}

		return nil
	}); transactionErr != nil {
		errString := transactionErr.Error()
		splt := strings.Split(errString, ":")
		if len(splt) == 2 {
			return LogicResult{}, appErrors.NewValidationError(map[string]string{
				splt[0]: splt[1],
			})
		}

		return LogicResult{}, appErrors.NewApplicationError(transactionErr).AddError("updateMapVariable.Logic", nil)
	}

	var groups []declarations.Group
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT g.name, g.id FROM %s AS g INNER JOIN %s AS vg ON vg.group_id = g.id AND vg.variable_id = ?", (declarations.Group{}).TableName(), (declarations.VariableGroup{}).TableName()), c.model.ItemID).Scan(&groups)
	if res.Error != nil {
		return LogicResult{}, appErrors.NewDatabaseError(res.Error)
	}

	return LogicResult{
		Variable: updated,
		Groups:   groups,
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
	logBuilder.Add("updateListItemByID", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
