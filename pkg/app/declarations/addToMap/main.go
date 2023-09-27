package deleteVariable

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
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

func (c Main) Logic() (interface{}, error) {
	var m declarations.Map
	if res := storage.Gorm().Where("name = ? AND project_id = ?", c.model.Name, c.model.ProjectID).Select("ID").First(&m); res.Error != nil {
		return nil, appErrors.NewNotFoundError(res.Error).AddError("addToMap.Logic", nil)
	}

	if c.model.Entry.Groups == nil {
		c.model.Entry.Groups = []string{}
	}

	mapNode := declarations.NewMapVariable(m.ID, c.model.Entry.Name, c.model.Entry.Behaviour, c.model.Entry.Metadata, c.model.Entry.Groups, c.model.Entry.Value)
	if res := storage.Gorm().Create(&mapNode); res.Error != nil {
		return nil, appErrors.NewApplicationError(errors.New(fmt.Sprintf("Entry with name '%s' already exists", c.model.Entry.Name))).AddError("addToMap.Logic", nil)
	}

	return nil, nil
}

func (c Main) Handle() (interface{}, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	if err := c.Authorize(); err != nil {
		return nil, err
	}

	_, err := c.Logic()

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func New(model Model) pkg.Job[Model, interface{}, interface{}] {
	return Main{model: model}
}
