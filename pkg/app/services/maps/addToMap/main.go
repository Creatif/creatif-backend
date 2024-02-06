package addToMap

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
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
		if err := validateGroupsExist(c.model.ProjectID, c.model.Entry.Groups); err != nil {
			return appErrors.NewValidationError(map[string]string{
				"exists": err.Error(),
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
		c.logBuilder.Add("addToList", err.Error())
		return LogicModel{}, appErrors.NewApplicationError(err).AddError("addToList.Logic", nil)
	}

	var m declarations.Map
	if res := storage.Gorm().Where(fmt.Sprintf("project_id = ? AND (id = ? OR short_id = ?)"), c.model.ProjectID, c.model.Name, c.model.Name).Select("ID", "name").First(&m); res.Error != nil {
		c.logBuilder.Add("addToList", res.Error.Error())
		return LogicModel{}, appErrors.NewNotFoundError(res.Error).AddError("addToList.Logic", nil)
	}

	if c.model.Entry.Groups == nil {
		c.model.Entry.Groups = []string{}
	}

	variable := declarations.NewMapVariable(m.ID, localeID, c.model.Entry.Name, c.model.Entry.Behaviour, c.model.Entry.Metadata, c.model.Entry.Groups, c.model.Entry.Value)
	var refs []declarations.Reference
	if transactionError := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&variable); res.Error != nil {
			c.logBuilder.Add("addToList", res.Error.Error())

			return errors.New(fmt.Sprintf("Map with name '%s' already exists.", c.model.Entry.Name))
		}

		if len(c.model.References) > 0 {
			references, err := shared.CreateDeclarationReferences(c.model.References, m.ID, variable.ID, "map", c.model.ProjectID)
			if err != nil {
				return err
			}

			tx.Create(&references)

			refs = references
		}

		return nil
	}); transactionError != nil {
		return LogicModel{}, transactionError
	}

	return LogicModel{
		Variable:   variable,
		References: refs,
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
