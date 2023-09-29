package updateVariable

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
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

func (c Main) Logic() (declarations.Variable, error) {
	var existing declarations.Variable
	storage.Gorm().Where("name = ? AND project_id = ?", c.model.Name, c.model.ProjectID).First(&existing)
	if err := storage.GetBy((declarations.Variable{}).TableName(), "name", c.model.Name, &existing, "id"); err != nil {
		return declarations.Variable{}, appErrors.NewNotFoundError(err).AddError("updateVariable.Logic", nil)
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

	var updated declarations.Variable
	if res := storage.Gorm().Model(&updated).Clauses(clause.Returning{Columns: []clause.Column{
		{Name: "id"},
		{Name: "project_id"},
		{Name: "name"},
		{Name: "behaviour"},
		{Name: "metadata"},
		{Name: "value"},
		{Name: "groups"},
		{Name: "created_at"},
		{Name: "updated_at"},
	}}).Where("id = ?", existing.ID).Select(c.model.Fields).Updates(existing); res.Error != nil {
		return declarations.Variable{}, appErrors.NewApplicationError(res.Error).AddError("updateVariable.Logic", nil)
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

	return newView(model), nil
}

func New(model Model) pkg.Job[Model, View, declarations.Variable] {
	return Main{model: model}
}
