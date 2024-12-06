package updateMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/app/services/shared/fileProcessor"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
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

	if sdk.Includes(c.model.Fields, "name") {
		if err := validateUniqueName(c.model.MapName, c.model.VariableName, c.model.Values.Name, c.model.ProjectID); err != nil {
			return err
		}
	}

	if sdk.Includes(c.model.Fields, "behaviour") {
		if err := validateBehaviour(c.model.MapName, c.model.ProjectID, c.model.VariableName, c.model.Values.Groups); err != nil {
			return err
		}
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
	var m declarations.Map
	if res := storage.Gorm().Where(
		fmt.Sprintf("(id = ? OR short_id = ?) AND project_id = ?"),
		c.model.MapName,
		c.model.MapName,
		c.model.ProjectID).
		Select("id", "name").First(&m); res.Error != nil {

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return LogicResult{}, appErrors.NewNotFoundError(res.Error).AddError("updateMapVariable.Logic", nil)
		}

		return LogicResult{}, appErrors.NewDatabaseError(res.Error).AddError("updateMapVariable.Logic", nil)
	}

	var existing declarations.MapVariable
	if res := storage.Gorm().Where(fmt.Sprintf("(id = ? OR short_id = ?) AND map_id = ?"),
		c.model.VariableName,
		c.model.VariableName,
		m.ID).
		First(&existing); res.Error != nil {

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return LogicResult{}, appErrors.NewNotFoundError(res.Error).AddError("updateMapVariable.Logic", nil)
		}

		return LogicResult{}, appErrors.NewDatabaseError(res.Error).AddError("updateMapVariable.Logic", nil)
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

	var updated declarations.MapVariable
	if err := storage.Transaction(func(tx *gorm.DB) error {
		var files []declarations.File
		if res := tx.Raw(fmt.Sprintf("SELECT * FROM %s WHERE project_id = ? AND map_id = ?", (declarations.File{}).TableName()), c.model.ProjectID, c.model.VariableName).Scan(&files); res.Error != nil {
			return res.Error
		}

		newValue, err := fileProcessor.UpdateFiles(
			c.model.ProjectID,
			c.model.Values.Value,
			c.model.ImagePaths,
			files,
			func(fileSystemFilePath, path, mimeType, extension, fileName string) (string, error) {
				image := declarations.NewFile(
					c.model.ProjectID,
					nil,
					&c.model.VariableName,
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
		}}).Where("id = ?", existing.ID).Updates(existing); res.Error != nil {
			return res.Error
		}

		if sdk.Includes(c.model.Fields, "groups") {
			res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE variable_id = ?", (declarations.VariableGroup{}).TableName()), c.model.VariableName)
			if res.Error != nil {
				return res.Error
			}

			if c.model.Values.Groups != nil && len(c.model.Values.Groups) > 0 {
				newGroup := declarations.NewVariableGroup(c.model.VariableName, c.model.Values.Groups)
				if res := tx.Create(&newGroup); res.Error != nil {
					return res.Error
				}
			}
		}

		if sdk.Includes(c.model.Fields, "connections") {
			newValue, newConnections, err := connections.RecreateConnections(tx, c.model.ProjectID, existing.ID, "map", c.model.Connections, existing.Value)
			if err != nil {
				return err
			}

			existing.Value = newValue

			if len(newConnections) != 0 {
				if res := tx.Create(&newConnections); res.Error != nil {
					return res.Error
				}
			}
		}

		return nil
	}); err != nil {
		errString := err.Error()
		splt := strings.Split(errString, ":")
		if len(splt) == 2 {
			return LogicResult{}, appErrors.NewValidationError(map[string]string{
				splt[0]: splt[1],
			})
		}

		return LogicResult{}, appErrors.NewApplicationError(err).AddError("updateMapVariable.Logic", nil)
	}

	groups := make([]declarations.Group, 0)
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT g.name, g.id FROM %s AS g INNER JOIN %s AS vg ON vg.variable_id = ? AND g.id = ANY(vg.groups)", (declarations.Group{}).TableName(), (declarations.VariableGroup{}).TableName()), c.model.VariableName).Scan(&groups)
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
