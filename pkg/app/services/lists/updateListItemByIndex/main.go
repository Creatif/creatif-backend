package updateListItemByIndex

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
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
		return declarations.ListVariable{}, appErrors.NewApplicationError(err).AddError("updateListItemByIndex.Logic", nil)
	}

	offset := c.model.ItemIndex
	var existing declarations.ListVariable
	if res := storage.Gorm().
		Raw(fmt.Sprintf(`
			SELECT lv.id
			FROM %s AS lv INNER JOIN %s AS l
			ON l.project_id = ? AND l.name = ? AND lv.list_id = l.id AND l.locale_id = ?
			ORDER BY lv.index ASC
			OFFSET ? LIMIT 1`, (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()), c.model.ProjectID, c.model.ListName, localeID, offset).
		Scan(&existing); res.Error != nil || res.RowsAffected == 0 {
		if res.RowsAffected == 0 {
			return declarations.ListVariable{}, appErrors.NewNotFoundError(res.Error).AddError("updateListItemByIndex.Logic", nil)
		}

		return declarations.ListVariable{}, appErrors.NewDatabaseError(res.Error).AddError("updateListItemByIndex.Logic", nil)
	}

	fmt.Println(existing.ID)

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
		return declarations.ListVariable{}, appErrors.NewApplicationError(res.Error).AddError("updateListItemByIndex.Logic", nil)
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
