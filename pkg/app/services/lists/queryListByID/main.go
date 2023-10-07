package queryListByID

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
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
	var variable declarations.ListVariable
	res := storage.Gorm().
		Raw(fmt.Sprintf(`
			SELECT lv.id, lv.name, lv.index, lv.short_id, lv.metadata, lv.value, lv.groups, lv.created_at, lv.updated_at
			FROM %s AS lv INNER JOIN %s AS l
			ON l.project_id = ? AND l.name = ? AND lv.list_id = l.id AND lv.id = ?`,
			(declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()), c.model.ProjectID, c.model.Name, c.model.ID).
		Scan(&variable)

	if res.Error != nil {
		return declarations.ListVariable{}, appErrors.NewDatabaseError(res.Error).AddError("queryListByID.Logic", nil)
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, appErrors.NewNotFoundError(res.Error).AddError("queryListByID.Logic", nil)
	}

	return variable, nil
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

func New(model Model) pkg.Job[Model, View, declarations.ListVariable] {
	return Main{model: model}
}
