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
		return nil, appErrors.NewNotFoundError(res.Error).AddError("removeMapEntry.Logic", nil)
	}

	res := storage.Gorm().Where("map_id = ? AND name = ?", m.ID, c.model.EntryName).Delete(&declarations.MapVariable{})
	if res.Error != nil {
		return nil, appErrors.NewNotFoundError(res.Error).AddError("removeMapEntry.Logic", nil)
	}

	if res.RowsAffected == 0 {
		return nil, appErrors.NewNotFoundError(errors.New(fmt.Sprintf("Variable with name '%s' not found.", c.model.EntryName))).AddError("removeMapEntry.Logic", nil)
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
