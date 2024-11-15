package updateListItemByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/app/services/shared/fileProcessor"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
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
		var images []declarations.File
		if res := tx.Raw(fmt.Sprintf("SELECT * FROM %s WHERE project_id = ? AND list_id = ?", (declarations.File{}).TableName()), c.model.ProjectID, c.model.ItemID).Scan(&images); res.Error != nil {
			return res.Error
		}

		newValue, err := fileProcessor.UpdateFiles(
			c.model.ProjectID,
			c.model.Values.Value,
			c.model.ImagePaths,
			images,
			func(fileSystemFilePath, path, mimeType, extension, fileName string) (string, error) {
				image := declarations.NewFile(
					c.model.ProjectID,
					&c.model.ItemID,
					nil,
					fileSystemFilePath,
					path,
					mimeType,
					extension,
					fileName,
				)

				if res := tx.Create(&image); res.Error != nil {
					return "", res.Error
				}

				return image.ID, nil
			},
			func(imageId, fieldName string) error {
				if fieldName != "" {
					if res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ? AND field_name = ?", (declarations.File{}).TableName()), imageId, fieldName); res.Error != nil {
						return res.Error
					}

					return nil
				}

				if res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", (declarations.File{}).TableName()), imageId); res.Error != nil {
					return res.Error
				}

				return nil
			},
		)

		if err != nil {
			return err
		}

		existing.Value = newValue

		if res := tx.Model(&updated).Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
			{Name: "name"},
			{Name: "behaviour"},
			{Name: "metadata"},
			{Name: "locale_id"},
			{Name: "value"},
			{Name: "created_at"},
			{Name: "updated_at"},
		}}).Where("id = ?", existing.ID).Updates(&existing); res.Error != nil {

			return appErrors.NewApplicationError(res.Error).AddError("updateListItemByID.Logic", nil)
		}

		if sdk.Includes(c.model.Fields, "groups") {
			if res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE variable_id = ?", (declarations.VariableGroup{}).TableName()), c.model.ItemID); res.Error != nil {
				return res.Error
			}

			if c.model.Values.Groups != nil && len(c.model.Values.Groups) > 0 {
				newGroup := declarations.NewVariableGroup(c.model.ItemID, c.model.Values.Groups)
				if res := tx.Create(&newGroup); res.Error != nil {
					return res.Error
				}
			}
		}

		if sdk.Includes(c.model.Fields, "references") {
			conns := sdk.Map(c.model.References, func(idx int, value shared.UpdateReference) connections.Connection {
				return connections.Connection{
					Path:          value.Name,
					StructureType: value.StructureType,
					VariableID:    value.VariableID,
				}
			})

			newValue, newConnections, err := connections.RecreateConnections(c.model.ProjectID, existing.ID, "list", conns, existing.Value)
			if err != nil {
				return err
			}

			existing.Value = newValue

			if res := tx.Create(&newConnections); res.Error != nil {
				return res.Error
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
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT g.name, g.id FROM %s AS g INNER JOIN %s AS vg ON vg.variable_id = ? AND g.id = ANY(vg.groups)", (declarations.Group{}).TableName(), (declarations.VariableGroup{}).TableName()), c.model.ItemID).Scan(&groups)
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, View, LogicResult] {
	return Main{model: model, auth: auth}
}
