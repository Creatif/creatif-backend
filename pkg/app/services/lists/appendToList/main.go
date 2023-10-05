package appendToList

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
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

func (c Main) Logic() (declarations.List, error) {
	var list declarations.List
	if err := storage.GetBy((declarations.List{}).TableName(), "name", c.model.Name, &list, "id"); err != nil {
		return declarations.List{}, appErrors.NewNotFoundError(err).AddError("appendToList.Logic", nil)
	}

	listVariables := make([]declarations.ListVariable, len(c.model.Variables))
	for i := 0; i < len(c.model.Variables); i++ {
		v := c.model.Variables[i]
		listVariables[i] = declarations.NewListVariable(list.ID, v.Name, v.Behaviour, v.Metadata, v.Groups, v.Value)
	}

	if res := storage.Gorm().Create(&listVariables); res.Error != nil {
		return declarations.List{}, res.Error
	}

	return list, nil
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

func New(model Model) pkg.Job[Model, View, declarations.List] {
	return Main{model: model}
}