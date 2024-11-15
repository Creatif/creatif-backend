package addToList

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
	"github.com/tidwall/sjson"
	"gorm.io/gorm"
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
(m.id = ? OR m.short_id = ?) AND m.project_id = ? AND 
mv.list_id = m.id AND mv.name = ? AND mv.locale_id = ?
`, (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName())

	var entry declarations.ListVariable
	res := storage.Gorm().Raw(sql, c.model.Name, c.model.Name, c.model.ProjectID, c.model.Entry.Name, entryLocaleId).Scan(&entry)
	if res.Error != nil {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Variable with name '%s' and locale '%s' for list with ID '%s' already exists.", c.model.Entry.Name, c.model.Entry.Locale, c.model.Name),
		})
	}

	if res.RowsAffected != 0 {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Variable with name '%s' and locale '%s' for list with ID '%s' already exists.", c.model.Entry.Name, c.model.Entry.Locale, c.model.Name),
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
		return LogicModel{}, appErrors.NewApplicationError(err).AddError("addToList.Logic", nil)
	}

	var m declarations.List
	if res := storage.Gorm().Where(fmt.Sprintf("project_id = ? AND (id = ? OR name = ? OR short_id = ?)"), c.model.ProjectID, c.model.Name, c.model.Name, c.model.Name).Select("ID", "short_id", "name").First(&m); res.Error != nil {
		return LogicModel{}, appErrors.NewNotFoundError(res.Error).AddError("addToList.Logic", nil)
	}

	if c.model.Entry.Groups == nil {
		c.model.Entry.Groups = []string{}
	}

	highestIndex, err := getHighestIndex(m.ID)
	if err != nil {
		return LogicModel{}, appErrors.NewApplicationError(err).AddError("addToMap.Logic", nil)
	}

	variable := declarations.NewListVariable(m.ID, localeID, c.model.Entry.Name, c.model.Entry.Behaviour, c.model.Entry.Metadata, c.model.Entry.Value)
	variable.Index = highestIndex + 1024

	var refs []declarations.Reference
	if transactionError := storage.Transaction(func(tx *gorm.DB) error {
		if len(c.model.ImagePaths) != 0 {
			newValue, err := fileProcessor.UploadFiles(
				c.model.ProjectID,
				c.model.Entry.Value,
				c.model.ImagePaths,
				func(fileSystemFilePath, path, mimeType, extension, fileName string) (string, error) {
					image := declarations.NewFile(
						c.model.ProjectID,
						&variable.ID,
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
			)

			if err != nil {
				return err
			}

			variable.Value = newValue
		}

		if len(c.model.References) > 0 {
			conns := sdk.Map(c.model.References, func(idx int, value shared.Reference) connections.Connection {
				return connections.Connection{
					Path:          value.Name,
					StructureType: value.StructureType,
					VariableID:    value.VariableID,
				}
			})

			for _, c := range conns {
				newValue, err := sjson.DeleteBytes(variable.Value, c.Path)
				if err != nil {
					return err
				}

				variable.Value = newValue
			}

			if err := connections.CheckConnectionsIntegrity(conns); err != nil {
				return err
			}

			created := sdk.Map(conns, func(idx int, value connections.Connection) declarations.Connection {
				return declarations.NewConnection(
					c.model.ProjectID,
					value.Path,
					variable.ID,
					"list",
					value.VariableID,
					value.StructureType,
				)
			})

			if res := tx.Create(&created); res.Error != nil {
				return res.Error
			}

			// 1. transform into connections
			// 2. check integrity
			// 3. save the connections
		}

		if res := tx.Create(&variable); res.Error != nil {
			return errors.New(fmt.Sprintf("List item with name '%s' already exists.", c.model.Entry.Name))
		}

		if len(c.model.Entry.Groups) > 0 {
			newGroup := declarations.NewVariableGroup(variable.ID, c.model.Entry.Groups)
			if res := tx.Create(&newGroup); res.Error != nil {
				return res.Error
			}
		}

		if len(c.model.References) > 0 {
			references, err := shared.CreateDeclarationReferences(c.model.References, m.ID, variable.ID, "list", c.model.ProjectID)
			if err != nil {
				return err
			}

			if res := tx.Create(&references); res.Error != nil {
				return res.Error
			}

			refs = references
		}

		return nil
	}); transactionError != nil {
		return LogicModel{}, appErrors.NewApplicationError(transactionError).AddError("addToList.Logic", nil)
	}

	return LogicModel{
		Variable:   variable,
		References: refs,
		Groups:     c.model.Entry.Groups,
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
