package updateMapVariable

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strings"
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
	var m declarations.Map
	if res := storage.Gorm().Where("name = ? AND project_id = ?", c.model.Name, c.model.ProjectID).First(&m); res.Error != nil {
		return LogicResult{}, appErrors.NewNotFoundError(res.Error).AddError("updateMapVariable.Logic", nil)
	}

	pqGroups := pq.StringArray{}
	if c.model.Entry.Groups == nil {
		pqGroups = pq.StringArray{}
	} else {
		for _, v := range c.model.Entry.Groups {
			pqGroups = append(pqGroups, v)
		}
	}

	returningFields := []string{"id", "name", "behaviour", "metadata", "groups", "value", "created_at", "updated_at"}

	var model declarations.MapVariable
	if res := storage.Gorm().Raw(fmt.Sprintf(
		"UPDATE %s SET behaviour = ?, metadata = ?, groups = ?, value = ? WHERE name = ? AND map_id = ? RETURNING %s", (declarations.MapVariable{}).TableName(), strings.Join(returningFields, ",")),
		c.model.Entry.Behaviour,
		c.model.Entry.Metadata,
		pqGroups,
		c.model.Entry.Value,
		c.model.Entry.Name,
		m.ID,
	).Scan(&model); res.Error != nil || res.RowsAffected == 0 {
		return LogicResult{}, appErrors.NewNotFoundError(errors.New("Could not update map")).AddError("updateMapVariable.Logic", nil)
	}

	return LogicResult{
		Map:   m,
		Entry: model,
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
