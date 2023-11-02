package replaceListItem

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("replaceListItem", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("replaceListItem", "Validated")
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
		c.logBuilder.Add("replaceListItem", err.Error())
		return declarations.ListVariable{}, appErrors.NewApplicationError(err).AddError("replaceListItem.Logic", nil)
	}

	if c.model.Variable.Groups == nil {
		c.model.Variable.Groups = []string{}
	}

	listAndItem, err := queryListAndItem(c.model.ProjectID, c.model.Name, c.model.ID, c.model.ShortID, c.model.ItemID, c.model.ItemShortID)
	if err != nil {
		c.logBuilder.Add("replaceListItem", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return declarations.ListVariable{}, appErrors.NewNotFoundError(err).AddError("replaceListItem.Logic", nil)
		}
	}

	listItem := declarations.NewListVariable(listAndItem.ListID, localeID, c.model.Variable.Name, c.model.Variable.Behaviour, c.model.Variable.Metadata, c.model.Variable.Groups, c.model.Variable.Value)
	listItem.Index = listAndItem.ItemIndex
	if err := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Where("list_id = ? AND id = ?", listAndItem.ListID, listAndItem.ItemID).Delete(&declarations.ListVariable{}); res.Error != nil {
			return res.Error
		}

		res := tx.Model(&listItem).Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
			{Name: "name"},
			{Name: "short_id"},
			{Name: "behaviour"},
			{Name: "metadata"},
			{Name: "value"},
			{Name: "groups"},
			{Name: "created_at"},
			{Name: "updated_at"},
		}}).Create(&listItem)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	}); err != nil {
		c.logBuilder.Add("replaceListItem", err.Error())
		return declarations.ListVariable{}, appErrors.NewDatabaseError(err).AddError("createList.Logic", nil)
	}

	return listItem, nil
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.ListVariable] {
	logBuilder.Add("replaceListItem", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
