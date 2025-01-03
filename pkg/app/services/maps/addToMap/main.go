package addToMap

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/app/services/shared/fileProcessor"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
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

	if len(c.model.Entry.Groups) > 0 {
		count, err := shared.ValidateGroupsExist(c.model.ProjectID, c.model.Entry.Groups)
		if err != nil {
			return appErrors.NewValidationError(map[string]string{
				"groupsExist": err.Error(),
			})
		}

		if count+len(c.model.Entry.Groups) > 20 {
			return appErrors.NewValidationError(map[string]string{
				"maximumGroups": fmt.Sprintf("You are trying to add %d more groups but you already have %d assigned to this item. Maximum number of groups per item is 20", len(c.model.Entry.Groups), count),
			})
		}
	}

	entryLocaleId, _ := locales.GetIDWithAlpha(c.model.Entry.Locale)

	sql := fmt.Sprintf(`
SELECT mv.id FROM %s AS mv 
INNER JOIN %s AS m ON 
(m.id = ? OR m.name = ? OR m.short_id = ?) AND m.project_id = ? AND 
mv.map_id = m.id AND mv.name = ? AND mv.locale_id = ?
`, (declarations.MapVariable{}).TableName(), (declarations.Map{}).TableName())

	var entry declarations.MapVariable
	res := storage.Gorm().Raw(sql, c.model.Name, c.model.Name, c.model.Name, c.model.ProjectID, c.model.Entry.Name, entryLocaleId).Scan(&entry)
	if res.Error != nil {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Variable with name '%s' and locale '%s' for map with ID '%s' already exists.", c.model.Entry.Name, c.model.Entry.Locale, c.model.Name),
		})
	}

	if res.RowsAffected != 0 {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Variable with name '%s' and locale '%s' for map with ID '%s' already exists.", c.model.Entry.Name, c.model.Entry.Locale, c.model.Name),
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

func (c Main) Logic() (LogicModel, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Entry.Locale)
	if err != nil {
		return LogicModel{}, appErrors.NewApplicationError(err).AddError("addToMap.Logic", nil)
	}

	var m declarations.Map
	if res := storage.Gorm().Where(fmt.Sprintf("project_id = ? AND (id = ? OR short_id = ?)"), c.model.ProjectID, c.model.Name, c.model.Name).Select("ID", "name").First(&m); res.Error != nil {
		return LogicModel{}, appErrors.NewNotFoundError(res.Error).AddError("addToMap.Logic", nil)
	}

	if c.model.Entry.Groups == nil {
		c.model.Entry.Groups = []string{}
	}

	highestIndex, err := getHighestIndex(m.ID)
	if err != nil {
		return LogicModel{}, appErrors.NewApplicationError(err).AddError("addToMap.Logic", nil)
	}

	variable := declarations.NewMapVariable(m.ID, localeID, c.model.Entry.Name, c.model.Entry.Behaviour, c.model.Entry.Metadata, c.model.Entry.Value)
	variable.Index = highestIndex + 1024
	var conns []declarations.Connection
	if transactionError := storage.Transaction(func(tx *gorm.DB) error {
		if len(c.model.ImagePaths) != 0 {
			newValue, err := fileProcessor.UploadFiles(
				c.model.ProjectID,
				c.model.Entry.Value,
				c.model.ImagePaths,
				func(fileSystemFilePath, path, mimeType, extension, fileName string) (string, error) {
					image := declarations.NewFile(
						c.model.ProjectID,
						nil,
						&variable.ID,
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
			)

			if err != nil {
				return err
			}

			variable.Value = newValue
		}

		if res := tx.Create(&variable); res.Error != nil {
			return errors.New(fmt.Sprintf("Map with name '%s' already exists.", c.model.Entry.Name))
		}

		if len(c.model.Entry.Groups) > 0 {
			newGroup := declarations.NewVariableGroup(variable.ID, c.model.Entry.Groups)
			if res := tx.Create(&newGroup); res.Error != nil {
				return res.Error
			}
		}

		if len(c.model.Connections) > 0 {
			newValue, newConnections, err := connections.CreateConnections(
				tx,
				c.model.ProjectID,
				variable.ID,
				"map",
				c.model.Connections,
				variable.Value,
			)

			if err != nil {
				return err
			}

			variable.Value = newValue

			if res := tx.Create(&newConnections); res.Error != nil {
				return res.Error
			}

			conns = newConnections
		}

		return nil
	}); transactionError != nil {
		errString := transactionError.Error()
		splt := strings.Split(errString, ":")
		if len(splt) == 2 {
			return LogicModel{}, appErrors.NewValidationError(map[string]string{
				splt[0]: splt[1],
			})
		}

		return LogicModel{}, appErrors.NewApplicationError(transactionError)
	}

	return LogicModel{
		Variable:    variable,
		Connections: conns,
		Groups:      c.model.Entry.Groups,
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, View, LogicModel] {
	return Main{model: model, auth: auth}
}
