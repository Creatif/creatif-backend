package queryListByID

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"gorm.io/gorm"
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
	var list declarations.List
	res := storage.Gorm().Where("project_id = ? AND name = ?", c.model.ProjectID, c.model.Name).Select("ID").First(&list)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.ListVariable{}, appErrors.NewNotFoundError(res.Error).AddError("queryListByIndex.Logic", nil)
		}

		return declarations.ListVariable{}, appErrors.NewDatabaseError(res.Error).AddError("queryListByIndex.Logic", nil)
	}

	var variable declarations.ListVariable
	res = storage.Gorm().Where("list_id = ? AND id = ?", list.ID, c.model.ID).First(&variable)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.ListVariable{}, appErrors.NewNotFoundError(res.Error).AddError("queryListByIndex.Logic", nil)
		}

		return declarations.ListVariable{}, appErrors.NewDatabaseError(res.Error).AddError("queryListByIndex.Logic", nil)
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