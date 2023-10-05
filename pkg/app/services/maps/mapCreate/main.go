package mapCreate

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
	newMap := declarations.NewMap(c.model.ProjectID, c.model.Name)

	names := make([]map[string]string, 0)
	if err := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&newMap); res.Error != nil {
			return res.Error
		}

		domainEntries := make([]declarations.MapVariable, len(c.model.Entries))
		entries := c.model.Entries
		for i, entry := range entries {
			if entry.Type == "variable" {
				m := entry.Model.(VariableModel)

				domainEntries[i] = declarations.NewMapVariable(
					newMap.ID,
					m.Name,
					m.Behaviour,
					m.Metadata,
					m.Groups,
					m.Value,
				)
			}
		}

		if res := tx.Create(&domainEntries); res.Error != nil {
			return res.Error
		}

		for _, d := range domainEntries {
			if d.ID != "" {
				names = append(names, map[string]string{
					"name": d.Name,
					"type": "variable",
				})
			}
		}

		return nil
	}); err != nil {
		return LogicResult{}, appErrors.NewDatabaseError(err).AddError("mapCreate.Logic", nil)
	}

	return LogicResult{
		ID:        newMap.ID,
		ProjectID: newMap.ProjectID,
		Name:      newMap.Name,
		Variables: names,
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