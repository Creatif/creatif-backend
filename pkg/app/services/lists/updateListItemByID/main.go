package updateListItemByID

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Main struct {
	model Model
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
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
		return declarations.ListVariable{}, appErrors.NewNotFoundError(err).AddError("updateListItemByID.Logic", nil)
	}
	var list declarations.List
	if res := storage.Gorm().Where("name = ? AND project_id = ? AND locale_id = ?", c.model.ListName, c.model.ProjectID, localeID).Select("id").First(&list); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.ListVariable{}, appErrors.NewNotFoundError(res.Error).AddError("updateListItemByID.Logic", nil)
		}

		return declarations.ListVariable{}, appErrors.NewDatabaseError(res.Error).AddError("updateListItemByID.Logic", nil)
	}

	var existing declarations.ListVariable
	if res := storage.Gorm().Where("id = ? AND list_id = ? AND locale_id = ?", c.model.ItemID, list.ID, localeID).First(&existing); res.Error != nil {
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

func New(model Model) pkg.Job[Model, View, declarations.ListVariable] {
	return Main{model: model}
}
