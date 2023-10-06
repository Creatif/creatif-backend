package switchByIndex

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
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

func (c Main) Logic() (LogicResult, error) {
	var to, from declarations.ListVariable
	if err := storage.Transaction(func(tx *gorm.DB) error {
		source, err := queryVariableByIndex(c.model.ProjectID, c.model.Name, c.model.Source)
		if err != nil {
			return err
		}
		destination, err := queryVariableByIndex(c.model.ProjectID, c.model.Name, c.model.Destination)
		if err != nil {
			return err
		}

		newToVariable, newFromVariable, err := handleUpdate(source, destination)
		if err != nil {
			return err
		}

		to = newToVariable
		from = newFromVariable

		return nil
	}); err != nil {
		return LogicResult{}, appErrors.NewDatabaseError(err)
	}

	return LogicResult{
		To:   to,
		From: from,
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

func New(model Model) pkg.Job[Model, View, LogicResult] {
	return Main{model: model}
}
