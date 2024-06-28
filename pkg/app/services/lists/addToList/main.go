package addToList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/app/services/shared/fileProcessor"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("addToList", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	c.logBuilder.Add("addToList", "Validated.")

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
		c.logBuilder.Add("addToList", err.Error())
		return LogicModel{}, appErrors.NewApplicationError(err).AddError("addToList.Logic", nil)
	}

	var m declarations.List
	if res := storage.Gorm().Where(fmt.Sprintf("project_id = ? AND (id = ? OR name = ? OR short_id = ?)"), c.model.ProjectID, c.model.Name, c.model.Name, c.model.Name).Select("ID", "short_id", "name").First(&m); res.Error != nil {
		c.logBuilder.Add("addToList", res.Error.Error())
		return LogicModel{}, appErrors.NewNotFoundError(res.Error).AddError("addToList.Logic", nil)
	}

	if c.model.Entry.Groups == nil {
		c.model.Entry.Groups = []string{}
	}

	highestIndex, err := getHighestIndex(m.ID)
	if err != nil {
		c.logBuilder.Add("addToMap", err.Error())
		return LogicModel{}, appErrors.NewApplicationError(err).AddError("addToMap.Logic", nil)
	}

	variable := declarations.NewListVariable(m.ID, localeID, c.model.Entry.Name, c.model.Entry.Behaviour, c.model.Entry.Metadata, c.model.Entry.Value)
	variable.Index = highestIndex + 1024

	var refs []declarations.Reference
	if transactionError := storage.Transaction(func(tx *gorm.DB) error {
		newValue, _, err := fileProcessor.UploadFiles(
			c.model.ProjectID,
			c.model.Entry.Value,
			c.model.ImagePaths,
			func(fileSystemFilePath, path, mimeType, extension string) (string, error) {
				image := declarations.NewImage(
					c.model.ProjectID,
					variable.ID,
					fileSystemFilePath,
					path,
					mimeType,
					extension,
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
		// before creating a variable, create images, ensures failure if image upload fails

		if res := tx.Create(&variable); res.Error != nil {
			c.logBuilder.Add("addToList", res.Error.Error())

			return errors.New(fmt.Sprintf("List item with name '%s' already exists.", c.model.Entry.Name))
		}

		if len(c.model.Entry.Groups) > 0 {
			groups := make([]declarations.VariableGroup, len(c.model.Entry.Groups))
			for i, g := range c.model.Entry.Groups {
				groups[i] = declarations.NewVariableGroup(g, variable.ID, c.model.Entry.Groups)
			}

			if res := tx.Create(&groups); res.Error != nil {
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, LogicModel] {
	logBuilder.Add("addToList", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
